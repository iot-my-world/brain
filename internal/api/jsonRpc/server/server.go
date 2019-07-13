package server

import (
	"github.com/iot-my-world/brain/internal/api/jsonRpc/service/provider"
)

type Server interface {
	Start() error
	SecureStart() error
	RegisterServiceProvider(provider.Provider) error
}
