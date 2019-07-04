package recordHandler

import (
	brainException "github.com/iot-my-world/brain/exception"
	brainRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/pkg/recordHandler/exception"
	"github.com/iot-my-world/brain/search/criterion"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/query"
	"github.com/iot-my-world/brain/security/claims"
	zx303StatusReading "github.com/iot-my-world/brain/tracker/zx303/reading/status"
	zx303StatusReadingRecordHandlerException "github.com/iot-my-world/brain/tracker/zx303/reading/status/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainZX303StatusReadingRecordHandler brainRecordHandler.RecordHandler,
) *RecordHandler {

	return &RecordHandler{
		recordHandler: brainZX303StatusReadingRecordHandler,
	}
}

type CreateRequest struct {
	ZX303StatusReading zx303StatusReading.Reading
}

type CreateResponse struct {
	ZX303StatusReading zx303StatusReading.Reading
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
		Entity: &request.ZX303StatusReading,
	}, &createResponse); err != nil {
		return nil, zx303StatusReadingRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdReading, ok := createResponse.Entity.(*zx303StatusReading.Reading)
	if !ok {
		return nil, zx303StatusReadingRecordHandlerException.Create{Reasons: []string{"could not cast created entity to zx303StatusReading"}}
	}

	return &CreateResponse{
		ZX303StatusReading: *createdReading,
	}, nil
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	ZX303StatusReading zx303StatusReading.Reading
}

func (r *RecordHandler) Retrieve(request *RetrieveRequest) (*RetrieveResponse, error) {
	retrievedZX303StatusReading := zx303StatusReading.Reading{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedZX303StatusReading,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, zx303StatusReadingRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &RetrieveResponse{
		ZX303StatusReading: retrievedZX303StatusReading,
	}, nil
}

type UpdateRequest struct {
	Claims             claims.Claims
	Identifier         identifier.Identifier
	ZX303StatusReading zx303StatusReading.Reading
}

type UpdateResponse struct{}

func (r *RecordHandler) Update(request *UpdateRequest) (*UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.ZX303StatusReading,
	}, &updateResponse); err != nil {
		return nil, zx303StatusReadingRecordHandlerException.Update{Reasons: []string{err.Error()}}
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
		return nil, zx303StatusReadingRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &DeleteResponse{}, nil
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []zx303StatusReading.Reading
	Total   int
}

func (r *RecordHandler) Collect(request *CollectRequest) (*CollectResponse, error) {
	var collectedZX303StatusReading []zx303StatusReading.Reading
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedZX303StatusReading,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, zx303StatusReadingRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedZX303StatusReading == nil {
		collectedZX303StatusReading = make([]zx303StatusReading.Reading, 0)
	}

	return &CollectResponse{
		Records: collectedZX303StatusReading,
		Total:   collectResponse.Total,
	}, nil
}
