package api

import (
	"health_monitor/cpu"
	"health_monitor/disk"
	"health_monitor/live"
	"health_monitor/ram"
	"health_monitor/target"
)

func CPUStatus() cpu.Info {
	return cpu.GetStatus()
}

func DiskStatus() map[string]disk.PartitionInfo {
	return disk.GetStatus()
}

func LiveStatus() live.Status {
	return live.GetStatus()
}

func RAMStatus() ram.Info {
	return ram.GetStatus()
}

func TargetStatus() map[string]target.Status {
	return target.GetStatus()
}
