package generator

import "gitlab.com/iotTracker/brain/communication/email"

type Generator interface {
	Generate(request *GenerateRequest) (*GenerateResponse, error)
}

type GenerateRequest struct {
	Data email.Data
}

type GenerateResponse struct {
	Email email.Email
}
