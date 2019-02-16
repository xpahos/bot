package users

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xpahos/bot/ctx"
)

func PrepareCommandUsers(msg *tgbotapi.MessageConfig, action map[string]int, username *string) {
	msg.Text = ctx.UsersActionMenuText
	msg.ReplyMarkup = ctx.UsersActionMenu
	action[*username] = ctx.ActionManageUserActionMenu
}
