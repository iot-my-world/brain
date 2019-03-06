package user

import (
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/user"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	RecordHandler userRecordHandler.RecordHandler
}

func New(recordHandler userRecordHandler.RecordHandler) *adaptor {
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
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createUserResponse := userRecordHandler.CreateResponse{}
	if err := s.RecordHandler.Create(
		&userRecordHandler.CreateRequest{
			Claims: claims,
			User:   request.User,
		}, &createUserResponse); err != nil {
		return err
	}

	response.User = createUserResponse.User

	return nil
}

type RetrieveRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type RetrieveResponse struct {
	User user.User `json:"user" bson:"user"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrieveUserResponse := userRecordHandler.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&userRecordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: id,
		},
		&retrieveUserResponse); err != nil {
		return err
	}

	response.User = retrieveUserResponse.User

	return nil
}

type UpdateRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
	User       user.User                           `json:"user"`
}

type UpdateResponse struct {
	User user.User `json:"user"`
}

func (s *adaptor) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	updateUserResponse := userRecordHandler.UpdateResponse{}
	if err := s.RecordHandler.Update(
		&userRecordHandler.UpdateRequest{
			Claims:     claims,
			Identifier: id,
		},
		&updateUserResponse); err != nil {
		return err
	}

	response.User = updateUserResponse.User

	return nil
}

type DeleteRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type DeleteResponse struct {
	User user.User `json:"user"`
}

func (s *adaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	deleteUserResponse := userRecordHandler.DeleteResponse{}
	if err := s.RecordHandler.Delete(
		&userRecordHandler.DeleteRequest{
			Identifier: id,
		},
		&deleteUserResponse); err != nil {
		return err
	}

	response.User = deleteUserResponse.User

	return nil
}

type ValidateRequest struct {
	User   user.User  `json:"user"`
	Method api.Method `json:"method"`
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid `json:"reasonsInvalid"`
}

func (s *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {

	validateUserResponse := userRecordHandler.ValidateResponse{}
	if err := s.RecordHandler.Validate(
		&userRecordHandler.ValidateRequest{
			User:   request.User,
			Method: request.Method,
		},
		&validateUserResponse); err != nil {
		return err
	}

	response.ReasonsInvalid = validateUserResponse.ReasonsInvalid

	return nil
}
