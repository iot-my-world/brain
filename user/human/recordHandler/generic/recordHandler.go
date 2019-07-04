package recordHandler

import (
	"github.com/iot-my-world/brain/log"
	brainRecordHandler "github.com/iot-my-world/brain/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/recordHandler/exception"
	humanUser "github.com/iot-my-world/brain/user/human"
	userRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler"
	userRecordHandlerException "github.com/iot-my-world/brain/user/human/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainUserRecordHandler brainRecordHandler.RecordHandler,
) userRecordHandler.RecordHandler {

	if brainUserRecordHandler == nil {
		log.Fatal(userRecordHandlerException.RecordHandlerNil{}.Error())
	}
	return &RecordHandler{
		recordHandler: brainUserRecordHandler,
	}
}

func (r *RecordHandler) Create(request *userRecordHandler.CreateRequest) (*userRecordHandler.CreateResponse, error) {
	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.User,
	}, &createResponse); err != nil {
		return nil, userRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdUser, ok := createResponse.Entity.(*humanUser.User)
	if !ok {
		return nil, userRecordHandlerException.Create{Reasons: []string{"could not cast created entity to user"}}
	}

	return &userRecordHandler.CreateResponse{
		User: *createdUser,
	}, nil
}

func (r *RecordHandler) Retrieve(request *userRecordHandler.RetrieveRequest) (*userRecordHandler.RetrieveResponse, error) {
	retrievedUser := humanUser.User{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedUser,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, userRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &userRecordHandler.RetrieveResponse{
		User: retrievedUser,
	}, nil
}

func (r *RecordHandler) Update(request *userRecordHandler.UpdateRequest) (*userRecordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.User,
	}, &updateResponse); err != nil {
		return nil, userRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &userRecordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *userRecordHandler.DeleteRequest) (*userRecordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, userRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &userRecordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *userRecordHandler.CollectRequest) (*userRecordHandler.CollectResponse, error) {
	var collectedUsers []humanUser.User
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedUsers,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, userRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedUsers == nil {
		collectedUsers = make([]humanUser.User, 0)
	}

	return &userRecordHandler.CollectResponse{
		Records: collectedUsers,
		Total:   collectResponse.Total,
	}, nil
}
