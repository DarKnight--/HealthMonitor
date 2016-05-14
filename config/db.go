package config

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"

	"health_monitor/utils"
)

// Database holds the active sqlite connection
var Database *sql.DB

func dbInit() {
	var DBDriver string
	var err error
	sql.Register(DBDriver, &sqlite3.SQLiteDriver{})
	Database, err = sql.Open(DBDriver, ConfigVars.DBFile)
	if err != nil {
		utils.Perror("Failed to create the handle")
		utils.Perror(err.Error())
	}
	if err = Database.Ping(); err != nil {
		utils.Perror("Failed to keep connection alive")
		utils.Perror(err.Error())
	}
}
