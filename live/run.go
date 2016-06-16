package live

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

// Status holds the status of the internet connectivity after the scan
type Status struct {
	Normal bool `json:"normal"`
}

var (
	liveStatus Status
	logFile    *os.File
	conf       *Config
)

// Live is the driver function of this module for monitor
func Live(status <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		logFileName = path.Join(setup.ConfigVars.HomeDir, "live.log")
		err         error
		x           bool
	)

	logFile = utils.OpenLogFile(logFileName)
	defer logFile.Close()

	utils.ModuleLogs(logFile, "Running with "+conf.Profile+" profile")
	liveStatus.Normal = true
	Default := conf.CheckByHEAD

	utils.ModuleLogs(logFile, "Default scan mode set to checkByHead")
	if x, err = conf.CheckByDNS(); x {
		utils.ModuleLogs(logFile, "checkByDNS successful, setting it to default.")
		Default = conf.CheckByDNS
	}
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "Error in checkByDNS")
	}

	if x, err = conf.Ping(); x {
		utils.ModuleLogs(logFile, "Ping scan successful, setting it to default.")
		Default = conf.Ping
	}
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "Error in Ping")
	}

	Default()
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

func internetCheck(defaultCheck func() (bool, error), conf *Config) {
	var (
		err error
		x   bool
	)
	if x, err = defaultCheck(); x {
		liveStatus.Normal = true
		return
	}
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "")
	}

	for i := 0; i < 3; i++ {
		time.Sleep(time.Duration(conf.RecheckThreshold) * time.Millisecond / 5)
		if x, err = conf.CheckByHEAD(); x {
			liveStatus.Normal = true
			return
		}
		if err != nil {
			utils.ModuleError(logFile, err.Error(), "")
		}
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

func GetConfJSON() []byte {
	data, err := json.Marshal(LoadConfig())
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the conf struct")
	}
	return data
}

func Init() {
	conf = LoadConfig()
}
