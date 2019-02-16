package form

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/storage"
	"time"
)

func ProcessKeyboardFormHN24(db *sql.DB, msg *tgbotapi.MessageConfig, update *tgbotapi.Update, action map[string]int) {
	username := update.Message.From.UserName
	message := update.Message.Text
	now := time.Now()

	if storage.FormUpdateHN24(db, &now, &message) {
		msg.Text = "Показания доски H2D(цифрами или 0)"
		action[username] = ctx.ActionManageFormH2D
	} else {
		msg.Text = ctx.FormHN24Text
		action[username] = ctx.ActionManageFormHN24
	}
}

func ProcessKeyboardFormH2D(db *sql.DB, msg *tgbotapi.MessageConfig, update *tgbotapi.Update, action map[string]int) {
	username := update.Message.From.UserName
	message := update.Message.Text
	now := time.Now()

	if storage.FormUpdateH2D(db, &now, &message) {
		msg.Text = "Показания доски HST(цифрами или 0)"
		action[username] = ctx.ActionManageFormHST
	} else {
		msg.Text = "Показания доски H2D(цифрами или 0)"
		action[username] = ctx.ActionManageFormH2D
	}
}

func ProcessKeyboardFormHST(db *sql.DB, msg *tgbotapi.MessageConfig, update *tgbotapi.Update, action map[string]int) {
	username := update.Message.From.UserName
	message := update.Message.Text
	now := time.Now()

	if storage.FormUpdateHST(db, &now, &message) {
		msg.Text = ctx.FormWeatherChangesText
		msg.ReplyMarkup = ctx.FormWeatherChanges
		action[username] = ctx.ActionManageFormWeatherChanges
	} else {
		msg.Text = "Показания доски HST(цифрами или 0)"
		action[username] = ctx.ActionManageFormHST
	}
}

func ProcessKeyboardFormWeatherChanges(db *sql.DB, msg *tgbotapi.MessageConfig, update *tgbotapi.Update, action map[string]int) {
	username := update.Message.From.UserName
	message := update.Message.Text
	now := time.Now()

	if storage.FormUpdateWeatherChanges(db, &now, &username, &message) {
		msg.Text = ctx.FormWeatherChangesAdditionalText
		msg.ReplyMarkup = ctx.YesNoMenu
		action[username] = ctx.ActionManageFormWeatherChangesAdditional
	} else {
		msg.Text = ctx.FormWeatherChangesText
		msg.ReplyMarkup = ctx.FormWeatherChanges
		action[username] = ctx.ActionManageFormWeatherChanges
	}
}

func ProcessKeyboardFormComments(db *sql.DB, msg *tgbotapi.MessageConfig, update *tgbotapi.Update, action map[string]int) {
	username := update.Message.From.UserName
	message := update.Message.Text
	now := time.Now()

	if storage.FormUpdateComments(db, &now, &message) {
		msg.Text = ctx.FormAvalancheForecastAlpText
		msg.ReplyMarkup = ctx.FormAvalancheForecast
		action[username] = ctx.ActionManageFormAvalancheForecastAlp
	} else {
		msg.Text = ctx.FormCommentsText
		action[username] = ctx.ActionManageFormComments
	}
}

func ProcessKeyboardFormDecline(db *sql.DB, msg *tgbotapi.MessageConfig, update *tgbotapi.Update, action map[string]int) {
	username := update.Message.From.UserName
	message := update.Message.Text
	now := time.Now()

	if storage.FormDecline(db, &now, &username, &message) {
		msg.Text = "Комментарий добавлен"
		action[username] = ctx.ActionNone
	} else {
		msg.Text = "Неудалось внести данные"
		action[username] = ctx.ActionNone
	}
}

func ProcessKeyboardFormArchive(db *sql.DB, msg *tgbotapi.MessageConfig, update *tgbotapi.Update, action map[string]int) {
	username := update.Message.From.UserName
	date, err := time.Parse("02 Jan 2006", update.Message.Text)

	if err != nil {
		msg.Text = "Неверный формат даты"
	} else {
		msg.ParseMode = "markdown"
		if storage.FormIsCompleted(db, &date) {
			msg.Text = generateTextReport(db, &date)
		} else {
			msg.Text = "Отчет еще не закончен"
		}
	}

	action[username] = ctx.ActionNone
}