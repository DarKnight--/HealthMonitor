package notify

import (
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendGrid(subject string, contentString string) error {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(strings.Split(conf.SendEmailTo, "@")[0], conf.SendEmailTo)
	content := mail.NewContent("text/plain", contentString)
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(conf.SendgridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	return err
}
