package api

import (
	"health_monitor/disk"
	"health_monitor/live"
	"health_monitor/utils"
)

var (
	//StatusFunc is a map of all the function which gives json object of module status
	StatusFunc   map[string]func() []byte
	ConfFunc     map[string]func() []byte
	ConfSaveFunc map[string]func([]byte) error
	ControlChan  chan utils.Status
)

func init() {
	StatusFunc = make(map[string]func() []byte)
	StatusFunc["live"] = live.GetStatusJSON
	StatusFunc["disk"] = disk.GetStatusJSON

	ConfFunc = make(map[string]func() []byte)
	ConfFunc["live"] = live.GetConfJSON
	ConfFunc["disk"] = disk.GetConfJSON

	ConfSaveFunc = make(map[string]func([]byte) error)
	ConfSaveFunc["live"] = live.SaveConfig
	ConfSaveFunc["disk"] = disk.SaveConfig
}

// GetStatusJSON will return json string of the status of module provided as a parameter
func GetStatusJSON(module string) []byte {
	return StatusFunc[module]()
}

// GetConfJSON will return json string of the config of module provided as a parameter
func GetConfJSON(module string) []byte {
	return ConfFunc[module]()
}

func SaveConfig(module string, data []byte) error {
	return ConfSaveFunc[module](data)
}
