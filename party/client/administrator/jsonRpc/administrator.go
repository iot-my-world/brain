package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	clientAdministrator "github.com/iot-my-world/brain/party/client/administrator"
	clientAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/party/client/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) clientAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) ValidateCreateRequest(request *clientAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *clientAdministrator.CreateRequest) (*clientAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	clientCreateResponse := clientAdministratorJsonRpcAdaptor.CreateResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		clientAdministrator.CreateService,
		clientAdministratorJsonRpcAdaptor.CreateRequest{
			Client: request.Client,
		},
		&clientCreateResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &clientAdministrator.CreateResponse{Client: clientCreateResponse.Client}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *clientAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *clientAdministrator.UpdateAllowedFieldsRequest) (*clientAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	clientUpdateAllowedFieldsResponse := clientAdministratorJsonRpcAdaptor.UpdateAllowedFieldsResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		clientAdministrator.UpdateAllowedFieldsService,
		clientAdministratorJsonRpcAdaptor.UpdateAllowedFieldsRequest{
			Client: request.Client,
		},
		&clientUpdateAllowedFieldsResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &clientAdministrator.UpdateAllowedFieldsResponse{
		Client: clientUpdateAllowedFieldsResponse.Client,
	}, nil
}
