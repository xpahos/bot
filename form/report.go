package form

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xpahos/bot/ctx"
	"github.com/xpahos/bot/helpers"

	"github.com/xpahos/bot/storage"

	"database/sql"
	"fmt"
	"time"

	"github.com/google/logger"
)

func PrepareCommandArchive(db *sql.DB, msg *tgbotapi.MessageConfig, action map[string]int, username *string) {
	dateStrings := storage.FormGetDatesList(db, 16)
	buttons := make([]tgbotapi.KeyboardButton, 0, len(dateStrings))

	for _, dateStr := range dateStrings {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		buttons = append(buttons, tgbotapi.NewKeyboardButton(date.Format("02 Jan 2006")))
	}

	buttonLen := len(buttons)
	rows := make([][]tgbotapi.KeyboardButton, 0, 1+buttonLen/4)
	for i := 0; i < buttonLen; i = i + 4 {
		end := helpers.Min(i+4, buttonLen)
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(buttons[i:end]...))
	}

	menu := tgbotapi.NewReplyKeyboard(rows...)

	msg.Text = "Выберите дату или введите в свободной форме"
	msg.ReplyMarkup = menu
	action[*username] = ctx.ActionManageFormArchive
}

func PrepareCommandReport(db *sql.DB, msg *tgbotapi.MessageConfig, action map[string]int) {
	now := time.Now()

	msg.ParseMode = "markdown"
	if storage.FormIsCompleted(db, &now) {
		msg.Text = generateTextReport(db, &now)
	} else {
		msg.Text = "Отчет еще не закончен"
	}
}

func generateTextReport(db *sql.DB, day *time.Time) string {
	logger.Info("Report")
	report := fmt.Sprintf("Отчет за *%s*:\n\n", day.Format("02 Jan 2006"))
	data, err := storage.FormGetOne(db, day)

	if err != nil {
		report += "Ошибка получения данных"
		return report
	}

	// TODO(serejkus): take care about allocations here
	var changes string
	iter := storage.FormGetWeatherChangesList(db, day)
	for i, change := range iter {
		changes += change
		if i+1 != len(iter) {
			changes += ", "
		}
	}

	report += fmt.Sprintf("*Отчет подготовил*:  `%s`\n", data.Username)
	report += fmt.Sprintf("*Ветровой перенос за 24 часа*:  `%s`\n", ctx.FormWindBlowingMappingText[data.WindBlowing])
	report += fmt.Sprintf("*Общая погодная тенденция*:  `%s`\n", ctx.FormWeatherTrendMappingText[data.WeatherTrend])
	report += fmt.Sprintf("*Показания доски HN24*:  `%d`\n", data.Hn24)
	report += fmt.Sprintf("*Показания доски H2D*:  `%d`\n", data.H2d)
	report += fmt.Sprintf("*Показания доски HST*:  `%d`\n", data.Hst)
	report += fmt.Sprintf("*Ощутиемые изменения*:  `%s`\n", changes)
	report += fmt.Sprintf("*Дополнительный комментарий*:  `%s`\n", data.Comments)
	report += fmt.Sprintf("*Лавинный прогноз в альпийской зоне*:  `%s`\n", data.AvalancheForAlp)
	report += fmt.Sprintf("*Уверенность в лавинном прогнозе в альпийской зоне*:  `%s`\n", data.AvalancheConfAlp)
	report += fmt.Sprintf("*Лавинный прогноз в зоне деревьев*:  `%s`\n", data.AvalancheForTree)
	report += fmt.Sprintf("*Уверенность в лавинном прогнозе в зоне деревье*:  `%s`\n", data.AvalancheConfAlp)
	report += fmt.Sprintf("*Лавинный прогноз ниже зоны деревьев*:  `%s`\n", data.AvalancheForBTree)
	report += fmt.Sprintf("*Уверенность в лавинном прогнозе ниже деревье*:  `%s`\n", data.AvalancheConfAlp)
	report += "\n"

	for i, problem := range storage.FormGetProblemList(db, day) {
		var buf string
		iter := strings.Split(string(problem.ProblemLocation), ",")
		for i, k := range iter {
			buf += ctx.FormProblemLocationMappingText[k]
			if i+1 != len(iter) {
				buf += ", "
			}
		}

		report += fmt.Sprintf("*Проблема %d*:\n*Тип проблемы*:  `%s`\n*Экспозиция*:  `%s`\n*Вероятность возникновения*:  `%s`\n*Размер проблемы*:  `%s`\n\n",
			i, ctx.FormProblemTypeMappingText[problem.ProblemType], buf,
			ctx.FormProblemLikelyHoodMappingText[problem.ProblemLikelyHood], problem.ProblemSize)
	}

	report += "*Форму подтвердили*:\n"
	for _, status := range storage.FormGetStatusList(db, day, true) {
		report += fmt.Sprintf("`%v`\n", status.Username)
	}
	report += "\n"

	report += "*Форму отклонили*:\n"
	for _, status := range storage.FormGetStatusList(db, day, false) {
		report += fmt.Sprintf("`%v` - `%v`\n", status.Username, string(status.Comment))
	}
	report += "\n"

	return report
}
