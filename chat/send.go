package chat

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

// Send checks result of sending message to bot and logs errors.
func Send(bot *tgbotapi.BotAPI, message tgbotapi.Chattable) {
	_, err := bot.Send(message)
	if err != nil {
		// TODO(serejkus): might be good idea to add previous stack frame info.
		logger.Errorf("failed to send message %+v: %v", message, err)
	}
}
