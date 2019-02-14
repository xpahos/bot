package settings

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/storage"
)

const (
	NOTIFY_ON  = "NOTIFY_ON"
	NOTIFY_OFF = "NOTIFY_OFF"
)

func PrepareCommandMenu(db *sql.DB, msg *tgbotapi.MessageConfig, action map[string]int, username *string) {
	isOn := storage.UsersIsOnNotifications(db, username)

	notificationButtonText := "Включить оповещения"
	notificationButtonAction := NOTIFY_ON

	if isOn {
		notificationButtonText = "Выключить оповещения"
		notificationButtonAction = NOTIFY_OFF
	}

	menu := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(notificationButtonText, notificationButtonAction),
		},
	)

	msg.Text = ctx.SettingsActionMenuText
	msg.ReplyMarkup = menu
	action[*username] = ctx.ActionManageSettingsActionMenu
}

func ProcessInlineSettingsMenu(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int) {
	logger.Infof("Settings menu")
	username := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	switch message {
	case NOTIFY_ON:
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
		if storage.UsersUpdateNotifications(db, &username, true) {
			msg.Text = "Уведомления включены"
		} else {
			msg.Text = "Неудалось включить уведомления"
		}
		bot.Send(msg)
		actionStateMap[username] = ctx.ActionNone
	case NOTIFY_OFF:
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
		if storage.UsersUpdateNotifications(db, &username, false) {
			msg.Text = "Уведомления выключены"
		} else {
			msg.Text = "Неудалось выключить уведомления"
		}
		bot.Send(msg)
		actionStateMap[username] = ctx.ActionNone
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.SettingsActionMenuText, message),
	)

	bot.Send(msg)
}