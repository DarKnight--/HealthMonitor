package disk

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/owtf/health_monitor/notify"
	"github.com/owtf/health_monitor/setup"
	"github.com/owtf/health_monitor/utils"
)

type (
	// PartitionStatus holds the status of the partition after the scan
	PartitionStatus struct {
		Inode int `json:"inode"`
		Space int `json:"space"`
	}
	// PartitionInfo holds the information of partition's status, contants and
	// stats after the scan
	PartitionInfo struct {
		Status PartitionStatus `json:"partition_status"`
		Stats  PartitionStats  `json:"partition_stats"`
		Const  PartitionConst  `json:"partition_consts"`
	}
)

var (
	diskInfo   map[string]PartitionInfo
	partition  []string
	logFile    *os.File
	conf       *Config
	diskAction *BasicCleaner
)

// Disk is driver funcion for the health_monitor to monitor disk
func Disk(status <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	logFileName := path.Join(setup.ConfigVars.HomeDir, "disk.log")
	logFile = utils.OpenLogFile(logFileName)
	defer logFile.Close()

	utils.ModuleLogs(logFile, "Running with "+conf.Profile+" profile")
	diskAction = NewBasicCleaner(setup.OSVarient)
	if diskAction.cleaner == nil {
		utils.ModuleError(logFile, "disk cleaner not implemented for this OS",
			"Please raise a issue on github")
	}
	partition = conf.GetDisk()
	diskInfo = make(map[string]PartitionInfo)
	loadPartitionConst()
	checkDisk(conf)

	for {
		select {
		case <-status:
			utils.ModuleLogs(logFile, "Received signal to turn off. Signing off")
			return
		case <-time.After(time.Millisecond * time.Duration(conf.RecheckThreshold)):
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
		printStatusLog(directory, tempStatus.Inode, diskInfo[directory].Status.Inode, "inode")
		printStatusLog(directory, tempStatus.Space, diskInfo[directory].Status.Space, "space")
		diskInfo[directory] = PartitionInfo{Status: tempStatus, Stats: tempStat, Const: diskInfo[directory].Const}
		utils.ModuleLogs(logFile, "Stats for mount "+directory+" :")
		utils.ModuleLogs(logFile, fmt.Sprintf("Inodes: \t Total: %d \t Free: %d",
			diskInfo[directory].Const.TotalInodes, tempStat.FreeInodes))
		utils.ModuleLogs(logFile, fmt.Sprintf("Blocks: \t Total: %d \t Free: %d",
			diskInfo[directory].Const.TotalBlocks, tempStat.FreeBlocks))
	}
}

// GetStatus function is getter function for the diskStatus to send status
// of disk monitor.
func GetStatus() map[string]PartitionInfo {
	return diskInfo
}

// GetStatusJSON function returns the json string of the diskInfo struct
func GetStatusJSON() []byte {
	data, err := json.Marshal(diskInfo)
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the diskStatus struct")
	}
	return data
}

func loadPartitionConst() {
	for _, directory := range partition {
		diskInfo[directory] = PartitionInfo{Status: PartitionStatus{Inode: 1, Space: 1}, Stats: PartitionStats{},
			Const: SetPartitionConst(directory)}
	}
}

func printStatusLog(directory string, status int, lastStatus int, types string) {
	switch status {
	case -1:
		utils.ModuleError(logFile, fmt.Sprintf("Unable to retrieve the informtaion about %s mount point",
			directory), "Check the mount point provided")
		if lastStatus != -1 {
			errorMsg := fmt.Sprintf("Disk %s for %s mount point can't be scanned due to error.", types, directory)
			notify.SendDesktopAlert("OWTF - Health Monitor", errorMsg, notify.Normal, "")
		}
	case 1:
		utils.ModuleLogs(logFile, fmt.Sprintf("Mount point %s %s status : OK",
			directory, types))

	case 2:
		utils.ModuleLogs(logFile, fmt.Sprintf("Mount point %s %s status : WARN",
			directory, types))
		if lastStatus < 2 {
			errorMsg := fmt.Sprintf("Disk %s for %s mount point status is above warning limit.", types, directory)
			notify.SendDesktopAlert("OWTF - Health Monitor", errorMsg, notify.Normal, "")
		}
	case 3:
		utils.ModuleLogs(logFile, fmt.Sprintf("Mount point %s %s status : Danger",
			directory, types))
		if lastStatus != 3 {
			errorMsg := fmt.Sprintf("Disk %s for %s mount point status is critical", types, directory)
			notify.SendDesktopAlert("OWTF - Health Monitor", errorMsg, notify.Critical, "")
			notify.SendEmailAlert("[OWTF-HEALTH-MONITOR]Error in disk module", errorMsg)
			BasicAction(directory)
		}
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
		utils.CheckConf(logFile, setup.MainLogFile, "disk", &setup.UserModuleState.Profile, setup.Disk)
	}
}

// BasicAction takes required action on directory to free space and inodes
func BasicAction(directory string) {
	if diskAction.cleaner == nil {
		return
	}
	if directory == "/" || strings.Contains(directory, os.Getenv("HOME")) {
		basicCleanup(*diskAction)
	}
}

// CleanTrash cleans the trash folder
func CleanTrash() error {
	if diskAction.cleaner == nil {
		return fmt.Errorf("Currently %s OS is not supported", diskAction.os)
	}
	return diskAction.cleaner.EmptyTrash()
}

// CleanPackageManagerCache cleans the package manager's cache directory
func CleanPackageManagerCache() error {
	if diskAction.cleaner == nil {
		return fmt.Errorf("Currently %s OS is not supported", diskAction.os)
	}
	return diskAction.cleaner.RemovePackageManagerCache()
}
