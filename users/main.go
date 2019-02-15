package users

import (
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/helpers"
	"github.com/xpahos/bot/storage"

	"database/sql"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

func ProcessInlineUserActionMenu(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Info("UserActionMenu")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "ADD":
		showMenu = false

		actionStateMap[userName] = ctx.ActionManageUserAdd
		helpers.ShowNextQuestionReply(bot, update, ctx.UserAddText, nil)
	case "DELETE":
		showMenu = false

		users := storage.UsersGetList(db)
		buttons := []tgbotapi.KeyboardButton{}
		for _, user := range users {
			buttons = append(buttons, tgbotapi.NewKeyboardButton(user))
		}

		rows := [][]tgbotapi.KeyboardButton{}
		buttonLen := len(buttons)
		for i := 0; i < buttonLen; i = i + 4 {
			end := helpers.Min(i+4, buttonLen)
			rows = append(rows, tgbotapi.NewKeyboardButtonRow(buttons[i:end]...))
		}
		userListMenu := tgbotapi.NewReplyKeyboard(rows...)

		actionStateMap[userName] = ctx.ActionManageUserDelete
		helpers.ShowNextQuestionReply(bot, update, ctx.UsersDeleteText, &userListMenu)
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.UsersActionMenuText, message),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormActionMenu
	}

	bot.Send(msg)
}
