package duty

import (
	"github.com/xpahos/bot/chat"
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/helpers"
	"github.com/xpahos/bot/storage"

	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

func ProcessInlineDutyActionMenu(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Infof("ProcessInlineDutyActionMenu")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	if message == "VIEW" {
		showMenu = false

		start := time.Now().AddDate(0, 0, -16)
		end := time.Now().AddDate(0, 0, 16)

		duties := storage.DutyGetList(db, &start, &end)

		results := make([]string, 1 + len(duties)) // 1 for header
		results[0] = fmt.Sprintf("Дежурства с *%v* по *%v*\n", start.Format("02 Jan 2006"), end.Format("02 Jan 2006"))
		for i, duty := range duties {
			results[i+1] = fmt.Sprintf("*%v* - `%v`", duty.Date, duty.User)
		}

		result := strings.Join(results, "\n")
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, result)
		msg.ParseMode = "markdown"
		chat.Send(bot, msg)

		actionStateMap[userName] = ctx.ActionNone
	} else {
		occupied := make(map[string]bool)
		start := time.Now()
		end := time.Now().AddDate(0, 0, 15)

		for _, duty := range storage.DutyGetList(db, &start, &end) {
			occupied[duty.Date] = true
		}

		after := time.Now()
		buttons := make([]tgbotapi.InlineKeyboardButton, 16)
		for i := range buttons {
			var userString string
			dbString := after.Format("2006-01-02")

			if occupied[dbString] {
				userString = "[" + after.Format("02 Jan 2006") + "]"
			} else {
				userString = after.Format("02 Jan 2006")
			}

			buttons[i] = tgbotapi.NewInlineKeyboardButtonData(userString, dbString)
			after = after.AddDate(0, 0, 1)
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
			helpers.ShowNextQuestionInline(bot, update, ctx.DutyDateText, &dateListMenu)
		case "DELETE":
			showMenu = false

			actionStateMap[userName] = ctx.ActionManageDutyDelete
			helpers.ShowNextQuestionInline(bot, update, ctx.DutyDateText, &dateListMenu)
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

	chat.Send(bot, msg)
}

func ProcessInlineDutyEdit(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int, add bool) {
	logger.Infof("InlineDutyEdit")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data
	now := time.Now()

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	if add {
		// If date is not occupied then add duty
		_, err := storage.DutyGetOne(db, &now)
		if (err == nil || err == sql.ErrNoRows) && storage.DutyAddOne(db, &message, &userName) {
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
	chat.Send(bot, msg)

	editMsg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.DutyDateText, message),
	)
	chat.Send(bot, editMsg)

	actionStateMap[userName] = ctx.ActionNone
}
