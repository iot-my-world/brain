package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	zx303DeviceAdministrator "github.com/iot-my-world/brain/tracker/zx303/administrator"
	zx303DeviceAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/tracker/zx303/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) zx303DeviceAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Create(request *zx303DeviceAdministrator.CreateRequest) (*zx303DeviceAdministrator.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) UpdateAllowedFields(request *zx303DeviceAdministrator.UpdateAllowedFieldsRequest) (*zx303DeviceAdministrator.UpdateAllowedFieldsResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) ValidateHeartbeatRequest(request *zx303DeviceAdministrator.HeartbeatRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.ZX303Identifier == nil {
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

func (a *administrator) Heartbeat(request *zx303DeviceAdministrator.HeartbeatRequest) (*zx303DeviceAdministrator.HeartbeatResponse, error) {
	if err := a.ValidateHeartbeatRequest(request); err != nil {
		return nil, err
	}

	// create wrapped identifier
	wrappedDeviceIdentifier, err := wrappedIdentifier.Wrap(request.ZX303Identifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// perform heartbeat
	loginResponse := zx303DeviceAdministratorJsonRpcAdaptor.HeartbeatResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"ZX303DeviceAdministrator.Heartbeat",
		zx303DeviceAdministratorJsonRpcAdaptor.HeartbeatRequest{
			WrappedZX303Identifier: *wrappedDeviceIdentifier,
		},
		&loginResponse,
	); err != nil {
		return nil, err
	}

	return nil, nil
}
