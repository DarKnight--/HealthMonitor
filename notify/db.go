package notify

import (
	"encoding/json"
	"fmt"

	"github.com/owtf/health_monitor/setup"
	"github.com/owtf/health_monitor/utils"
)

//LoadConfig load the config of the module from the db
func LoadConfig() *Config {
	var conf = new(Config)
	err := setup.Database.QueryRow("SELECT * FROM Alert WHERE profile=?",
		setup.UserModuleState.Profile).Scan(&conf.Profile, &conf.SendgridAPIKey,
		&conf.ElasticMailKey, &conf.ElasticMainUName, &conf.MailjetPublicKey,
		&conf.MailjetSecretKey, &conf.SendEmailTo, &conf.MailgunDomain, &conf.MailgunPrivateKey,
		&conf.MailgunPublicKey, &conf.SendDesktopNotific, &conf.MailOptionToUse, &conf.IconPath)
	if err != nil {
		utils.ModuleError(setup.DBLogFile, "Error while quering from databse", err.Error())
		return nil // TODO better to have fallback call to default profile
	}
	return conf
}

func saveData(newConf *Config) error {
	_, err := setup.Database.Exec(`INSERT OR REPLACE INTO Alert VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		newConf.Profile, newConf.SendgridAPIKey, newConf.ElasticMailKey,
		newConf.ElasticMainUName, newConf.MailjetPublicKey, newConf.MailjetSecretKey,
		newConf.SendEmailTo, newConf.MailgunDomain, newConf.MailgunPrivateKey,
		newConf.MailgunPublicKey, newConf.SendDesktopNotific, newConf.MailOptionToUse, newConf.IconPath)
	if err != nil {
		utils.ModuleError(setup.DBLogFile, "Module: alert, Unable to insert/update profile",
			err.Error())
		return err
	}
	utils.ModuleLogs(setup.DBLogFile, fmt.Sprintf("Module: alert, Updated/Inserted the config of %s profile in db",
		newConf.Profile))
	return nil
}

//SaveConfig save the config of the module to the database
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
		utils.ModuleError(setup.DBLogFile, "Module: alert, Unable to decode obtained json.", err.Error())
		return err
	}
	conf = newConfig
	return saveData(newConfig)
}
