package password

import (
	"github.com/iot-my-world/brain/communication/email"
	humanUser "github.com/iot-my-world/brain/user/human"
	"net/mail"
)

type Data struct {
	URLToken string
	User     humanUser.User
}

func (d Data) Details() email.Details {
	return email.Details{
		Subject: "Set Password",
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
