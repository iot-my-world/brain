package email

type Type string

const SetPassword Type = "SetPassword"

type Email interface {
	Body() string
	Type() Type
}
