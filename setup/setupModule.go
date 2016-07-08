package setup

import (
	"bytes"
	"log"
	"os"

	"health_monitor/utils"

	"github.com/BurntSushi/toml"
)

var (
	//ModulesStatus holds the running status of all the modules of monitor
	ModulesStatus struct {
		Profile string
		Live    bool
		Target  bool
		Disk    bool
		RAM     bool
		CPU     bool
	}
)

func loadStatus() {
	if _, err := os.Stat(ConfigVars.ModuleInfoFilePath); os.IsNotExist(err) {
		utils.ModuleError(MainLogFile, "The module status file is missing.", "Creating one with default settings")
		initStatus()
		return
	}
	_, err := toml.DecodeFile(ConfigVars.ModuleInfoFilePath, &ModulesStatus) // Read the module status file
	if err != nil {
		utils.ModuleError(MainLogFile, "The module status file is corrupt, creating one with default values", err.Error())
		initStatus()
	} else if ModulesStatus.Profile == "" { //TODO add check to ensure profile exists in db
		utils.ModuleError(MainLogFile, "The module status file does not contain profile or profile does not exists", "Creating one with default values")
		initStatus()
	}
}

func initStatus() {
	ModulesStatus.Profile = "default"
	ModulesStatus.Live = true
	ModulesStatus.Target = true
	ModulesStatus.Disk = true
	ModulesStatus.RAM = true
	ModulesStatus.CPU = true

	SaveStatus()
}

//SaveStatus saves the status of all the modules to disk
func SaveStatus() {
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(ModulesStatus)
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
