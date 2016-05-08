package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/mattn/go-sqlite3"

	"health_monitor/utils"
)

func dbInit() {
	var DBDriver string
	file, err := os.OpenFile(Logs.HealthMonitorLog, os.O_RDWR|os.O_CREATE|
		os.O_APPEND, 0666)
	utils.PLogError(err)
	defer file.Close()
	log.SetOutput(file)
	sql.Register(DBDriver, &sqlite3.SQLiteDriver{})
	database, err := sql.Open(DBDriver, ConfigVars.DBFile)
	if err != nil {
		log.Println("Failed to create the handle")
	}
	if err2 := database.Ping(); err2 != nil {
		log.Println("Failed to keep connection alive")
	}

}
