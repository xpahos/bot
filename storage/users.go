package storage

import (
	"database/sql"
	"errors"

	"github.com/xpahos/bot/ctx"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

func UsersCheckTrusted(db *sql.DB, cache map[string]bool, update *tgbotapi.Update) bool {
	var userName string
	var chatID int64

	if update.CallbackQuery != nil {
		userName = update.CallbackQuery.From.UserName
		chatID = update.CallbackQuery.Message.Chat.ID
	} else if update.Message != nil {
		userName = update.Message.From.UserName
		chatID = update.Message.Chat.ID
	}

	userCache, err := db.Prepare("SELECT username, chat_id FROM users WHERE username = ?")
	if err != nil {
		logger.Errorf("Users get error: %v", err)
		return false
	}
	defer userCache.Close()

	var user string
	var dbChatID *int64
	err = userCache.QueryRow(userName).Scan(&user, &dbChatID)
	if err != nil {
		logger.Errorf("Users get error: %v", err)
		return false
	}

	// something strange
	if user != userName {
		return false
	}

	if dbChatID == nil {
		userChatID, err := db.Prepare("UPDATE users SET chat_id = ? WHERE username = ?")
		if err != nil {
			logger.Errorf("Users get chat_id error: %v", err)
		} else {
			_, err = userChatID.Exec(chatID, userName)

			if err != nil {
				logger.Errorf("Users get chat_id error: %v", err)
			}
		}
	}

	cache[userName] = true
	return true
}

func UsersAddOne(db *sql.DB, user *string) bool {
	userInsert, err := db.Prepare("INSERT INTO users(username) VALUES(?)")
	if err != nil {
		logger.Errorf("Users add error: %v", err)
		return false
	}
	defer userInsert.Close()

	_, err = userInsert.Exec(*user)
	if err != nil {
		logger.Errorf("Users add error: %v", err)
		return false
	}

	return true

}

func UsersDeleteOne(db *sql.DB, user *string) bool {
	if ctx.TrustedUsers[*user] {
		logger.Errorf("Someone tried to delete trusted user %v", *user)
		return false
	}

	userDelete, err := db.Prepare("DELETE FROM users WHERE username = ?")
	if err != nil {
		logger.Errorf("Users delete error: %v", err)
		return false
	}
	defer userDelete.Close()

	_, err = userDelete.Exec(*user)
	if err != nil {
		logger.Errorf("Users delete error: %v", err)
		return false
	}

	return true
}

func UsersGetList(db *sql.DB) []string {
	userSelect, err := db.Query("SELECT username FROM users ORDER BY username LIMIT 16")
	if err != nil {
		logger.Errorf("Users get error: %v", err)
		return make([]string, 0)
	}
	defer userSelect.Close()

	var buf string
	result := make([]string, 0, 9)
	for userSelect.Next() {
		err = userSelect.Scan(&buf)
		if err != nil {
			logger.Errorf("Users get error: %v", err)
			return make([]string, 0)
		}

		if len(buf) != 0 {
			result = append(result, buf)
		}
	}

	return result
}

func UsersGetChatIDList(db *sql.DB) []int64 {
	userSelect, err := db.Query("SELECT chat_id FROM users WHERE notifications = true AND chat_id IS NOT NULL")
	if err != nil {
		logger.Errorf("Users get error: %v", err)
		return make([]int64, 0)
	}
	defer userSelect.Close()

	var buf int64
	result := make([]int64, 0, 9)
	for userSelect.Next() {
		err = userSelect.Scan(&buf)
		if err != nil {
			logger.Errorf("Users get error: %v", err)
			return make([]int64, 0)
		}

		result = append(result, buf)
	}

	return result
}

func UsersGetOneNotificationInfo(db *sql.DB, user *string) (ctx.SettingsNotificationInfoStruct, error) {
	var result ctx.SettingsNotificationInfoStruct

	userSelect, err := db.Prepare("SELECT notifications, time_start, time_end, time_zone FROM users WHERE username = ?")
	if err != nil {
		logger.Errorf("Users notifications get error: %v", err)
		return result, errors.New("DB error")
	}
	defer userSelect.Close()

	err = userSelect.QueryRow(*user).Scan(&result.IsOn, &result.TimeStart, &result.TimeEnd, &result.TimeZone)
	if err != nil {
		logger.Errorf("Users notifications get error: %v", err)
		return result, errors.New("DB error")
	}

	return result, nil
}

func UsersUpdateNotifications(db *sql.DB, user *string, flag bool) bool {
	userChatID, err := db.Prepare("UPDATE users SET notifications = ? WHERE username = ?")
	if err != nil {
		logger.Errorf("Users set notifications error: %v", err)
		return false
	}

	_, err = userChatID.Exec(flag, *user)
	if err != nil {
		logger.Errorf("Users set notifications error: %v", err)
		return false
	}

	return true
}

func UsersUpdateNotificationsTime(db *sql.DB, user *string, val int, flag int) bool {
	var field string

	switch flag {
	case ctx.SettingsTimeStartUpdate:
		field = "time_start"
	case ctx.SettingsTimeEndUpdate:
		field = "time_end"
	case ctx.SettingsTimeZoneUpdate:
		field = "time_zone"
	default:
		logger.Errorf("Incorrect flag for notifications time update")
		return false
	}

	query, err := db.Prepare("UPDATE users SET " + field + " = ? WHERE username = ?")
	if err != nil {
		logger.Errorf("Users set %s field error: %s", field, err)
		return false
	}

	_, err = query.Exec(val, *user)
	if err != nil {
		logger.Errorf("Users set %s field error: %s", field, err)
		return false
	}

	return true
}

func UsersGetAllNotifiable(db *sql.DB) []ctx.UsersNotificationDurationStruct {
	userSelect, err := db.Query(`SELECT
       		username, chat_id, time_start, time_end, time_zone
		FROM users WHERE notifications = true AND chat_id IS NOT NULL`)
	if err != nil {
		logger.Errorf("Users get error: %v", err)
		return make([]ctx.UsersNotificationDurationStruct, 0)
	}
	defer userSelect.Close()

	var buf ctx.UsersNotificationDurationStruct
	result := make([]ctx.UsersNotificationDurationStruct, 0, 9)
	for userSelect.Next() {
		err = userSelect.Scan(&buf.Username, &buf.ChatID, &buf.TimeStart, &buf.TimeEnd, &buf.TimeZone)
		if err != nil {
			logger.Errorf("Users get error: %v", err)
			return make([]ctx.UsersNotificationDurationStruct, 0)
		}

		result = append(result, buf)
	}

	return result
}

func UsersGetOneNotifiable(db *sql.DB, user *string) (ctx.UsersNotificationDurationStruct, error) {
	var result ctx.UsersNotificationDurationStruct

	userSelect, err := db.Prepare(`SELECT
       		username, chat_id, time_start, time_end, time_zone
		FROM users WHERE notifications = true AND chat_id IS NOT NULL AND username = ?`)
	if err != nil {
		logger.Errorf("Users notifiable get error: %v", err)
		return result, errors.New("Empty")
	}
	defer userSelect.Close()

	err = userSelect.QueryRow(*user).Scan(&result.Username, &result.ChatID, &result.TimeStart, &result.TimeEnd, &result.TimeZone)
	if err != nil {
		logger.Errorf("Users notifiable get error: %v", err)
		return result, errors.New("Empty")
	}

	return result, nil
}
