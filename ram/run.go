package ram

import (
	"encoding/json"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"health_monitor/setup"
	"health_monitor/utils"
)

type (
	Status struct {
		Normal bool
	}

	Info struct {
		Status Status
		Stats  MemoryStat
		Consts MemoryConst
	}
)

var (
	ramInfo Info
	logFile *os.File
	conf    *Config
)

func Ram(status <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var logFileName = path.Join(setup.ConfigVars.HomeDir, "ram.log")

	logFile = utils.OpenLogFile(logFileName)
	defer logFile.Close()

	utils.ModuleLogs(logFile, "Running with "+conf.Profile+" profile")
	ramInfo.Consts.InitMemoryConst()
	for {
		select {
		case <-status:
			utils.ModuleLogs(logFile, "Recieved signal to turn off. Signing off")
			return
		case <-time.After(time.Millisecond * time.Duration(conf.RecheckThreshold)):
			checkRam()
			runtime.Gosched()
		}
	}
}

func checkRam() {
	ramInfo.Stats.LoadMemoryStats()
	if ramInfo.Stats.FreePhysical < conf.RamWarningLimit {
		ramInfo.Status.Normal = false
		utils.ModuleLogs(logFile, "Ram is being used over the warning limit")
	} else {
		ramInfo.Status.Normal = true
		utils.ModuleLogs(logFile, "Ram usage is normal")
	}
}

func GetStatus() Info {
	return ramInfo
}

func GetConfJSON() []byte {
	data, err := json.Marshal(LoadConfig())
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the conf struct")
	}
	return data
}

func GetStatusJSON() []byte {
	data, err := json.Marshal(ramInfo)
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the ramInfo struct")
	}
	return data
}

func Init() {
	conf = LoadConfig()
}
