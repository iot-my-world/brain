package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/action"
	messageAdministrator "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/administrator"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/administrator/exception"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/validator"
)

type administrator struct {
	sigfoxBackendDataCallbackMessageValidator     validator.Validator
	sigfoxBackendDataCallbackMessageRecordHandler recordHandler.RecordHandler
}

func New(
	sigfoxBackendDataCallbackMessageValidator validator.Validator,
	sigfoxBackendDataCallbackMessageRecordHandler recordHandler.RecordHandler,
) messageAdministrator.Administrator {
	return &administrator{
		sigfoxBackendDataCallbackMessageValidator:     sigfoxBackendDataCallbackMessageValidator,
		sigfoxBackendDataCallbackMessageRecordHandler: sigfoxBackendDataCallbackMessageRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *messageAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		sigfoxBackendDataCallbackMessageValidateResponse, err := a.sigfoxBackendDataCallbackMessageValidator.Validate(&validator.ValidateRequest{
			Claims:  request.Claims,
			Message: request.Message,
			Action:  action.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating message message: "+err.Error())
		} else {
			if len(sigfoxBackendDataCallbackMessageValidateResponse.ReasonsInvalid) > 0 {
				for _, reason := range sigfoxBackendDataCallbackMessageValidateResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("message message invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
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

	createResponse, err := a.sigfoxBackendDataCallbackMessageRecordHandler.Create(&recordHandler.CreateRequest{
		Message: request.Message,
	})
	if err != nil {
		return nil, exception.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &messageAdministrator.CreateResponse{
		Message: createResponse.Message,
	}, nil
}
