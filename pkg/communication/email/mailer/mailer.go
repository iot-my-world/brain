package mailer

import (
	email2 "github.com/iot-my-world/brain/pkg/communication/email"
)

type AuthInfo struct {
	Identity string
	Username string
	Password string
	Host     string
}

type Mailer interface {
	Send(request *SendRequest) (*SendResponse, error)
}

type SendRequest struct {
	Email email2.Email
}

type SendResponse struct {
}
