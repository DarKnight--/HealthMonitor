package notify

type Config struct {
	Profile           string
	SendgridAPIKey    string
	ElasticMailKey    string
	ElasticMainUName  string
	MailjetPublicKey  string
	MailjetSecretKey  string
	SendEmailTo       string
	MailgunDomain     string
	MailgunPrivateKey string
	MailgunPublicKey  string
	IconPath          string
}

var (
	config    Config
	fromName  = "OWTF Health Monitor"
	fromEmail = "alerts_health_monitor@owasp-owtf.org"
)
