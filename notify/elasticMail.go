package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func elasticMail(subject string, body string) {
	apiUrl := "https://api.elasticemail.com/mailer/send"
	form := url.Values{}
	form.Add("username", config.ElasticMainUName)
	form.Add("api_key", config.ElasticMailKey)
	form.Add("from", fromEmail)
	form.Add("from_name", fromName)
	form.Add("to", config.SendEmailTo)
	form.Add("subject", subject)
	form.Add("body_text", body)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiUrl, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	fmt.Println(resp.Status)
}
