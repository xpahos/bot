package helpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xpahos/bot/chat"
)

func ShowNextQuestionInline(bot *tgbotapi.BotAPI, update *tgbotapi.Update, text string, menu *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	if menu != nil {
		msg.ReplyMarkup = menu
	} else {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}
	chat.Send(bot, msg)
}

func ShowNextQuestionReply(bot *tgbotapi.BotAPI, update *tgbotapi.Update, text string, menu *tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	if menu != nil {
		msg.ReplyMarkup = menu
	} else {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}
	chat.Send(bot, msg)
}
