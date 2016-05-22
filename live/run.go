package live

import (
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"health_monitor/config"
	"health_monitor/utils"
)

// Status holds the status of the internet connectivity after the scan
type Status struct {
	Normal bool
}

var (
	liveStatus Status
	logFile    *os.File
)

func loadData() *Config {
	var l Config
	err := config.Database.QueryRow("SELECT * FROM Live WHERE profile=?",
		config.ConfigVars.Profile).Scan(&l.profile, &l.headURL, &l.recheckThreshold,
		&l.pingThreshold, &l.headThreshold, &l.pingAddress, &l.pingProtocol)
	if err != nil {
		return nil // TODO better to have fallback call to default profile
	}
	return &l
}

// Live is the driver function of this module for monitor
func Live(status chan utils.Status, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		logFileName = path.Join(config.ConfigVars.HomeDir, "live.log")
		err         error
		live        *Config
	)

	logFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		utils.PLogError(err)
	}
	defer logFile.Close()

	live = loadData()
	utils.ModuleLogs(logFile, "Loaded "+live.profile+" profile successfully")
	liveStatus.Normal = true
	Default := live.checkByHEAD

	utils.ModuleLogs(logFile, "Default scan mode set to checkByHead")
	if live.checkByDNS() {
		utils.ModuleLogs(logFile, "checkByDNS successful, setting it to default.")
		Default = live.checkByDNS
	}
	if live.ping() {
		utils.ModuleLogs(logFile, "Ping scan successful, setting it to default.")
		Default = live.ping
	}
	Default()
	printStatusLog()

	for {
		select {
		case signal := <-status:
			if signal.Module == 1 && signal.Run == false {
				return
			}

		case <-time.After(time.Millisecond * time.Duration(live.recheckThreshold)):
			internetCheck(Default, live)
			printStatusLog()
			runtime.Gosched()
		}
	}
}

// GetLiveStatus function is getter funtion for the liveStatus to send status
// of internet connectivity monitor.
func GetLiveStatus() Status {
	return liveStatus
}

func internetCheck(defaultCheck func() bool, live *Config) {
	if defaultCheck() {
		liveStatus.Normal = true
		return
	}

	for i := 0; i < 3; i++ {
		time.Sleep(time.Duration(live.recheckThreshold) * time.Millisecond / 5)
		if live.checkByHEAD() {
			liveStatus.Normal = true
			return
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
