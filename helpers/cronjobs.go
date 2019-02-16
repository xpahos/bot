package helpers

import (
	"database/sql"
	"github.com/google/logger"
	"time"

	"github.com/xpahos/bot/storage"
)

const constSleepTime  = 50*60

func CronJobCheckDuty(notifies chan<- int64, db *sql.DB) {
	cronJobSleepRoundTime(true)

	for {
		for _, user := range storage.UsersGetAllNotifiable(db) {
			now := time.Now().UTC().Add(time.Duration(user.TimeZone) * time.Hour)

			if now.Hour() >= user.TimeStart && now.Hour() < user.TimeEnd {
				notifies <- user.ChatID
				logger.Infof("Duty notification was scheduled to send to %s", user.Username)
			}
		}
		cronJobSleepRoundTime(false)
	}
}

func CronJobCheckReport(notifies chan<- int64, db *sql.DB) {
	cronJobSleepRoundTime(true)

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
		cronJobSleepRoundTime(false)
	}
}

/*
	Round time on start
*/
func cronJobSleepRoundTime(skip bool) {
	_, min, sec := time.Now().Clock()
	sleepTime := (60-sec) + (60 - min)*60

	// Sleep only if current minutes less then 10
	if !skip || (skip && sleepTime < constSleepTime) {
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}
