package storage

import (
	"github.com/xpahos/bot/ctx"

	"database/sql"
	"strings"
	"time"

	"github.com/google/logger"
)

func FormInitRecord(db *sql.DB, day *time.Time, user *string) bool {
	formInsert, err := db.Prepare("INSERT INTO form(date, username) VALUES(?, ?)")
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}
	defer formInsert.Close()

	_, err = formInsert.Exec(day.Format("2006-01-02"), *user)
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}

	return true
}

func FormDeleteRecord(db *sql.DB, day *time.Time) bool {
	formDelete, err := db.Prepare("DELETE FROM form WHERE date = ?")
	if err != nil {
		logger.Infof("Form delete error: %v", err)
		return false
	}
	defer formDelete.Close()

	_, err = formDelete.Exec(day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("Form delete error: %v", err)
		return false
	}

	formDelete, err = db.Prepare("DELETE FROM form_problems WHERE date = ?")
	if err != nil {
		logger.Infof("Form delete error: %v", err)
		return false
	}
	defer formDelete.Close()

	_, err = formDelete.Exec(day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("Form delete error: %v", err)
		return false
	}

	formDelete, err = db.Prepare("DELETE FROM form_weather_changes WHERE date = ?")
	if err != nil {
		logger.Infof("Form delete error: %v", err)
		return false
	}
	defer formDelete.Close()

	_, err = formDelete.Exec(day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("Form delete error: %v", err)
		return false
	}

	return true
}

func FormAddProblem(db *sql.DB, day *time.Time, user *string, ctx *ctx.FormProblemStruct) bool {
	formInsert, err := db.Prepare("INSERT INTO form_problems(date, username, type, location, likelyhood, size) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}
	defer formInsert.Close()

	var location []string
	for key := range ctx.ProblemLocation {
		location = append(location, key)
	}

	_, err = formInsert.Exec(day.Format("2006-01-02"), *user, ctx.ProblemType, strings.Join(location, ","), ctx.ProblemLikelyHood, ctx.ProblemSize)
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}

	return true
}

func FormUpdateWindBlowing(db *sql.DB, day *time.Time, msg *string) bool {
	formUpdate, err := db.Prepare("UPDATE form SET wind_blowing = ? WHERE date = ?")
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}
	defer formUpdate.Close()

	_, err = formUpdate.Exec(*msg, day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}

	return true
}

func FormUpdateWeatherTrend(db *sql.DB, day *time.Time, msg *string) bool {
	formUpdate, err := db.Prepare("UPDATE form SET weather_trend = ? WHERE date = ?")
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}
	defer formUpdate.Close()

	_, err = formUpdate.Exec(*msg, day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}

	return true
}

func FormUpdateHN24(db *sql.DB, day *time.Time, msg *string) bool {
	formUpdate, err := db.Prepare("UPDATE form SET hn24 = ? WHERE date = ?")
	if err != nil {
		logger.Infof("%v", err)
		return false
	}
	defer formUpdate.Close()

	_, err = formUpdate.Exec(msg, day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("%v", err)
		return false
	}

	return true
}

func FormUpdateH2D(db *sql.DB, day *time.Time, msg *string) bool {
	formUpdate, err := db.Prepare("UPDATE form SET h2d = ? WHERE date = ?")
	if err != nil {
		logger.Infof("%v", err)
		return false
	}
	defer formUpdate.Close()

	_, err = formUpdate.Exec(msg, day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("%v", err)
		return false
	}

	return true
}

func FormUpdateHST(db *sql.DB, day *time.Time, msg *string) bool {
	formUpdate, err := db.Prepare("UPDATE form SET hst = ? WHERE date = ?")
	if err != nil {
		logger.Infof("%v", err)
		return false
	}
	defer formUpdate.Close()

	_, err = formUpdate.Exec(msg, day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("%v", err)
		return false
	}

	return true
}

func FormUpdateWeatherChanges(db *sql.DB, day *time.Time, userName *string, msg *string) bool {
	formUpdate, err := db.Prepare("INSERT INTO form_weather_changes(date, username, changes) VALUES(?, ?, ?)")
	if err != nil {
		logger.Infof("%v", err)
		return false
	}
	defer formUpdate.Close()

	_, err = formUpdate.Exec(day.Format("2006-01-02"), *userName, *msg)
	if err != nil {
		logger.Infof("%v", err)
		return false
	}

	return true
}

func FormUpdateComments(db *sql.DB, day *time.Time, msg *string) bool {
	formUpdate, err := db.Prepare("UPDATE form SET comments = ? WHERE date = ?")
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}
	defer formUpdate.Close()

	_, err = formUpdate.Exec(*msg, day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}

	return true
}

func FormUpdateAvalanche(db *sql.DB, day *time.Time, msg *string, zone int) bool {
	var field string

	switch zone {
	case ctx.Alp:
		field = "avalanche_forecast_alp"
	case ctx.Tree:
		field = "avalanche_forecast_tree"
	case ctx.BTree:
		field = "avalanche_forecast_btree"
	}

	formUpdate, err := db.Prepare("UPDATE form SET " + field + " = ? WHERE date = ?")
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}
	defer formUpdate.Close()

	_, err = formUpdate.Exec(msg, day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}

	return true
}

func FormComplete(db *sql.DB, day *time.Time) bool {
	formUpdate, err := db.Prepare("UPDATE form SET completed = true WHERE date = ?")
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}
	defer formUpdate.Close()

	_, err = formUpdate.Exec(day.Format("2006-01-02"))
	if err != nil {
		logger.Infof("Form insert error: %v", err)
		return false
	}

	return true
}

func FormIsCompleted(db *sql.DB, day *time.Time) bool {
	formFilled, err := db.Prepare("SELECT completed FROM form WHERE date = ?")
	if err != nil {
		return false
	}
	defer formFilled.Close()

	filled := false
	err = formFilled.QueryRow(day.Format("2006-01-02")).Scan(&filled)
	if err != nil {
		return false
	}

	return filled
}

func FormConfirm(db *sql.DB, day *time.Time, userName *string) bool {
	formStatus, err := db.Prepare("SELECT count(id) FROM form_status WHERE date = ? and userName = ?")
	if err != nil {
		logger.Infof("Form confirm error: %v", err)
		return false
	}
	defer formStatus.Close()

	count := 0
	err = formStatus.QueryRow(day.Format("2006-01-02"), *userName).Scan(&count)
	if err != nil {
		logger.Infof("Form confirm error: %v", err)
		return false
	}

	if count == 0 {
		formStatus, err = db.Prepare("INSERT INTO form_status(date, username, status) VALUES(?, ?, 'YES')")
		if err != nil {
			logger.Infof("Form confirm error: %v", err)
			return false
		}

		_, err = formStatus.Exec(day.Format("2006-01-02"), *userName)
		if err != nil {
			logger.Infof("Form confirm error: %v", err)
			return false
		}
	} else {
		formStatus, err = db.Prepare("UPDATE form_status SET status = 'YES' WHERE date = ? AND username = ?")
		if err != nil {
			logger.Infof("Form confirm error: %v", err)
			return false
		}

		_, err = formStatus.Exec(day.Format("2006-01-02"), *userName)
		if err != nil {
			logger.Infof("Form confirm error: %v", err)
			return false
		}
	}

	return true
}

func FormDecline(db *sql.DB, day *time.Time, userName *string, msg *string) bool {
	formStatus, err := db.Prepare("SELECT count(id) FROM form_status WHERE date = ? and username = ?")
	if err != nil {
		logger.Infof("Form decline error: %v", err)
		return false
	}
	defer formStatus.Close()

	count := 0
	err = formStatus.QueryRow(day.Format("2006-01-02"), *userName).Scan(&count)
	if err != nil {
		logger.Infof("Form decline error: %v", err)
		return false
	}

	if count == 0 {
		formStatus, err = db.Prepare("INSERT INTO form_status(date, username, status, comment) VALUES(?, ?, 'NO', ?)")
		if err != nil {
			logger.Infof("Form decline error: %v", err)
			return false
		}

		_, err = formStatus.Exec(day.Format("2006-01-02"), *userName, *msg)
		if err != nil {
			logger.Infof("Form decline error: %v", err)
			return false
		}
	} else {
		formStatus, err = db.Prepare("UPDATE form_status SET status = 'NO', comment = ? WHERE date = ? AND username = ?")
		if err != nil {
			logger.Infof("Form decline error: %v", err)
			return false
		}

		_, err = formStatus.Exec(*msg, day.Format("2006-01-02"), *userName)
		if err != nil {
			logger.Infof("Form decline error: %v", err)
			return false
		}
	}

	return true
}

func FormGetWeatherChangesList(db *sql.DB, day *time.Time) []string {
	formSelect, err := db.Prepare("SELECT changes FROM form_weather_changes WHERE date = ?")
	if err != nil {
		logger.Errorf("Form get weather changes error: %v", err)
		return nil
	}
	defer formSelect.Close()

	rows, err := formSelect.Query(day.Format("2006-01-02"))
	if err != nil {
		logger.Errorf("Form get weather changes error: %v", err)
		return nil
	}
	defer rows.Close()

	var buf string
	result := make([]string, 0, 1)
	for rows.Next() {
		err = rows.Scan(&buf)
		if err != nil {
			logger.Errorf("Form get weather changes error: %v", err)
			return nil
		}

		result = append(result, buf)
	}

	if err := rows.Err(); err != nil {
		logger.Errorf("GetWeatherChangesList got error: %v", err)
		return nil
	}

	return result
}

func FormGetProblemList(db *sql.DB, day *time.Time) []ctx.FormProblemReadOnlyStruct {
	formSelect, err := db.Prepare("SELECT type, location, likelyhood, size FROM form_problems WHERE date = ?")
	if err != nil {
		logger.Errorf("Form get error: %v", err)
		return nil
	}
	defer formSelect.Close()

	rows, err := formSelect.Query(day.Format("2006-01-02"))
	if err != nil {
		logger.Errorf("Form get error: %v", err)
		return nil
	}
	defer rows.Close()

	var buf ctx.FormProblemReadOnlyStruct
	result := make([]ctx.FormProblemReadOnlyStruct, 0, 14)
	for rows.Next() {
		err = rows.Scan(&buf.ProblemType, &buf.ProblemLocation, &buf.ProblemLikelyHood, &buf.ProblemSize)
		if err != nil {
			logger.Errorf("Form get error: %v", err)
			return nil
		}

		result = append(result, buf)
	}

	if err := rows.Err(); err != nil {
		logger.Errorf("GetProblemList got error: %v", err)
		return nil
	}

	return result
}

func FormGetStatusList(db *sql.DB, day *time.Time, confirm bool) []ctx.FormStatusStruct {
	var status string

	if confirm {
		status = "YES"
	} else {
		status = "NO"
	}

	formSelect, err := db.Prepare("SELECT username, comment FROM form_status WHERE date = ? AND status = ?")
	if err != nil {
		logger.Errorf("Form get error: %v", err)
		return nil
	}
	defer formSelect.Close()

	rows, err := formSelect.Query(day.Format("2006-01-02"), status)
	if err != nil {
		logger.Errorf("Form get error: %v", err)
		return nil
	}
	defer rows.Close()

	var buf ctx.FormStatusStruct
	var comment sql.NullString
	result := make([]ctx.FormStatusStruct, 0, 14)
	for rows.Next() {
		err = rows.Scan(&buf.Username, &comment)
		if err != nil {
			logger.Errorf("Form get error: %v", err)
			return nil
		}

		if comment.Valid {
			buf.Comment = comment.String
		}

		result = append(result, buf)
	}

	if err := rows.Err(); err != nil {
		logger.Errorf("GetStatusList got error: %v", err)
		return nil
	}

	return result
}

func FormGetOne(db *sql.DB, day *time.Time) (ctx.FormStruct, error) {
	var form ctx.FormStruct

	formSelect, err := db.Prepare(`SELECT 
        username, wind_blowing, weather_trend, hn24, h2d, hst, comments, avalanche_forecast_alp, avalanche_forecast_tree, avalanche_forecast_btree
        FROM form WHERE date = ?`)
	if err != nil {
		logger.Errorf("Form get error: %v", err)
		return form, err
	}
	defer formSelect.Close()

	err = formSelect.QueryRow(day.Format("2006-01-02")).Scan(&form.Username, &form.WindBlowing, &form.WeatherTrend, &form.Hn24, &form.H2d, &form.Hst,
		&form.Comments, &form.AvalancheAlp, &form.AvalancheTree, &form.AvalancheBTree)
	if err != nil {
		logger.Errorf("Form get error: %v", err)
		return form, err
	}

	return form, nil
}

func FormGetDatesList(db *sql.DB, count int) []string {
	query, err := db.Prepare(`SELECT date FROM form WHERE completed = TRUE ORDER BY date DESC LIMIT ?`)
	if err != nil {
		logger.Errorf("Form dates get error: %v", err)
		return nil
	}
	defer query.Close()

	rows, err := query.Query(count)
	if err != nil {
		logger.Errorf("Form dates get error: %v", err)
		return nil
	}
	defer rows.Close()

	var buf string
	result := make([]string, 0, 0)
	for rows.Next() {
		err = rows.Scan(&buf)
		if err != nil {
			logger.Errorf("Form get error: %v", err)
			return nil
		}

		result = append(result, buf)
	}

	if err := rows.Err(); err != nil {
		logger.Errorf("GetDatesList got error: %v", err)
		return nil
	}

	return result
}