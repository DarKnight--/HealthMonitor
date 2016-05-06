package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"

	"health_monitor/utils"
)

var (
	dbParameters struct {
		User     string
		Password string
		Database string
		Address  string
		Port     string
	}

	DBInstance *pg.DB //database instance for use in other modules
)

func loadDBParams() {
	_, err := toml.DecodeFile(DB.DBConfigFile, &dbParameters)
	utils.PDBFileError(err) // TODO if this function is not used elsewhere remove it from utils
}

func init() {
	file, err := os.OpenFile(Logs.HealthMonitorLog, os.O_RDWR|os.O_CREATE|
		os.O_APPEND, 0666)
	utils.PLogError(err)
	defer file.Close()
	log.SetOutput(file)
	loadDBParams()
	DBInstance = pg.Connect(&pg.Options{
		User:     dbParameters.User,
		Password: dbParameters.Password,
		Database: dbParameters.Database,
		Addr:     dbParameters.Address + ":" + dbParameters.Port,
	})
}
