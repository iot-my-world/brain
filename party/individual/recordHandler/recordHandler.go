package recordHandler

import (
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/party/individual"
	individualRecordHandlerException "github.com/iot-my-world/brain/party/individual/recordHandler/exception"
	brainRecordHandler "github.com/iot-my-world/brain/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/recordHandler/exception"
	"github.com/iot-my-world/brain/search/criterion"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/query"
	"github.com/iot-my-world/brain/security/claims"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainIndividualRecordHandler brainRecordHandler.RecordHandler,
) *RecordHandler {

	return &RecordHandler{
		recordHandler: brainIndividualRecordHandler,
	}
}

type CreateRequest struct {
	Individual individual.Individual
}

type CreateResponse struct {
	Individual individual.Individual
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
		Entity: &request.Individual,
	}, &createResponse); err != nil {
		return nil, individualRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdIndividual, ok := createResponse.Entity.(*individual.Individual)
	if !ok {
		return nil, individualRecordHandlerException.Create{Reasons: []string{"could not cast created entity to individual"}}
	}

	return &CreateResponse{
		Individual: *createdIndividual,
	}, nil
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Individual individual.Individual
}

func (r *RecordHandler) Retrieve(request *RetrieveRequest) (*RetrieveResponse, error) {
	retrievedIndividual := individual.Individual{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedIndividual,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, individualRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &RetrieveResponse{
		Individual: retrievedIndividual,
	}, nil
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Individual individual.Individual
}

type UpdateResponse struct{}

func (r *RecordHandler) Update(request *UpdateRequest) (*UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Individual,
	}, &updateResponse); err != nil {
		return nil, individualRecordHandlerException.Update{Reasons: []string{err.Error()}}
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
		return nil, individualRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &DeleteResponse{}, nil
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []individual.Individual
	Total   int
}

func (r *RecordHandler) Collect(request *CollectRequest) (*CollectResponse, error) {
	var collectedIndividual []individual.Individual
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedIndividual,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, individualRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedIndividual == nil {
		collectedIndividual = make([]individual.Individual, 0)
	}

	return &CollectResponse{
		Records: collectedIndividual,
		Total:   collectResponse.Total,
	}, nil
}
