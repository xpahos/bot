package ctx

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	FormActionMenuText               = "Выберите действие с отчетом: "
	FormWindBlowingText              = "Ветровой перенос за последние 24 часа"
	FormWeatherTrendText             = "Общая погодная тенденция"
	FormHN24Text                     = "Показания доски HN24(цифрами или 0)"
	FormCommentsText                 = "Комментарий в свободной форме"
	FormAvalancheForecastAlpText     = "Лавинный прогноз в альпийской зоне(1 - 5): "
	FormAvalancheForecastTreeText    = "Лавинный прогноз в зоне деревьев(1 - 5): "
	FormAvalancheForecastBTreeText   = "Лавинный прогноз в зоне ниже деревьев(1 - 5): "
	FormAvalancheConfidenceAlpText   = "Уверенность в лавинном прогнозе в альпийской зоне(1 - не уверен, 5 - уверен): "
	FormAvalancheConfidenceTreeText  = "Уверенность в лавинном прогнозе в зоне деревьев(1 - не уверен, 5 - уверен): "
	FormAvalancheConfidenceBTreeText = "Уверенность в лавинном прогнозе в зоне ниже деревьев(1 - не уверен, 5 - уверен): "
	FormCompletedText                = "Отчет завершен"
	FormWeatherChangesText           = "Ощутимые изменения(выберите или введите произвольный вариант)"
	FormWeatherChangesAdditionalText = "Добавить дополнительные изменения?"
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

var FormWindBlowing = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Слабый", "LOW"),
		tgbotapi.NewInlineKeyboardButtonData("Сильный", "MEDIUM"),
		tgbotapi.NewInlineKeyboardButtonData("Очень сильный", "HIGH"),
	},
)

var FormWeatherTrend = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Ухудшается", "WORSE"),
		tgbotapi.NewInlineKeyboardButtonData("Не изменяется", "SAME"),
		tgbotapi.NewInlineKeyboardButtonData("Улучшается", "BETTER"),
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
