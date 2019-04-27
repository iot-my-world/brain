package jsonRpc

import (
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	brainException "gitlab.com/iotTracker/brain/exception"
	zx303TaskAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/task/administrator"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) zx303TaskAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Create(request *zx303TaskAdministrator.CreateRequest) (*zx303TaskAdministrator.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}
