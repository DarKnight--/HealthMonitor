package api

import (
	"health_monitor/cpu"
	"health_monitor/disk"
	"health_monitor/live"
	"health_monitor/ram"
	"health_monitor/target"
)

// CPUStatus returns status and information of the CPU module.
// Information contains CPU usage.
func CPUStatus() cpu.Info {
	return cpu.GetStatus()
}

// DiskStatus returns status and information of the disk module.
// Information contains free and total inode + disk blocks.
func DiskStatus() map[string]disk.PartitionInfo {
	return disk.GetStatus()
}

// LiveStatus returns status of the live module.
func LiveStatus() live.Status {
	return live.GetStatus()
}

// RAMStatus returns status and information of the ram module.
// Information contains free RAM.
func RAMStatus() ram.Info {
	return ram.GetStatus()
}

// TargetStatus returns status of the target module
func TargetStatus() map[string]target.Status {
	return target.GetStatus()
}
