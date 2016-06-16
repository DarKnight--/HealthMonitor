package setup

import (
	"os"
	"path"
)

var config = `# Config file for OWTF-HealthMonitor

# Config file should contain absolute paths or file relative to $HOME directory

# HomeDir is the directory where all the logs and config file will reside
HomeDir = ".owtf_monitor/"

# DBFile is the sqlite database file absolute path
DBFile = ".owtf_monitor/config/monitor.db"

# OWTFAddress is the address of OWTF API
OWTFAddress = "http://127.0.0.1:8009"

# Name of the file where information about modules is kept.
ModuleInfoFilePath = ".owtf_monitor/config/status.conf"

# Port for the web server
Port = "8080"
`

func setupConfig() {
	var baseDir = path.Join(os.Getenv("HOME"), ".owtf_monitor")
	var configDir = path.Join(baseDir, "config")
	var configFile = path.Join(configDir, "config.toml")

	// Update current config variables
	ConfigVars.HomeDir = baseDir
	ConfigVars.DBFile = path.Join(configDir, "monitor.db")
	ConfigVars.OWTFAddress = "http://127.0.0.1:8009"
	ConfigVars.Port = "8080"
	ConfigVars.ModuleInfoFilePath = path.Join(configDir, "status.conf")

	_, err := os.Stat(configDir)
	if err != nil {
		// Create the config directory as it does not exists.
		os.MkdirAll(configDir, 0777)
	}

	file, _ := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()

	file.WriteString(config)

	// Complete initialisation process
	dbInit()
	loadStatus()
}
