package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	readingAction "github.com/iot-my-world/brain/tracker/tk102/reading/action"
	readingAdministrator "github.com/iot-my-world/brain/tracker/tk102/reading/administrator"
	readingAdministratorException "github.com/iot-my-world/brain/tracker/tk102/reading/administrator/exception"
	readingRecordHandler "github.com/iot-my-world/brain/tracker/tk102/reading/recordHandler"
	readingValidator "github.com/iot-my-world/brain/tracker/tk102/reading/validator"
)

type administrator struct {
	readingRecordHandler readingRecordHandler.RecordHandler
	readingValidator     readingValidator.Validator
}

func New(
	readingRecordHandler readingRecordHandler.RecordHandler,
	readingValidator readingValidator.Validator,
) readingAdministrator.Administrator {
	return &administrator{
		readingRecordHandler: readingRecordHandler,
		readingValidator:     readingValidator,
	}
}

func (a *administrator) ValidateCreateRequest(request *readingAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	validateReadingResponse, err := a.readingValidator.Validate(&readingValidator.ValidateRequest{
		Reading: request.Reading,
		Action:  readingAction.Create,
	})
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "error validating reading: "+err.Error())
	}
	if len(validateReadingResponse.ReasonsInvalid) > 0 {
		for _, reason := range validateReadingResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("reading invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *readingAdministrator.CreateRequest) (*readingAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.readingRecordHandler.Create(&readingRecordHandler.CreateRequest{
		Reading: request.Reading,
	})
	if err != nil {
		return nil, readingAdministratorException.ReadingCreation{Reason: err.Error()}
	}

	return &readingAdministrator.CreateResponse{
		Reading: createResponse.Reading,
	}, nil
}

func (a *administrator) ValidateCreateBulkRequest(request *readingAdministrator.CreateBulkRequest) error {
	reasonsInvalid := make([]string, 0)

	for readingIdx := range request.Readings {
		validateReadingResponse, err := a.readingValidator.Validate(&readingValidator.ValidateRequest{
			Reading: request.Readings[readingIdx],
			Action:  readingAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating reading: "+err.Error())
		}
		if len(validateReadingResponse.ReasonsInvalid) > 0 {
			for _, reason := range validateReadingResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("reading invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) CreateBulk(request *readingAdministrator.CreateBulkRequest) (*readingAdministrator.CreateBulkResponse, error) {
	if err := a.ValidateCreateBulkRequest(request); err != nil {
		return nil, err
	}

	createBulkResponse, err := a.readingRecordHandler.CreateBulk(&readingRecordHandler.CreateBulkRequest{
		Readings: request.Readings,
	})
	if err != nil {
		return nil, readingAdministratorException.BulkReadingCreation{Reason: err.Error()}
	}

	return &readingAdministrator.CreateBulkResponse{
		Readings: createBulkResponse.Readings,
	}, nil
}
