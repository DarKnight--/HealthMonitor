package setup

import (
	"os"

	"github.com/owtf/health_monitor/utils"
)

//Live saves the default config of live module to db
func Live() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Live(
		profile  		CHAR(50) PRIMARY KEY NOT NULL,
		head_url 		CHAR(50) NOT NULL,
		recheck_threshold   	INT NOT NULL,
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

//Disk saves the default config of disk module to db
func Disk() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Disk(
		profile			CHAR(50) PRIMARY KEY NOT NULL,
		space_w_limit		INT NOT NULL,
		space_d_limit		INT NOT NULL,
		inode_w_limit		INT NOT NULL,
		inode_d_limit		INT NOT NULL,
		recheck_threshold 	INT NOT NULL,
		compression_output	CHAR(500) NOT NULL,
		disk			CHAR(500) NOT NULL
		);`)
	_, err := Database.Exec(`INSERT OR REPLACE INTO Disk VALUES ("default", 2000, 1000, 2000, 1000, 5000,
			"/tmp/owtf", "/,` + os.Getenv("HOME") + `");`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to insert value to Disk table", err.Error())
		return
	}
	utils.ModuleLogs(DBLogFile, "Inserted default values to Disk table")
}

//RAM saves the default config of ram module to db
func RAM() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Ram(
		profile			CHAR(50) PRIMARY KEY NOT NULL,
		ram_w_limit		INT NOT NULL,
		recheck_threshold 	INT NOT NULL
		);`)
	_, err := Database.Exec(`INSERT OR REPLACE INTO Ram VALUES ("default", 95, 5000);`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to insert value to Ram table", err.Error())
		return
	}
	utils.ModuleLogs(DBLogFile, "Inserted default values to Ram table")
}

//CPU saves the default config of cpu module to db
func CPU() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS CPU(
		profile			CHAR(50) PRIMARY KEY NOT NULL,
		cpu_w_limit		INT NOT NULL,
		recheck_threshold 	INT NOT NULL
		);`)
	_, err := Database.Exec(`INSERT OR REPLACE INTO CPU VALUES ("default", 95, 5000);`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to insert value to CPU table", err.Error())
		return
	}
	utils.ModuleLogs(DBLogFile, "Inserted default values to CPU table")
}

//Target saves the default config of target module to db
func Target() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Target(
		profile			CHAR(50) PRIMARY KEY NOT NULL,
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

//TargetHash saves the fuzzy hash of the target response
func TargetHash() {
	_, err := Database.Exec(`CREATE TABLE IF NOT EXISTS TargetHash(
		url		CHAR(50) PRIMARY KEY NOT NULL,
		hash		CHAR(300) NOT NULL,
		Timestamp 	DATETIME DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to create TargetHash table", err.Error())
		return
	}
}

//Alert saves the default config of alert module to db
func Alert() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Alert(
		profile			CHAR(50) PRIMARY KEY NOT NULL,
		sg_api_key		CHAR(100),
		em_api_key		CHAR(100),
		em_api_uname		CHAR(100),
		mj_api_pubkey		CHAR(100),
		mj_api_prikey		CHAR(100),
		send_mail_to		CHAR(50),
		mg_api_domain		CHAR(100),
		mg_api_pubkey		CHAR(100),
		mg_api_prikey		CHAR(100),
		max_email_retry		CHAR(100),
		send_desktop_notific	INTEGER NOT NULL,
		mail_to_use		CHAR(100),
		icon_path		CHAR(100)
		);`)
	// TODO set iconPath
	_, err := Database.Exec(`INSERT OR REPLACE INTO Alert VALUES ("default", "",
		 "", "", "", "", "", "", "", "", 3, 1, "", "");`)
	if err != nil {
		utils.ModuleError(DBLogFile, "Unable to insert value to Alert table", err.Error())
		return
	}
	utils.ModuleLogs(DBLogFile, "Inserted default values to Alert table")
}

func setupDB() {
	Live()
	Target()
	Disk()
	RAM()
	CPU()
	TargetHash()
	Alert()
	return
}
