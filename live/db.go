package live

import (
	"encoding/json"
	"fmt"

	"health_monitor/setup"
	"health_monitor/utils"
)

func LoadConfig() *Config {
	var conf *Config = new(Config)
	err := setup.Database.QueryRow("SELECT * FROM Live WHERE profile=?",
		setup.ModulesStatus.Profile).Scan(&conf.Profile, &conf.HeadURL,
		&conf.RecheckThreshold, &conf.PingThreshold, &conf.HeadThreshold,
		&conf.PingAddress, &conf.PingProtocol)
	if err != nil {
		utils.ModuleError(setup.DBLogFile, "Error while quering from databse", err.Error())
		return nil // TODO better to have fallback call to default profile
	}
	return conf
}

func saveData(newConf *Config) error {
	_, err := setup.Database.Exec(`INSERT OR REPLACE INTO Live VALUES(?,?,?,?,?,?,?)`,
		newConf.Profile, newConf.HeadURL, newConf.RecheckThreshold,
		newConf.PingThreshold, newConf.HeadThreshold, newConf.PingAddress, newConf.PingProtocol)
	if err != nil {
		utils.ModuleError(setup.DBLogFile, "Module: live, Unable to insert/update profile",
			err.Error())
		return err
	}
	utils.ModuleLogs(setup.DBLogFile, fmt.Sprintf("Module: live, Updated/Inserted the config of %s profile in db",
		newConf.Profile))
	return nil
}

func SaveConfig(data []byte, profile string) error {
	if len(data) == 0 {
		if profile != conf.Profile {
			conf.Profile = profile
			return saveData(conf)
		}
		return nil
	}
	var newConfig = new(Config)
	err := json.Unmarshal(data, newConfig)
	if err != nil {
		utils.ModuleError(setup.DBLogFile, "Module: live, Unable to decode obtained json.", err.Error())
		return err
	}
	conf = newConfig
	return saveData(newConfig)
}
