package forgotPassword

import (
	"bytes"
	"github.com/go-errors/errors"
	"gitlab.com/iotTracker/brain/user"
	"text/template"
)

type Data struct {
	URLToken string
	User     user.User
}

func GenerateEmail(data Data) (string, error) {
	// parse html file as template
	emailTemplate, err := template.ParseFiles("email/template/forgotPassword/template.html")
	if err != nil {
		return "", errors.New("failed to parse file: " + err.Error())
	}

	var emailBytes bytes.Buffer
	if err := emailTemplate.Execute(&emailBytes, data); err != nil {
		return "", errors.New("failed to execute template: " + err.Error())
	}

	return emailBytes.String(), nil
}
