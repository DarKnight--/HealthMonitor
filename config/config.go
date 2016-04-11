package config

import (
	"os"
	"path"
	"log"
)

var (
	HomeDir string = os.Getenv("HOME")
	DB struct{
		DBConfigFile string 
	}
	
	testConfig bool = loadConfig()
	errorMsg string = `The config file is missing. Please run the setup script
						to setup the OWTF Health Monitor.`
)

func loadConfig() bool{
	DB.DBConfigFile = path.Join(os.Getenv("HOME"), ".owtfMonitor", ".config" ,
									"db_config.json")
	if _, err := os.Stat(DB.DBConfigFile); os.IsNotExist(err){
		log.SetOutput(os.Stdout)
		log.Println(errorMsg)
		os.Exit(1)
	}
	return true
}