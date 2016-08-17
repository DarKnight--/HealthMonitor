package notify

import (
	"encoding/json"

	"github.com/owtf/health_monitor/setup"
	"github.com/owtf/health_monitor/utils"
)

// Config holds all the necessary parameters required by the module
type Config struct {
	Profile               string
	SendgridAPIKey        string
	ElasticMailKey        string
	ElasticMainUName      string
	MailjetPublicKey      string
	MailjetSecretKey      string
	SendEmailTo           string
	MailgunDomain         string
	MailgunPrivateKey     string
	MailgunPublicKey      string
	DesktopNoticSupported bool
	SendDesktopNotific    bool
	MailOptionToUse       string
	IconPath              string
}

var (
	conf         *Config
	fromName     = "OWTF Health Monitor"
	fromEmail    = "alerts_health_monitor@owasp-owtf.org"
	desktopAlert *DesktopAlert
)

// Init is the initialization function of the module
func Init() {
	conf = LoadConfig()
	if conf == nil {
		utils.CheckConf(setup.MainLogFile, setup.MainLogFile, "alert", &setup.UserModuleState.Profile, setup.Alert)
	}
	desktopAlert = nil
	conf.DesktopNoticSupported = false
	if CheckDesktopAlertSupport() {
		desktopAlert = NewDesktopAlert(conf.IconPath)
		conf.DesktopNoticSupported = true
	}
}

// SendEmailAlert sends the email notification if enabled using specified client
func SendEmailAlert(subject string, body string) {
	switch conf.MailOptionToUse {
	case "sendgrid":
		sendGrid(subject, body)
	case "mailgun":
		mailGun(subject, body)
	case "elasticemail":
		elasticMail(subject, body)
	}
}

// SendDesktopAlert sends the desktop notification if enabled and required packages
// are installed.
func SendDesktopAlert(subject string, summary string, urgent MessageImportance, iconPath string) {
	if iconPath == "" {
		iconPath = conf.IconPath
	}
	if conf.SendDesktopNotific && conf.DesktopNoticSupported {
		desktopAlert.Push(subject, summary, iconPath, urgent)
	}
}

//GetConfJSON returns the json byte array of the module's config
func GetConfJSON() []byte {
	data, err := json.Marshal(LoadConfig())
	if err != nil {
		utils.ModuleError(setup.MainLogFile, err.Error(), "[!] Check the conf struct(alert module)")
	}
	return data
}
