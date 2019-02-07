package form

import (
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/storage"

	"database/sql"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

func ProcessInlineFormProblemMenu(bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int, formProblemMap map[string]*ctx.FormProblemStruct) {
	logger.Info("FormProblemMenu")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true
	skipNext := false

	switch message {
	case "Y":
		var tmp ctx.FormProblemStruct
		formProblemMap[userName] = &tmp
		showMenu = false
	default:
		showMenu = false
		skipNext = true
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%s %s", ctx.FormProblemMenuText, message),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormProblemType
	}

	bot.Send(msg)

	if skipNext {
		actionStateMap[userName] = ctx.ActionManageFormComments
		ShowNextQuestion(bot, update, ctx.FormCommentsText, nil)
	} else {
		actionStateMap[userName] = ctx.ActionManageFormProblemType
		ShowNextQuestion(bot, update, ctx.FormProblemTypeText, &ctx.FormProblemType)
	}
}

func ProcessInlineFormType(bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int, formProblemMap map[string]*ctx.FormProblemStruct) {
	logger.Info("FormType")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "DRY_LOOSE", "STORM_SLAB", "WIND_SLAB", "PERS_SLAB", "DEEP_PERS_SLAB", "WET_LOOSE", "WET_SLAB", "CORN_FALL", "GLIDE":
		if formProblemMap[userName] == nil {
			var tmp ctx.FormProblemStruct
			formProblemMap[userName] = &tmp
		}

		showMenu = false
		formProblemMap[userName].ProblemType = message
		logger.Info(message)
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%s %s", ctx.FormProblemTypeText, message),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormProblemType
	}

	bot.Send(msg)

	if !showMenu {
		actionStateMap[userName] = ctx.ActionManageFormProblemLocation
		ShowNextQuestion(bot, update, ctx.FormProblemLocationText, &ctx.FormProblemLocation)
	}
}

func ProcessInlineFormLocations(bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int, formProblemMap map[string]*ctx.FormProblemStruct) {
	logger.Info("Locations")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "CLEAR":
		if formProblemMap[userName] != nil {
			formProblemMap[userName].ProblemLocation = nil
		}
	case "DONE":
		showMenu = false
	case "_":
		break
	default:
		if formProblemMap[userName] == nil {
			var tmp ctx.FormProblemStruct
			formProblemMap[userName] = &tmp
		}
		if formProblemMap[userName].ProblemLocation == nil {
			formProblemMap[userName].ProblemLocation = make(map[string]bool)
		}
		formProblemMap[userName].ProblemLocation[message] = true
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%s %v", ctx.FormProblemLocationText, formProblemMap[userName].ProblemLocation),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormProblemLocation
	}

	bot.Send(msg)

	if !showMenu {
		actionStateMap[userName] = ctx.ActionManageFormProblemLikelyHood
		ShowNextQuestion(bot, update, ctx.FormProblemLikelyHoodText, &ctx.FormProblemLikelyHood)
	}
}

func ProcessInlineFormLikelyHood(bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int, formProblemMap map[string]*ctx.FormProblemStruct) {
	logger.Info("LikelyHood")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "UNLIKELY", "LIKELY", "CERTAIN":
		if formProblemMap[userName] == nil {
			var tmp ctx.FormProblemStruct
			formProblemMap[userName] = &tmp
		}

		showMenu = false

		formProblemMap[userName].ProblemLikelyHood = message
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.FormProblemLikelyHoodText, formProblemMap[userName].ProblemLikelyHood),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormProblemLikelyHood
	}

	bot.Send(msg)

	if !showMenu {
		actionStateMap[userName] = ctx.ActionManageFormProblemSize
		ShowNextQuestion(bot, update, ctx.FormProblemSizeText, &ctx.FormProblemSize)
	}
}

func ProcessInlineFormSize(db *sql.DB, bot *tgbotapi.BotAPI, update *tgbotapi.Update, actionStateMap map[string]int, formProblemMap map[string]*ctx.FormProblemStruct) {
	logger.Info("FormSize")
	userName := update.CallbackQuery.From.UserName
	message := update.CallbackQuery.Data

	showMenu := true

	switch message {
	case "1", "2", "3", "4", "5":
		if formProblemMap[userName] == nil {
			var tmp ctx.FormProblemStruct
			formProblemMap[userName] = &tmp
		}

		showMenu = false
		formProblemMap[userName].ProblemSize = message
	}

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("%v %v", ctx.FormProblemSizeText, formProblemMap[userName].ProblemSize),
	)

	if showMenu {
		msg.ReplyMarkup = &ctx.FormProblemSize
	}

	bot.Send(msg)

	if !showMenu {
		now := time.Now()
		storage.FormAddProblem(db, &now, &userName, formProblemMap[userName])
		actionStateMap[userName] = ctx.ActionManageFormProblemMenu
		ShowNextQuestion(bot, update, ctx.FormProblemMenuText, &ctx.YesNoMenu)
	}
}
