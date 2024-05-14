package Email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

const verificationHTMLPath = "verificationHtml/default.html"

type Sender struct {
	template *template.Template
	server   string
	author   string
	auth     smtp.Auth
	host     string
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	ApiHost  string
}

type Content struct {
	Author           string
	Target           string
	VerificationLink string
}

func New(config Config) (*Sender, error) {
	bytes, err := os.ReadFile(verificationHTMLPath)
	if err != nil {
		return nil, err
	}

	html := string(bytes)

	tmpl, err := template.New("VerificationLink").Parse(html)
	if err != nil {
		return nil, err
	}

	server := fmt.Sprintf("%s:%d", config.Host, config.Port)
	auth := smtp.PlainAuth("", config.User, config.Password, config.Host)
	sender := Sender{
		template: tmpl,
		server:   server,
		author:   config.User,
		auth:     auth,
		host:     config.ApiHost,
	}
	return &sender, nil
}

func (sender *Sender) SendVerify(targetEmail string, key string) error {
	content := Content{
		Author:           sender.author,
		Target:           targetEmail,
		VerificationLink: "https://" + sender.host + "/#register/" + key,
	}

	var result bytes.Buffer
	err := sender.template.Execute(&result, content)
	if err != nil {
		return err
	}

	err = smtp.SendMail(sender.server, sender.auth, sender.author, []string{targetEmail}, result.Bytes())
	if err != nil {
		return err
	}

	return nil
}
