package notify

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
	config       Config
	fromName     = "OWTF Health Monitor"
	fromEmail    = "alerts_health_monitor@owasp-owtf.org"
	desktopAlert *DesktopAlert
)

func Init() {
	desktopAlert = nil
	config.DesktopNoticSupported = false
	if CheckDesktopAlertSupport() {
		desktopAlert = DesktopAlertBuilder("OWTF Health Monitor", config.IconPath)
		config.DesktopNoticSupported = true
	}
}

func SendEmailAlert(subject string, body string) {
	switch config.MailOptionToUse {
	case "sendgrid":
		sendGrid(subject, body)
	case "mailgun":
		mailGun(subject, body)
	case "elasticemail":
		elasticMail(subject, body)
	}
}

func SendDesktopAlert(subject string, summary string, urgent int, iconPath string) {
	if iconPath == "" {
		iconPath = config.IconPath
	}
	if config.SendDesktopNotific && config.DesktopNoticSupported {
		desktopAlert.Push(subject, summary, iconPath, urgent)
	}
}
