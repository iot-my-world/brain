package recordHandler

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/party/client"
	clientRecordHandler "github.com/iot-my-world/brain/party/client/recordHandler"
	clientRecordHandlerException "github.com/iot-my-world/brain/party/client/recordHandler/exception"
	brainRecordHandler "github.com/iot-my-world/brain/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainClientRecordHandler brainRecordHandler.RecordHandler,
) clientRecordHandler.RecordHandler {

	if brainClientRecordHandler == nil {
		log.Fatal(clientRecordHandlerException.RecordHandlerNil{}.Error())
	}
	return &RecordHandler{
		recordHandler: brainClientRecordHandler,
	}
}

func (r *RecordHandler) Create(request *clientRecordHandler.CreateRequest) (*clientRecordHandler.CreateResponse, error) {
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

	return &clientRecordHandler.CreateResponse{
		Client: *createdClient,
	}, nil
}

func (r *RecordHandler) Retrieve(request *clientRecordHandler.RetrieveRequest) (*clientRecordHandler.RetrieveResponse, error) {
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

	return &clientRecordHandler.RetrieveResponse{
		Client: retrievedClient,
	}, nil
}

func (r *RecordHandler) Update(request *clientRecordHandler.UpdateRequest) (*clientRecordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Client,
	}, &updateResponse); err != nil {
		return nil, clientRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &clientRecordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *clientRecordHandler.DeleteRequest) (*clientRecordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, clientRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &clientRecordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *clientRecordHandler.CollectRequest) (*clientRecordHandler.CollectResponse, error) {
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

	return &clientRecordHandler.CollectResponse{
		Records: collectedClients,
		Total:   collectResponse.Total,
	}, nil
}
