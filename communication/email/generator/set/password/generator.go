package password

import (
	"fmt"
	"gitlab.com/iotTracker/brain/communication/email"
	emailGenerator "gitlab.com/iotTracker/brain/communication/email/generator"
	"gitlab.com/iotTracker/brain/log"
	"html/template"
)

type generator struct {
	emailTemplate *template.Template
}

func New(
	pathToTemplateFolder string,
) emailGenerator.Generator {

	// parse html file as template
	// email/template/set/password/template.html
	emailTemplate, err := template.ParseFiles(fmt.Sprintf("%s/%s", pathToTemplateFolder, "set/password/template.html"))
	if err != nil {
		log.Fatal("failed to parse file: " + err.Error())
	}

	return &generator{
		emailTemplate: emailTemplate,
	}
}

func (g *generator) Generate() (email.Email, error) {
	return Email{}, nil
}
