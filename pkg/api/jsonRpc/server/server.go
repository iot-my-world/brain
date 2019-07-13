package server

import (
	provider2 "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
)

type Server interface {
	Start() error
	SecureStart() error
	RegisterServiceProvider(provider2.Provider) error
}
