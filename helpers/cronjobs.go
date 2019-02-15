package helpers

import (
	"database/sql"
	"github.com/google/logger"
	"time"

	"github.com/xpahos/bot/storage"
)

func CronJobCheckDuty(notifies chan<- int64, db *sql.DB) {
	cronJobRoundTime()

	for {
		for _, user := range storage.UsersGetAllNotifiable(db) {
			now := time.Now().UTC().Add(time.Duration(user.TimeZone) * time.Hour)

			if now.Hour() >= user.TimeStart && now.Hour() < user.TimeEnd {
				notifies <- user.ChatID
				logger.Infof("Duty notification was scheduled to send to %s", user.Username)
			}
		}

		time.Sleep(time.Hour)
	}
}

func CronJobCheckReport(notifies chan<- int64, db *sql.DB) {
	cronJobRoundTime()

	for {
		now := time.Now()
		username, err := storage.DutyGetOne(db, &now)

		if err == nil {
			user, err := storage.UsersGetOneNotifiable(db, &username)

			if err != nil {
				logger.Errorf("Report notification was not sent because user %s turned off notifications", username)
			} else {
				now := time.Now().UTC().Add(time.Duration(user.TimeZone) * time.Hour)
				if now.Hour() >= user.TimeStart && now.Hour() < user.TimeEnd {
					notifies <- user.ChatID
					logger.Infof("Report notification was scheduled to send to %s", user.Username)
				}
			}
		} else {
			logger.Errorf("Report notification was not sent because of no duty")
		}
		time.Sleep(time.Hour)
	}
}

/*
	Round time on start
*/
func cronJobRoundTime() {
	now := time.Now()
	sleepTime := (60-now.Second())*60 + (60 - now.Minute())

	// Sleep only if current minutes less then 10
	if sleepTime > 10*60 {
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}
