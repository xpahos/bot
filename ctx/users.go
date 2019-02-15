package ctx

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var TrustedUsers = map[string]bool{"xpahos": true, "Khalmax": true, "snow_niki": true}

const (
	UsersActionMenuText    = "Выберите действие с пользователем: "
	UserAddText            = "Введите имя пользователя: "
	UsersAddConfirmText    = "Пользователь добавлен"
	UsersDeleteText        = "Выберите пользователя. Выведены первые 16 пользователей. Если в списке нет нужного пользователя укажите его вручную."
	UsersDeleteConfirmText = "Пользователь удален"
)

type UsersNotificationDurationStruct struct {
	Username  string
	ChatID    int64
	TimeStart int
	TimeEnd   int
	TimeZone  int
}

var UsersActionMenu = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Добавить", "ADD"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "DELETE"),
	},
)
