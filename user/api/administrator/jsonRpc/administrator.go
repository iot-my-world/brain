package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	apiUserAdministrator "github.com/iot-my-world/brain/user/api/administrator"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) apiUserAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Create(request *apiUserAdministrator.CreateRequest) (*apiUserAdministrator.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) UpdateAllowedFields(request *apiUserAdministrator.UpdateAllowedFieldsRequest) (*apiUserAdministrator.UpdateAllowedFieldsResponse, error) {
	return nil, brainException.NotImplemented{}
}
