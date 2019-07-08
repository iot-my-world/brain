package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/device/sigbug/action"
	sigbugAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator"
	"github.com/iot-my-world/brain/pkg/device/sigbug/administrator/exception"
	"github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	"github.com/iot-my-world/brain/pkg/device/sigbug/validator"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
)

type administrator struct {
	sigbugDeviceValidator validator.Validator
	sigbugRecordHandler   recordHandler.RecordHandler
}

func New(
	sigbugDeviceValidator validator.Validator,
	sigbugRecordHandler recordHandler.RecordHandler,
) sigbugAdministrator.Administrator {
	return &administrator{
		sigbugDeviceValidator: sigbugDeviceValidator,
		sigbugRecordHandler:   sigbugRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *sigbugAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		sigbugDeviceValidateResponse, err := a.sigbugDeviceValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			Sigbug: request.Sigbug,
			Action: action.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating sigbug device: "+err.Error())
		} else {
			if len(sigbugDeviceValidateResponse.ReasonsInvalid) > 0 {
				for _, reason := range sigbugDeviceValidateResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("sigbug device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *sigbugAdministrator.CreateRequest) (*sigbugAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.sigbugRecordHandler.Create(&recordHandler.CreateRequest{
		Sigbug: request.Sigbug,
	})
	if err != nil {
		return nil, exception.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &sigbugAdministrator.CreateResponse{
		Sigbug: createResponse.Sigbug,
	}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *sigbugAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// device must be valid
		validationResponse, err := a.sigbugDeviceValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			Action: action.UpdateAllowedFields,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating device: "+err.Error())
		} else {
			if len(validationResponse.ReasonsInvalid) > 0 {
				for _, reason := range validationResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("sigbug device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *sigbugAdministrator.UpdateAllowedFieldsRequest) (*sigbugAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the device
	deviceRetrieveResponse, err := a.sigbugRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Sigbug.Id},
	})
	if err != nil {
		return nil, exception.DeviceRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the device

	// update the device
	_, err = a.sigbugRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Sigbug.Id},
		Sigbug:     deviceRetrieveResponse.Sigbug,
	})
	if err != nil {
		return nil, exception.DeviceUpdate{Reasons: []string{err.Error()}}
	}

	return &sigbugAdministrator.UpdateAllowedFieldsResponse{
		Sigbug: deviceRetrieveResponse.Sigbug,
	}, nil
}
