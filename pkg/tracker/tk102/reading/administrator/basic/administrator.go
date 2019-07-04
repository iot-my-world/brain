package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/tracker/tk102/reading/action"
	administrator2 "github.com/iot-my-world/brain/pkg/tracker/tk102/reading/administrator"
	"github.com/iot-my-world/brain/pkg/tracker/tk102/reading/administrator/exception"
	"github.com/iot-my-world/brain/pkg/tracker/tk102/reading/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/tk102/reading/validator"
)

type administrator struct {
	readingRecordHandler recordHandler.RecordHandler
	readingValidator     validator.Validator
}

func New(
	readingRecordHandler recordHandler.RecordHandler,
	readingValidator validator.Validator,
) administrator2.Administrator {
	return &administrator{
		readingRecordHandler: readingRecordHandler,
		readingValidator:     readingValidator,
	}
}

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	validateReadingResponse, err := a.readingValidator.Validate(&validator.ValidateRequest{
		Reading: request.Reading,
		Action:  action.Create,
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

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.readingRecordHandler.Create(&recordHandler.CreateRequest{
		Reading: request.Reading,
	})
	if err != nil {
		return nil, exception.ReadingCreation{Reason: err.Error()}
	}

	return &administrator2.CreateResponse{
		Reading: createResponse.Reading,
	}, nil
}

func (a *administrator) ValidateCreateBulkRequest(request *administrator2.CreateBulkRequest) error {
	reasonsInvalid := make([]string, 0)

	for readingIdx := range request.Readings {
		validateReadingResponse, err := a.readingValidator.Validate(&validator.ValidateRequest{
			Reading: request.Readings[readingIdx],
			Action:  action.Create,
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

func (a *administrator) CreateBulk(request *administrator2.CreateBulkRequest) (*administrator2.CreateBulkResponse, error) {
	if err := a.ValidateCreateBulkRequest(request); err != nil {
		return nil, err
	}

	createBulkResponse, err := a.readingRecordHandler.CreateBulk(&recordHandler.CreateBulkRequest{
		Readings: request.Readings,
	})
	if err != nil {
		return nil, exception.BulkReadingCreation{Reason: err.Error()}
	}

	return &administrator2.CreateBulkResponse{
		Readings: createBulkResponse.Readings,
	}, nil
}
