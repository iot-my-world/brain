package password

import "gitlab.com/iotTracker/brain/communication/email"

type Email struct {
}

func (e Email) Body() string {
	return ""
}

func (e Email) Type() email.Type {
	return email.SetPassword
}
