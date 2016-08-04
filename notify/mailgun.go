package notify

import (
	"fmt"

	"github.com/mailgun/mailgun-go"
)

func mailGun(subject string, body string) {
	mg := mailgun.NewMailgun(conf.MailgunDomain, conf.MailgunPrivateKey, conf.MailgunPublicKey)
	m := mg.NewMessage(
		fmt.Sprintf("%s %s", fromName, fromEmail),
		subject,
		body,
		conf.SendEmailTo,
	)
	_, id, err := mg.Send(m)
	fmt.Println(id)
	fmt.Println(err)
}
