package notify

import (
	"fmt"

	"github.com/mailgun/mailgun-go"
)

func mailGun(subject string, body string) {
	mg := mailgun.NewMailgun(config.MailgunDomain, config.MailgunPrivateKey, config.MailgunPublicKey)
	m := mg.NewMessage(
		fmt.Sprintf("%s %s", fromName, fromEmail),
		subject,
		body,
		config.SendEmailTo,
	)
	_, id, err := mg.Send(m)
	fmt.Println(id)
	fmt.Println(err)
}
