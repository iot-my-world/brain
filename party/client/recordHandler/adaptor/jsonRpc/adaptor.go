package client

import (
	"gitlab.com/iotTracker/brain/party/client"
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/validate"
	"net/http"
)

type adaptor struct {
	RecordHandler client.RecordHandler
}

func New(recordHandler client.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type CreateRequest struct {
	Client client.Client `json:"client"`
}

type CreateResponse struct {
	Client client.Client `json:"client"`
}

func (s *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	createClientResponse := client.CreateResponse{}

	if err := s.RecordHandler.Create(
		&client.CreateRequest{
			Client: request.Client,
		},
		&createClientResponse); err != nil {
		return err
	}

	response.Client = createClientResponse.Client

	return nil
}

type RetrieveRequest struct {
	Identifier search.WrappedIdentifier `json:"identifier"`
}

type RetrieveResponse struct {
	Client client.Client `json:"client" bson:"client"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrieveClientResponse := client.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&client.RetrieveRequest{
			Identifier: id,
		},
		&retrieveClientResponse); err != nil {
		return err
	}

	response.Client = retrieveClientResponse.Client

	return nil
}

type UpdateRequest struct {
	Identifier search.WrappedIdentifier `json:"identifier"`
	Client     client.Client            `json:"client"`
}

type UpdateResponse struct {
	Client client.Client `json:"client"`
}

func (s *adaptor) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	updateClientResponse := client.UpdateResponse{}
	if err := s.RecordHandler.Update(
		&client.UpdateRequest{
			Identifier: id,
		},
		&updateClientResponse); err != nil {
		return err
	}

	response.Client = updateClientResponse.Client

	return nil
}

type DeleteRequest struct {
	Identifier search.WrappedIdentifier `json:"identifier"`
}

type DeleteResponse struct {
	Client client.Client `json:"client"`
}

func (s *adaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	deleteClientResponse := client.DeleteResponse{}
	if err := s.RecordHandler.Delete(
		&client.DeleteRequest{
			Identifier: id,
		},
		&deleteClientResponse); err != nil {
		return err
	}

	response.Client = deleteClientResponse.Client

	return nil
}

type ValidateRequest struct {
	Client client.Client `json:"client"`
}

type ValidateResponse struct {
	ReasonsInvalid []validate.ReasonInvalid `json:"reasonsInvalid"`
}

func (s *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {

	validateClientResponse := client.ValidateResponse{}
	if err := s.RecordHandler.Validate(
		&client.ValidateRequest{
			Client: request.Client,
		},
		&validateClientResponse); err != nil {
		return err
	}

	response.ReasonsInvalid = validateClientResponse.ReasonsInvalid

	return nil
}
