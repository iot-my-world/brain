package password

import (
	"gitlab.com/iotTracker/brain/communication/email"
	"gitlab.com/iotTracker/brain/user"
	"net/mail"
)

type Data struct {
	URLToken string
	User     user.User
}

func (d Data) Details() email.Details {
	return email.Details{
		Subject: "Set Password",
		To: mail.Address{
			Name:    d.User.Name,
			Address: d.User.EmailAddress,
		},
		From: mail.Address{
			Name:    "SpotNav Team",
			Address: "noreply@spotnav.co.za",
		},
	}
}
