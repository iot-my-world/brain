package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/party"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	"github.com/iot-my-world/brain/pkg/party/administrator/exception"
	clientAdministrator "github.com/iot-my-world/brain/pkg/party/client/administrator"
	"github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	clientRecordHandlerException "github.com/iot-my-world/brain/pkg/party/client/recordHandler/exception"
	companyAdministrator "github.com/iot-my-world/brain/pkg/party/company/administrator"
	companyRecordHandler "github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	companyRecordHandlerException "github.com/iot-my-world/brain/pkg/party/company/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/party/registrar"
	systemRecordHandler "github.com/iot-my-world/brain/pkg/party/system/recordHandler"
	systemRecordHandlerException "github.com/iot-my-world/brain/pkg/party/system/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
)

type administrator struct {
	clientRecordHandler  recordHandler.RecordHandler
	companyRecordHandler companyRecordHandler.RecordHandler
	systemRecordHandler  systemRecordHandler.RecordHandler
	systemClaims         *humanUserLoginClaims.Login
	companyAdministrator companyAdministrator.Administrator
	clientAdministrator  clientAdministrator.Administrator
	partyRegistrar       registrar.Registrar
}

func New(
	clientRecordHandler recordHandler.RecordHandler,
	companyRecordHandler companyRecordHandler.RecordHandler,
	systemRecordHandler systemRecordHandler.RecordHandler,
	systemClaims *humanUserLoginClaims.Login,
	companyAdministrator companyAdministrator.Administrator,
	clientAdministrator clientAdministrator.Administrator,
	partyRegistrar registrar.Registrar,
) partyAdministrator.Administrator {
	return &administrator{
		clientRecordHandler:  clientRecordHandler,
		companyRecordHandler: companyRecordHandler,
		systemRecordHandler:  systemRecordHandler,
		systemClaims:         systemClaims,
		companyAdministrator: companyAdministrator,
		clientAdministrator:  clientAdministrator,
		partyRegistrar:       partyRegistrar,
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
				return nil, exception.NotFound{}
			default:
				return nil, exception.PartyRetrieval{Reasons: []string{err.Error()}}
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
				return nil, exception.NotFound{}
			default:
				return nil, exception.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.PartyType = party.Company
		response.Party = companyRecordHandlerRetrieveResponse.Company

	case party.Client:
		clientRecordHandlerRetrieveResponse, err := a.clientRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Claims.PartyDetails().PartyId,
		})
		if err != nil {
			switch err.(type) {
			case clientRecordHandlerException.NotFound:
				return nil, exception.NotFound{}
			default:
				return nil, exception.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.PartyType = party.Client
		response.Party = clientRecordHandlerRetrieveResponse.Client

	default:
		return nil, exception.InvalidParty{Reasons: []string{string(request.Claims.PartyDetails().PartyType)}}
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
				return nil, exception.NotFound{}
			default:
				return nil, exception.PartyRetrieval{Reasons: []string{err.Error()}}
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
				return nil, exception.NotFound{}
			default:
				return nil, exception.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.Party = companyRecordHandlerRetrieveResponse.Company

	case party.Client:
		clientRecordHandlerRetrieveResponse, err := a.clientRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Identifier,
		})
		if err != nil {
			switch err.(type) {
			case clientRecordHandlerException.NotFound:
				return nil, exception.NotFound{}
			default:
				return nil, exception.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.Party = clientRecordHandlerRetrieveResponse.Client

	default:
		return nil, exception.InvalidParty{Reasons: []string{string(request.Claims.PartyDetails().PartyType)}}
	}

	return &response, nil
}

func (a *administrator) ValidateCreateAndInviteCompanyRequest(request *partyAdministrator.CreateAndInviteCompanyRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) CreateAndInviteCompany(request *partyAdministrator.CreateAndInviteCompanyRequest) (*partyAdministrator.CreateAndInviteCompanyResponse, error) {
	if err := a.ValidateCreateAndInviteCompanyRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// set company parent details to system
	request.Company.ParentPartyType = a.systemClaims.PartyType
	request.Company.ParentId = a.systemClaims.PartyId

	// create company via company administrator
	createResponse, err := a.companyAdministrator.Create(&companyAdministrator.CreateRequest{
		Claims:  a.systemClaims,
		Company: request.Company,
	})
	if err != nil {
		err = exception.CreateAndInviteCompany{Reasons: []string{"company create", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// invite company admin user
	inviteResponse, err := a.partyRegistrar.InviteCompanyAdminUser(&registrar.InviteCompanyAdminUserRequest{
		Claims: a.systemClaims,
		CompanyIdentifier: id.Identifier{
			Id: createResponse.Company.Id,
		},
	})
	if err != nil {
		err = exception.CreateAndInviteCompany{Reasons: []string{"invite company admin user", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &partyAdministrator.CreateAndInviteCompanyResponse{RegistrationURLToken: inviteResponse.URLToken}, nil
}

func (a *administrator) ValidateCreateAndInviteClientRequest(request *partyAdministrator.CreateAndInviteClientRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) CreateAndInviteClient(request *partyAdministrator.CreateAndInviteClientRequest) (*partyAdministrator.CreateAndInviteClientResponse, error) {
	if err := a.ValidateCreateAndInviteClientRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// set client parent details to system
	request.Client.ParentPartyType = a.systemClaims.PartyType
	request.Client.ParentId = a.systemClaims.PartyId

	// create client via client administrator
	createResponse, err := a.clientAdministrator.Create(&clientAdministrator.CreateRequest{
		Claims: a.systemClaims,
		Client: request.Client,
	})
	if err != nil {
		err = exception.CreateAndInviteClient{Reasons: []string{"client create", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// invite client admin user
	inviteResponse, err := a.partyRegistrar.InviteClientAdminUser(&registrar.InviteClientAdminUserRequest{
		Claims: a.systemClaims,
		ClientIdentifier: id.Identifier{
			Id: createResponse.Client.Id,
		},
	})
	if err != nil {
		err = exception.CreateAndInviteClient{Reasons: []string{"invite client admin user", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &partyAdministrator.CreateAndInviteClientResponse{RegistrationURLToken: inviteResponse.URLToken}, nil
}
