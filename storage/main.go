package storage

import (
	"database/sql"

	"github.com/google/logger"
)

func InitDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
        username VARCHAR(255) NOT NULL UNIQUE PRIMARY KEY,
        chat_id BIGINT,
        notifications BOOL DEFAULT false
    ) ENGINE=INNODB CHARSET=utf8;`)

	if err != nil {
		logger.Errorf("DB init error: %v", err)
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS duty (
        date DATE NOT NULL UNIQUE PRIMARY KEY,
        username VARCHAR(255) NOT NULL
    ) ENGINE=INNODB CHARSET=utf8;`)

	if err != nil {
		logger.Errorf("DB init error: %v", err)
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS form (
        date DATE NOT NULL UNIQUE PRIMARY KEY,
        username VARCHAR(255) NOT NULL,
        completed BOOL,
        wind_blowing ENUM('LOW', 'MEDIUM', 'HIGH'),
        weather_trend ENUM('WORSE', 'SAME', 'BETTER'),
        hn24 INT,
        h2d INT,
        hst INT,
        comments TEXT,
        avalanche_forecast_alp ENUM('1', '2', '3', '4', '5'),
        avalanche_forecast_tree ENUM('1', '2', '3', '4', '5'),
        avalanche_forecast_btree ENUM('1', '2', '3', '4', '5')
    ) ENGINE=INNODB CHARSET=utf8;`)

	if err != nil {
		logger.Errorf("DB init error: %v", err)
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS form_weather_changes (
        id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
        date DATE NOT NULL,
        username VARCHAR(255),
        changes VARCHAR(255)
    ) ENGINE=INNODB CHARSET=utf8;`)

	if err != nil {
		logger.Errorf("DB init error: %v", err)
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS form_problems (
        id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
        date DATE NOT NULL,
        username VARCHAR(255),
        type ENUM("DRY_LOOSE", "STORM_SLAB", "WIND_SLAB", "PERS_SLAB", "DEEP_PERS_SLAB", "WET_LOOSE", "WET_SLAB", "CORN_FALL", "GLIDE"),
        location SET("NW1", "NW2", "NW3", "N1", "N2", "N3", "NE1", "NE2", "NE3", "W1", "W2", "W3", "E1", "E2", "E3", "SW1", "SW2", "SW3", "S1", "S2", "S3", "SE1", "SE2", "SE3"),
        likelyhood ENUM("UNLIKELY", "LIKELY", "CERTAIN"),
        size ENUM('1', '2', '3', '4', '5')
    ) ENGINE=INNODB CHARSET=utf8;`)

	if err != nil {
		logger.Errorf("DB init error: %v", err)
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS form_status (
        id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
        date DATE NOT NULL,
        username VARCHAR(255) NOT NULL,
        status ENUM('YES', 'NO') NOT NULL,
        comment TEXT
    ) ENGINE=INNODB CHARSET=utf8;`)

	if err != nil {
		logger.Errorf("DB init error: %v", err)
		return err
	}

	return nil
}
