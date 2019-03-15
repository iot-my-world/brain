package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	partyAdministrator "gitlab.com/iotTracker/brain/party/administrator"
	partyHandlerException "gitlab.com/iotTracker/brain/party/administrator/exception"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
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

func (bh *basicHandler) ValidateGetMyPartyRequest(request *partyAdministrator.GetMyPartyRequest) error {
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

func (bh *basicHandler) GetMyParty(request *partyAdministrator.GetMyPartyRequest, response *partyAdministrator.GetMyPartyResponse) error {
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

func (bh *basicHandler) ValidateRetrievePartyRequest(request *partyAdministrator.RetrievePartyRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}
	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}
	if !party.IsValidType(request.PartyType) {
		reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("party type '%s' is invalid", string(request.PartyType)))
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (bh *basicHandler) RetrieveParty(request *partyAdministrator.RetrievePartyRequest, response *partyAdministrator.RetrievePartyResponse) error {
	if err := bh.ValidateRetrievePartyRequest(request); err != nil {
		return err
	}

	switch request.PartyType {
	case party.System:
		systemRecordHandlerRetrieveResponse := systemRecordHandler.RetrieveResponse{}
		if err := bh.systemRecordHandler.Retrieve(&systemRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Identifier,
		}, &systemRecordHandlerRetrieveResponse); err != nil {
			return partyHandlerException.PartyRetrieval{Reasons: []string{err.Error()}}
		}
		response.PartyType = party.System
		response.Party = systemRecordHandlerRetrieveResponse.System

	case party.Company:
		companyRecordHandlerRetrieveResponse := companyRecordHandler.RetrieveResponse{}
		if err := bh.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Identifier,
		}, &companyRecordHandlerRetrieveResponse); err != nil {
			return partyHandlerException.PartyRetrieval{Reasons: []string{err.Error()}}
		}
		response.PartyType = party.Company
		response.Party = companyRecordHandlerRetrieveResponse.Company

	case party.Client:
		clientRecordHandlerRetrieveResponse := clientRecordHandler.RetrieveResponse{}
		if err := bh.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Identifier,
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
