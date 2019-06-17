package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	sf001DeviceAdministrator "github.com/iot-my-world/brain/tracker/sf001/administrator"
	sf001DeviceAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/tracker/sf001/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) sf001DeviceAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Create(request *sf001DeviceAdministrator.CreateRequest) (*sf001DeviceAdministrator.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) UpdateAllowedFields(request *sf001DeviceAdministrator.UpdateAllowedFieldsRequest) (*sf001DeviceAdministrator.UpdateAllowedFieldsResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) ValidateHeartbeatRequest(request *sf001DeviceAdministrator.HeartbeatRequest) error {
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

func (a *administrator) Heartbeat(request *sf001DeviceAdministrator.HeartbeatRequest) (*sf001DeviceAdministrator.HeartbeatResponse, error) {
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
