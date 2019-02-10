package form

import (
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/storage"
	"github.com/xpahos/bot/helpers"

	"database/sql"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

func ProcessInlineFormActionMenu(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Info("Action")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data
	now := time.Now()

	showMenu := true

	switch message {
	case "ADD":
		showMenu = false

		if storage.FormInitRecord(db, &now, &userName) {
			actionStateMap[userName] = ctx.ActionManageFormWindBlowing
			helpers.ShowNextQuestionInline(bot, update, ctx.FormWindBlowingText, &ctx.FormWindBlowing)
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
			logger.Infof("User %s deleted report for %s", userName, now.Format("02 Jan 2006"))
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
	logger.Info("WindBlowing")
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
		helpers.ShowNextQuestionInline(bot, update, ctx.FormWeatherTrendText, &ctx.FormWeatherTrend)
	}
}

func ProcessInlineFormWeatherTrend(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Info("WeatherTrend")
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
		helpers.ShowNextQuestionInline(bot, update, ctx.FormHN24Text, nil)
	}
}

func ProcessInlineFormWeatherChangesAdditional(bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
    logger.Info("WeatherChangesAdditional")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.FormWeatherChangesAdditionalText, message),
	)

	bot.Send(msg)

    switch message {
        case "Y":
		    actionStateMap[userName] = ctx.ActionManageFormWeatherChanges
            helpers.ShowNextQuestionReply(bot, update, ctx.FormWeatherChangesText, &ctx.FormWeatherChanges)
        default:
		    actionStateMap[userName] = ctx.ActionManageFormProblemMenu
            helpers.ShowNextQuestionInline(bot, update, ctx.FormProblemMenuText, &ctx.YesNoMenu)
    }
}

func ProcessInlineFormAvalancheForecast(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int, zone int, notify chan<- string) {
	logger.Info("AvalancheForecast")
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
			helpers.ShowNextQuestionInline(bot, update, ctx.FormAvalancheForecastTreeText, &ctx.FormAvalancheForecast)
		} else if zone == ctx.Tree {
			actionStateMap[userName] = ctx.ActionManageFormAvalancheForecastBTree
			helpers.ShowNextQuestionInline(bot, update, ctx.FormAvalancheForecastBTreeText, &ctx.FormAvalancheForecast)
		} else if zone == ctx.BTree {
			storage.FormComplete(db, &now)
			helpers.ShowNextQuestionInline(bot, update, ctx.FormCompletedText, nil)
			if notify != nil {
				notify <- userName
			}
			actionStateMap[userName] = ctx.ActionNone
		}
	}
}
