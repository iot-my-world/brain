package user

import (
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	RecordHandler user.RecordHandler
}

func New(recordHandler user.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type CreateRequest struct {
	User user.User `json:"user"`
}

type CreateResponse struct {
	User user.User `json:"user"`
}

func (s *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	createUserResponse := user.CreateResponse{}

	if err := s.RecordHandler.Create(
		&user.CreateRequest{
			User: request.User,
		},
		&createUserResponse); err != nil {
		return err
	}

	response.User = createUserResponse.User

	return nil
}

type RetrieveRequest struct {
	Identifier search.WrappedIdentifier `json:"identifier"`
}

type RetrieveResponse struct {
	User user.User `json:"user" bson:"user"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrieveUserResponse := user.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&user.RetrieveRequest{
			Identifier: id,
		},
		&retrieveUserResponse); err != nil {
		return err
	}

	response.User = retrieveUserResponse.User

	return nil
}

type UpdateRequest struct {
	Identifier search.WrappedIdentifier `json:"identifier"`
	User       user.User                `json:"user"`
}

type UpdateResponse struct {
	User user.User `json:"user"`
}

func (s *adaptor) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	updateUserResponse := user.UpdateResponse{}
	if err := s.RecordHandler.Update(
		&user.UpdateRequest{
			Identifier: id,
		},
		&updateUserResponse); err != nil {
		return err
	}

	response.User = updateUserResponse.User

	return nil
}

type DeleteRequest struct {
	Identifier search.WrappedIdentifier `json:"identifier"`
}

type DeleteResponse struct {
	User user.User `json:"user"`
}

func (s *adaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	deleteUserResponse := user.DeleteResponse{}
	if err := s.RecordHandler.Delete(
		&user.DeleteRequest{
			Identifier: id,
		},
		&deleteUserResponse); err != nil {
		return err
	}

	response.User = deleteUserResponse.User

	return nil
}

type ValidateRequest struct {
	User user.User `json:"user"`
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid `json:"reasonsInvalid"`
}

func (s *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {

	validateUserResponse := user.ValidateResponse{}
	if err := s.RecordHandler.Validate(
		&user.ValidateRequest{
			User: request.User,
		},
		&validateUserResponse); err != nil {
		return err
	}

	response.ReasonsInvalid = validateUserResponse.ReasonsInvalid

	return nil
}
