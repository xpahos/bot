package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/duty"
	"github.com/xpahos/bot/form"
	"github.com/xpahos/bot/storage"
	"github.com/xpahos/bot/users"
	"github.com/xpahos/bot/util"
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

	notifyReport := make(chan string, 1)
	go util.NotifyNewReport(notifyReport, bot, db)

	var actionStateMap map[string]int = make(map[string]int)
	var trustedUsersCache map[string]bool = make(map[string]bool)
	var formProblemMap map[string]*ctx.FormProblemStruct = make(map[string]*ctx.FormProblemStruct)

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		var userName string
		var chatId int64
		var msgId int
		if update.CallbackQuery != nil {
			userName = update.CallbackQuery.From.UserName
			chatId = update.CallbackQuery.Message.Chat.ID
			msgId = -1
		} else if update.Message != nil {
			userName = update.Message.From.UserName
			chatId = update.Message.Chat.ID
			msgId = update.Message.MessageID
		}

		if !trustedUsersCache[userName] && !storage.UsersCheckTrusted(db, trustedUsersCache, &update) {
			msg := tgbotapi.NewMessage(chatId, "Вы не авторизованы для выполнения этой операции")
			if msgId != -1 {
				msg.ReplyToMessageID = msgId
			}
			bot.Send(msg)
			continue
		}

		if update.CallbackQuery != nil {
			logger.Infof("Inline: %v %v %v", update.CallbackQuery.Data, userName, actionStateMap[userName])
			switch actionStateMap[userName] {
			case ctx.ActionManageFormActionMenu:
				form.ProcessInlineFormActionMenu(db, bot, &update, actionStateMap)
			case ctx.ActionManageFormWindBlowing:
				form.ProcessInlineFormWindBlowing(db, bot, &update, actionStateMap)
			case ctx.ActionManageFormWeatherTrend:
				form.ProcessInlineFormWeatherTrend(db, bot, &update, actionStateMap)
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
			}

			continue
		}

		message := update.Message.Text
		now := time.Now()

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		//msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		logger.Infof("[%v] %v %v", userName, message, update.Message.IsCommand())
		if !update.Message.IsCommand() {
			switch actionStateMap[userName] {
			case ctx.ActionManageFormHN24:
				if storage.FormUpdateHN24(db, &now, &message) {
					msg.Text = "Показания доски H2D(цифрами или 0)"
					actionStateMap[userName] = ctx.ActionManageFormH2D
				} else {
					msg.Text = ctx.FormHN24Text
					actionStateMap[userName] = ctx.ActionManageFormHN24
				}
			case ctx.ActionManageFormH2D:
				if storage.FormUpdateH2D(db, &now, &message) {
					msg.Text = "Показания доски HST(цифрами или 0)"
					actionStateMap[userName] = ctx.ActionManageFormHST
				} else {
					msg.Text = "Показания доски H2D(цифрами или 0)"
					actionStateMap[userName] = ctx.ActionManageFormH2D
				}
			case ctx.ActionManageFormHST:
				if storage.FormUpdateHST(db, &now, &message) {
					msg.Text = "Ощутимые изменения(выберите или введите произвольный вариант)"
					msg.ReplyMarkup = ctx.FormWeatherChanges
					actionStateMap[userName] = ctx.ActionManageFormWeatherChanges
				} else {
					msg.Text = "Показания доски HST(цифрами или 0)"
					actionStateMap[userName] = ctx.ActionManageFormHST
				}
			case ctx.ActionManageFormWeatherChanges:
				if storage.FormUpdateWeatherChanges(db, &now, &message) {
					msg.Text = ctx.FormProblemMenuText
					msg.ReplyMarkup = ctx.YesNoMenu
					actionStateMap[userName] = ctx.ActionManageFormProblemMenu
				} else {
					msg.Text = "Ощутимые изменения(выбрать или ввести произвольный вариант)"
					msg.ReplyMarkup = ctx.FormWeatherChanges
					actionStateMap[userName] = ctx.ActionManageFormWeatherChanges
				}
			case ctx.ActionManageFormComments:
				if storage.FormUpdateComments(db, &now, &message) {
					msg.Text = ctx.FormAvalancheForecastAlpText
					msg.ReplyMarkup = ctx.FormAvalancheForecast
					actionStateMap[userName] = ctx.ActionManageFormAvalancheForecastAlp
				} else {
					msg.Text = ctx.FormCommentsText
					actionStateMap[userName] = ctx.ActionManageFormComments
				}
			case ctx.ActionManageFormDeclineComment:
				if storage.FormDecline(db, &now, &userName, &message) {
					msg.Text = "Комментарий добавлен"
					actionStateMap[userName] = ctx.ActionNone
				} else {
					msg.Text = "Неудалось внести данные"
					actionStateMap[userName] = ctx.ActionNone
				}
			case ctx.ActionManageUserAdd:
				if storage.UsersAddOne(db, &message) {
					msg.Text = "Пользователь добавлен"
					logger.Infof("User %v added user %v", userName, message)
				} else {
					msg.Text = "Пользователь уже существует или его имя длинее 255 символов"
				}
				actionStateMap[userName] = ctx.ActionNone
			case ctx.ActionManageUserDelete:
				if storage.UsersDeleteOne(db, &message) {
					msg.Text = "Пользователь удален"
					logger.Infof("User %v deleted user %v", userName, message)
				} else {
					msg.Text = "Не удалось удалить пользователя"
				}
				actionStateMap[userName] = ctx.ActionNone
			default:
				actionStateMap[userName] = ctx.ActionNone
				msg.Text = "Неизвестная команда"
			}
		} else {
			if actionStateMap[userName] != ctx.ActionNone &&
				actionStateMap[userName] != ctx.ActionManageFormActionMenu &&
				actionStateMap[userName] != ctx.ActionManageUserActionMenu &&
				actionStateMap[userName] != ctx.ActionManageDutyActionMenu {
				msg.Text = "Предыдущее действие не завершено"
			} else {
				switch update.Message.Command() {
				case "help":
					msg.ParseMode = "markdown"
					msg.Text = ctx.HelpText
				case "form":
					duty := storage.DutyGetOne(db, &now)
					if duty != userName {
						msg.Text = fmt.Sprintf("Сегодня дежурный %v", duty)
						actionStateMap[userName] = ctx.ActionNone
					} else {
						msg.Text = ctx.FormActionMenuText
						msg.ReplyMarkup = ctx.FormActionMenu
						actionStateMap[userName] = ctx.ActionManageFormActionMenu
					}
				case "confirm":
					duty := storage.DutyGetOne(db, &now)
					if duty == userName {
						msg.Text = "Вы не можете подтверждать свои отчеты"
					} else {
						if storage.FormIsCompleted(db, &now) {
							storage.FormConfirm(db, &now, &userName)
							msg.Text = "Отчет подтвержден"
						} else {
							msg.Text = "Отчет еще не закончен"
						}
					}
					actionStateMap[userName] = ctx.ActionNone
				case "decline":
					duty := storage.DutyGetOne(db, &now)
					if duty == userName {
						msg.Text = "Вы не можете подтверждать свои отчеты"
						actionStateMap[userName] = ctx.ActionNone
					} else {
						if storage.FormIsCompleted(db, &now) {
							msg.Text = "Введите доплнительный комментарий"
							actionStateMap[userName] = ctx.ActionManageFormDeclineComment
						} else {
							msg.Text = "Отчет еще не закончен"
							actionStateMap[userName] = ctx.ActionNone
						}
					}
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
					actionStateMap[userName] = ctx.ActionManageUserActionMenu
				case "duty":
					msg.Text = ctx.DutyActionMenuText
					msg.ReplyMarkup = ctx.DutyActionMenu
					actionStateMap[userName] = ctx.ActionManageDutyActionMenu
				default:
					msg.Text = "Неизвестная команда"
				}
			}
		}
		bot.Send(msg)
	}
}
