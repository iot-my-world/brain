package backendRecordHandler

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	brainRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/pkg/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/sigfox/backend"
	backendRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler"
	backendRecordHandlerException "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler/exception"
)

type RecordHandler struct {
	backendRecordHandler brainRecordHandler.RecordHandler
}

func New(
	brainBackendRecordHandler brainRecordHandler.RecordHandler,
) backendRecordHandler.RecordHandler {

	return &RecordHandler{
		backendRecordHandler: brainBackendRecordHandler,
	}
}

type CreateRequest struct {
	Backend backend.Backend
}

type CreateResponse struct {
	Backend backend.Backend
}

func (r *RecordHandler) ValidateCreateRequest(request *backendRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *RecordHandler) Create(request *backendRecordHandler.CreateRequest) (*backendRecordHandler.CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.backendRecordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.Backend,
	}, &createResponse); err != nil {
		return nil, backendRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdDevice, ok := createResponse.Entity.(*backend.Backend)
	if !ok {
		return nil, backendRecordHandlerException.Create{Reasons: []string{"could not cast created entity to backend"}}
	}

	return &backendRecordHandler.CreateResponse{
		Backend: *createdDevice,
	}, nil
}

func (r *RecordHandler) Retrieve(request *backendRecordHandler.RetrieveRequest) (*backendRecordHandler.RetrieveResponse, error) {
	retrievedBackend := backend.Backend{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedBackend,
	}
	if err := r.backendRecordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, backendRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &backendRecordHandler.RetrieveResponse{
		Backend: retrievedBackend,
	}, nil
}

func (r *RecordHandler) Update(request *backendRecordHandler.UpdateRequest) (*backendRecordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.backendRecordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Backend,
	}, &updateResponse); err != nil {
		return nil, backendRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &backendRecordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *backendRecordHandler.DeleteRequest) (*backendRecordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.backendRecordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, backendRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &backendRecordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *backendRecordHandler.CollectRequest) (*backendRecordHandler.CollectResponse, error) {
	var collectedBackend []backend.Backend
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedBackend,
	}
	err := r.backendRecordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, backendRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedBackend == nil {
		collectedBackend = make([]backend.Backend, 0)
	}

	return &backendRecordHandler.CollectResponse{
		Records: collectedBackend,
		Total:   collectResponse.Total,
	}, nil
}
