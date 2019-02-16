package main

import (
	"flag"
	"github.com/xpahos/bot/chat"
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/duty"
	"github.com/xpahos/bot/form"
	"github.com/xpahos/bot/helpers"
	"github.com/xpahos/bot/settings"
	"github.com/xpahos/bot/storage"
	"github.com/xpahos/bot/users"
	"os"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var logPath = flag.String("log", "bot.log", "Log path")
var verbose = flag.Bool("verbose", false, "Print info level logs to stdout")

func DispatchMessage(db *sql.DB, bot *tgbotapi.BotAPI, action map[string]int, formProblemMap map[string]*ctx.FormProblemStruct, notifyReport chan<- ctx.NotifyNewReportStruct, updates <-chan tgbotapi.Update) {
	for update := range updates {
		if update.CallbackQuery != nil {
			username := update.CallbackQuery.From.UserName

			actionStateIdx := action[username]
			logger.Infof("Inline: %s %s %d", update.CallbackQuery.Data, username, actionStateIdx)
			switch actionStateIdx {
			case ctx.ActionManageFormActionMenu:
				form.ProcessInlineFormActionMenu(db, bot, &update, action)
			case ctx.ActionManageFormWindBlowing:
				form.ProcessInlineFormWindBlowing(db, bot, &update, action)
			case ctx.ActionManageFormWeatherTrend:
				form.ProcessInlineFormWeatherTrend(db, bot, &update, action)
			case ctx.ActionManageFormWeatherChangesAdditional:
				form.ProcessInlineFormWeatherChangesAdditional(bot, &update, action)
			case ctx.ActionManageFormProblemMenu:
				form.ProcessInlineFormProblemMenu(bot, &update, action, formProblemMap)
			case ctx.ActionManageFormProblemType:
				form.ProcessInlineFormType(bot, &update, action, formProblemMap)
			case ctx.ActionManageFormProblemLocation:
				form.ProcessInlineFormLocations(bot, &update, action, formProblemMap)
			case ctx.ActionManageFormProblemLikelyHood:
				form.ProcessInlineFormLikelyHood(bot, &update, action, formProblemMap)
			case ctx.ActionManageFormProblemSize:
				form.ProcessInlineFormSize(db, bot, &update, action, formProblemMap)
			case ctx.ActionManageFormAvalancheForecastAlp:
				form.ProcessInlineFormAvalanche(db, bot, &update, action, ctx.AlpForecast, nil)
			case ctx.ActionManageFormAvalancheForecastTree:
				form.ProcessInlineFormAvalanche(db, bot, &update, action, ctx.TreeForecast, nil)
			case ctx.ActionManageFormAvalancheForecastBTree:
				form.ProcessInlineFormAvalanche(db, bot, &update, action, ctx.BTreeForecast, nil)
			case ctx.ActionManageFormAvalancheConfidenceAlp:
				form.ProcessInlineFormAvalanche(db, bot, &update, action, ctx.AlpConfidence, nil)
			case ctx.ActionManageFormAvalancheConfidenceTree:
				form.ProcessInlineFormAvalanche(db, bot, &update, action, ctx.TreeConfidence, nil)
			case ctx.ActionManageFormAvalancheConfidenceBTree:
				form.ProcessInlineFormAvalanche(db, bot, &update, action, ctx.BTreeConfidence, notifyReport)
			case ctx.ActionManageUserActionMenu:
				users.ProcessInlineUserActionMenu(db, bot, &update, action)
			case ctx.ActionManageDutyActionMenu:
				duty.ProcessInlineDutyActionMenu(db, bot, &update, action)
			case ctx.ActionManageDutyAdd:
				duty.ProcessInlineDutyEdit(db, bot, &update, action, true)
			case ctx.ActionManageDutyDelete:
				duty.ProcessInlineDutyEdit(db, bot, &update, action, false)
			case ctx.ActionManageSettingsActionMenu:
				settings.ProcessInlineSettingsMenu(db, bot, &update, action)
			default:
				logger.Errorf("unknown action state index %d", actionStateIdx)
			}
		} else {
			username := update.Message.From.UserName
			message := update.Message.Text

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			//msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			logger.Infof("[%s] %s %v", username, message, update.Message.IsCommand())
			if !update.Message.IsCommand() {
				switch action[username] {
				case ctx.ActionManageFormHN24:
					form.ProcessKeyboardFormHN24(db, &msg, &update, action)
				case ctx.ActionManageFormH2D:
					form.ProcessKeyboardFormH2D(db, &msg, &update, action)
				case ctx.ActionManageFormHST:
					form.ProcessKeyboardFormHST(db, &msg, &update, action)
				case ctx.ActionManageFormWeatherChanges:
					form.ProcessKeyboardFormWeatherChanges(db, &msg, &update, action)
				case ctx.ActionManageFormComments:
					form.ProcessKeyboardFormComments(db, &msg, &update, action)
				case ctx.ActionManageFormDeclineComment:
					form.ProcessKeyboardFormDecline(db, &msg, &update, action)
				case ctx.ActionManageFormArchive:
					form.ProcessKeyboardFormArchive(db, &msg, &update, action)
				case ctx.ActionManageUserAdd:
					users.ProcessKeyboardUserAdd(db, &msg, &update, action)
				case ctx.ActionManageUserDelete:
					users.ProcessKeyboardUserDelete(db, &msg, &update, action)
				case ctx.ActionManageSettingsTimeStart, ctx.ActionManageSettingsTimeEnd, ctx.ActionManageSettingsTimeZone:
					settings.ProcessKeyboardSettingsTime(db, &msg, &update, action)
				default:
					msg.Text = "Неизвестная команда"
					action[username] = ctx.ActionNone
				}
			} else {
				if action[username] != ctx.ActionNone &&
					action[username] != ctx.ActionManageFormActionMenu &&
					action[username] != ctx.ActionManageUserActionMenu &&
					action[username] != ctx.ActionManageDutyActionMenu &&
					action[username] != ctx.ActionManageSettingsActionMenu {
					msg.Text = "Предыдущее действие не завершено"
				} else {
					switch update.Message.Command() {
					case "help":
						msg.ParseMode = "markdown"
						msg.Text = ctx.HelpText
					case "form":
						form.PrepareCommandForm(db, &msg, action, &username)
					case "confirm":
						form.PrepareCommandConfirm(db, &msg, action, &username)
					case "decline":
						form.PrepareCommandDecline(db, &msg, action, &username)
					case "archive":
						form.PrepareCommandArchive(db, &msg, action, &username)
					case "report":
						form.PrepareCommandReport(db, &msg, action)
					case "users":
						users.PrepareCommandUsers(&msg, action, &username)
					case "duty":
						duty.PrepareCommandDuty(&msg, action, &username)
					case "settings":
						settings.PrepareCommandMenu(db, &msg, action, &username)
					default:
						msg.Text = "Неизвестная команда"
					}
				}
			}
			chat.Send(bot, msg)
		}
	}
}

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

	var (
		actionStateMap    = make(map[string]int)
		trustedUsersCache = make(map[string]bool)
		formProblemMap    = make(map[string]*ctx.FormProblemStruct)
	)

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

	// Routine for async processing updates
	updatesChannel := make(chan tgbotapi.Update, 1)
	go DispatchMessage(db, bot, actionStateMap, formProblemMap, notifyReport, updatesChannel)

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
			chat.Send(bot, msg)
			if err != nil {
				logger.Errorf("failed to send message %+v: %v", msg, err)
			}
			continue
		}

		updatesChannel <- update

	}
}
