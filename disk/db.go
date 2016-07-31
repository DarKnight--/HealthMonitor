package disk

import (
	"encoding/json"
	"fmt"

	"health_monitor/setup"
	"health_monitor/utils"
)

//LoadConfig load the config of the module from the db
func LoadConfig() *Config {
	var conf = new(Config)
	err := setup.Database.QueryRow("SELECT * FROM Disk WHERE profile=?",
		setup.UserModuleState.Profile).Scan(&conf.Profile, &conf.SpaceWarningLimit,
		&conf.SpaceDangerLimit, &conf.InodeWarningLimit, &conf.InodeDangerLimit,
		&conf.RecheckThreshold, &conf.Disks)
	if err != nil {
		utils.ModuleError(setup.DBLogFile, "Module: disk, Error while quering from databse", err.Error())
		return nil // TODO better to have fallback call to default profile
	}
	return conf
}

func saveData(newConf *Config) error {
	_, err := setup.Database.Exec(`INSERT OR REPLACE INTO Disk VALUES(?,?,?,?,?,?,?)`,
		newConf.Profile, newConf.SpaceWarningLimit, newConf.SpaceDangerLimit,
		newConf.InodeWarningLimit, newConf.InodeDangerLimit, newConf.RecheckThreshold,
		newConf.Disks)
	if err != nil {
		utils.ModuleError(setup.DBLogFile, "Module: disk, Unable to insert/update profile", err.Error())
		return err
	}
	utils.ModuleLogs(setup.DBLogFile, fmt.Sprintf("Module: disk, Updated/Inserted the %s profile in db",
		newConf.Profile))
	return nil
}

//SaveConfig save the config of the module to the database
func SaveConfig(data []byte, profile string) error {
	if data == nil {
		if profile != conf.Profile {
			conf.Profile = profile
			return saveData(conf)
		}
		return nil
	}
	var newConfig = new(Config)
	err := json.Unmarshal(data, newConfig)
	if err != nil {
		utils.ModuleError(setup.DBLogFile, "Module: disk, Unable to decode obtained json.", err.Error())
		return err
	}
	conf = newConfig
	return saveData(newConfig)
}
