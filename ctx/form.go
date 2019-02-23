package ctx

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	FormActionMenuText               = "Выберите действие с отчетом: "
	FormWindBlowingText              = "Ветровой перенос за последние 24 часа: "
	FormWeatherTrendText             = "Общая погодная тенденция: "
	FormHN24Text                     = "Показания доски HN24(цифрами или 0): "
	FormCommentsText                 = "Комментарий в свободной форме: "
	FormAvalancheForecastAlpText     = "Лавинный прогноз в альпийской зоне(1 - 5): "
	FormAvalancheForecastTreeText    = "Лавинный прогноз в зоне деревьев(1 - 5): "
	FormAvalancheForecastBTreeText   = "Лавинный прогноз в зоне ниже деревьев(1 - 5): "
	FormAvalancheConfidenceAlpText   = "Уверенность в лавинном прогнозе в альпийской зоне(1 - не уверен, 5 - уверен): "
	FormAvalancheConfidenceTreeText  = "Уверенность в лавинном прогнозе в зоне деревьев(1 - не уверен, 5 - уверен): "
	FormAvalancheConfidenceBTreeText = "Уверенность в лавинном прогнозе в зоне ниже деревьев(1 - не уверен, 5 - уверен): "
	FormCompletedText                = "Отчет завершен"
	FormWeatherChangesText           = "Ощутимые изменения(выберите или введите произвольный вариант): "
	FormWeatherChangesAdditionalText = "Добавить дополнительные ощутимые изменения? "
)

type FormStruct struct {
	Username           string
	WindBlowing        string
	WeatherTrend       string
	Hn24               int
	H2d                int
	Hst                int
	Comments           string
	AvalancheForAlp    string
	AvalancheForTree   string
	AvalancheForBTree  string
	AvalancheConfAlp   string
	AvalancheConfTree  string
	AvalancheConfBTree string
}

type FormStatusStruct struct {
	Username string
	Comment  string
}

var FormActionMenu = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Добавить", "ADD"),
		tgbotapi.NewInlineKeyboardButtonData("Редактировать", "EDIT"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "DELETE"),
	},
)

var FormWindBlowingMappingText = map[string]string{
	"LOW":    "Слабый",
	"MEDIUM": "Сильный",
	"HIGH":   "Очень сильный",
}

var FormWindBlowing = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(FormWindBlowingMappingText["LOW"], "LOW"),
		tgbotapi.NewInlineKeyboardButtonData(FormWindBlowingMappingText["MEDIUM"], "MEDIUM"),
		tgbotapi.NewInlineKeyboardButtonData(FormWindBlowingMappingText["HIGH"], "HIGH"),
	},
)

var FormWeatherTrendMappingText = map[string]string{
	"WORSE":  "Ухудшается",
	"SAME":   "Не изменяется",
	"BETTER": "Улучшается",
}

var FormWeatherTrend = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(FormWeatherTrendMappingText["WORSE"], "WORSE"),
		tgbotapi.NewInlineKeyboardButtonData(FormWeatherTrendMappingText["SAME"], "SAME"),
		tgbotapi.NewInlineKeyboardButtonData(FormWeatherTrendMappingText["BETTER"], "BETTER"),
	},
)

var FormWeatherChanges = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Перенос"),
		tgbotapi.NewKeyboardButton("Переход через ноль"),
		tgbotapi.NewKeyboardButton("Свежий снег"),
		tgbotapi.NewKeyboardButton("Теплая ночь"),
	),
)

var FormAvalancheForecast = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("1", "1"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
	},
)
