package recordHandler

import (
	brainException "github.com/iot-my-world/brain/exception"
	brainRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/pkg/recordHandler/exception"
	sf0012 "github.com/iot-my-world/brain/pkg/tracker/sf001"
	"github.com/iot-my-world/brain/pkg/tracker/sf001/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/sf001/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainSF001RecordHandler brainRecordHandler.RecordHandler,
) recordHandler.RecordHandler {

	return &RecordHandler{
		recordHandler: brainSF001RecordHandler,
	}
}

type CreateRequest struct {
	SF001 sf0012.SF001
}

type CreateResponse struct {
	SF001 sf0012.SF001
}

func (r *RecordHandler) ValidateCreateRequest(request *recordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *RecordHandler) Create(request *recordHandler.CreateRequest) (*recordHandler.CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.SF001,
	}, &createResponse); err != nil {
		return nil, exception.Create{Reasons: []string{err.Error()}}
	}
	createdDevice, ok := createResponse.Entity.(*sf0012.SF001)
	if !ok {
		return nil, exception.Create{Reasons: []string{"could not cast created entity to sf001"}}
	}

	return &recordHandler.CreateResponse{
		SF001: *createdDevice,
	}, nil
}

func (r *RecordHandler) Retrieve(request *recordHandler.RetrieveRequest) (*recordHandler.RetrieveResponse, error) {
	retrievedSF001 := sf0012.SF001{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedSF001,
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

	return &recordHandler.RetrieveResponse{
		SF001: retrievedSF001,
	}, nil
}

func (r *RecordHandler) Update(request *recordHandler.UpdateRequest) (*recordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.SF001,
	}, &updateResponse); err != nil {
		return nil, exception.Update{Reasons: []string{err.Error()}}
	}

	return &recordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *recordHandler.DeleteRequest) (*recordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, exception.Delete{Reasons: []string{err.Error()}}
	}

	return &recordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *recordHandler.CollectRequest) (*recordHandler.CollectResponse, error) {
	var collectedSF001 []sf0012.SF001
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedSF001,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, exception.Collect{Reasons: []string{err.Error()}}
	}

	if collectedSF001 == nil {
		collectedSF001 = make([]sf0012.SF001, 0)
	}

	return &recordHandler.CollectResponse{
		Records: collectedSF001,
		Total:   collectResponse.Total,
	}, nil
}
