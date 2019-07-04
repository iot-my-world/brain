package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	clientAdministrator "github.com/iot-my-world/brain/party/client/administrator"
	clientAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/party/client/administrator/adaptor/jsonRpc"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
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

func (a *administrator) ValidateDeleteRequest(request *clientAdministrator.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.ClientIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "client identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Delete(request *clientAdministrator.DeleteRequest) (*clientAdministrator.DeleteResponse, error) {
	if err := a.ValidateDeleteRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// wrap identifier
	id, err := wrappedIdentifier.Wrap(request.ClientIdentifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	response := clientAdministratorJsonRpcAdaptor.DeleteResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		clientAdministrator.DeleteService,
		clientAdministratorJsonRpcAdaptor.DeleteRequest{
			ClientIdentifier: *id,
		},
		&response); err != nil {
		return nil, err
	}

	return &clientAdministrator.DeleteResponse{}, nil
}
