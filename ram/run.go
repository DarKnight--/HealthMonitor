package ram

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/owtf/health_monitor/notify"
	"github.com/owtf/health_monitor/owtf"
	"github.com/owtf/health_monitor/setup"
	"github.com/owtf/health_monitor/utils"
)

type (
	// Status holds the status of the RAM after the scan
	Status struct {
		Normal bool
	}
	//Info holds the information of RAM's status, contants and
	// stats after the scan
	Info struct {
		Status Status
		Stats  MemoryStat
		Consts MemoryConst
	}
)

var (
	ramInfo    Info
	logFile    *os.File
	conf       *Config
	lastStatus Status
)

//RAM is the driver function of this module for monitor
func RAM(status <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var logFileName = path.Join(setup.ConfigVars.HomeDir, "ram.log")

	logFile = utils.OpenLogFile(logFileName)
	defer logFile.Close()

	utils.ModuleLogs(logFile, "Running with "+conf.Profile+" profile")
	// initialization of the constants for the module
	conf.InitMemoryConst(&ramInfo.Consts)
	ramInfo.Status.Normal = true
	checkRAM()

	for {
		select {
		case <-status:
			utils.ModuleLogs(logFile, "Received signal to turn off. Signing off")
			return
		case <-time.After(time.Millisecond * time.Duration(conf.RecheckThreshold)):
			checkRAM()
			runtime.Gosched()
		}
	}
}

func percent()

func checkRAM() {
	// storing the value for detecting the change in the status
	lastStatus.Normal = ramInfo.Status.Normal
	conf.LoadMemoryStats(&ramInfo.Stats)

	if ramInfo.Stats.FreePhysical < (100-conf.RAMWarningLimit)*ramInfo.Consts.TotalPhysical/100 {
		ramInfo.Status.Normal = false
		if lastStatus.Normal {
			errorMsg := fmt.Sprintf("Ram usage is above warning limit. Ram usage = %d", utils.Percent(
				ramInfo.Stats.FreePhysical, ramInfo.Consts.TotalPhysical))
			notify.SendEmailAlert("[OWTF-HEALTH-MONITOR]Error in RAM module", errorMsg)
			notify.SendDesktopAlert("OWTF - Health Monitor", "RAM usage is above warn limit.", notify.Critical, "")
			owtf.PauseOWTF(logFile)
		}
		utils.ModuleLogs(logFile, "Ram is being used over the warning limit")
	} else {
		ramInfo.Status.Normal = true
		utils.ModuleLogs(logFile, "Ram usage is normal")
	}
}

// GetStatus function is getter function for the ramInfo to send status
// of ram monitor
func GetStatus() Info {
	return ramInfo
}

//GetConfJSON returns the json byte array of the module's config
func GetConfJSON() []byte {
	data, err := json.Marshal(LoadConfig())
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the conf struct")
	}
	return data
}

// GetStatusJSON function returns the json string of the ramInfo struct
func GetStatusJSON() []byte {
	data, err := json.Marshal(ramInfo)
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the ramInfo struct")
	}
	return data
}

//Init is the initialization function of the module
func Init() {
	conf = LoadConfig()
	if conf == nil {
		utils.CheckConf(logFile, setup.MainLogFile, "ram", &setup.UserModuleState.Profile, setup.RAM)
	}
}
