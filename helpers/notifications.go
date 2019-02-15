package helpers

import (
	"fmt"
	"github.com/xpahos/bot/chat"

	"github.com/xpahos/bot/ctx"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

func NotifyNewReport(notifies <-chan ctx.NotifyNewReportStruct, bot *tgbotapi.BotAPI) {
	for notify := range notifies {
		logger.Infof("Sending notification about report from %s", notify.Username)
		msg := tgbotapi.NewMessage(notify.ChatID, fmt.Sprintf("Готов отчет пользователя %s", notify.Username))
		chat.Send(bot, msg)
	}
}

func NotifyNoDuty(notifies <-chan int64, bot *tgbotapi.BotAPI) {
	for chatID := range notifies {
		logger.Infof("Sending notification about no duty %d", chatID)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Пожалуйста, выберите завтрашнего дежурного"))
		chat.Send(bot, msg)
	}
}

func NotifyNoReport(notifies <-chan int64, bot *tgbotapi.BotAPI) {
	for chatID := range notifies {
		logger.Infof("Sending notification about no report %d", chatID)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Пожалуйста, заполните форму"))
		chat.Send(bot, msg)
	}
}
