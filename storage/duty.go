package storage

import (
	"github.com/xpahos/bot/ctx"

	"database/sql"
	"time"

	"github.com/google/logger"
)

func DutyGetList(db *sql.DB, start *time.Time, end *time.Time) []ctx.DutyInfo {
	dutySelect, err := db.Prepare("SELECT username, date FROM duty WHERE date BETWEEN ? AND ?")
	if err != nil {
		logger.Errorf("Duties get error: %v", err)
		return nil
	}
	defer dutySelect.Close()

	rows, err := dutySelect.Query(start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		logger.Errorf("Duties get error: %v", err)
		return nil
	}
	defer rows.Close()

	var buf ctx.DutyInfo
	result := make([]ctx.DutyInfo, 0, 14)
	for rows.Next() {
		err = rows.Scan(&buf.User, &buf.Date)
		if err != nil {
			logger.Errorf("Duties get error: %v", err)
			return nil
		}

		result = append(result, buf)
	}

	if err := rows.Err(); err != nil {
		logger.Errorf("Duties got error: %v", err)
		return nil
	}

	return result
}

func DutyAddOne(db *sql.DB, dateString *string, user *string) bool {
	dutyInsert, err := db.Prepare("INSERT INTO duty(date, username) VALUES(?, ?)")
	if err != nil {
		logger.Errorf("Duties add error: %v", err)
		return false
	}
	defer dutyInsert.Close()

	_, err = dutyInsert.Exec(*dateString, *user)
	if err != nil {
		logger.Errorf("Duties add error: %v", err)
		return false
	}

	return true
}

func DutyDeleteOne(db *sql.DB, dateString *string) bool {
	dutyDelete, err := db.Prepare("DELETE FROM duty WHERE date = ?")
	if err != nil {
		logger.Errorf("Duties delete error: %v", err)
		return false
	}
	defer dutyDelete.Close()

	_, err = dutyDelete.Exec(*dateString)
	if err != nil {
		logger.Errorf("Duties delete error: %v", err)
		return false
	}

	return true
}

func DutyGetOne(db *sql.DB, day *time.Time) (string, error) {
	dutySelect, err := db.Prepare("SELECT username FROM duty WHERE date = ?")
	if err != nil {
		logger.Errorf("Duties get error: %v", err)
		return "Empty", err
	}
	defer dutySelect.Close()

	var result string
	err = dutySelect.QueryRow(day.Format("2006-01-02")).Scan(&result)
	if err != nil {
		logger.Errorf("Duties get error: %v", err)
		return "Empty", err
	}

	return result, nil
}
