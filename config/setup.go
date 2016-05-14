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
		setupDB()
	}

	file, _ := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE, 0666)

	file.WriteString(config)

	// Complete initialisation process
	dbInit()
	logFile.Close()
}

func setupDB() {
	return
}
