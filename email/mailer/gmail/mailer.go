package gmail

import (
	"gitlab.com/iotTracker/brain/email/mailer"
	"net/smtp"
	"github.com/jpoehls/gophermail"
	"net/mail"
)

type gmailMailer struct {
	authInfo mailer.AuthInfo
}

func New(
	authInfo mailer.AuthInfo,
) *gmailMailer {
	return &gmailMailer{
		authInfo: authInfo,
	}
}

func (gm *gmailMailer) Send(request *mailer.SendRequest, response *mailer.SendResponse) error {

	// Set up authentication information.
	auth := smtp.PlainAuth(
		gm.authInfo.Identity,
		gm.authInfo.Username,
		gm.authInfo.Password,
		gm.authInfo.Host,
	)

	msg := &gophermail.Message{
		To:       []mail.Address{{Address: request.To}},
		HTMLBody: request.Body,
		From:     mail.Address{Name: "SpotNav Team", Address: request.From},
		Subject:  request.Subject,
	}
	err := gophermail.SendMail("smtp.gmail.com:587", auth, msg)
	if err != nil {
		return err
	}

	return nil
}
