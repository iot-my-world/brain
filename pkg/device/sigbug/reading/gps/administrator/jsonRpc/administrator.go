package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	messageAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/administrator"
	messageAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) messageAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) ValidateCreateRequest(request *messageAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *messageAdministrator.CreateRequest) (*messageAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	messageCreateResponse := messageAdministratorJsonRpcAdaptor.CreateResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		messageAdministrator.CreateService,
		messageAdministratorJsonRpcAdaptor.CreateRequest{
			Reading: request.Reading,
		},
		&messageCreateResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &messageAdministrator.CreateResponse{Reading: messageCreateResponse.Reading}, nil
}
