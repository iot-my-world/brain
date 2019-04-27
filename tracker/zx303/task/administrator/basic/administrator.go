package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	zx303TaskAction "gitlab.com/iotTracker/brain/tracker/zx303/task/action"
	zx303TaskAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/task/administrator"
	zx303TaskAdministratorException "gitlab.com/iotTracker/brain/tracker/zx303/task/administrator/exception"
	zx303TaskRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/task/recordHandler"
	zx303TaskValidator "gitlab.com/iotTracker/brain/tracker/zx303/task/validator"
)

type administrator struct {
	zx303TaskValidator     zx303TaskValidator.Validator
	zx303TaskRecordHandler *zx303TaskRecordHandler.RecordHandler
}

func New(
	zx303TaskValidator zx303TaskValidator.Validator,
	zx303TaskRecordHandler *zx303TaskRecordHandler.RecordHandler,
) zx303TaskAdministrator.Administrator {
	return &administrator{
		zx303TaskValidator:     zx303TaskValidator,
		zx303TaskRecordHandler: zx303TaskRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *zx303TaskAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		zx303DeviceValidateResponse, err := a.zx303TaskValidator.Validate(&zx303TaskValidator.ValidateRequest{
			Claims:    request.Claims,
			ZX303Task: request.ZX303Task,
			Action:    zx303TaskAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating zx303 task: "+err.Error())
		}
		if len(zx303DeviceValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range zx303DeviceValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("zx303 task invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *zx303TaskAdministrator.CreateRequest) (*zx303TaskAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.zx303TaskRecordHandler.Create(&zx303TaskRecordHandler.CreateRequest{
		ZX303Task: request.ZX303Task,
	})
	if err != nil {
		return nil, zx303TaskAdministratorException.ZX303TaskCreation{Reasons: []string{err.Error()}}
	}

	return &zx303TaskAdministrator.CreateResponse{
		ZX303Task: createResponse.ZX303Task,
	}, nil
}
