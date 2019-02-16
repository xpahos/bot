package duty

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xpahos/bot/ctx"
)

func PrepareCommandDuty(msg *tgbotapi.MessageConfig, action map[string]int, username *string)  {
	msg.Text = ctx.DutyActionMenuText
	msg.ReplyMarkup = ctx.DutyActionMenu
	action[*username] = ctx.ActionManageDutyActionMenu
}
