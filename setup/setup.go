package setup

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"

	"health_monitor/utils"
)

var (
	// ConfigVars will hold necessary variables loaded from config file
	ConfigVars struct {
		HomeDir     string
		DBFile      string
		OWTFAddress string
		Profile     string
		Port        string
	}
	// HealthMonitorLog holds the path to main log file
	HealthMonitorLog string
	logFile          *os.File
)

func init() {
	var err error
	var basePath = path.Join(os.Getenv("HOME"), ".owtf_monitor")
	// The necessary config file required by health_monitor
	var configFile = path.Join(basePath, "config", "config.toml")
	HealthMonitorLog = path.Join(basePath, "monitor.log")
	os.Mkdir(basePath, 0777)
	logFile, err = os.OpenFile(HealthMonitorLog, os.O_RDWR|os.O_CREATE|
		os.O_APPEND, 0666)
	utils.PLogError(err)
	defer logFile.Close()
	log.SetOutput(logFile)
	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		log.Println("The config file is missing. Creating one with default settings")
		setupConfig()
		return
	}

	_, err = toml.DecodeFile(configFile, &ConfigVars) // Read the config file
	if err != nil {
		log.Println(err)
		log.Println("The config file is corrupt, creating one with default values")
		setupConfig()
		return
	}

	// Update the values if relative path is used
	ConfigVars.HomeDir = utils.GetPath(ConfigVars.HomeDir)
	ConfigVars.DBFile = utils.GetPath(ConfigVars.DBFile)

	dbInit()
	logFile.Close()
}
