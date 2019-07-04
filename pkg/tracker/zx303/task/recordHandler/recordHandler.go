package recordHandler

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	brainRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/pkg/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/query"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainZX303TaskRecordHandler brainRecordHandler.RecordHandler,
) *RecordHandler {

	return &RecordHandler{
		recordHandler: brainZX303TaskRecordHandler,
	}
}

type CreateRequest struct {
	ZX303Task task.Task
}

type CreateResponse struct {
	ZX303Task task.Task
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
		Entity: &request.ZX303Task,
	}, &createResponse); err != nil {
		return nil, exception.Create{Reasons: []string{err.Error()}}
	}
	createdReading, ok := createResponse.Entity.(*task.Task)
	if !ok {
		return nil, exception.Create{Reasons: []string{"could not cast created entity to zx303Task"}}
	}

	return &CreateResponse{
		ZX303Task: *createdReading,
	}, nil
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	ZX303Task task.Task
}

func (r *RecordHandler) Retrieve(request *RetrieveRequest) (*RetrieveResponse, error) {
	retrievedZX303Task := task.Task{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedZX303Task,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, exception.NotFound{}
		default:
			return nil, err
		}
	}

	return &RetrieveResponse{
		ZX303Task: retrievedZX303Task,
	}, nil
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	ZX303Task  task.Task
}

type UpdateResponse struct{}

func (r *RecordHandler) Update(request *UpdateRequest) (*UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.ZX303Task,
	}, &updateResponse); err != nil {
		return nil, exception.Update{Reasons: []string{err.Error()}}
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
		return nil, exception.Delete{Reasons: []string{err.Error()}}
	}

	return &DeleteResponse{}, nil
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []task.Task
	Total   int
}

func (r *RecordHandler) Collect(request *CollectRequest) (*CollectResponse, error) {
	var collectedZX303Task []task.Task
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedZX303Task,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, exception.Collect{Reasons: []string{err.Error()}}
	}

	if collectedZX303Task == nil {
		collectedZX303Task = make([]task.Task, 0)
	}

	return &CollectResponse{
		Records: collectedZX303Task,
		Total:   collectResponse.Total,
	}, nil
}