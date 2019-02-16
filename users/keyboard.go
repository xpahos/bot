package users

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/storage"
)

func ProcessKeyboardUserAdd(db *sql.DB, msg *tgbotapi.MessageConfig, update *tgbotapi.Update, action map[string]int) {
	username := update.Message.From.UserName
	message := update.Message.Text

	if storage.UsersAddOne(db, &message) {
		msg.Text = "Пользователь добавлен"
		logger.Infof("User %s added user %s", username, message)
	} else {
		msg.Text = "Пользователь уже существует или его имя длинее 255 символов"
	}
	action[username] = ctx.ActionNone
}

func ProcessKeyboardUserDelete(db *sql.DB, msg *tgbotapi.MessageConfig, update *tgbotapi.Update, action map[string]int) {
	username := update.Message.From.UserName
	message := update.Message.Text

	if storage.UsersDeleteOne(db, &message) {
		msg.Text = "Пользователь удален"
		logger.Infof("User %s deleted user %s", username, message)
	} else {
		msg.Text = "Не удалось удалить пользователя"
	}
	action[username] = ctx.ActionNone
}
