package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/duty"
	"github.com/xpahos/bot/form"
	"github.com/xpahos/bot/helpers"
	"github.com/xpahos/bot/settings"
	"github.com/xpahos/bot/storage"
	"github.com/xpahos/bot/users"
)

var logPath = flag.String("log", "bot.log", "Log path")
var verbose = flag.Bool("verbose", false, "Print info level logs to stdout")

func main() {
	flag.Parse()

	fd, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
		return
	}
	defer fd.Close()
	defer logger.Init("Chat Bot", *verbose, true, fd).Close()

	db, err := sql.Open("mysql", os.Getenv("DB"))
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatalf("Failed to ping database: %v", err)
		return
	}

	db.SetConnMaxLifetime(10 * time.Second)

	err = storage.InitDB(db)
	if err != nil {
		logger.Fatalf("Failed to init database: %v", err)
		return
	}

	logger.Info("Starting bot")

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		logger.Fatalf("Bot auth problem: %v", err)
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1

	// Run all goroutines for notifications
	notifyReport := make(chan ctx.NotifyNewReportStruct, 1)
	notifyNoDuty := make(chan int64, 1)
	notifyNoReport := make(chan int64, 1)
	go helpers.NotifyNewReport(notifyReport, bot)
	go helpers.NotifyNoDuty(notifyNoDuty, bot)
	go helpers.NotifyNoReport(notifyNoReport, bot)

	// Run all cronjobs
	go helpers.CronJobCheckDuty(notifyNoDuty, db)
	go helpers.CronJobCheckReport(notifyNoReport, db)

	var (
		actionStateMap    = make(map[string]int)
		trustedUsersCache = make(map[string]bool)
		formProblemMap    = make(map[string]*ctx.FormProblemStruct)
	)

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		logger.Fatalf("failed to get updates: %v", err)
		return
	}
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		var (
			username string
			chatID   int64
			msgID    int
		)
		if update.CallbackQuery != nil {
			username = update.CallbackQuery.From.UserName
			chatID = update.CallbackQuery.Message.Chat.ID
			msgID = -1
		} else if update.Message != nil {
			username = update.Message.From.UserName
			chatID = update.Message.Chat.ID
			msgID = update.Message.MessageID
		}

		if !trustedUsersCache[username] && !storage.UsersCheckTrusted(db, trustedUsersCache, &update) {
			msg := tgbotapi.NewMessage(chatID, "Вы не авторизованы для выполнения этой операции")
			if msgID != -1 {
				msg.ReplyToMessageID = msgID
			}
			bot.Send(msg)
			continue
		}

		if update.CallbackQuery != nil {
			actionStateIdx := actionStateMap[username]
			logger.Infof("Inline: %s %s %d", update.CallbackQuery.Data, username, actionStateIdx)
			switch actionStateIdx {
			case ctx.ActionManageFormActionMenu:
				form.ProcessInlineFormActionMenu(db, bot, &update, actionStateMap)
			case ctx.ActionManageFormWindBlowing:
				form.ProcessInlineFormWindBlowing(db, bot, &update, actionStateMap)
			case ctx.ActionManageFormWeatherTrend:
				form.ProcessInlineFormWeatherTrend(db, bot, &update, actionStateMap)
			case ctx.ActionManageFormWeatherChangesAdditional:
				form.ProcessInlineFormWeatherChangesAdditional(bot, &update, actionStateMap)
			case ctx.ActionManageFormProblemMenu:
				form.ProcessInlineFormProblemMenu(bot, &update, actionStateMap, formProblemMap)
			case ctx.ActionManageFormProblemType:
				form.ProcessInlineFormType(bot, &update, actionStateMap, formProblemMap)
			case ctx.ActionManageFormProblemLocation:
				form.ProcessInlineFormLocations(bot, &update, actionStateMap, formProblemMap)
			case ctx.ActionManageFormProblemLikelyHood:
				form.ProcessInlineFormLikelyHood(bot, &update, actionStateMap, formProblemMap)
			case ctx.ActionManageFormProblemSize:
				form.ProcessInlineFormSize(db, bot, &update, actionStateMap, formProblemMap)
			case ctx.ActionManageFormAvalancheForecastAlp:
				form.ProcessInlineFormAvalancheForecast(db, bot, &update, actionStateMap, ctx.Alp, nil)
			case ctx.ActionManageFormAvalancheForecastTree:
				form.ProcessInlineFormAvalancheForecast(db, bot, &update, actionStateMap, ctx.Tree, nil)
			case ctx.ActionManageFormAvalancheForecastBTree:
				form.ProcessInlineFormAvalancheForecast(db, bot, &update, actionStateMap, ctx.BTree, notifyReport)
			case ctx.ActionManageUserActionMenu:
				users.ProcessInlineUserActionMenu(db, bot, &update, actionStateMap)
			case ctx.ActionManageDutyActionMenu:
				duty.ProcessInlineDutyActionMenu(db, bot, &update, actionStateMap)
			case ctx.ActionManageDutyAdd:
				duty.ProcessInlineDutyEdit(db, bot, &update, actionStateMap, true)
			case ctx.ActionManageDutyDelete:
				duty.ProcessInlineDutyEdit(db, bot, &update, actionStateMap, false)
			case ctx.ActionManageSettingsActionMenu:
				settings.ProcessInlineSettingsMenu(db, bot, &update, actionStateMap)
			default:
				logger.Infof("unknown action state index %d", actionStateIdx)
			}

			continue
		}

		message := update.Message.Text
		now := time.Now()

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		//msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		logger.Infof("[%s] %s %v", username, message, update.Message.IsCommand())
		if !update.Message.IsCommand() {
			switch actionStateMap[username] {
			case ctx.ActionManageFormHN24:
				if storage.FormUpdateHN24(db, &now, &message) {
					msg.Text = "Показания доски H2D(цифрами или 0)"
					actionStateMap[username] = ctx.ActionManageFormH2D
				} else {
					msg.Text = ctx.FormHN24Text
					actionStateMap[username] = ctx.ActionManageFormHN24
				}
			case ctx.ActionManageFormH2D:
				if storage.FormUpdateH2D(db, &now, &message) {
					msg.Text = "Показания доски HST(цифрами или 0)"
					actionStateMap[username] = ctx.ActionManageFormHST
				} else {
					msg.Text = "Показания доски H2D(цифрами или 0)"
					actionStateMap[username] = ctx.ActionManageFormH2D
				}
			case ctx.ActionManageFormHST:
				if storage.FormUpdateHST(db, &now, &message) {
					msg.Text = ctx.FormWeatherChangesText
					msg.ReplyMarkup = ctx.FormWeatherChanges
					actionStateMap[username] = ctx.ActionManageFormWeatherChanges
				} else {
					msg.Text = "Показания доски HST(цифрами или 0)"
					actionStateMap[username] = ctx.ActionManageFormHST
				}
			case ctx.ActionManageFormWeatherChanges:
				if storage.FormUpdateWeatherChanges(db, &now, &username, &message) {
					msg.Text = ctx.FormWeatherChangesAdditionalText
					msg.ReplyMarkup = ctx.YesNoMenu
					actionStateMap[username] = ctx.ActionManageFormWeatherChangesAdditional
				} else {
					msg.Text = ctx.FormWeatherChangesText
					msg.ReplyMarkup = ctx.FormWeatherChanges
					actionStateMap[username] = ctx.ActionManageFormWeatherChanges
				}
			case ctx.ActionManageFormComments:
				if storage.FormUpdateComments(db, &now, &message) {
					msg.Text = ctx.FormAvalancheForecastAlpText
					msg.ReplyMarkup = ctx.FormAvalancheForecast
					actionStateMap[username] = ctx.ActionManageFormAvalancheForecastAlp
				} else {
					msg.Text = ctx.FormCommentsText
					actionStateMap[username] = ctx.ActionManageFormComments
				}
			case ctx.ActionManageFormDeclineComment:
				if storage.FormDecline(db, &now, &username, &message) {
					msg.Text = "Комментарий добавлен"
					actionStateMap[username] = ctx.ActionNone
				} else {
					msg.Text = "Неудалось внести данные"
					actionStateMap[username] = ctx.ActionNone
				}
			case ctx.ActionManageUserAdd:
				if storage.UsersAddOne(db, &message) {
					msg.Text = "Пользователь добавлен"
					logger.Infof("User %s added user %s", username, message)
				} else {
					msg.Text = "Пользователь уже существует или его имя длинее 255 символов"
				}
				actionStateMap[username] = ctx.ActionNone
			case ctx.ActionManageUserDelete:
				if storage.UsersDeleteOne(db, &message) {
					msg.Text = "Пользователь удален"
					logger.Infof("User %s deleted user %s", username, message)
				} else {
					msg.Text = "Не удалось удалить пользователя"
				}
				actionStateMap[username] = ctx.ActionNone
			case ctx.ActionManageFormArchive:
				date, err := time.Parse("02 Jan 2006", message)

				if err != nil {
					msg.Text = "Неверный формат даты"
				} else {
					msg.ParseMode = "markdown"
					if storage.FormIsCompleted(db, &date) {
						msg.Text = form.GenerateTextReport(db, &date)
					} else {
						msg.Text = "Отчет еще не закончен"
					}
				}

				actionStateMap[username] = ctx.ActionNone
			case ctx.ActionManageSettingsTimeStart, ctx.ActionManageSettingsTimeEnd, ctx.ActionManageSettingsTimeZone:
				settings.ProcessKeyboardSettingsTime(db, &msg, &update, actionStateMap)
			default:
				msg.Text = "Неизвестная команда"
				actionStateMap[username] = ctx.ActionNone
			}
		} else {
			if actionStateMap[username] != ctx.ActionNone &&
				actionStateMap[username] != ctx.ActionManageFormActionMenu &&
				actionStateMap[username] != ctx.ActionManageUserActionMenu &&
				actionStateMap[username] != ctx.ActionManageDutyActionMenu &&
				actionStateMap[username] != ctx.ActionManageSettingsActionMenu {
				msg.Text = "Предыдущее действие не завершено"
			} else {
				switch update.Message.Command() {
				case "help":
					msg.ParseMode = "markdown"
					msg.Text = ctx.HelpText
				case "form":
					duty, err := storage.DutyGetOne(db, &now)
					if err == nil {
						msg.Text = "Не выбран дежурный"
						actionStateMap[username] = ctx.ActionNone
					} else {
						if duty != username {
							msg.Text = fmt.Sprintf("Сегодня дежурный %s", duty)
							actionStateMap[username] = ctx.ActionNone
						} else {
							msg.Text = ctx.FormActionMenuText
							msg.ReplyMarkup = ctx.FormActionMenu
							actionStateMap[username] = ctx.ActionManageFormActionMenu
						}
					}
				case "confirm":
					duty, err := storage.DutyGetOne(db, &now)
					if err == nil {
						msg.Text = "Не выбран дежурный"
					} else {
						if duty == username {
							msg.Text = "Вы не можете подтверждать свои отчеты"
						} else {
							if storage.FormIsCompleted(db, &now) {
								storage.FormConfirm(db, &now, &username)
								msg.Text = "Отчет подтвержден"
							} else {
								msg.Text = "Отчет еще не закончен"
							}
						}
					}
					actionStateMap[username] = ctx.ActionNone
				case "decline":
					duty, err := storage.DutyGetOne(db, &now)
					if err == nil {
						msg.Text = "Не выбран дежурный"
						actionStateMap[username] = ctx.ActionNone
					} else {
						if duty == username {
							msg.Text = "Вы не можете подтверждать свои отчеты"
							actionStateMap[username] = ctx.ActionNone
						} else {
							if storage.FormIsCompleted(db, &now) {
								msg.Text = "Введите доплнительный комментарий"
								actionStateMap[username] = ctx.ActionManageFormDeclineComment
							} else {
								msg.Text = "Отчет еще не закончен"
								actionStateMap[username] = ctx.ActionNone
							}
						}
					}
				case "archive":
					after := time.Now().AddDate(0, 0, -16)
					buttons := []tgbotapi.KeyboardButton{}
					for i := 0; i < 16; i++ {
						after = after.AddDate(0, 0, 1)
						buttons = append(buttons, tgbotapi.NewKeyboardButton(after.Format("02 Jan 2006")))
					}

					dateListMenu := tgbotapi.NewReplyKeyboard(
						tgbotapi.NewKeyboardButtonRow(buttons[0:4]...),
						tgbotapi.NewKeyboardButtonRow(buttons[4:8]...),
						tgbotapi.NewKeyboardButtonRow(buttons[8:12]...),
						tgbotapi.NewKeyboardButtonRow(buttons[12:16]...),
					)

					msg.Text = "Выберите дату или введите в свободной форме"
					msg.ReplyMarkup = dateListMenu
					actionStateMap[username] = ctx.ActionManageFormArchive
				case "report":
					msg.ParseMode = "markdown"
					if storage.FormIsCompleted(db, &now) {
						msg.Text = form.GenerateTextReport(db, &now)
					} else {
						msg.Text = "Отчет еще не закончен"
					}
				case "users":
					msg.Text = ctx.UsersActionMenuText
					msg.ReplyMarkup = ctx.UsersActionMenu
					actionStateMap[username] = ctx.ActionManageUserActionMenu
				case "duty":
					msg.Text = ctx.DutyActionMenuText
					msg.ReplyMarkup = ctx.DutyActionMenu
					actionStateMap[username] = ctx.ActionManageDutyActionMenu
				case "settings":
					settings.PrepareCommandMenu(db, &msg, actionStateMap, &username)
				default:
					msg.Text = "Неизвестная команда"
				}
			}
		}
		bot.Send(msg)
	}
}
