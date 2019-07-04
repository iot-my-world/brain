package generator

import (
	email2 "github.com/iot-my-world/brain/pkg/communication/email"
)

type Generator interface {
	Generate(request *GenerateRequest) (*GenerateResponse, error)
}

type GenerateRequest struct {
	Data email2.Data
}

type GenerateResponse struct {
	Email email2.Email
}
