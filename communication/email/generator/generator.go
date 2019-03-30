package generator

import "gitlab.com/iotTracker/brain/communication/email"

type Generator interface {
	Generate() (email.Email, error)
}
