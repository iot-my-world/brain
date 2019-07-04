package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	administrator2 "github.com/iot-my-world/brain/pkg/user/api/administrator"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) administrator2.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) UpdateAllowedFields(request *administrator2.UpdateAllowedFieldsRequest) (*administrator2.UpdateAllowedFieldsResponse, error) {
	return nil, brainException.NotImplemented{}
}
