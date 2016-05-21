package disk

import (
	"log"
	"os"
	"sync"
	"time"

	"health_monitor/config"
	"health_monitor/utils"
)

// Status holds the status of the disk after the scan
type Status struct {
	SpaceNormal  bool
	SpaceWarning bool
	InodeNormal  bool
	InodeWarning bool
}

var (
	diskStatus map[string]Status
	partition  []string
)

func loadData() *Config {
	var conf Config
	err := config.Database.QueryRow("SELECT * FROM Disk WHERE profile=?",
		config.ConfigVars.Profile).Scan(&conf.profile, &conf.spaceWarningLimit,
		&conf.spaceDangerLimit, &conf.inodeWarningLimit, &conf.inodeDangerLimit,
		&conf.recheckThreshold, &conf.disks)
	if err != nil {
		return nil // TODO better to have fallback call to default profile
	}
	return &conf
}

// Disk is driver funcion for the health_monitor to monitor disk
func Disk(status chan utils.Status, wg *sync.WaitGroup) {
	defer wg.Done()
	var conf *Config
	conf = loadData()
	log.SetOutput(os.Stdout)
	partition = conf.getDisk()
	for {
		select {
		case signal := <-status:
			if signal.Module == 2 && signal.Run == false {
				return
			}

		case <-time.After(time.Millisecond * time.Duration(conf.recheckThreshold)):
			checkDisk(conf)
		}
	}
}

func checkDisk(conf *Config) {
	for _, path := range partition {
		var tempData Status
		tempData.InodeNormal, tempData.InodeWarning = conf.inodesInfo(path)
		tempData.SpaceNormal, tempData.SpaceWarning = conf.diskInfo(path)
		diskStatus["path"] = tempData
	}
}

// GetDiskStatus function is getter funtion for the diskStatus to send status
// of disk monitor.
func GetDiskStatus() map[string]Status {
	return diskStatus
}
