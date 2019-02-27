package form

import (
	"github.com/xpahos/bot/chat"
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/helpers"
	"github.com/xpahos/bot/storage"

	"database/sql"
	"fmt"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

func ProcessInlineFormActionMenu(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Info("Action")
	username := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data
	now := time.Now()

	showMenu := true

	switch message {
	case "ADD":
		showMenu = false

		if storage.FormInitRecord(db, &now, &username) {
			actionStateMap[username] = ctx.ActionManageFormWindBlowing
			helpers.ShowNextQuestionInline(bot, update, ctx.FormWindBlowingText, &ctx.FormWindBlowing)
		} else {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Отчет уже существует")
			chat.Send(bot, msg)
			actionStateMap[username] = ctx.ActionNone
		}
	case "DELETE":
		showMenu = false
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
		if storage.FormDeleteRecord(db, &now) {
			msg.Text = "Отчет за сегодняшний день удален"
			logger.Infof("User %s deleted report for %s", username, now.Format("02 Jan 2006"))
		} else {
			msg.Text = "Не удалось удалить отчет за сегодня"
		}
		chat.Send(bot, msg)
		actionStateMap[username] = ctx.ActionNone
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%s %s", ctx.FormActionMenuText, message),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormActionMenu
	}

	chat.Send(bot, msg)
}

func ProcessInlineFormWindBlowing(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Info("WindBlowing")
	username := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "LOW", "MEDIUM", "HIGH":
		showMenu = false
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%s %s", ctx.FormWindBlowingText, ctx.FormWindBlowingMappingText[message]),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormWindBlowing
	}

	chat.Send(bot, msg)

	if !showMenu {
		now := time.Now()
		storage.FormUpdateWindBlowing(db, &now, &message)
		actionStateMap[username] = ctx.ActionManageFormWeatherTrend
		helpers.ShowNextQuestionInline(bot, update, ctx.FormWeatherTrendText, &ctx.FormWeatherTrend)
	}
}

func ProcessInlineFormWeatherTrend(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Info("WeatherTrend")
	username := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "WORSE", "SAME", "BETTER":
		showMenu = false
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%s %s", ctx.FormWeatherTrendText, ctx.FormWeatherTrendMappingText[message]),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormWeatherTrend
	}

	chat.Send(bot, msg)

	if !showMenu {
		now := time.Now()
		storage.FormUpdateWeatherTrend(db, &now, &message)
		actionStateMap[username] = ctx.ActionManageFormHN24
		helpers.ShowNextQuestionInline(bot, update, ctx.FormHN24Text, nil)
	}
}

func ProcessInlineFormWeatherChangesAdditional(bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Info("WeatherChangesAdditional")
	username := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%s %s", ctx.FormWeatherChangesAdditionalText, message),
	)

	chat.Send(bot, msg)

	switch message {
	case "Y":
		actionStateMap[username] = ctx.ActionManageFormWeatherChanges
		helpers.ShowNextQuestionReply(bot, update, ctx.FormWeatherChangesText, &ctx.FormWeatherChanges)
	default:
		actionStateMap[username] = ctx.ActionManageFormProblemMenu
		helpers.ShowNextQuestionInline(bot, update, ctx.FormProblemMenuText, &ctx.YesNoMenu)
	}
}

func ProcessInlineFormAvalanche(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int, zone int, notify chan<- ctx.NotifyNewReportStruct) {
	logger.Info("AvalancheForecast")
	username := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "1", "2", "3", "4", "5":
		showMenu = false
	}

	var text string

	switch zone {
	case ctx.AlpForecast:
		text = fmt.Sprintf("%s %s", ctx.FormAvalancheForecastAlpText, message)
	case ctx.TreeForecast:
		text = fmt.Sprintf("%s %s", ctx.FormAvalancheForecastTreeText, message)
	case ctx.BTreeForecast:
		text = fmt.Sprintf("%s %s", ctx.FormAvalancheForecastBTreeText, message)
	case ctx.AlpConfidence:
		text = fmt.Sprintf("%s %s", ctx.FormAvalancheConfidenceAlpText, message)
	case ctx.TreeConfidence:
		text = fmt.Sprintf("%s %s", ctx.FormAvalancheConfidenceTreeText, message)
	case ctx.BTreeConfidence:
		text = fmt.Sprintf("%s %s", ctx.FormAvalancheConfidenceBTreeText, message)
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		text,
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormAvalancheForecast
	}

	chat.Send(bot, msg)

	if !showMenu {
		now := time.Now()
		storage.FormUpdateAvalanche(db, &now, &message, zone)
		switch zone {
		case ctx.AlpForecast:
			actionStateMap[username] = ctx.ActionManageFormAvalancheConfidenceAlp
			helpers.ShowNextQuestionInline(bot, update, ctx.FormAvalancheConfidenceAlpText, &ctx.FormAvalancheForecast)
		case ctx.AlpConfidence:
			actionStateMap[username] = ctx.ActionManageFormAvalancheForecastTree
			helpers.ShowNextQuestionInline(bot, update, ctx.FormAvalancheForecastTreeText, &ctx.FormAvalancheForecast)
		case ctx.TreeForecast:
			actionStateMap[username] = ctx.ActionManageFormAvalancheConfidenceTree
			helpers.ShowNextQuestionInline(bot, update, ctx.FormAvalancheConfidenceTreeText, &ctx.FormAvalancheForecast)
		case ctx.TreeConfidence:
			actionStateMap[username] = ctx.ActionManageFormAvalancheForecastBTree
			helpers.ShowNextQuestionInline(bot, update, ctx.FormAvalancheForecastBTreeText, &ctx.FormAvalancheForecast)
		case ctx.BTreeForecast:
			actionStateMap[username] = ctx.ActionManageFormAvalancheConfidenceBTree
			helpers.ShowNextQuestionInline(bot, update, ctx.FormAvalancheConfidenceBTreeText, &ctx.FormAvalancheForecast)
		case ctx.BTreeConfidence:
			storage.FormComplete(db, &now)
			if notify != nil {
				for _, chatID := range storage.UsersGetChatIDList(db) {
					notify <- ctx.NotifyNewReportStruct{Username: username, ChatID: chatID}
				}
			}
			helpers.ShowNextQuestionInline(bot, update, ctx.FormCompletedText, nil)
			actionStateMap[username] = ctx.ActionNone
		}
	}
}
