package config

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	
	"gopkg.in/pg.v4"
	
	"HealthMonitor/utils"
)

var (
	dbParameters struct{
		User		string
		Password		string
		Database	string
		Address		string
		Port		string	
	}
	
	DBInstance *pg.DB //database instance for use in program
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
	loadDBParams()
	DBInstance = pg.Connect(&pg.Options{
		User: 		dbParameters.User,
		Password: 	dbParameters.Password,
		Database: 	dbParameters.Database,
		Addr:		dbParameters.Address + ":" + dbParameters.Port,	
	})
	fmt.Println(dbParameters.User)
	fmt.Println(dbParameters.Password)
	fmt.Println(dbParameters.Database)
	fmt.Println(dbParameters.Address)
	fmt.Println(dbParameters.Port)
}