package gmail

import (
	emailMailer "github.com/iot-my-world/brain/communication/email/mailer"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/jpoehls/gophermail"
	"net/smtp"
)

type mailer struct {
	authInfo emailMailer.AuthInfo
}

func New(
	authInfo emailMailer.AuthInfo,
) emailMailer.Mailer {
	return &mailer{
		authInfo: authInfo,
	}
}

func (m *mailer) ValidateSendRequest(request *emailMailer.SendRequest) error {
	reasonsInvalid := make([]string, 0)

	for _, toAddress := range request.Email.Details.To {
		if toAddress.Address == "" {
			reasonsInvalid = append(reasonsInvalid, "to email address blank")
		}

		if toAddress.Name == "" {
			reasonsInvalid = append(reasonsInvalid, "to name blank")
		}
	}

	if request.Email.Details.From.Address == "" {
		reasonsInvalid = append(reasonsInvalid, "from email address blank")
	}

	if request.Email.Details.From.Name == "" {
		reasonsInvalid = append(reasonsInvalid, "from name blank")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (m *mailer) Send(request *emailMailer.SendRequest) (*emailMailer.SendResponse, error) {
	if err := m.ValidateSendRequest(request); err != nil {
		return nil, err
	}

	// Set up authentication information.
	auth := smtp.PlainAuth(
		m.authInfo.Identity,
		m.authInfo.Username,
		m.authInfo.Password,
		m.authInfo.Host,
	)

	msg := &gophermail.Message{
		To:       request.Email.Details.To,
		HTMLBody: request.Email.Body,
		From:     request.Email.Details.From,
		Subject:  request.Email.Details.Subject,
	}
	err := gophermail.SendMail("smtp.gmail.com:587", auth, msg)
	if err != nil {
		return nil, err
	}

	return &emailMailer.SendResponse{}, nil
}
