package api

import (
	"health_monitor/disk"
	"health_monitor/live"
)

var (
	//StatusFunc is a map of all the function which gives json object of module status
	StatusFunc map[string]func() []byte
)

func init() {
	StatusFunc = make(map[string]func() []byte)
	StatusFunc["live"] = live.GetStatusJSON
	StatusFunc["disk"] = disk.GetStatusJSON
}

// GetStatusJSON will return json string of the module provided as a parameter
func GetStatusJSON(module string) []byte {
	return StatusFunc[module]()
}
