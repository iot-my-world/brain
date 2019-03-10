package basic

import (
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	partyHandler "gitlab.com/iotTracker/brain/party/handler"
	brainException "gitlab.com/iotTracker/brain/exception"
	partyHandlerException "gitlab.com/iotTracker/brain/party/handler/exception"
	"gitlab.com/iotTracker/brain/party"
)

type basicHandler struct {
	clientRecordHandler  clientRecordHandler.RecordHandler
	companyRecordHandler companyRecordHandler.RecordHandler
	systemRecordHandler  systemRecordHandler.RecordHandler
}

func New(
	clientRecordHandler clientRecordHandler.RecordHandler,
	companyRecordHandler companyRecordHandler.RecordHandler,
	systemRecordHandler systemRecordHandler.RecordHandler,
) *basicHandler {
	return &basicHandler{
		clientRecordHandler:  clientRecordHandler,
		companyRecordHandler: companyRecordHandler,
		systemRecordHandler:  systemRecordHandler,
	}
}

func (bh *basicHandler) ValidateGetMyPartyRequest(request *partyHandler.GetMyPartyRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (bh *basicHandler) GetMyParty(request *partyHandler.GetMyPartyRequest, response *partyHandler.GetMyPartyResponse) error {
	if err := bh.ValidateGetMyPartyRequest(request); err != nil {
		return err
	}

	switch request.Claims.PartyDetails().PartyType {
	case party.System:
		systemRecordHandlerRetrieveResponse := systemRecordHandler.RetrieveResponse{}
		if err := bh.systemRecordHandler.Retrieve(&systemRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Claims.PartyDetails().PartyId,
		}, &systemRecordHandlerRetrieveResponse); err != nil {
			return partyHandlerException.PartyRetrieval{Reasons: []string{err.Error()}}
		}
		response.PartyType = party.System
		response.Party = systemRecordHandlerRetrieveResponse.System

	case party.Company:
		companyRecordHandlerRetrieveResponse := companyRecordHandler.RetrieveResponse{}
		if err := bh.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Claims.PartyDetails().PartyId,
		}, &companyRecordHandlerRetrieveResponse); err != nil {
			return partyHandlerException.PartyRetrieval{Reasons: []string{err.Error()}}
		}
		response.PartyType = party.Company
		response.Party = companyRecordHandlerRetrieveResponse.Company

	case party.Client:
		clientRecordHandlerRetrieveResponse := clientRecordHandler.RetrieveResponse{}
		if err := bh.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Claims.PartyDetails().PartyId,
		}, &clientRecordHandlerRetrieveResponse); err != nil {
			return partyHandlerException.PartyRetrieval{Reasons: []string{err.Error()}}
		}
		response.PartyType = party.Client
		response.Party = clientRecordHandlerRetrieveResponse.Client

	default:
		return partyHandlerException.InvalidParty{Reasons: []string{string(request.Claims.PartyDetails().PartyType)}}
	}

	return nil
}
