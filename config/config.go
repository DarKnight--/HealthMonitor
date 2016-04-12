package config

import (
	"fmt"
	"os"
	"path"
)

var (
	HomeDir string = os.Getenv("HOME") //home directory path
	
	DB struct{	//structure containing database related data
		DBConfigFile string 
	}
	
	testConfig bool = loadConfig()
	errorMsg string = `The config file is missing. Please run the setup script
						to setup the OWTF Health Monitor.`
)

// This function will iniailise all the configuration variable defined
func loadConfig() bool{
	DB.DBConfigFile = path.Join(os.Getenv("HOME"), ".owtfMonitor", "config" ,
									"db_config.json")
	if _, err := os.Stat(DB.DBConfigFile); os.IsNotExist(err){
		fmt.Println(errorMsg)
		os.Exit(1)
	}
	return true
}