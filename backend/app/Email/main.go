package Email

import (
	"fmt"
	"net/smtp"
	"path/filepath"
	"simple_account/app/Email/ChangeEmailMail"
	"simple_account/app/Email/EmailOwnerMail"
	"simple_account/app/Email/EmailSender"
	"simple_account/app/Email/RegisterMail"
	"simple_account/app/Email/TimedKeys"
	"time"
)

type Manager struct {
	Templates *Templates
	TimedKeys *TimedKeys.TimedKeys
}

type Templates struct {
	Register    *RegisterMail.Mail
	ChangeEmail *ChangeEmailMail.Mail
	EmailOwner  *EmailOwnerMail.Mail
}

type Config struct {
	FilesPath           string
	Host                string
	Port                int
	User                string
	Password            string
	WebsiteHost         string
	VerificationDuraion time.Duration
}

func New(config Config) (*Manager, error) {
	server := fmt.Sprintf("%s:%d", config.Host, config.Port)
	auth := smtp.PlainAuth("", config.User, config.Password, config.Host)

	sender := &EmailSender.Sender{
		Server: server,
		Author: config.User,
		Auth:   auth,
		Host:   config.WebsiteHost,
	}

	registerTemplatePath := filepath.Join(config.FilesPath, "register")
	register, err := RegisterMail.New(sender, registerTemplatePath)
	if err != nil {
		return nil, err
	}

	changeEmailTemplatePath := filepath.Join(config.FilesPath, "change_email")
	changeEmail, err := ChangeEmailMail.New(sender, changeEmailTemplatePath)
	if err != nil {
		return nil, err
	}

	verificationTemplatePath := filepath.Join(config.FilesPath, "email_owner")
	emailOwner, err := EmailOwnerMail.New(sender, verificationTemplatePath)
	if err != nil {
		return nil, err
	}

	templates := Templates{
		Register:    register,
		ChangeEmail: changeEmail,
		EmailOwner:  emailOwner,
	}

	timedKeys := TimedKeys.New(config.VerificationDuraion)

	manager := &Manager{
		Templates: &templates,
		TimedKeys: timedKeys,
	}

	return manager, nil
}
