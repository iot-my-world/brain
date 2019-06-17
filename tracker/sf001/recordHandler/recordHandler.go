package recordHandler

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	brainRecordHandler "gitlab.com/iotTracker/brain/recordHandler"
	brainRecordHandlerException "gitlab.com/iotTracker/brain/recordHandler/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/sf001"
	sf001RecordHandlerException "gitlab.com/iotTracker/brain/tracker/sf001/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainSF001RecordHandler brainRecordHandler.RecordHandler,
) *RecordHandler {

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
		Entity: &request.SF001,
	}, &createResponse); err != nil {
		return nil, sf001RecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdDevice, ok := createResponse.Entity.(*sf001.SF001)
	if !ok {
		return nil, sf001RecordHandlerException.Create{Reasons: []string{"could not cast created entity to sf001"}}
	}

	return &CreateResponse{
		SF001: *createdDevice,
	}, nil
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	SF001 sf001.SF001
}

func (r *RecordHandler) Retrieve(request *RetrieveRequest) (*RetrieveResponse, error) {
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

	return &RetrieveResponse{
		SF001: retrievedSF001,
	}, nil
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	SF001      sf001.SF001
}

type UpdateResponse struct{}

func (r *RecordHandler) Update(request *UpdateRequest) (*UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.SF001,
	}, &updateResponse); err != nil {
		return nil, sf001RecordHandlerException.Update{Reasons: []string{err.Error()}}
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
		return nil, sf001RecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &DeleteResponse{}, nil
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []sf001.SF001
	Total   int
}

func (r *RecordHandler) Collect(request *CollectRequest) (*CollectResponse, error) {
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

	return &CollectResponse{
		Records: collectedSF001,
		Total:   collectResponse.Total,
	}, nil
}
