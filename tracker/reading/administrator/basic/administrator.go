package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	readingAdministrator "gitlab.com/iotTracker/brain/tracker/reading/administrator"
	readingAdministratorException "gitlab.com/iotTracker/brain/tracker/reading/administrator/exception"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
	readingValidator "gitlab.com/iotTracker/brain/tracker/reading/validator"
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
