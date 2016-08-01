package api

import (
	"encoding/json"

	"health_monitor/cpu"
	"health_monitor/disk"
	"health_monitor/live"
	"health_monitor/ram"
	"health_monitor/setup"
	"health_monitor/target"
	"health_monitor/utils"
)

var (
	//StatusFunc is a map of all the function which gives json byte array of module status
	StatusFunc map[string]func() []byte
	//ConfFunc is the map of all the function which gives json byte array of module config
	ConfFunc map[string]func() []byte
	//ConfSaveFunc is the map of all the function which save the module config to database
	ConfSaveFunc map[string]func([]byte, string) error
)

func init() {
	StatusFunc = make(map[string]func() []byte)
	StatusFunc["live"] = live.GetStatusJSON
	StatusFunc["target"] = target.GetStatusJSON
	StatusFunc["disk"] = disk.GetStatusJSON
	StatusFunc["ram"] = ram.GetStatusJSON
	StatusFunc["cpu"] = cpu.GetStatusJSON

	ConfFunc = make(map[string]func() []byte)
	ConfFunc["live"] = live.GetConfJSON
	ConfFunc["target"] = target.GetConfJSON
	ConfFunc["disk"] = disk.GetConfJSON
	ConfFunc["ram"] = ram.GetConfJSON
	ConfFunc["cpu"] = cpu.GetConfJSON

	ConfSaveFunc = make(map[string]func([]byte, string) error)
	ConfSaveFunc["live"] = live.SaveConfig
	ConfSaveFunc["target"] = target.SaveConfig
	ConfSaveFunc["disk"] = disk.SaveConfig
	ConfSaveFunc["ram"] = ram.SaveConfig
	ConfSaveFunc["cpu"] = cpu.SaveConfig
}

// GetStatusJSON will return json string of the status of module provided as a parameter
func GetStatusJSON(module string) []byte {
	return StatusFunc[module]()
}

// GetConfJSON will return json string of the config of module provided as a parameter
func GetConfJSON(module string) []byte {
	return ConfFunc[module]()
}

//SaveConfig saves the config obtained to the database and load it
func SaveConfig(module string, data []byte) error {
	profile := getProfile(data)
	err := ConfSaveFunc[module](data, profile)
	if profile == setup.UserModuleState.Profile {
		return err
	}
	for _, function := range ConfSaveFunc {
		err := function(nil, profile)
		if err != nil {
			return err
		}
	}
	setup.UserModuleState.Profile = profile
	return nil
}

func getProfile(data []byte) string {
	var Temp struct {
		Profile string
	}
	json.Unmarshal(data, &Temp)
	return Temp.Profile
}

//ChangeModuleStatus sends the signal to main function about the changing the status of module
func ChangeModuleStatus(module string, status bool) {
	switch module {
	case "live":
		setup.UserModuleState.Live = status
	case "target":
		if status {
			utils.AddOWTFModuleDependence()
		} else {
			utils.RemoveOWTFModuleDependence()
		}
		setup.UserModuleState.Target = status
	case "disk":
		setup.UserModuleState.Disk = status
	case "inode":
		setup.UserModuleState.Disk = status
	case "ram":
		setup.UserModuleState.RAM = status
	case "cpu":
		setup.UserModuleState.CPU = status
	default:
		return
	}
	utils.SendModuleStatus(module, status)
}

//ModuleStatus return the running status of the given module.
func ModuleStatus(module string) bool {
	switch module {
	case "live":
		return setup.InternalModuleState.Live
	case "target":
		return setup.InternalModuleState.Target
	case "disk":
		return setup.InternalModuleState.Disk
	case "inode":
		return setup.InternalModuleState.Disk
	case "ram":
		return setup.InternalModuleState.RAM
	case "cpu":
		return setup.InternalModuleState.CPU
	default:
		return false
	}
}
