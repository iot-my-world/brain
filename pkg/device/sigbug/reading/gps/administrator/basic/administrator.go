package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	sigbugGPSReadingAction "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/action"
	sigbugGPSReadingAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/administrator"
	sigbugGPSReadingAdministratorException "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/administrator/exception"
	sigbugGPSReadingRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/recordHandler"
	sigbugGPSReadingValidator "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/validator"
)

type administrator struct {
	sigfoxBackendDataCallbackReadingValidator sigbugGPSReadingValidator.Validator
	sigbugGPSReadingRecordHandler             sigbugGPSReadingRecordHandler.RecordHandler
}

func New(
	sigfoxBackendDataCallbackReadingValidator sigbugGPSReadingValidator.Validator,
	sigbugGPSReadingRecordHandler sigbugGPSReadingRecordHandler.RecordHandler,
) messageAdministrator.Administrator {
	return &administrator{
		sigfoxBackendDataCallbackReadingValidator: sigfoxBackendDataCallbackReadingValidator,
		sigbugGPSReadingRecordHandler:             sigbugGPSReadingRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *messageAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		sigfoxBackendDataCallbackReadingValidateResponse, err := a.sigfoxBackendDataCallbackReadingValidator.Validate(&sigbugGPSReadingValidator.ValidateRequest{
			Claims:  request.Claims,
			Reading: request.Reading,
			Action:  sigbugGPSReadingAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating reading reading: "+err.Error())
		} else {
			if len(sigfoxBackendDataCallbackReadingValidateResponse.ReasonsInvalid) > 0 {
				for _, reason := range sigfoxBackendDataCallbackReadingValidateResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("reading invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *messageAdministrator.CreateRequest) (*messageAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.sigbugGPSReadingRecordHandler.Create(&sigbugGPSReadingRecordHandler.CreateRequest{
		Reading: request.Reading,
	})
	if err != nil {
		return nil, sigbugGPSReadingAdministratorException.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &messageAdministrator.CreateResponse{
		Reading: createResponse.Reading,
	}, nil
}
