package password

import (
	email2 "github.com/iot-my-world/brain/pkg/communication/email"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	"net/mail"
)

type Data struct {
	URLToken string
	User     humanUser.User
}

func (d Data) Details() email2.Details {
	return email2.Details{
		Subject: "Set Password",
		To: []mail.Address{{
			Name:    d.User.Name,
			Address: d.User.EmailAddress,
		}},
		From: mail.Address{
			Name:    "IOT My World Team",
			Address: "iotmywordteam@gmail.com",
		},
	}
}
