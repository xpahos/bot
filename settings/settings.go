package settings

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
	"github.com/xpahos/bot/chat"
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/storage"
	"strconv"
)

const (
	NOTIFY_ON  = "NOTIFY_ON"
	NOTIFY_OFF = "NOTIFY_OFF"
	TIME_START = "TIME_START"
	TIME_END   = "TIME_END"
	TIME_ZONE  = "TIME_ZONE"
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
		chat.Send(bot, msg)
		actionStateMap[username] = ctx.ActionNone
	case NOTIFY_OFF:
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
		if storage.UsersUpdateNotifications(db, &username, false) {
			msg.Text = "Уведомления выключены"
		} else {
			msg.Text = "Неудалось выключить уведомления"
		}
		chat.Send(bot, msg)
		actionStateMap[username] = ctx.ActionNone
	case TIME_START:
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx.SettingsTimeStartText)
		chat.Send(bot, msg)
		actionStateMap[username] = ctx.ActionManageSettingsTimeStart
	case TIME_END:
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx.SettingsTimeEndText)
		chat.Send(bot, msg)
		actionStateMap[username] = ctx.ActionManageSettingsTimeEnd
	case TIME_ZONE:
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx.SettingsTimeZoneText)
		chat.Send(bot, msg)
		actionStateMap[username] = ctx.ActionManageSettingsTimeZone
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.SettingsActionMenuText, message),
	)

	chat.Send(bot, msg)
}

func ProcessKeyboardSettingsTime(db *sql.DB, msg *tgbotapi.MessageConfig, update *tgbotapi.Update, actionStateMap map[string]int) {
	username := update.Message.From.UserName
	message, err := strconv.Atoi(update.Message.Text)

	switch actionStateMap[username] {
	case ctx.ActionManageSettingsTimeStart:
		if err != nil || (message < 0 || message > 24) {
			msg.Text = "Неверный формат времени. Допустимые значения от 0 до 24"
		} else {
			if storage.UsersUpdateNotificationsTime(db, &username, message, ctx.SettingsTimeStartUpdate) {
				msg.Text = fmt.Sprintf("Уведомления будут приходить с %d", message)
			} else {
				msg.Text = "Неудалось изменить время начала уведомлений"
			}
		}
	case ctx.ActionManageSettingsTimeEnd:
		if err != nil || (message < 0 || message > 24) {
			msg.Text = "Неверный формат времени. Допустимые значения от 0 до 24"
		} else {
			if storage.UsersUpdateNotificationsTime(db, &username, message, ctx.SettingsTimeEndUpdate) {
				msg.Text = fmt.Sprintf("Уведомления будут приходить до %d", message)
			} else {
				msg.Text = "Неудалось изменить время окончания уведомлений"
			}
		}
	case ctx.ActionManageSettingsTimeZone:
		if err != nil || (message < -12 || message > 14) {
			msg.Text = "Неверный формат временной зоны. Допустимые значения от -12 до 14"
		} else {
			if storage.UsersUpdateNotificationsTime(db, &username, message, ctx.SettingsTimeZoneUpdate) {
				msg.Text = fmt.Sprintf("Уведомления будут приходить во временной зоне %d", message)
			} else {
				msg.Text = "Неудалось изменить временную зону"
			}
		}
	default:
		msg.Text = "Неверная комманда"
	}
	actionStateMap[username] = ctx.ActionNone
}
