package setup

import (
	"bytes"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	ModulesStatus struct {
		Profile string
		Live    bool
		Target  bool
		Disk    bool
		Ram     bool
	}
)

func loadStatus() {
	if _, err := os.Stat(ConfigVars.ModuleInfoFilePath); os.IsNotExist(err) {
		log.Println("The module status file is missing. Creating one with default settings")
		initStatus()
		return
	}
	_, err := toml.DecodeFile(ConfigVars.ModuleInfoFilePath, &ModulesStatus) // Read the module status file
	if err != nil {
		log.Println(err)
		log.Println("The module status file is corrupt, creating one with default values")
		initStatus()
		return
	}
}

func initStatus() {
	ModulesStatus.Profile = "default"
	ModulesStatus.Live = true
	ModulesStatus.Target = true
	ModulesStatus.Disk = true
	ModulesStatus.Ram = true

	SaveStatus()
}

func SaveStatus() {
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(ModulesStatus)
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
