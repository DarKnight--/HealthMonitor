package config

import (
	"path"
)

var (
	// Logs will hold all the log files names used by health_monitor
	Logs struct {
		HealthMonitorLog string
	}
)

// loadLogs function will initialise the logs variable
func logsInit() bool {
	Logs.HealthMonitorLog = path.Join(ConfigVars.HomeDir, "monitor.log")
	return true
}
