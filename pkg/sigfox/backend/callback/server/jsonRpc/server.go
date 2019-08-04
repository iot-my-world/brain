package jsonRpc

import (
	"encoding/hex"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	sigfoxBackendCallbackServer "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server"
	sigfoxBackendCallbackServerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server/adaptor/jsonRpc"
)

type server struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) sigfoxBackendCallbackServer.Server {
	return &server{
		jsonRpcClient: jsonRpcClient,
	}
}

func (s *server) ValidateHandleDataMessageRequest(request *sigfoxBackendCallbackServer.HandleDataMessageRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (s *server) HandleDataMessage(request *sigfoxBackendCallbackServer.HandleDataMessageRequest) (*sigfoxBackendCallbackServer.HandleDataMessageResponse, error) {
	if err := s.ValidateHandleDataMessageRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	handleDataMessageResponse := sigfoxBackendCallbackServerJsonRpcAdaptor.HandleDataMessageResponse{}
	if err := s.jsonRpcClient.JsonRpcRequest(
		sigfoxBackendCallbackServer.HandleDataMessageService,
		sigfoxBackendCallbackServerJsonRpcAdaptor.HandleDataMessageRequest{
			DeviceId: request.Message.DeviceId,
			Data:     hex.EncodeToString(request.Message.Data),
		},
		&handleDataMessageResponse); err != nil {
		return nil, err
	}

	return &sigfoxBackendCallbackServer.HandleDataMessageResponse{}, nil
}
