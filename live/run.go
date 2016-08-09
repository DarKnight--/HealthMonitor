package live

import (
	"encoding/json"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"health_monitor/notify"
	"health_monitor/owtf"
	"health_monitor/setup"
	"health_monitor/utils"
)

// Status holds the status of the internet connectivity after the scan
type Status struct {
	Normal bool `json:"normal"`
}

var (
	liveStatus Status
	lastStatus Status
	logFile    *os.File
	conf       *Config
)

// Live is the driver function of this module for monitor
func Live(status <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		logFileName = path.Join(setup.ConfigVars.HomeDir, "live.log")
		err         error
	)

	logFile = utils.OpenLogFile(logFileName)
	defer logFile.Close()

	utils.ModuleLogs(logFile, "Running with "+conf.Profile+" profile")
	liveStatus.Normal = true
	Default := conf.CheckByHEAD

	utils.ModuleLogs(logFile, "Default scan mode set to checkByHead")
	if err = conf.CheckByDNS(); err == nil {
		utils.ModuleLogs(logFile, "checkByDNS successful, setting it to default.")
		Default = conf.CheckByDNS
	} else {
		utils.ModuleError(logFile, err.Error(), "Error in checkByDNS")
	}

	if err = conf.Ping(); err == nil {
		utils.ModuleLogs(logFile, "Ping scan successful, setting it to default.")
		Default = conf.Ping
	} else {
		utils.ModuleError(logFile, err.Error(), "Error in Ping")
	}

	internetCheck(Default, conf)
	printStatusLog()

	for {
		select {
		case <-status:
			utils.ModuleLogs(logFile, "Recieved signal to turn off. Signing off")
			return
		case <-time.After(time.Millisecond * time.Duration(conf.RecheckThreshold)):
			internetCheck(Default, conf)
			printStatusLog()
			runtime.Gosched()
		case <-utils.LiveEmergency:
			internetCheck(Default, conf)
			runtime.Gosched()
		}
	}
}

// GetStatus function is getter funtion for the liveStatus to send status
// of internet connectivity monitor.
func GetStatus() Status {
	return liveStatus
}

// GetStatusJSON function retuns the json string of the liveStatus struct
func GetStatusJSON() []byte {
	data, err := json.Marshal(liveStatus)
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the liveStatus struct")
	}
	return data
}

func internetCheck(defaultCheck func() error, conf *Config) {
	var err error
	lastStatus.Normal = liveStatus.Normal
	if err = defaultCheck(); err == nil {
		liveStatus.Normal = true
		owtf.ResumeOWTF(logFile)
		return
	}
	liveStatus.Normal = false
	utils.ModuleError(logFile, err.Error(), "")

	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond * 500)
		if err = conf.CheckByHEAD(); err == nil {
			liveStatus.Normal = true
			return
		}
		utils.ModuleError(logFile, err.Error(), "")
	}
	if lastStatus.Normal {
		owtf.PauseOWTF(logFile)
		notify.SendDesktopAlert("OWTF - Health Monitor", "Your internet connection is down", notify.CRITICAL, "")
	}
	liveStatus.Normal = false
}

func printStatusLog() {
	if liveStatus.Normal {
		utils.ModuleLogs(logFile, "Scan successful, Status : Up")
	} else {
		utils.ModuleLogs(logFile, "Scan successful, Status : Down")
	}
}

//GetConfJSON returns the json byte array of the module's config
func GetConfJSON() []byte {
	data, err := json.Marshal(LoadConfig())
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the conf struct")
	}
	return data
}

//Init is the initialization function of the module
func Init() {
	conf = LoadConfig()
	if conf == nil {
		utils.CheckConf(logFile, setup.MainLogFile, "live", &setup.UserModuleState.Profile, setup.Live)
	}
}
