package disk

import (
	"encoding/json"
	"fmt"

	"health_monitor/setup"
	"health_monitor/utils"
)

func loadData() *Config {
	var conf *Config = new(Config)
	err := setup.Database.QueryRow("SELECT * FROM Disk WHERE profile=?",
		setup.ConfigVars.Profile).Scan(&conf.Profile, &conf.SpaceWarningLimit,
		&conf.SpaceDangerLimit, &conf.InodeWarningLimit, &conf.InodeDangerLimit,
		&conf.RecheckThreshold, &conf.Disks)
	if err != nil {
		utils.ModuleError(logFile, "Error while quering from databse", err.Error())
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
		utils.ModuleError(logFile, "Unable to insert/update profile", err.Error())
		return err
	}
	utils.ModuleLogs(logFile, fmt.Sprintf("Updated/Inserted the %s profile in db",
		newConf.Profile))
	return nil
}

func SaveConfig(data []byte) error {
	var newConfig = new(Config)
	err := json.Unmarshal(data, newConfig)
	if err != nil {
		utils.ModuleError(logFile, "Unable to decode obtained json.", err.Error())
		return err
	}
	return saveData(newConfig)
}
