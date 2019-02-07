package storage

import (
	"database/sql"

	"github.com/xpahos/bot/ctx"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
)

func UsersCheckTrusted(db *sql.DB, cache map[string]bool, update *tgbotapi.Update) bool {
	var userName string
	var chatId int64

	if update.CallbackQuery != nil {
		userName = update.CallbackQuery.From.UserName
		chatId = update.CallbackQuery.Message.Chat.ID
	} else if update.Message != nil {
		userName = update.Message.From.UserName
		chatId = update.Message.Chat.ID
	}

	userCache, err := db.Prepare("SELECT username, chat_id FROM users WHERE username = ?")
	if err != nil {
		logger.Errorf("Users get error: %v", err)
		return false
	}
	defer userCache.Close()

	var user string
	var dbChatId *int64
	err = userCache.QueryRow(userName).Scan(&user, &dbChatId)
	if err != nil {
		logger.Errorf("Users get error: %v", err)
		return false
	}

	// something strange
	if user != userName {
		return false
	}

	if dbChatId == nil {
		userChatId, err := db.Prepare("UPDATE users SET chat_id = ? WHERE username = ?")
		if err != nil {
			logger.Errorf("Users get chat_id error: %v", err)
		} else {
			_, err = userChatId.Exec(chatId, userName)

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
		return make([]string, 0, 0)
	}
	defer userSelect.Close()

	var buf string
	var result []string = make([]string, 0, 9)
	for userSelect.Next() {
		err = userSelect.Scan(&buf)
		if err != nil {
			logger.Errorf("Users get error: %v", err)
			return make([]string, 0, 0)
		}

		if len(buf) != 0 {
			result = append(result, buf)
		}
	}

	return result
}

func UsersGetChatIdList(db *sql.DB) []int64 {
	userSelect, err := db.Query("SELECT chat_id FROM users")
	if err != nil {
		logger.Errorf("Users get error: %v", err)
		return make([]int64, 0, 0)
	}
	defer userSelect.Close()

	var buf *int64
	var result []int64 = make([]int64, 0, 9)
	for userSelect.Next() {
		err = userSelect.Scan(&buf)
		if err != nil {
			logger.Errorf("Users get error: %v", err)
			return make([]int64, 0, 0)
		}

		if buf != nil {
			result = append(result, *buf)
		}
	}

	return result
}
