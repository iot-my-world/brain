package recordHandler

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/party/client"
	clientRecordHandlerException "github.com/iot-my-world/brain/party/client/recordHandler/exception"
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
	brainClientRecordHandler brainRecordHandler.RecordHandler,
) *RecordHandler {

	if brainClientRecordHandler == nil {
		log.Fatal(clientRecordHandlerException.RecordHandlerNil{}.Error())
	}
	return &RecordHandler{
		recordHandler: brainClientRecordHandler,
	}
}

type CreateRequest struct {
	Client client.Client
}

type CreateResponse struct {
	Client client.Client
}

func (r *RecordHandler) Create(request *CreateRequest) (*CreateResponse, error) {
	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.Client,
	}, &createResponse); err != nil {
		return nil, clientRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdClient, ok := createResponse.Entity.(*client.Client)
	if !ok {
		return nil, clientRecordHandlerException.Create{Reasons: []string{"could not cast created entity to client"}}
	}

	return &CreateResponse{
		Client: *createdClient,
	}, nil
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Client client.Client
}

func (r *RecordHandler) Retrieve(request *RetrieveRequest) (*RetrieveResponse, error) {
	retrievedClient := client.Client{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedClient,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, clientRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &RetrieveResponse{
		Client: retrievedClient,
	}, nil
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Client     client.Client
}

type UpdateResponse struct{}

func (r *RecordHandler) Update(request *UpdateRequest) (*UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Client,
	}, &updateResponse); err != nil {
		return nil, clientRecordHandlerException.Update{Reasons: []string{err.Error()}}
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
		return nil, clientRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &DeleteResponse{}, nil
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []client.Client
	Total   int
}

func (r *RecordHandler) Collect(request *CollectRequest) (*CollectResponse, error) {
	var collectedClients []client.Client
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedClients,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, clientRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	return &CollectResponse{
		Records: collectedClients,
		Total:   collectResponse.Total,
	}, nil
}
