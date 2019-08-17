package handler

import (
	"github.com/iot-my-world/brain/pkg/security/claims"
	sigfoxBackendDataDataCallbackReading "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type Handler interface {
	Handle(*HandleRequest) error
	WantReading(sigfoxBackendDataDataCallbackReading.Reading) bool
}

type HandleRequest struct {
	Claims      claims.Claims
	DataReading sigfoxBackendDataDataCallbackReading.Reading
}
