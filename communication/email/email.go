package email

import "net/mail"

type Type string

const SetPassword Type = "SetPassword"

type Email struct {
	Body    string
	Details Details
}

type Details struct {
	Subject string
	To      []mail.Address
	From    mail.Address
}

type Data interface {
	Details() Details
}
