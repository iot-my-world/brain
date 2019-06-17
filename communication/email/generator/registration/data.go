package registration

import (
	"gitlab.com/iotTracker/brain/communication/email"
	humanUser "gitlab.com/iotTracker/brain/user/human"
	"net/mail"
)

type Data struct {
	URLToken string
	User     humanUser.User
}

func (d Data) Details() email.Details {
	return email.Details{
		Subject: "Welcome To SpotNav",
		To: []mail.Address{{
			Name:    d.User.Name,
			Address: d.User.EmailAddress,
		}},
		From: mail.Address{
			Name:    "SpotNav Team",
			Address: "noreply@spotnav.co.za",
		},
	}
}
