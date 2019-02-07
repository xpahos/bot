package ctx

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type FormProblemStruct struct {
	ProblemType       string
	ProblemLocation   map[string]bool
	ProblemLikelyHood string
	ProblemSize       string
}

type FormProblemReadOnlyStruct struct {
	ProblemType       string
	ProblemLocation   []byte
	ProblemLikelyHood string
	ProblemSize       string
}

const (
	FormProblemMenuText       = "Добавить преблемные места?"
	FormProblemTypeText       = "Выберите тип проблемы:"
	FormProblemLocationText   = "Выберите экспозицию (С-З - направление, 1 - альпийская зона, 2 - зона деревьев, 3 - ниже зоны деревьев)\nТекущее значение:"
	FormProblemLikelyHoodText = "Вероятность схода лавины:"
	FormProblemSizeText       = "Размер потенциальной лавины(1 - не может засыпать человека, 2 - может засыпать человека, 3 - может уничтожить дом, 4 и 5 - может уничтожить часть или все поселение):"
)

var FormProblemType = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Славы", "DRY_LOOSE"),
		tgbotapi.NewInlineKeyboardButtonData("Свежие доски", "STORM_SLAB"),
		tgbotapi.NewInlineKeyboardButtonData("Ветряные доски", "WIND_SLAB"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Поверхностные доски", "PERS_SLAB"),
		tgbotapi.NewInlineKeyboardButtonData("Глубинные доски", "DEEP_PERS_SLAB"),
		tgbotapi.NewInlineKeyboardButtonData("Мокрые славы", "WET_LOOSE"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Мокрые доски", "WET_SLAB"),
		tgbotapi.NewInlineKeyboardButtonData("Падение карнизов", "CORN_FALL"),
		tgbotapi.NewInlineKeyboardButtonData("Снежные платки", "GLIDE"),
	},
)

var FormProblemLocation = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("С-З 3", "NW3"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("С 3", "N3"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("С-В 3", "NE3"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("С-З 2", "NW2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("С 2", "N2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("С-В 2", "NE2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("С-З 1", "NW1"),
		tgbotapi.NewInlineKeyboardButtonData("С 1", "NW1"),
		tgbotapi.NewInlineKeyboardButtonData("С-В 1", "NE1"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("З 3", "W3"),
		tgbotapi.NewInlineKeyboardButtonData("З 2", "W2"),
		tgbotapi.NewInlineKeyboardButtonData("З 1", "W1"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("В 1", "E1"),
		tgbotapi.NewInlineKeyboardButtonData("В 2", "E2"),
		tgbotapi.NewInlineKeyboardButtonData("В 3", "E3"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("Ю-З 1", "SW1"),
		tgbotapi.NewInlineKeyboardButtonData("Ю 1", "S1"),
		tgbotapi.NewInlineKeyboardButtonData("Ю-В 1", "SE1"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("Ю-З 2", "SW2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("Ю 2", "S2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("Ю-В 2", "SE2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Ю-З 3", "SW3"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("Ю 3", "S3"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData("Ю-В 3", "SE3"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Завершить", "DONE"),
		tgbotapi.NewInlineKeyboardButtonData("Очистить", "CLEAR"),
	},
)

var FormProblemLikelyHood = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Маловероятно", "UNLIKELY"),
		tgbotapi.NewInlineKeyboardButtonData("Вероятно", "LIKELY"),
		tgbotapi.NewInlineKeyboardButtonData("Точно", "CERTAIN"),
	},
)

var FormProblemSize = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("1", "1"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
	},
)
