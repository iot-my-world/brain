package recordHandler

import (
	"github.com/iot-my-world/brain/internal/log"
	brainRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/pkg/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/user/human"
	"github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	"github.com/iot-my-world/brain/pkg/user/human/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainUserRecordHandler brainRecordHandler.RecordHandler,
) recordHandler.RecordHandler {

	if brainUserRecordHandler == nil {
		log.Fatal(exception.RecordHandlerNil{}.Error())
	}
	return &RecordHandler{
		recordHandler: brainUserRecordHandler,
	}
}

func (r *RecordHandler) Create(request *recordHandler.CreateRequest) (*recordHandler.CreateResponse, error) {
	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.User,
	}, &createResponse); err != nil {
		return nil, exception.Create{Reasons: []string{err.Error()}}
	}
	createdUser, ok := createResponse.Entity.(*human.User)
	if !ok {
		return nil, exception.Create{Reasons: []string{"could not cast created entity to user"}}
	}

	return &recordHandler.CreateResponse{
		User: *createdUser,
	}, nil
}

func (r *RecordHandler) Retrieve(request *recordHandler.RetrieveRequest) (*recordHandler.RetrieveResponse, error) {
	retrievedUser := human.User{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedUser,
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
		User: retrievedUser,
	}, nil
}

func (r *RecordHandler) Update(request *recordHandler.UpdateRequest) (*recordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.User,
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
	var collectedUsers []human.User
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedUsers,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, exception.Collect{Reasons: []string{err.Error()}}
	}

	if collectedUsers == nil {
		collectedUsers = make([]human.User, 0)
	}

	return &recordHandler.CollectResponse{
		Records: collectedUsers,
		Total:   collectResponse.Total,
	}, nil
}
