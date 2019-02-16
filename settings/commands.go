package settings

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/storage"
	"strconv"
)

func PrepareCommandMenu(db *sql.DB, msg *tgbotapi.MessageConfig, action map[string]int, username *string) {
	info, err := storage.UsersGetOneNotificationInfo(db, username)

	if err != nil {
		msg.Text = "Неизвестная ошибка"
		action[*username] = ctx.ActionNone
		return
	}

	notificationButtonText := "Включить оповещения"
	notificationButtonAction := NOTIFY_ON

	if info.IsOn {
		notificationButtonText = "Выключить оповещения"
		notificationButtonAction = NOTIFY_OFF
	}

	notificationButtonTimeStart := fmt.Sprintf("С %d", info.TimeStart)
	notificationButtonTimeEnd := fmt.Sprintf("До %d", info.TimeEnd)
	notificationButtonTimeZone := "Зона "

	if info.TimeZone > 0 {
		notificationButtonTimeZone += "+"
	}
	notificationButtonTimeZone += strconv.Itoa(info.TimeZone)

	menu := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(notificationButtonText, notificationButtonAction),
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(notificationButtonTimeStart, TIME_START),
			tgbotapi.NewInlineKeyboardButtonData(notificationButtonTimeEnd, TIME_END),
			tgbotapi.NewInlineKeyboardButtonData(notificationButtonTimeZone, TIME_ZONE),
		},
	)

	msg.Text = ctx.SettingsActionMenuText
	msg.ReplyMarkup = menu
	action[*username] = ctx.ActionManageSettingsActionMenu
}
