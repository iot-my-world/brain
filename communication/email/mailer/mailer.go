package mailer

import "gitlab.com/iotTracker/brain/communication/email"

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
	Email email.Email
}

type SendResponse struct {
}
