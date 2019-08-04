package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	sigbugAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator"
	sigbugAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/administrator/adaptor/jsonRpc"
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

func (a *administrator) ValidateCreateRequest(request *sigbugAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *sigbugAdministrator.CreateRequest) (*sigbugAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	sigbugCreateResponse := sigbugAdministratorJsonRpcAdaptor.CreateResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		sigbugAdministrator.CreateService,
		sigbugAdministratorJsonRpcAdaptor.CreateRequest{
			Sigbug: request.Sigbug,
		},
		&sigbugCreateResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &sigbugAdministrator.CreateResponse{Sigbug: sigbugCreateResponse.Sigbug}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *sigbugAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *sigbugAdministrator.UpdateAllowedFieldsRequest) (*sigbugAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	sigbugUpdateAllowedFieldsResponse := sigbugAdministratorJsonRpcAdaptor.UpdateAllowedFieldsResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		sigbugAdministrator.UpdateAllowedFieldsService,
		sigbugAdministratorJsonRpcAdaptor.UpdateAllowedFieldsRequest{
			Sigbug: request.Sigbug,
		},
		&sigbugUpdateAllowedFieldsResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &sigbugAdministrator.UpdateAllowedFieldsResponse{
		Sigbug: sigbugUpdateAllowedFieldsResponse.Sigbug,
	}, nil
}

func (a *administrator) LastMessageUpdate(request *sigbugAdministrator.LastMessageUpdateRequest) (*sigbugAdministrator.LastMessageUpdateResponse, error) {
	return nil, brainException.NotImplemented{}
}
