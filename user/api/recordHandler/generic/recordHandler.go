package recordHandler

import (
	brainException "github.com/iot-my-world/brain/exception"
	brainRecordHandler "github.com/iot-my-world/brain/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/recordHandler/exception"
	apiUser "github.com/iot-my-world/brain/user/api"
	apiUserRecordHandler "github.com/iot-my-world/brain/user/api/recordHandler"
	apiUserRecordHandlerException "github.com/iot-my-world/brain/user/api/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainAPIUserRecordHandler brainRecordHandler.RecordHandler,
) apiUserRecordHandler.RecordHandler {

	return &RecordHandler{
		recordHandler: brainAPIUserRecordHandler,
	}
}

func (r *RecordHandler) ValidateCreateRequest(request *apiUserRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *RecordHandler) Create(request *apiUserRecordHandler.CreateRequest) (*apiUserRecordHandler.CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.User,
	}, &createResponse); err != nil {
		return nil, apiUserRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdApiUser, ok := createResponse.Entity.(*apiUser.User)
	if !ok {
		return nil, apiUserRecordHandlerException.Create{Reasons: []string{"could not cast created entity to api user"}}
	}

	return &apiUserRecordHandler.CreateResponse{
		User: *createdApiUser,
	}, nil
}

func (r *RecordHandler) Retrieve(request *apiUserRecordHandler.RetrieveRequest) (*apiUserRecordHandler.RetrieveResponse, error) {
	retrievedUser := apiUser.User{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedUser,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, apiUserRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &apiUserRecordHandler.RetrieveResponse{
		User: retrievedUser,
	}, nil
}

func (r *RecordHandler) Update(request *apiUserRecordHandler.UpdateRequest) (*apiUserRecordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.User,
	}, &updateResponse); err != nil {
		return nil, apiUserRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &apiUserRecordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *apiUserRecordHandler.DeleteRequest) (*apiUserRecordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, apiUserRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &apiUserRecordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *apiUserRecordHandler.CollectRequest) (*apiUserRecordHandler.CollectResponse, error) {
	var collectedUser []apiUser.User
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedUser,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, apiUserRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedUser == nil {
		collectedUser = make([]apiUser.User, 0)
	}

	return &apiUserRecordHandler.CollectResponse{
		Records: collectedUser,
		Total:   collectResponse.Total,
	}, nil
}
