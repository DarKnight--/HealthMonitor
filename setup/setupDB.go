package setup

import (
	"os"

	"health_monitor/utils"
)

func SetupLive() {
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

func SetupDisk() {
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

func SetupRAM() {
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

func SetupCPU() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS CPU(
		profile				CHAR(50) PRIMARY KEY NOT NULL,
		cpu_w_limit			INT NOT NULL,
		recheck_threshold 	INT NOT NULL
		);`)
	_, err := Database.Exec(`INSERT OR REPLACE INTO CPU VALUES ("default", 95, 5000);`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to insert value to CPU table", err.Error())
		return
	}
	utils.ModuleLogs(DBLogFile, "Inserted default values to CPU table")
}

func SetupTarget() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Target(
		profile				CHAR(50) PRIMARY KEY NOT NULL,
		fuzzy_threshold		INT NOT NULL,
		recheck_threshold 	INT NOT NULL
		);`)
	_, err := Database.Exec(`INSERT OR REPLACE INTO Target VALUES ("default", 50, 5000);`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to insert value to Target table", err.Error())
		return
	}
	utils.ModuleLogs(DBLogFile, "Inserted default values to Target table")
}

func SetupTargetHash() {
	_, err := Database.Exec(`CREATE TABLE IF NOT EXISTS TargetHash(
		url			CHAR(50) PRIMARY KEY NOT NULL,
		hash		CHAR(300) NOT NULL,
		Timestamp 	DATETIME DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to create TargetHash table", err.Error())
		return
	}
}

func setupDB() {
	SetupLive()
	SetupTarget()
	SetupDisk()
	SetupRAM()
	SetupCPU()
	SetupTargetHash()
	return
}
