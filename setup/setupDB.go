package setup

import (
	"os"

	"health_monitor/utils"
)

func setupLive() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Live(
		profile  			CHAR(50) PRIMARY KEY NOT NULL,
		head_url 			CHAR(50) NOT NULL,
		recheck_threshold   INT NOT NULL,
		ping_threshold		INT NOT NULL,
		head_threshold		INT NOT NULL,
		ping_address		CHAR(50) NOT NULL,
		ping_protocol		CHAR(10)
		);`)
	_, err := Database.Exec(`INSERT OR REPLACE INTO Live VALUES (
	"default", "https://google.com", 30000, 4000, 4000,"8.8.8.8", "tcp");`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to insert value to Live table", err.Error())
		return
	}
	utils.ModuleLogs(DBLogFile, "Inserted default values to Live table")
}

func setupDisk() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Disk(
		profile				CHAR(50) PRIMARY KEY NOT NULL,
		space_w_limit		INT NOT NULL,
		space_d_limit		INT NOT NULL,
		inode_w_limit		INT NOT NULL,
		inode_d_limit		INT NOT NULL,
		recheck_threshold 	INT NOT NULL,
		disk				CHAR(500) NOT NULL
		);`)
	_, err := Database.Exec(`INSERT OR REPLACE INTO Disk VALUES ("default", 2000, 1000, 2000, 1000, 5000,
			"/,` + os.Getenv("HOME") + `");`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to insert value to Disk table", err.Error())
		return
	}
	utils.ModuleLogs(DBLogFile, "Inserted default values to Disk table")
}

func setupRAM() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Ram(
		profile				CHAR(50) PRIMARY KEY NOT NULL,
		ram_w_limit			INT NOT NULL,
		recheck_threshold 	INT NOT NULL
		);`)
	_, err := Database.Exec(`INSERT OR REPLACE INTO Ram VALUES ("default", 90000, 5000);`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to insert value to Ram table", err.Error())
		return
	}
	utils.ModuleLogs(DBLogFile, "Inserted default values to Ram table")
}

func setupCPU() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS CPU(
		profile				CHAR(50) PRIMARY KEY NOT NULL,
		ram_w_limit			INT NOT NULL,
		recheck_threshold 	INT NOT NULL
		);`)
	_, err := Database.Exec(`INSERT OR REPLACE INTO CPU VALUES ("default", 90000, 5000);`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to insert value to CPU table", err.Error())
		return
	}
	utils.ModuleLogs(DBLogFile, "Inserted default values to CPU table")
}

func setupDB() {
	setupLive()
	setupDisk()
	setupRAM()
	setupCPU()
	return
}
