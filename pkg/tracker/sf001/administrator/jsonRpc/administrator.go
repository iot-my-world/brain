package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	administrator2 "github.com/iot-my-world/brain/pkg/tracker/sf001/administrator"
	sf001DeviceAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/trackercker/sf001/administrator/adaptor/jsonRpc"
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

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) UpdateAllowedFields(request *administrator2.UpdateAllowedFieldsRequest) (*administrator2.UpdateAllowedFieldsResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) ValidateHeartbeatRequest(request *administrator2.HeartbeatRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.SF001Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if !a.jsonRpcClient.LoggedIn() {
		reasonsInvalid = append(reasonsInvalid, "json rpc client is not logged in")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Heartbeat(request *administrator2.HeartbeatRequest) (*administrator2.HeartbeatResponse, error) {
	if err := a.ValidateHeartbeatRequest(request); err != nil {
		return nil, err
	}

	// create wrapped identifier
	wrappedDeviceIdentifier, err := wrappedIdentifier.Wrap(request.SF001Identifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// perform heartbeat
	loginResponse := sf001DeviceAdministratorJsonRpcAdaptor.HeartbeatResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"SF001DeviceAdministrator.Heartbeat",
		sf001DeviceAdministratorJsonRpcAdaptor.HeartbeatRequest{
			WrappedSF001Identifier: *wrappedDeviceIdentifier,
		},
		&loginResponse,
	); err != nil {
		return nil, err
	}

	return nil, nil
}