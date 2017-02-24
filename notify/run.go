package notify

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/owtf/health_monitor/setup"
	"github.com/owtf/health_monitor/utils"
)

var (
	conf         *Config
	fromName     = "OWTF Health Monitor"
	fromEmail    = "alerts_health_monitor@owasp-owtf.org"
	desktopAlert *DesktopAlert
	emailQueue   *utils.Queue
)

type email struct {
	Subject string
	Body    string
}

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

// Notify is driver funcion for the health_monitor to send email
func Notify() {
	emailQueue = utils.NewQueue()
	var err error
	for {
		select {
		case <-time.After(time.Second):
			for emailQueue.Len() != 0 {
				emailData := reflect.ValueOf(emailQueue.Poll())
				i := 0
				for ; i < conf.MaxEmailRetry; i++ {
					err = sendEmail(emailData.FieldByName("Subject").Interface().(string),
						emailData.FieldByName("Body").Interface().(string))
					if err == nil {
						break
					}
				}
				if i == conf.MaxEmailRetry {
					utils.ModuleError(setup.MainLogFile, err.Error(), "[!] Unable to send email")
				}
			}
		}
	}
}

// SendEmailAlert sends the email notification if enabled using specified client
func SendEmailAlert(subject string, body string) {
	emailQueue.Push(email{Subject: subject, Body: body})
}

//GetConfJSON returns the json byte array of the module's config
func GetConfJSON() []byte {
	data, err := json.Marshal(LoadConfig())
	if err != nil {
		utils.ModuleError(setup.MainLogFile, err.Error(), "[!] Check the conf struct(alert module)")
	}
	return data
}
