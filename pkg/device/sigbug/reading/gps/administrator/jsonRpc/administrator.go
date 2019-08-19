package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	sigbugGPSReadingAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/administrator"
	sigbugGPSReadingAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) sigbugGPSReadingAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) ValidateCreateRequest(request *sigbugGPSReadingAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *sigbugGPSReadingAdministrator.CreateRequest) (*sigbugGPSReadingAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	messageCreateResponse := sigbugGPSReadingAdministratorJsonRpcAdaptor.CreateResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		sigbugGPSReadingAdministrator.CreateService,
		sigbugGPSReadingAdministratorJsonRpcAdaptor.CreateRequest{
			Reading: request.Reading,
		},
		&messageCreateResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &sigbugGPSReadingAdministrator.CreateResponse{Reading: messageCreateResponse.Reading}, nil
}
