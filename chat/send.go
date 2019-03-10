package chat

import (
	"fmt"
	"runtime"
	
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

// Send checks result of sending message to bot and logs errors.
func Send(bot *tgbotapi.BotAPI, message tgbotapi.Chattable) {
	_, err := bot.Send(message)
	if err != nil {
		logger.Errorf("%sfailed to send message %+v: %v", caller(), message, err)
	}
}

func caller() string {
	_, file, line, ok := runtime.Caller(2) // 2: 1 for Send, 1 for Send's caller
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d: ", file, line)
}
