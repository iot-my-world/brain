package sigfoxBackendDataCallbackMessageRecordHandler

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	brainRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/pkg/recordHandler/exception"
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	sigfoxBackendDataCallbackMessageRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler"
	sigfoxBackendDataCallbackMessageRecordHandlerException "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler/exception"
)

type RecordHandler struct {
	sigfoxBackendDataCallbackMessageRecordHandler brainRecordHandler.RecordHandler
}

func New(
	brainMessageRecordHandler brainRecordHandler.RecordHandler,
) sigfoxBackendDataCallbackMessageRecordHandler.RecordHandler {

	return &RecordHandler{
		sigfoxBackendDataCallbackMessageRecordHandler: brainMessageRecordHandler,
	}
}

type CreateRequest struct {
	Message sigfoxBackendDataCallbackMessage.Message
}

type CreateResponse struct {
	Message sigfoxBackendDataCallbackMessage.Message
}

func (r *RecordHandler) ValidateCreateRequest(request *sigfoxBackendDataCallbackMessageRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *RecordHandler) Create(request *sigfoxBackendDataCallbackMessageRecordHandler.CreateRequest) (*sigfoxBackendDataCallbackMessageRecordHandler.CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.sigfoxBackendDataCallbackMessageRecordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.Message,
	}, &createResponse); err != nil {
		return nil, sigfoxBackendDataCallbackMessageRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdMessage, ok := createResponse.Entity.(*sigfoxBackendDataCallbackMessage.Message)
	if !ok {
		return nil, sigfoxBackendDataCallbackMessageRecordHandlerException.Create{Reasons: []string{"could not cast created entity to message"}}
	}

	return &sigfoxBackendDataCallbackMessageRecordHandler.CreateResponse{
		Message: *createdMessage,
	}, nil
}

func (r *RecordHandler) Retrieve(request *sigfoxBackendDataCallbackMessageRecordHandler.RetrieveRequest) (*sigfoxBackendDataCallbackMessageRecordHandler.RetrieveResponse, error) {
	retrievedMessage := sigfoxBackendDataCallbackMessage.Message{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedMessage,
	}
	if err := r.sigfoxBackendDataCallbackMessageRecordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, sigfoxBackendDataCallbackMessageRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &sigfoxBackendDataCallbackMessageRecordHandler.RetrieveResponse{
		Message: retrievedMessage,
	}, nil
}

func (r *RecordHandler) Update(request *sigfoxBackendDataCallbackMessageRecordHandler.UpdateRequest) (*sigfoxBackendDataCallbackMessageRecordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.sigfoxBackendDataCallbackMessageRecordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Message,
	}, &updateResponse); err != nil {
		return nil, sigfoxBackendDataCallbackMessageRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &sigfoxBackendDataCallbackMessageRecordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *sigfoxBackendDataCallbackMessageRecordHandler.DeleteRequest) (*sigfoxBackendDataCallbackMessageRecordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.sigfoxBackendDataCallbackMessageRecordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, sigfoxBackendDataCallbackMessageRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &sigfoxBackendDataCallbackMessageRecordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *sigfoxBackendDataCallbackMessageRecordHandler.CollectRequest) (*sigfoxBackendDataCallbackMessageRecordHandler.CollectResponse, error) {
	var collectedMessage []sigfoxBackendDataCallbackMessage.Message
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedMessage,
	}
	err := r.sigfoxBackendDataCallbackMessageRecordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, sigfoxBackendDataCallbackMessageRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedMessage == nil {
		collectedMessage = make([]sigfoxBackendDataCallbackMessage.Message, 0)
	}

	return &sigfoxBackendDataCallbackMessageRecordHandler.CollectResponse{
		Records: collectedMessage,
		Total:   collectResponse.Total,
	}, nil
}
