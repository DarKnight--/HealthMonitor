package config

import (
	"path"
	"os"
)
var(
	Logs struct{
		HealthMonitorLog string
	}
	testLogs bool = loadLogs()
)

// loadLogs function will initialise the logs variable
func loadLogs() bool {
	Logs.HealthMonitorLog = path.Join(os.Getenv("HOME"), "monitor.log")
	return true
}