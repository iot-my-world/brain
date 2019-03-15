package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	partyHandler "gitlab.com/iotTracker/brain/party/handler"
	partyHandlerException "gitlab.com/iotTracker/brain/party/handler/exception"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/claims/login"
)

type basicHandler struct {
	clientRecordHandler  clientRecordHandler.RecordHandler
	companyRecordHandler companyRecordHandler.RecordHandler
	systemRecordHandler  systemRecordHandler.RecordHandler
	userRecordHandler    userRecordHandler.RecordHandler
}

func New(
	clientRecordHandler clientRecordHandler.RecordHandler,
	companyRecordHandler companyRecordHandler.RecordHandler,
	systemRecordHandler systemRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
) *basicHandler {
	return &basicHandler{
		clientRecordHandler:  clientRecordHandler,
		companyRecordHandler: companyRecordHandler,
		systemRecordHandler:  systemRecordHandler,
		userRecordHandler:    userRecordHandler,
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

func (bh *basicHandler) ValidateGetMyUserRequest(request *partyHandler.GetMyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// claims must be login claims to be able to get user
		if request.Claims.Type() != claims.Login {
			reasonsInvalid = append(reasonsInvalid, "claims must be of type login")
		}
	}
	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (bh *basicHandler) GetMyUser(request *partyHandler.GetMyUserRequest, response *partyHandler.GetMyUserResponse) error {
	if err := bh.ValidateGetMyUserRequest(request); err != nil {
		return err
	}

	// parse the claims to login claims
	loginClaims, ok := request.Claims.(login.Login)
	if !ok {
		return partyHandlerException.InvalidClaims{Reasons: []string{"cannot assert login claims type"}}
	}

	// retrieve user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := bh.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	}, &userRetrieveResponse); err != nil {
		return partyHandlerException.PartyRetrieval{Reasons: []string{"user retrieval", err.Error()}}
	}

	response.User = userRetrieveResponse.User

	return nil
}

func (bh *basicHandler) ValidateRetrievePartyRequest(request *partyHandler.RetrievePartyRequest) error {
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

func (bh *basicHandler) RetrieveParty(request *partyHandler.RetrievePartyRequest, response *partyHandler.RetrievePartyResponse) error {
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
