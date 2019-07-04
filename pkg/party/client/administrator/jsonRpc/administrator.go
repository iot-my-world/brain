package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	administrator2 "github.com/iot-my-world/brain/pkg/party/client/administrator"
	"github.com/iot-my-world/brain/pkg/party/client/administrator/adaptor/jsonRpc"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
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

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	clientCreateResponse := jsonRpc.CreateResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.CreateService,
		jsonRpc.CreateRequest{
			Client: request.Client,
		},
		&clientCreateResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.CreateResponse{Client: clientCreateResponse.Client}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *administrator2.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *administrator2.UpdateAllowedFieldsRequest) (*administrator2.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	clientUpdateAllowedFieldsResponse := jsonRpc.UpdateAllowedFieldsResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.UpdateAllowedFieldsService,
		jsonRpc.UpdateAllowedFieldsRequest{
			Client: request.Client,
		},
		&clientUpdateAllowedFieldsResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.UpdateAllowedFieldsResponse{
		Client: clientUpdateAllowedFieldsResponse.Client,
	}, nil
}

func (a *administrator) ValidateDeleteRequest(request *administrator2.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.ClientIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "client identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Delete(request *administrator2.DeleteRequest) (*administrator2.DeleteResponse, error) {
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
		administrator2.DeleteService,
		jsonRpc.DeleteRequest{
			ClientIdentifier: *id,
		},
		&response); err != nil {
		return nil, err
	}

	return &administrator2.DeleteResponse{}, nil
}
