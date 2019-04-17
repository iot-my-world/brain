package jsonRpc

import (
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	brainException "gitlab.com/iotTracker/brain/exception"
	apiUserDeviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/apiUser/administrator"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) apiUserDeviceAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Create(request *apiUserDeviceAdministrator.CreateRequest) (*apiUserDeviceAdministrator.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) UpdateAllowedFields(request *apiUserDeviceAdministrator.UpdateAllowedFieldsRequest) (*apiUserDeviceAdministrator.UpdateAllowedFieldsResponse, error) {
	return nil, brainException.NotImplemented{}
}
