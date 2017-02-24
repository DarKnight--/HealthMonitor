package notify

import (
	"fmt"

	"github.com/mailgun/mailgun-go"
)

func mailGun(subject string, body string) error {
	mg := mailgun.NewMailgun(conf.MailgunDomain, conf.MailgunPrivateKey, conf.MailgunPublicKey)
	m := mg.NewMessage(
		fmt.Sprintf("%s %s", fromName, fromEmail),
		subject,
		body,
		conf.SendEmailTo,
	)
	_, _, err := mg.Send(m)
	return err
}
