package form

import (
    //"ctx"
    "storage"

    "fmt"
    "time"
    "database/sql"

    "github.com/google/logger"
)

func GenerateTextReport(db *sql.DB, day *time.Time) string {
    logger.Infof("Report")
    report := fmt.Sprintf("Отчет за *`%v`*:\n\n", day.Format("02 Jan 2006"))
    data, err := storage.FormGetOne(db, day)

    if err != nil {
        report += "Ошибка получения данных"
    } else {
        report += fmt.Sprintf("*Отчет подготовил*: `%v`\n", data.Username)
        report += fmt.Sprintf("*Ветровой перенос за 24 часа*: `%v`\n", data.WindBlowing)
        report += fmt.Sprintf("*Общая погодная тенденция*: `%v`\n", data.WeatherTrend)
        report += fmt.Sprintf("*Показания доски HN24*: `%v`\n", data.Hn24)
        report += fmt.Sprintf("*Показания доски H2D*: `%v`\n", data.H2d)
        report += fmt.Sprintf("*Показания доски HST*: `%v`\n", data.Hst)
        report += fmt.Sprintf("*Ощутиемые изменения*: `%v`\n", data.WeatherChanges)
        report += fmt.Sprintf("*Дополнительный комментарий*: `%v`\n", data.Comments)
        report += fmt.Sprintf("*Лавинный прогноз в альийской зоне*: `%v`\n", data.AvalancheAlp)
        report += fmt.Sprintf("*Лавинный прогноз в зоне деревьев*: `%v`\n", data.AvalancheTree)
        report += fmt.Sprintf("*Лавинный прогноз ниже зоны деревьев*: `%v`\n\n", data.AvalancheBTree)

        for i, problem := range storage.FormGetProblemList(db, day) {
            report += fmt.Sprintf("*Проблема %v*:\n*Тип проблемы*: `%v`\n*Экспозиция*: `%v`\n*Вероятность возникновения*: `%v`\n*Размер проблемы*: `%v`\n\n",
                i, problem.ProblemType, string(problem.ProblemLocation), problem.ProblemLikelyHood, problem.ProblemSize)
        }

        report += "*Форму подтвердили*:\n"
        for _, status := range storage.FormGetStatusList(db, day, true) {
            report += status.Username + "\n"
        }
        report += "\n"

        report += "*Форму отклонили*:\n"
        for _, status := range storage.FormGetStatusList(db, day, true) {
            report += status.Username + " - " + status.Comment + "\n"
        }
        report += "\n"

    }

    return report
}
