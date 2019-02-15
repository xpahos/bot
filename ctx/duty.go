package ctx

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type DutyInfo struct {
	User string
	Date string
}

const (
	DutyActionMenuText = "Выберите действие с расписанием: "
	DutyDateText       = "Выберите дату: "
)

var DutyActionMenu = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Добавить", "ADD"),
		tgbotapi.NewInlineKeyboardButtonData("Просмотреть", "VIEW"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "DELETE"),
	},
)
