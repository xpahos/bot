package util

import (
    "storage"

    "fmt"
    "database/sql"

    "github.com/google/logger"
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

func NotifyNewReport(username <-chan string, bot *tgbotapi.BotAPI, db *sql.DB) {
    for userName := range username {
        logger.Infof("Sending notification about report from %v", userName)
        for _, chatId := range storage.UsersGetChatIdList(db) {
            msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("Готов отчет пользователя %v", userName))
            bot.Send(msg)
        }
    }
}
