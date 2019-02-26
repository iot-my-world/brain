package system

import (
	"gitlab.com/iotTracker/brain/party/system"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"net/http"
)

type adaptor struct {
	RecordHandler systemRecordHandler.RecordHandler
}

func New(recordHandler systemRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type RetrieveRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type RetrieveResponse struct {
	System system.System `json:"system" bson:"system"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrieveSystemResponse := systemRecordHandler.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&systemRecordHandler.RetrieveRequest{
			Identifier: id,
		},
		&retrieveSystemResponse); err != nil {
		return err
	}

	response.System = retrieveSystemResponse.System

	return nil
}
