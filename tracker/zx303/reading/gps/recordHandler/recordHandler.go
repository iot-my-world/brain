package recordHandler

import (
	brainException "github.com/iot-my-world/brain/exception"
	brainRecordHandler "github.com/iot-my-world/brain/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/recordHandler/exception"
	"github.com/iot-my-world/brain/search/criterion"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/query"
	"github.com/iot-my-world/brain/security/claims"
	zx303GPSReading "github.com/iot-my-world/brain/tracker/zx303/reading/gps"
	zx303GPSReadingRecordHandlerException "github.com/iot-my-world/brain/tracker/zx303/reading/gps/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainZX303GPSReadingRecordHandler brainRecordHandler.RecordHandler,
) *RecordHandler {

	return &RecordHandler{
		recordHandler: brainZX303GPSReadingRecordHandler,
	}
}

type CreateRequest struct {
	ZX303GPSReading zx303GPSReading.Reading
}

type CreateResponse struct {
	ZX303GPSReading zx303GPSReading.Reading
}

func (r *RecordHandler) ValidateCreateRequest(request *CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *RecordHandler) Create(request *CreateRequest) (*CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.ZX303GPSReading,
	}, &createResponse); err != nil {
		return nil, zx303GPSReadingRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdReading, ok := createResponse.Entity.(*zx303GPSReading.Reading)
	if !ok {
		return nil, zx303GPSReadingRecordHandlerException.Create{Reasons: []string{"could not cast created entity to zx303GPSReading"}}
	}

	return &CreateResponse{
		ZX303GPSReading: *createdReading,
	}, nil
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	ZX303GPSReading zx303GPSReading.Reading
}

func (r *RecordHandler) Retrieve(request *RetrieveRequest) (*RetrieveResponse, error) {
	retrievedZX303GPSReading := zx303GPSReading.Reading{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedZX303GPSReading,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, zx303GPSReadingRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &RetrieveResponse{
		ZX303GPSReading: retrievedZX303GPSReading,
	}, nil
}

type UpdateRequest struct {
	Claims          claims.Claims
	Identifier      identifier.Identifier
	ZX303GPSReading zx303GPSReading.Reading
}

type UpdateResponse struct{}

func (r *RecordHandler) Update(request *UpdateRequest) (*UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.ZX303GPSReading,
	}, &updateResponse); err != nil {
		return nil, zx303GPSReadingRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &UpdateResponse{}, nil
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
}

func (r *RecordHandler) Delete(request *DeleteRequest) (*DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, zx303GPSReadingRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &DeleteResponse{}, nil
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []zx303GPSReading.Reading
	Total   int
}

func (r *RecordHandler) Collect(request *CollectRequest) (*CollectResponse, error) {
	var collectedZX303GPSReading []zx303GPSReading.Reading
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedZX303GPSReading,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, zx303GPSReadingRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedZX303GPSReading == nil {
		collectedZX303GPSReading = make([]zx303GPSReading.Reading, 0)
	}

	return &CollectResponse{
		Records: collectedZX303GPSReading,
		Total:   collectResponse.Total,
	}, nil
}
