package registration

import (
	"bytes"
	"fmt"
	"gitlab.com/iotTracker/brain/communication/email"
	emailGenerator "gitlab.com/iotTracker/brain/communication/email/generator"
	emailGeneratorException "gitlab.com/iotTracker/brain/communication/email/generator/exception"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"html/template"
)

type generator struct {
	emailTemplate *template.Template
}

func New(
	pathToTemplateFolder string,
) emailGenerator.Generator {

	emailTemplate, err := template.ParseFiles(fmt.Sprintf("%s/%s", pathToTemplateFolder, "registration/template.html"))
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
		return nil, emailGeneratorException.TemplateExecution{Reasons: []string{err.Error()}}
	}

	return &emailGenerator.GenerateResponse{
		Email: email.Email{
			Body:    emailBytes.String(),
			Details: request.Data.Details(),
		},
	}, nil
}
