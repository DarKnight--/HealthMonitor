package notify

// Config holds all the necessary parameters required by the module
type Config struct {
	Profile               string
	SendgridAPIKey        string
	ElasticMailKey        string
	ElasticMailUName      string
	MailjetPublicKey      string
	MailjetSecretKey      string
	SendEmailTo           string
	MailgunDomain         string
	MailgunPrivateKey     string
	MailgunPublicKey      string
	MaxEmailRetry         int
	DesktopNoticSupported bool
	SendDesktopNotific    bool
	MailOptionToUse       string
	IconPath              string
}

func sendEmail(subject string, body string) error {
	switch conf.MailOptionToUse {
	case "sendgrid":
		return sendGrid(subject, body)
	case "mailgun":
		return mailGun(subject, body)
	case "elasticemail":
		return elasticMail(subject, body)
	}
	return nil
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
