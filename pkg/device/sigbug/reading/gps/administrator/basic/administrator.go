package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	messageAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/administrator"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/action"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/administrator/exception"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/validator"
)

type administrator struct {
	sigfoxBackendDataCallbackReadingValidator validator.Validator
	sigbugGPSReadingRecordHandler             recordHandler.RecordHandler
}

func New(
	sigfoxBackendDataCallbackReadingValidator validator.Validator,
	sigbugGPSReadingRecordHandler recordHandler.RecordHandler,
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
		sigfoxBackendDataCallbackReadingValidateResponse, err := a.sigfoxBackendDataCallbackReadingValidator.Validate(&validator.ValidateRequest{
			Claims:  request.Claims,
			Reading: request.Reading,
			Action:  action.Create,
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

	createResponse, err := a.sigbugGPSReadingRecordHandler.Create(&recordHandler.CreateRequest{
		Reading: request.Reading,
	})
	if err != nil {
		return nil, exception.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &messageAdministrator.CreateResponse{
		Reading: createResponse.Reading,
	}, nil
}
