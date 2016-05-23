package config

import (
	"os"
	"path"
)

var config = `# Config file for OWTF-HealthMonitor

# Config file should contain absolute paths or file relative to $HOME directory

# HomeDir is the directory where all the logs and config file will reside
HomeDir = ".owtfMonitor/"

# DBFile is the sqlite database file absolute path
DBFile = ".owtfMonitor/config/monitor.db"

# OWTFAddress is the address of OWTF API
OWTFAddress = "http://127.0.0.1:8009"

# Name of the profile to set when monitor starts. It stores last used profile.
Profile = "default"

`

func setupConfig() {
	var baseDir = path.Join(os.Getenv("HOME"), ".owtfMonitor")
	var configDir = path.Join(baseDir, "config")
	var configFile = path.Join(configDir, "config.toml")

	// Update current config variables
	ConfigVars.HomeDir = baseDir
	ConfigVars.DBFile = path.Join(configDir, "monitor.db")
	ConfigVars.OWTFAddress = "http://127.0.0.1:8009"

	_, err := os.Stat(configDir)
	if err != nil {
		// Create the config directory as it does not exists.
		os.MkdirAll(configDir, 0777)
	}

	file, _ := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE, 0666)

	file.WriteString(config)

	// Complete initialisation process
	dbInit()
	logFile.Close()
}

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
	Database.Exec(`INSERT INTO Live VALUES (
	"default", "https://google.com", 30000, 4000, 4000,"8.8.8.8", "tcp");`)
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
	Database.Exec(`INSERT INTO Disk VALUES ("default", 2000, 1000, 2000, 1000, 5000,
			"/,` + os.Getenv("HOME") + `");`)
}

func setupDB() {
	setupLive()
	setupDisk()
	return
}
