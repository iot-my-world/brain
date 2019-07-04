package recordHandler

import (
	"github.com/iot-my-world/brain/internal/log"
	client2 "github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	"github.com/iot-my-world/brain/pkg/party/client/recordHandler/exception"
	brainRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/pkg/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainClientRecordHandler brainRecordHandler.RecordHandler,
) recordHandler.RecordHandler {

	if brainClientRecordHandler == nil {
		log.Fatal(exception.RecordHandlerNil{}.Error())
	}
	return &RecordHandler{
		recordHandler: brainClientRecordHandler,
	}
}

func (r *RecordHandler) Create(request *recordHandler.CreateRequest) (*recordHandler.CreateResponse, error) {
	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.Client,
	}, &createResponse); err != nil {
		return nil, exception.Create{Reasons: []string{err.Error()}}
	}
	createdClient, ok := createResponse.Entity.(*client2.Client)
	if !ok {
		return nil, exception.Create{Reasons: []string{"could not cast created entity to client"}}
	}

	return &recordHandler.CreateResponse{
		Client: *createdClient,
	}, nil
}

func (r *RecordHandler) Retrieve(request *recordHandler.RetrieveRequest) (*recordHandler.RetrieveResponse, error) {
	retrievedClient := client2.Client{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedClient,
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
		Client: retrievedClient,
	}, nil
}

func (r *RecordHandler) Update(request *recordHandler.UpdateRequest) (*recordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Client,
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
	var collectedClients []client2.Client
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedClients,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, exception.Collect{Reasons: []string{err.Error()}}
	}

	if collectedClients == nil {
		collectedClients = make([]client2.Client, 0)
	}

	return &recordHandler.CollectResponse{
		Records: collectedClients,
		Total:   collectResponse.Total,
	}, nil
}
