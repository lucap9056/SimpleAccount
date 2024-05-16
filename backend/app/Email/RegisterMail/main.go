package RegisterMail

import (
	"bytes"
	"simple_account/app/Email/EmailSender"
	"text/template"
)

type Mail struct {
	sender   *EmailSender.Sender
	template *template.Template
}

type content struct {
	Author           string
	Target           string
	VerificationLink string
}

func New(sender *EmailSender.Sender, filesPath string) (*Mail, error) {
	templ, err := EmailSender.ReadTemplate(filesPath)
	if err != nil {
		return nil, err
	}

	mail := Mail{
		sender:   sender,
		template: templ,
	}
	return &mail, nil
}

func (m *Mail) SendVerify(language string, targetEmail string, key string) error {

	content := content{
		Author:           m.sender.Author,
		Target:           targetEmail,
		VerificationLink: "https://" + m.sender.Host + "/#register/" + key,
	}
	var result bytes.Buffer
	err := m.template.ExecuteTemplate(&result, language, content)
	if err != nil {
		err := m.template.ExecuteTemplate(&result, "default", content)
		if err != nil {
			return err
		}
	}

	err = m.sender.Send(targetEmail, result.Bytes())
	if err != nil {
		return err
	}

	return nil
}
