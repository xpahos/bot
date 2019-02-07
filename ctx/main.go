package ctx

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	Alp   = 0
	Tree  = 1
	BTree = 2
)

const (
	ActionNone                             = 0
	ActionManageUserActionMenu             = 10
	ActionManageUserAdd                    = 15
	ActionManageUserDelete                 = 20
	ActionManageDutyActionMenu             = 25
	ActionManageDutyAdd                    = 30
	ActionManageDutyDelete                 = 35
	ActionManageFormActionMenu             = 40
	ActionManageFormWindBlowing            = 45
	ActionManageFormWeatherTrend           = 50
	ActionManageFormHN24                   = 55
	ActionManageFormH2D                    = 60
	ActionManageFormHST                    = 65
	ActionManageFormWeatherChanges         = 70
	ActionManageFormProblemMenu            = 75
	ActionManageFormProblemType            = 80
	ActionManageFormProblemLocation        = 85
	ActionManageFormProblemLikelyHood      = 90
	ActionManageFormProblemSize            = 95
	ActionManageFormComments               = 100
	ActionManageFormAvalancheForecastAlp   = 105
	ActionManageFormAvalancheForecastTree  = 110
	ActionManageFormAvalancheForecastBTree = 115
	ActionManageFormDeclineComment         = 120
)

var YesNoMenu = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Да", "Y"),
		tgbotapi.NewInlineKeyboardButtonData("Нет", "N"),
	},
)

const HelpText = `
*/form* - Управление отчетом за сегодняшний день

*/confirm* - Подтвердить отчет за сегодняшний день

*/decline* - Отклонить отчет за сегодняшний день

*/report* - Вывести отчет за сегодняшний день в текстовом формате

*/duty* - Работа с графиком дежурств. Добавление и удаление возможно только за предстоящие даты.
Просмотр показывает дежурства за неделю до и неделю после.

*/users* - Добавление и удаление пользователей, которые имеют доступ к боту.
Доступны операции добавления и удаления пользователей.`
