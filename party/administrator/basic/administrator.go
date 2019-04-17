package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	partyAdministrator "gitlab.com/iotTracker/brain/party/administrator"
	partyAdministratorException "gitlab.com/iotTracker/brain/party/administrator/exception"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	clientRecordHandlerException "gitlab.com/iotTracker/brain/party/client/recordHandler/exception"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	companyRecordHandlerException "gitlab.com/iotTracker/brain/party/company/recordHandler/exception"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	systemRecordHandlerException "gitlab.com/iotTracker/brain/party/system/recordHandler/exception"
)

type administrator struct {
	clientRecordHandler  *clientRecordHandler.RecordHandler
	companyRecordHandler *companyRecordHandler.RecordHandler
	systemRecordHandler  systemRecordHandler.RecordHandler
}

func New(
	clientRecordHandler *clientRecordHandler.RecordHandler,
	companyRecordHandler *companyRecordHandler.RecordHandler,
	systemRecordHandler systemRecordHandler.RecordHandler,
) partyAdministrator.Administrator {
	return &administrator{
		clientRecordHandler:  clientRecordHandler,
		companyRecordHandler: companyRecordHandler,
		systemRecordHandler:  systemRecordHandler,
	}
}

func (a *administrator) ValidateGetMyPartyRequest(request *partyAdministrator.GetMyPartyRequest) error {
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

func (a *administrator) GetMyParty(request *partyAdministrator.GetMyPartyRequest) (*partyAdministrator.GetMyPartyResponse, error) {
	if err := a.ValidateGetMyPartyRequest(request); err != nil {
		return nil, err
	}

	response := partyAdministrator.GetMyPartyResponse{}

	switch request.Claims.PartyDetails().PartyType {
	case party.System:
		systemRecordHandlerRetrieveResponse, err := a.systemRecordHandler.Retrieve(&systemRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Claims.PartyDetails().PartyId,
		})
		if err != nil {
			switch err.(type) {
			case systemRecordHandlerException.NotFound:
				return nil, partyAdministratorException.NotFound{}
			default:
				return nil, partyAdministratorException.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.PartyType = party.System
		response.Party = systemRecordHandlerRetrieveResponse.System

	case party.Company:
		companyRecordHandlerRetrieveResponse, err := a.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Claims.PartyDetails().PartyId,
		})
		if err != nil {
			switch err.(type) {
			case companyRecordHandlerException.NotFound:
				return nil, partyAdministratorException.NotFound{}
			default:
				return nil, partyAdministratorException.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.PartyType = party.Company
		response.Party = companyRecordHandlerRetrieveResponse.Company

	case party.Client:
		clientRecordHandlerRetrieveResponse, err := a.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Claims.PartyDetails().PartyId,
		})
		if err != nil {
			switch err.(type) {
			case clientRecordHandlerException.NotFound:
				return nil, partyAdministratorException.NotFound{}
			default:
				return nil, partyAdministratorException.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.PartyType = party.Client
		response.Party = clientRecordHandlerRetrieveResponse.Client

	default:
		return nil, partyAdministratorException.InvalidParty{Reasons: []string{string(request.Claims.PartyDetails().PartyType)}}
	}

	return &response, nil
}

func (a *administrator) ValidateRetrievePartyRequest(request *partyAdministrator.RetrievePartyRequest) error {
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

func (a *administrator) RetrieveParty(request *partyAdministrator.RetrievePartyRequest) (*partyAdministrator.RetrievePartyResponse, error) {
	if err := a.ValidateRetrievePartyRequest(request); err != nil {
		return nil, err
	}
	response := partyAdministrator.RetrievePartyResponse{}
	switch request.PartyType {
	case party.System:
		systemRecordHandlerRetrieveResponse, err := a.systemRecordHandler.Retrieve(&systemRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Identifier,
		})
		if err != nil {
			switch err.(type) {
			case systemRecordHandlerException.NotFound:
				return nil, partyAdministratorException.NotFound{}
			default:
				return nil, partyAdministratorException.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.Party = systemRecordHandlerRetrieveResponse.System

	case party.Company:
		companyRecordHandlerRetrieveResponse, err := a.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Identifier,
		})
		if err != nil {
			switch err.(type) {
			case companyRecordHandlerException.NotFound:
				return nil, partyAdministratorException.NotFound{}
			default:
				return nil, partyAdministratorException.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.Party = companyRecordHandlerRetrieveResponse.Company

	case party.Client:
		clientRecordHandlerRetrieveResponse, err := a.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Identifier,
		})
		if err != nil {
			switch err.(type) {
			case clientRecordHandlerException.NotFound:
				return nil, partyAdministratorException.NotFound{}
			default:
				return nil, partyAdministratorException.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.Party = clientRecordHandlerRetrieveResponse.Client

	default:
		return nil, partyAdministratorException.InvalidParty{Reasons: []string{string(request.Claims.PartyDetails().PartyType)}}
	}

	return &response, nil
}
