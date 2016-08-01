package setup

import (
	"bytes"
	"log"
	"os"

	"health_monitor/utils"

	"github.com/BurntSushi/toml"
)

// ModulesStatus holds the status of module
type ModulesStatus struct {
	Profile string
	Live    bool
	Target  bool
	Disk    bool
	RAM     bool
	CPU     bool
}

var (
	//UserModuleState holds the status of all the modules of monitor set by user
	UserModuleState = ModulesStatus{}
	//InternalModuleState holds the running status of all the modules of monitor
	InternalModuleState = ModulesStatus{}
)

func loadStatus() {
	if _, err := os.Stat(ConfigVars.ModuleInfoFilePath); os.IsNotExist(err) {
		utils.ModuleError(MainLogFile, "The module status file is missing.", "Creating one with default settings")
		initStatus()
		return
	}
	_, err := toml.DecodeFile(ConfigVars.ModuleInfoFilePath, &UserModuleState) // Read the module status file
	if err != nil {
		utils.ModuleError(MainLogFile, "The module status file is corrupt, creating one with default values", err.Error())
		initStatus()
	} else if UserModuleState.Profile == "" { //TODO add check to ensure profile exists in db
		utils.ModuleError(MainLogFile, "The module status file does not contain profile or profile does not exists", "Creating one with default values")
		initStatus()
	}
	InternalModuleState = UserModuleState
}

func initStatus() {
	UserModuleState.Profile = "default"
	UserModuleState.Live = true
	UserModuleState.Target = true
	UserModuleState.Disk = true
	UserModuleState.RAM = true
	UserModuleState.CPU = true
	SaveStatus()
}

//SaveStatus saves the status of all the modules to disk
func SaveStatus() {
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(UserModuleState)
	log.SetOutput(MainLogFile)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(ConfigVars.ModuleInfoFilePath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write(buffer.Bytes())
}
