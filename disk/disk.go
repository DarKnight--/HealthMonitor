package disk

// It will work only on linux

import (
	"strings"
	"syscall"
)

type (
	// Config holds all the necessary parameters required by the module
	Config struct {
		Profile           string
		SpaceWarningLimit int
		SpaceDangerLimit  int
		InodeWarningLimit int
		InodeDangerLimit  int
		RecheckThreshold  int
		Disks             string
	}

	// PartitionStats holds the data about the remaining inodes and blocks of the
	// mount points
	PartitionStats struct {
		FreeInodes int `json:"free_inodes"`
		FreeBlocks int `json:"free_blocks"`
	}
	// PartitionConst holds the constant data of the mount point
	PartitionConst struct {
		TotalInodes int `json:"total_inodes"`
		TotalBlocks int `json:"total_blocks"`
	}
)

// InodesInfo will return the status of the inodes availabe based on the
// Warning and danger limit set in the config.
// -1 will refer to error, 1 for above warning limit, 2 & 3 for warning and
// danger limit respectively
func (conf Config) InodesInfo(directory string, pStats *PartitionStats) int {
	var stat syscall.Statfs_t
	err := syscall.Statfs(directory, &stat)
	if err != nil {
		return -1
	}

	pStats.FreeInodes = int(stat.Ffree)

	return compareLimit(int(stat.Ffree), conf.InodeWarningLimit,
		conf.InodeDangerLimit)
}

// SpaceInfo will return the status of the space availabe based on the
// Warning and danger limit set in the config.
// -1 will refer to error, 1 for above warning limit, 2 & 3 for warning and
// danger limit respectively
func (conf Config) SpaceInfo(directory string, pStats *PartitionStats) int {
	var stat syscall.Statfs_t
	err := syscall.Statfs(directory, &stat)
	if err != nil {
		return -1
	}

	pStats.FreeBlocks = int(stat.Bfree)

	return compareLimit(int(stat.Bfree), conf.SpaceWarningLimit,
		conf.SpaceDangerLimit)
}

// GetDisk will return all the mount points set to monitor
func (conf Config) GetDisk() []string {
	return strings.Split(conf.Disks, ",")
}

func compareLimit(value int, wLimit int, dLimit int) int {
	if value > wLimit { // disk is safe and has enough space
		return 1
	}

	if value > dLimit { // Warning limit reached
		return 2
	}

	return 3 // disk is running out of inodes, signal to free them
}

// SetPartitionConst will return the TotalInodes and TotalBlocks for the given
// mount point
func SetPartitionConst(directory string) PartitionConst {
	var stat syscall.Statfs_t
	err := syscall.Statfs(directory, &stat)
	if err != nil {
		return PartitionConst{}
	}
	var partitionConst PartitionConst
	partitionConst.TotalBlocks = int(stat.Blocks)
	partitionConst.TotalInodes = int(stat.Files)
	return partitionConst
}
