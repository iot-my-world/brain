package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	backendAdministrator "github.com/iot-my-world/brain/pkg/sigfox/backend/administrator"
	backendAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/sigfox/backend/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) backendAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) ValidateCreateRequest(request *backendAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *backendAdministrator.CreateRequest) (*backendAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	backendCreateResponse := backendAdministratorJsonRpcAdaptor.CreateResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		backendAdministrator.CreateService,
		backendAdministratorJsonRpcAdaptor.CreateRequest{
			Backend: request.Backend,
		},
		&backendCreateResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &backendAdministrator.CreateResponse{Backend: backendCreateResponse.Backend}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *backendAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *backendAdministrator.UpdateAllowedFieldsRequest) (*backendAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	backendUpdateAllowedFieldsResponse := backendAdministratorJsonRpcAdaptor.UpdateAllowedFieldsResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		backendAdministrator.UpdateAllowedFieldsService,
		backendAdministratorJsonRpcAdaptor.UpdateAllowedFieldsRequest{
			Backend: request.Backend,
		},
		&backendUpdateAllowedFieldsResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &backendAdministrator.UpdateAllowedFieldsResponse{
		Backend: backendUpdateAllowedFieldsResponse.Backend,
	}, nil
}
