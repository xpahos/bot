package duty

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
	}
	bot.Send(msg)
}

func ProcessInlineDutyActionMenu(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Infof("123")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	if message == "VIEW" {
		showMenu = false

		start := time.Now().AddDate(0, 0, -16)
		end := time.Now().AddDate(0, 0, 16)

		duties := storage.DutyGetList(db, &start, &end)

		result := fmt.Sprintf("Дежурства с *%v* по *%v*\n\n", start.Format("02 Jan 2006"), end.Format("02 Jan 2006"))
		for _, duty := range duties {
			result = fmt.Sprintf("%v*%v* - %v\n", result, duty.Date, duty.User)
		}

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, result)
		msg.ParseMode = "markdown"
		bot.Send(msg)

		actionStateMap[userName] = ctx.ActionNone
	} else {
		occupied := make(map[string]bool)
		start := time.Now().AddDate(0, 0, 1)
		end := time.Now().AddDate(0, 0, 16)

		for _, duty := range storage.DutyGetList(db, &start, &end) {
			occupied[duty.Date] = true
		}

		after := time.Now()
		buttons := []tgbotapi.InlineKeyboardButton{}
		for i := 0; i < 16; i++ {
			after = after.AddDate(0, 0, 1)

			var userString string
			dbString := after.Format("2006-01-02")

			if occupied[dbString] {
				userString = "[" + after.Format("02 Jan 2006") + "]"
			} else {
				userString = after.Format("02 Jan 2006")
			}

			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(userString, dbString))
		}

		dateListMenu := tgbotapi.NewInlineKeyboardMarkup(
			buttons[0:4],
			buttons[4:8],
			buttons[8:12],
			buttons[12:16],
		)

		switch message {
		case "ADD":
			showMenu = false

			actionStateMap[userName] = ctx.ActionManageDutyAdd
			ShowNextQuestion(bot, update, ctx.DutyDateText, &dateListMenu)
		case "DELETE":
			showMenu = false

			actionStateMap[userName] = ctx.ActionManageDutyDelete
			ShowNextQuestion(bot, update, ctx.DutyDateText, &dateListMenu)
		}
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.DutyActionMenuText, message),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.DutyActionMenu
	}

	bot.Send(msg)
}

func ProcessInlineDutyEdit(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int, add bool) {
	logger.Infof("123")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	if add {
		if storage.DutyAddOne(db, &message, &userName) {
			msg.Text = "Дежурство добавлено"
			logger.Infof("User %v is duty on %v", userName, message)
		} else {
			msg.Text = "Произошла ошибка добавления дежурства. Неправильный формат даты или дата уже занята"
		}
	} else {
		if storage.DutyDeleteOne(db, &message) {
			msg.Text = "Дежурство удалено"
			logger.Infof("User %v is out of duty on %v", userName, message)
		} else {
			msg.Text = "Произошла ошибка удаления дежурства. Неправильный формат даты или дата уже удалена"
		}
	}
	bot.Send(msg)

	editMsg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.DutyActionMenuText, message),
	)
	bot.Send(editMsg)

	actionStateMap[userName] = ctx.ActionNone
}
