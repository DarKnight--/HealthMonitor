package disk

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"health_monitor/config"
	"health_monitor/utils"
)

type (
	// PartitionStatus holds the status of the partition after the scan
	PartitionStatus struct {
		Inode int
		Space int
	}
	// PartitionInfo holds the information of partition's status, contants and
	// stats after the scan
	PartitionInfo struct {
		Status PartitionStatus
		Stats  PartitionStats
		Const  PartitionConst
	}
)

var (
	diskInfo  map[string]PartitionInfo
	partition []string
	logFile   *os.File
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
	var (
		logFileName = path.Join(config.ConfigVars.HomeDir, "disk.log")
		err         error
		conf        *Config
	)

	logFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		utils.PLogError(err)
	}
	defer logFile.Close()

	conf = loadData()
	utils.ModuleLogs(logFile, "Loaded "+conf.profile+" profile successfully")
	partition = conf.GetDisk()
	diskInfo = make(map[string]PartitionInfo)
	loadPartitionConst()
	checkDisk(conf)

	for {
		select {
		case signal := <-status:
			if signal.Module == 2 && signal.Run == false {
				return
			}

		case <-time.After(time.Millisecond * time.Duration(conf.recheckThreshold)):
			checkDisk(conf)
			runtime.Gosched()
		}
	}
}

func checkDisk(conf *Config) {
	for _, directory := range partition {
		var tempStatus PartitionStatus
		var tempStat PartitionStats
		tempStatus.Inode = conf.InodesInfo(directory, &tempStat)
		tempStatus.Space = conf.SpaceInfo(directory, &tempStat)
		diskInfo[directory] = PartitionInfo{tempStatus, tempStat, diskInfo[directory].Const}
		printStatusLog(directory, tempStatus.Inode, "inode")
		printStatusLog(directory, tempStatus.Space, "space")
		utils.ModuleLogs(logFile, "Stats for mount "+directory+" :")
		utils.ModuleLogs(logFile, fmt.Sprintf("Inodes: \t Total: %d \t Free: %d",
			diskInfo[directory].Const.TotalInodes, tempStat.FreeInodes))
		utils.ModuleLogs(logFile, fmt.Sprintf("Blocks: \t Total: %d \t Free: %d",
			diskInfo[directory].Const.TotalBlocks, tempStat.FreeBlocks))
	}
}

// GetDiskStatus function is getter funtion for the diskStatus to send status
// of disk monitor.
func GetDiskStatus() map[string]PartitionInfo {
	return diskInfo
}

func loadPartitionConst() {
	for _, directory := range partition {
		diskInfo[directory] = PartitionInfo{PartitionStatus{}, PartitionStats{},
			SetPartitionConst(directory)}
	}
}

func printStatusLog(directory string, status int, types string) {
	switch status {
	case -1:
		utils.ModuleError(logFile, fmt.Sprintf("Unable to retrieve the informtaion about %s mount point",
			directory), "Check the mount point provided")

	case 1:
		utils.ModuleLogs(logFile, fmt.Sprintf("Mount point %s %s status : OK",
			directory, types))

	case 2:
		utils.ModuleLogs(logFile, fmt.Sprintf("Mount point %s %s status : WARN",
			directory, types))

	case 3:
		utils.ModuleLogs(logFile, fmt.Sprintf("Mount point %s %s status : Danger",
			directory, types))
	}
}
