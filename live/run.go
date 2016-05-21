package live

import (
	"log"
	"os"
	"sync"
	"time"

	"health_monitor/config"
	"health_monitor/utils"
)

// Status holds the status of the internet connectivity after the scan
type Status struct {
	Normal bool
}

var liveStatus Status

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
	var live *Config
	live = loadData()
	log.SetOutput(os.Stdout)
	liveStatus.Normal = true
	Default := live.checkByHEAD
	if live.checkByDNS() {
		Default = live.checkByDNS
	}
	if live.ping() {
		Default = live.ping
	}
	Default()
	for {
		select {
		case signal := <-status:
			if signal.Module == 1 && signal.Run == false {
				return
			}

		case <-time.After(time.Millisecond * time.Duration(live.recheckThreshold)):
			internetCheck(Default, live)
		}
	}
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

// GetLiveStatus function is getter funtion for the liveStatus to send status
// of internet connectivity monitor.
func GetLiveStatus() Status {
	return liveStatus
}
