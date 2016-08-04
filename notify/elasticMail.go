package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func elasticMail(subject string, body string) {
	apiURL := "https://api.elasticemail.com/mailer/send"
	form := url.Values{}
	form.Add("username", conf.ElasticMainUName)
	form.Add("api_key", conf.ElasticMailKey)
	form.Add("from", fromEmail)
	form.Add("from_name", fromName)
	form.Add("to", conf.SendEmailTo)
	form.Add("subject", subject)
	form.Add("body_text", body)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	fmt.Println(resp.Status)
}
