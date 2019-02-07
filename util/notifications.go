package util

import (
	"github.com/xpahos/bot/storage"

	"database/sql"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

func NotifyNewReport(username <-chan string, bot *tgbotapi.BotAPI, db *sql.DB) {
	for userName := range username {
		logger.Infof("Sending notification about report from %v", userName)
		for _, chatID := range storage.UsersGetChatIDList(db) {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Готов отчет пользователя %v", userName))
			bot.Send(msg)
		}
	}
}
