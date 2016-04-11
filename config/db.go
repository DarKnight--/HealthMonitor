package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	
	"HealthMonitor/utils"
)

var (
	dbParameters struct{
		User		string
		Passord		string
		Database	string	
	}
)

func loadDBParams(){
	file, err := ioutil.ReadFile(DB.DBConfigFile)
	utils.PDBFileError(err)
	err = json.Unmarshal(file, &dbParameters)
	utils.PDBFileError(err)
}

func init() {
	
	file, err := os.OpenFile(Logs.HealthMonitorLog , os.O_RDWR | os.O_CREATE | 
							os.O_APPEND, 0666)
	utils.PLogError(err)
	defer file.Close()
	log.SetOutput(file)	
}