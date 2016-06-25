package cpu

import (
	"encoding/json"
	"fmt"

	"health_monitor/setup"
	"health_monitor/utils"
)

//LoadConfig load the config of the module from the db
func LoadConfig() *Config {
	var conf = new(Config)
	err := setup.Database.QueryRow("SELECT * FROM CPU WHERE profile=?",
		setup.ModulesStatus.Profile).Scan(&conf.Profile, &conf.CPUWarningLimit,
		&conf.RecheckThreshold)
	if err != nil {
		utils.ModuleError(logFile, "Error while quering from databse", err.Error())
		return nil
	}
	return conf
}

func saveData(newConf *Config) error {
	_, err := setup.Database.Exec(`INSERT OR REPLACE INTO CPU VALUES(?,?,?)`,
		newConf.Profile, newConf.CPUWarningLimit, newConf.RecheckThreshold)
	if err != nil {
		utils.ModuleError(logFile, "Module: cpu, Unable to insert/update profile", err.Error())
		return err
	}
	utils.ModuleLogs(logFile, fmt.Sprintf("Module: cpu, Updated/Inserted the %s profile in db",
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
		utils.ModuleError(setup.DBLogFile, "Module: cpu, Unable to decode obtained json.", err.Error())
		return err
	}
	conf = newConfig
	return saveData(newConfig)
}
