package config

import (
	"os"
	"path"
)
var (
	Logs struct{
		HealthMonitorLog string
	}
	testLog bool =  loadLogParams()
)

func loadLogParams() bool {
	Logs.HealthMonitorLog = path.Join(os.Getenv("HOME"), ".owtfMonitor",
									"monitor.log")
	return true
}
									