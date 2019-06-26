package recordHandler

import (
	brainException "github.com/iot-my-world/brain/exception"
	brainRecordHandler "github.com/iot-my-world/brain/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/recordHandler/exception"
	"github.com/iot-my-world/brain/tracker/sf001"
	sf001RecordHandler "github.com/iot-my-world/brain/tracker/sf001/recordHandler"
	sf001RecordHandlerException "github.com/iot-my-world/brain/tracker/sf001/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainSF001RecordHandler brainRecordHandler.RecordHandler,
) sf001RecordHandler.RecordHandler {

	return &RecordHandler{
		recordHandler: brainSF001RecordHandler,
	}
}

type CreateRequest struct {
	SF001 sf001.SF001
}

type CreateResponse struct {
	SF001 sf001.SF001
}

func (r *RecordHandler) ValidateCreateRequest(request *sf001RecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *RecordHandler) Create(request *sf001RecordHandler.CreateRequest) (*sf001RecordHandler.CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.SF001,
	}, &createResponse); err != nil {
		return nil, sf001RecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdDevice, ok := createResponse.Entity.(*sf001.SF001)
	if !ok {
		return nil, sf001RecordHandlerException.Create{Reasons: []string{"could not cast created entity to sf001"}}
	}

	return &sf001RecordHandler.CreateResponse{
		SF001: *createdDevice,
	}, nil
}

func (r *RecordHandler) Retrieve(request *sf001RecordHandler.RetrieveRequest) (*sf001RecordHandler.RetrieveResponse, error) {
	retrievedSF001 := sf001.SF001{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedSF001,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, sf001RecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &sf001RecordHandler.RetrieveResponse{
		SF001: retrievedSF001,
	}, nil
}

func (r *RecordHandler) Update(request *sf001RecordHandler.UpdateRequest) (*sf001RecordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.SF001,
	}, &updateResponse); err != nil {
		return nil, sf001RecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &sf001RecordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *sf001RecordHandler.DeleteRequest) (*sf001RecordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, sf001RecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &sf001RecordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *sf001RecordHandler.CollectRequest) (*sf001RecordHandler.CollectResponse, error) {
	var collectedSF001 []sf001.SF001
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedSF001,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, sf001RecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedSF001 == nil {
		collectedSF001 = make([]sf001.SF001, 0)
	}

	return &sf001RecordHandler.CollectResponse{
		Records: collectedSF001,
		Total:   collectResponse.Total,
	}, nil
}
