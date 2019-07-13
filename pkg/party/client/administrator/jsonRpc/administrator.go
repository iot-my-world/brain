package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/internal/api/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	clientAdministrator "github.com/iot-my-world/brain/pkg/party/client/administrator"
	"github.com/iot-my-world/brain/pkg/party/client/administrator/adaptor/jsonRpc"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
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

	clientCreateResponse := jsonRpc.CreateResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		clientAdministrator.CreateService,
		jsonRpc.CreateRequest{
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

	clientUpdateAllowedFieldsResponse := jsonRpc.UpdateAllowedFieldsResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		clientAdministrator.UpdateAllowedFieldsService,
		jsonRpc.UpdateAllowedFieldsRequest{
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

	response := jsonRpc.DeleteResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		clientAdministrator.DeleteService,
		jsonRpc.DeleteRequest{
			ClientIdentifier: *id,
		},
		&response); err != nil {
		return nil, err
	}

	return &clientAdministrator.DeleteResponse{}, nil
}
