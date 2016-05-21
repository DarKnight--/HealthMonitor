package disk

import (
	"fmt"
	"strings"
	"syscall"

	"health_monitor/utils"
)

// Config holds all the necessary parameters required by the module
type Config struct {
	profile           string
	spaceWarningLimit int
	spaceDangerLimit  int
	inodeWarningLimit int
	inodeDangerLimit  int
	recheckThreshold  int
	disks             string
}

func (conf Config) inodesInfo(directory string) (bool, bool) { // It will work only on linux
	var stat syscall.Statfs_t
	err := syscall.Statfs(directory, &stat)
	if err != nil {
		utils.Perror(fmt.Sprintf("Unable to retrieve disk information about %s",
			directory))
		return false, false
	}

	return compareLimit(int(stat.Ffree), conf.inodeWarningLimit,
		conf.inodeDangerLimit)
}

func (conf Config) diskInfo(directory string) (bool, bool) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(directory, &stat)
	if err != nil {
		utils.Perror(fmt.Sprintf("Unable to retrieve disk information about %s",
			directory))
		return false, false
	}

	return compareLimit(int(stat.Bfree), conf.spaceWarningLimit,
		conf.spaceDangerLimit)
}

func (conf Config) getDisk() []string {
	return strings.Split(conf.disks, ",")
}

func compareLimit(value int, wLimit int, dLimit int) (bool, bool) {
	if value > wLimit { // disk is safe and has enough space
		return true, false
	}

	if value > dLimit { // Warning limit reached
		return true, true
	}

	return false, true // disk is running out of inodes, signal to free them
}
