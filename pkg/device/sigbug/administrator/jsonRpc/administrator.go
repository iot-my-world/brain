package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	sigbugAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) sigbugAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Create(request *sigbugAdministrator.CreateRequest) (*sigbugAdministrator.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) UpdateAllowedFields(request *sigbugAdministrator.UpdateAllowedFieldsRequest) (*sigbugAdministrator.UpdateAllowedFieldsResponse, error) {
	return nil, brainException.NotImplemented{}
}
