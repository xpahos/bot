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
	FormProblemMenuText       = "Добавить проблемные места?"
	FormProblemTypeText       = "Выберите тип проблемы:"
	FormProblemLocationText   = "Выберите экспозицию (С-З - направление, А - альпийская зона, Д - зона деревьев, Н - ниже зоны деревьев)\nТекущее значение:"
	FormProblemLikelyHoodText = "Вероятность схода лавины:"
	FormProblemSizeText       = "Размер потенциальной лавины(1 - не может засыпать человека, 2 - может засыпать человека, 3 - может уничтожить дом, 4 и 5 - может уничтожить часть или все поселение):"
)

var FormProblemTypeMappingText = map[string]string{
	"DRY_LOOSE":      "Слафы",
	"STORM_SLAB":     "Свежие доски",
	"WIND_SLAB":      "Ветровые доски",
	"PERS_SLAB":      "Поверхностные доски",
	"DEEP_PERS_SLAB": "Глубинные доски",
	"WET_LOOSE":      "Мокрые славы",
	"WET_SLAB":       "Мокрые доски",
	"CORN_FALL":      "Падение карнизов",
	"GLIDE":          "Снежные платки",
}

var FormProblemType = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(FormProblemTypeMappingText["DRY_LOOSE"], "DRY_LOOSE"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemTypeMappingText["STORM_SLAB"], "STORM_SLAB"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemTypeMappingText["WIND_SLAB"], "WIND_SLAB"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(FormProblemTypeMappingText["PERS_SLAB"], "PERS_SLAB"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemTypeMappingText["DEEP_PERS_SLAB"], "DEEP_PERS_SLAB"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemTypeMappingText["WET_LOOSE"], "WET_LOOSE"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(FormProblemTypeMappingText["WET_SLAB"], "WET_SLAB"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemTypeMappingText["CORN_FALL"], "CORN_FALL"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemTypeMappingText["GLIDE"], "GLIDE"),
	},
)

var FormProblemLocationMappingText = map[string]string{
	"NW3": "Северо-Запад ниже деревьев",
	"N3":  "Север ниже деревьев",
	"NE3": "Северо-Восток ниже деревьев",
	"NW2": "Северо-Запад деревья",
	"N2":  "Север деревья",
	"NE2": "Северо-Восток деревья",
	"NW1": "Северов-Запад альпийская зона",
	"N1":  "Север альпийская зона",
	"NE1": "Северо-Восток альпийская зона",
	"W3":  "Запад ниже деревьев",
	"W2":  "Запад деревья",
	"W1":  "Запад альпийская зона",
	"E1":  "Восток альпийская зона",
	"E2":  "Восток деревья",
	"E3":  "Восток ниже деревьев",
	"SW1": "Юго-Запад альпийская зона",
	"S1":  "Юг альпийская зона",
	"SE1": "Юго-Восток альпийская зона",
	"SW2": "Юго-Запад деревья",
	"S2":  "Юг деревья",
	"SE2": "Юго-Запад деревья",
	"SW3": "Юго-Запад ниже деревьев",
	"S3":  "Юг ниже деревьев",
	"SE3": "Юго-Восток ниже деревьев",
}

var FormProblemLocationMappingMenuText = map[string]string{
	"NW3":   "С-З Н",
	"N3":    "С Н",
	"NE3":   "С-В Н",
	"NW2":   "С-З Д",
	"N2":    "С Д",
	"NE2":   "С-В Д",
	"NW1":   "С-З А",
	"N1":    "С А",
	"NE1":   "С-В А",
	"W3":    "З Н",
	"W2":    "З Д",
	"W1":    "З А",
	"E1":    "В А",
	"E2":    "В Д",
	"E3":    "В Н",
	"SW1":   "Ю-З А",
	"S1":    "Ю А",
	"SE1":   "Ю-В А",
	"SW2":   "Ю-З Д",
	"S2":    "Ю Д",
	"SE2":   "Ю-В Д",
	"SW3":   "Ю-З Н",
	"S3":    "Ю Н",
	"SE3":   "Ю-В Н",
	"DONE":  "Завершить",
	"CLEAR": "Очистить",
}

var FormProblemLocation = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["NW3"], "NW3"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["N3"], "N3"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["NE3"], "NE3"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["NW2"], "NW2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["N2"], "N2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["NE2"], "NE2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["NW1"], "NW1"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["N1"], "N1"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["NE1"], "NE1"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["W3"], "W3"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["W2"], "W2"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["W1"], "W1"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["E1"], "E1"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["E2"], "E2"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["E3"], "E3"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["SW1"], "SW1"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["S1"], "S1"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["SE1"], "SE1"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["SW2"], "SW2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["S2"], "S2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["SE2"], "SE2"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["SW3"], "SW3"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["S3"], "S3"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(" ", "_"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["SE3"], "SE3"),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["DONE"], "DONE"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLocationMappingMenuText["CLEAR"], "CLEAR"),
	},
)

var FormProblemLikelyHoodMappingText = map[string]string{
	"UNLIKELY": "Маловероятно",
	"LIKELY":   "Вероятно",
	"CERTAIN":  "Точно",
}

var FormProblemLikelyHood = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLikelyHoodMappingText["UNLIKELY"], "UNLIKELY"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLikelyHoodMappingText["LIKELY"], "LIKELY"),
		tgbotapi.NewInlineKeyboardButtonData(FormProblemLikelyHoodMappingText["CERTAIN"], "CERTAIN"),
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
