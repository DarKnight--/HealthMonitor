package ram

import (
	"encoding/json"
	"fmt"

	"health_monitor/setup"
	"health_monitor/utils"
)

func LoadConfig() *Config {
	var conf *Config = new(Config)
	err := setup.Database.QueryRow("SELECT * FROM Ram WHERE profile=?",
		setup.ModulesStatus.Profile).Scan(&conf.Profile, &conf.RamWarningLimit,
		&conf.RecheckThreshold)
	if err != nil {
		utils.ModuleError(logFile, "Error while quering from databse", err.Error())
		return nil // TODO better to have fallback call to default profile
	}
	return conf
}

func saveData(newConf *Config) error {
	_, err := setup.Database.Exec(`INSERT OR REPLACE INTO Ram VALUES(?,?,?)`,
		newConf.Profile, newConf.RamWarningLimit, newConf.RecheckThreshold)
	if err != nil {
		utils.ModuleError(logFile, "Module: ram, Unable to insert/update profile", err.Error())
		return err
	}
	utils.ModuleLogs(logFile, fmt.Sprintf("Module: ram, Updated/Inserted the %s profile in db",
		newConf.Profile))
	return nil
}

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
		utils.ModuleError(setup.DBLogFile, "Module: ram, Unable to decode obtained json.", err.Error())
		return err
	}
	conf = newConfig
	return saveData(newConfig)
}
