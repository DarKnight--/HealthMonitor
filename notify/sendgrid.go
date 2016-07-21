package alert

import (
	"fmt"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendGrid(subject string, contentString string) {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(strings.Split(config.SendEmailTo, "@")[0], config.SendEmailTo)
	content := mail.NewContent("text/plain", contentString)
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(config.SendgridApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
