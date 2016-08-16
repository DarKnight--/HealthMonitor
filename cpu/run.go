package cpu

import (
	"encoding/json"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"health_monitor/notify"
	"health_monitor/setup"
	"health_monitor/utils"
)

type (
	// Status holds the status of the CPU after the scan
	Status struct {
		Normal bool
	}
	//Info holds the information of CPU's status and stats after the scan
	Info struct {
		Status Status
		Stats  Stat
	}
)

var (
	cpuInfo    Info
	logFile    *os.File
	conf       *Config
	lastStatus Status
)

//CPU is the driver function of this module for monitor
func CPU(status <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var logFileName = path.Join(setup.ConfigVars.HomeDir, "cpu.log")

	logFile = utils.OpenLogFile(logFileName)
	defer logFile.Close()

	utils.ModuleLogs(logFile, "Running with "+conf.Profile+" profile")
	conf.Init()
	time.Sleep(time.Second)
	cpuInfo.Status.Normal = true
	checkCPU()

	for {
		select {
		case <-status:
			utils.ModuleLogs(logFile, "Recieved signal to turn off. Signing off")
			return
		case <-time.After(time.Millisecond * time.Duration(conf.RecheckThreshold)):
			checkCPU()
			runtime.Gosched()
		}
	}
}

func checkCPU() {
	lastStatus.Normal = cpuInfo.Status.Normal
	conf.CPUUsage(&cpuInfo.Stats)
	if cpuInfo.Stats.CPUUsage < conf.CPUWarningLimit {
		cpuInfo.Status.Normal = true
		utils.ModuleLogs(logFile, "CPU usage is normal")
	} else {
		if lastStatus.Normal {
			notify.SendDesktopAlert("OWTF - Health Monitor", "CPU usage is above warn limit", notify.Critical, "")
		}
		cpuInfo.Status.Normal = false
		utils.ModuleLogs(logFile, "CPU is being used over the warning limit")
	}
}

// GetStatus function is getter funtion for the cpuInfo to send status
// of cpu monitor
func GetStatus() Info {
	return cpuInfo
}

//GetConfJSON returns the json byte array of the module's config
func GetConfJSON() []byte {
	data, err := json.Marshal(LoadConfig())
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the conf struct")
	}
	return data
}

// GetStatusJSON function retuns the json string of the cpuInfo struct
func GetStatusJSON() []byte {
	data, err := json.Marshal(cpuInfo)
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the cpuInfo struct")
	}
	return data
}

//Init is the initialization function of the module
func Init() {
	conf = LoadConfig()
	if conf == nil {
		utils.CheckConf(logFile, setup.MainLogFile, "cpu", &setup.UserModuleState.Profile, setup.CPU)
	}
}
