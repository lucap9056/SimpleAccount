package EmailSender

import (
	"net/smtp"
	"os"
	"path/filepath"
	"text/template"
)

type Sender struct {
	Server string
	Author string
	Auth   smtp.Auth
	Host   string
}

func (sender *Sender) Send(target string, content []byte) error {
	return smtp.SendMail(sender.Server, sender.Auth, sender.Author, []string{target}, content)
}

func ReadTemplate(dirPath string) (*template.Template, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	htmlFiles := []string{}
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dirPath, file.Name())
			htmlFiles = append(htmlFiles, filePath)
		}
	}
	return template.ParseFiles(htmlFiles...)
}
