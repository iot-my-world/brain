package password

import (
	"bytes"
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/communication/email"
	emailGenerator "github.com/iot-my-world/brain/pkg/communication/email/generator"
	"github.com/iot-my-world/brain/pkg/communication/email/generator/exception"
	"html/template"
)

type generator struct {
	emailTemplate *template.Template
}

func New(
	pathToTemplateFolder string,
) emailGenerator.Generator {

	emailTemplate, err := template.ParseFiles(fmt.Sprintf("%s/%s", pathToTemplateFolder, "set/password/template.html"))
	if err != nil {
		log.Fatal("failed to parse file: " + err.Error())
	}

	return &generator{
		emailTemplate: emailTemplate,
	}
}

func (g *generator) ValidateGenerateEmailRequest(request *emailGenerator.GenerateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Data == nil {
		reasonsInvalid = append(reasonsInvalid, "data is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (g *generator) Generate(request *emailGenerator.GenerateRequest) (*emailGenerator.GenerateResponse, error) {
	if err := g.ValidateGenerateEmailRequest(request); err != nil {
		return nil, err
	}

	var emailBytes bytes.Buffer
	if err := g.emailTemplate.Execute(&emailBytes, request.Data); err != nil {
		return nil, exception.TemplateExecution{Reasons: []string{err.Error()}}
	}

	return &emailGenerator.GenerateResponse{
		Email: email.Email{
			Body:    emailBytes.String(),
			Details: request.Data.Details(),
		},
	}, nil
}
