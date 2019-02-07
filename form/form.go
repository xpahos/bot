package form

import (
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/storage"

	"database/sql"
	"fmt"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

func ShowNextQuestion(bot *tgbotapi.BotAPI, update *tgbotapi.Update, text string, menu *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	if menu != nil {
		msg.ReplyMarkup = menu
	} else {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}
	bot.Send(msg)
}

func ProcessInlineFormActionMenu(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Infof("Action")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data
	now := time.Now()

	showMenu := true

	switch message {
	case "ADD":
		showMenu = false

		if storage.FormInitRecord(db, &now, &userName) {
			actionStateMap[userName] = ctx.ActionManageFormWindBlowing
			ShowNextQuestion(bot, update, ctx.FormWindBlowingText, &ctx.FormWindBlowing)
		} else {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Отчет уже существует")
			bot.Send(msg)
			actionStateMap[userName] = ctx.ActionNone
		}
	case "DELETE":
		showMenu = false
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
		if storage.FormDeleteRecord(db, &now) {
			msg.Text = "Отчет за сегодняшний день удален"
			logger.Infof("User %v deleted report for %v", userName, now.Format("02 Jan 2006"))
		} else {
			msg.Text = "Неудалось удалить отчет за сегодня"
		}
		bot.Send(msg)
		actionStateMap[userName] = ctx.ActionNone
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.FormActionMenuText, message),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormActionMenu
	}

	bot.Send(msg)
}

func ProcessInlineFormWindBlowing(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "LOW", "MEDIUM", "HIGH":
		showMenu = false
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.FormWindBlowingText, message),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormWindBlowing
	}

	bot.Send(msg)

	if !showMenu {
		now := time.Now()
		storage.FormUpdateWindBlowing(db, &now, &message)
		actionStateMap[userName] = ctx.ActionManageFormWeatherTrend
		ShowNextQuestion(bot, update, ctx.FormWeatherTrendText, &ctx.FormWeatherTrend)
	}

}

func ProcessInlineFormWeatherTrend(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Infof("123")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "WORSE", "SAME", "BETTER":
		showMenu = false
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.FormWeatherTrendText, message),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormWeatherTrend
	}

	bot.Send(msg)

	if !showMenu {
		now := time.Now()
		storage.FormUpdateWeatherTrend(db, &now, &message)
		actionStateMap[userName] = ctx.ActionManageFormHN24
		ShowNextQuestion(bot, update, ctx.FormHN24Text, nil)
	}

}

func ProcessInlineFormAvalancheForecast(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int, zone int, notify chan<- string) {
	logger.Infof("123")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "1", "2", "3", "4", "5":
		showMenu = false
	}

	var text string

	switch zone {
	case ctx.Alp:
		text = fmt.Sprintf("%v %v", ctx.FormAvalancheForecastAlpText, message)
	case ctx.Tree:
		text = fmt.Sprintf("%v %v", ctx.FormAvalancheForecastTreeText, message)
	case ctx.BTree:
		text = fmt.Sprintf("%v %v", ctx.FormAvalancheForecastBTreeText, message)
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		text,
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormAvalancheForecast
	}

	bot.Send(msg)

	if !showMenu {
		now := time.Now()
		storage.FormUpdateAvalanche(db, &now, &message, zone)
		if zone == ctx.Alp {
			actionStateMap[userName] = ctx.ActionManageFormAvalancheForecastTree
			ShowNextQuestion(bot, update, ctx.FormAvalancheForecastTreeText, &ctx.FormAvalancheForecast)
		} else if zone == ctx.Tree {
			actionStateMap[userName] = ctx.ActionManageFormAvalancheForecastBTree
			ShowNextQuestion(bot, update, ctx.FormAvalancheForecastBTreeText, &ctx.FormAvalancheForecast)
		} else if zone == ctx.BTree {
			storage.FormComplete(db, &now)
			ShowNextQuestion(bot, update, ctx.FormCompletedText, nil)
			if notify != nil {
				notify <- userName
			}
			actionStateMap[userName] = ctx.ActionNone
		}
	}
}
