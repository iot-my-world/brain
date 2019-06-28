package recordHandler

import (
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/party/individual"
	individualRecordHandler "github.com/iot-my-world/brain/party/individual/recordHandler"
	individualRecordHandlerException "github.com/iot-my-world/brain/party/individual/recordHandler/exception"
	brainRecordHandler "github.com/iot-my-world/brain/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainIndividualRecordHandler brainRecordHandler.RecordHandler,
) individualRecordHandler.RecordHandler {

	return &RecordHandler{
		recordHandler: brainIndividualRecordHandler,
	}
}

func (r *RecordHandler) ValidateCreateRequest(request *individualRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *RecordHandler) Create(request *individualRecordHandler.CreateRequest) (*individualRecordHandler.CreateResponse, error) {
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

	return &individualRecordHandler.CreateResponse{
		Individual: *createdIndividual,
	}, nil
}

func (r *RecordHandler) Retrieve(request *individualRecordHandler.RetrieveRequest) (*individualRecordHandler.RetrieveResponse, error) {
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

	return &individualRecordHandler.RetrieveResponse{
		Individual: retrievedIndividual,
	}, nil
}

func (r *RecordHandler) Update(request *individualRecordHandler.UpdateRequest) (*individualRecordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Individual,
	}, &updateResponse); err != nil {
		return nil, individualRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &individualRecordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *individualRecordHandler.DeleteRequest) (*individualRecordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, individualRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &individualRecordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *individualRecordHandler.CollectRequest) (*individualRecordHandler.CollectResponse, error) {
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

	return &individualRecordHandler.CollectResponse{
		Records: collectedIndividual,
		Total:   collectResponse.Total,
	}, nil
}
