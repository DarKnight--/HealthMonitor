package config

import (
	"fmt"
	"os"
	"path"
)

var (
	//home directory path
	HomeDir string = os.Getenv("HOME")

	//structure containing database related data
	DB struct {
		DBConfigFile string
	}

	testConfig bool   = loadConfig()
	errorMsg   string = "The config file is missing. Please run the setup script" +
		"to setup the OWTF Health Monitor."
)

// This function will iniailise all the configuration variable defined
func loadConfig() bool {
	DB.DBConfigFile = path.Join(os.Getenv("HOME"), ".owtfMonitor", "config",
		"db_config.toml")
	if _, err := os.Stat(DB.DBConfigFile); os.IsNotExist(err) {
		fmt.Println(errorMsg)
		os.Exit(1)
	}
	return true
}
