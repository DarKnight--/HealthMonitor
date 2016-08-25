package setup

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/owtf/health_monitor/utils"
)

var (
	// ConfigVars will hold necessary variables loaded from config file
	ConfigVars struct {
		HomeDir            string
		DBFile             string
		OWTFAddress        string
		ModuleInfoFilePath string
		Port               string
	}
	// HealthMonitorLog holds the path to main log file
	HealthMonitorLog string
	//OSVarient holds the os name of the current system
	OSVarient string
	//MainLogFile is the file pointer of the monitor.log file
	MainLogFile *os.File
)

func init() {
	var err error
	var basePath = path.Join(os.Getenv("HOME"), ".owtf_monitor")
	// The necessary config file required by health_monitor
	var configFile = path.Join(basePath, "config", "config.toml")
	HealthMonitorLog = path.Join(basePath, "monitor.log")
	os.Mkdir(basePath, 0777)
	MainLogFile = utils.OpenLogFile(HealthMonitorLog)
	log.SetOutput(MainLogFile)

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
	ConfigVars.ModuleInfoFilePath = utils.GetPath(ConfigVars.ModuleInfoFilePath)

	temp, err := exec.Command("lsb_release", "-is").Output()
	if err != nil {
		log.Println("Unable to get os info")
		log.Println(err)
	}
	OSVarient = strings.TrimSpace(string(temp))
	dbInit()
	loadStatus()
}
