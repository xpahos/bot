package form

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/storage"
	"time"
)

func PrepareCommandForm(db *sql.DB, msg *tgbotapi.MessageConfig, action map[string]int, username *string) {
	now := time.Now()

	duty, err := storage.DutyGetOne(db, &now)
	if err != nil {
		msg.Text = "Не выбран дежурный"
		action[*username] = ctx.ActionNone
		return
	}

	if duty != *username {
		msg.Text = fmt.Sprintf("Сегодня дежурный %s", duty)
		action[*username] = ctx.ActionNone
		return
	}

	msg.Text = ctx.FormActionMenuText
	msg.ReplyMarkup = ctx.FormActionMenu
	action[*username] = ctx.ActionManageFormActionMenu
}

func PrepareCommandConfirm(db *sql.DB, msg *tgbotapi.MessageConfig, action map[string]int, username *string) {
	action[*username] = ctx.ActionNone
	now := time.Now()

	duty, err := storage.DutyGetOne(db, &now)
	if err != nil {
		msg.Text = "Не выбран дежурный"
		return
	}

	if duty == *username {
		msg.Text = "Вы не можете подтверждать свои отчеты"
		return
	}

	if storage.FormIsCompleted(db, &now) {
		storage.FormConfirm(db, &now, username)
		msg.Text = "Отчет подтвержден"
		return
	}

	msg.Text = "Отчет еще не закончен"
}

func PrepareCommandDecline(db *sql.DB, msg *tgbotapi.MessageConfig, action map[string]int, username *string)  {
	now := time.Now()

	duty, err := storage.DutyGetOne(db, &now)
	if err != nil {
		msg.Text = "Не выбран дежурный"
		action[*username] = ctx.ActionNone
		return
	}

	if duty == *username {
		msg.Text = "Вы не можете подтверждать свои отчеты"
		action[*username] = ctx.ActionNone
		return
	}

	if storage.FormIsCompleted(db, &now) {
		msg.Text = "Введите дополнительный комментарий"
		action[*username] = ctx.ActionManageFormDeclineComment
		return
	}

	msg.Text = "Отчет еще не закончен"
	action[*username] = ctx.ActionNone
}
