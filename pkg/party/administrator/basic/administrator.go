package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/party"
	administrator2 "github.com/iot-my-world/brain/pkg/party/administrator"
	"github.com/iot-my-world/brain/pkg/party/administrator/exception"
	administrator3 "github.com/iot-my-world/brain/pkg/party/client/administrator"
	"github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	exception2 "github.com/iot-my-world/brain/pkg/party/client/recordHandler/exception"
	administrator4 "github.com/iot-my-world/brain/pkg/party/company/administrator"
	recordHandler2 "github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	exception4 "github.com/iot-my-world/brain/pkg/party/company/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/party/registrar"
	recordHandler3 "github.com/iot-my-world/brain/pkg/party/system/recordHandler"
	exception3 "github.com/iot-my-world/brain/pkg/party/system/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
)

type administrator struct {
	clientRecordHandler  recordHandler.RecordHandler
	companyRecordHandler recordHandler2.RecordHandler
	systemRecordHandler  recordHandler3.RecordHandler
	systemClaims         *humanUserLoginClaims.Login
	companyAdministrator administrator4.Administrator
	clientAdministrator  administrator3.Administrator
	partyRegistrar       registrar.Registrar
}

func New(
	clientRecordHandler recordHandler.RecordHandler,
	companyRecordHandler recordHandler2.RecordHandler,
	systemRecordHandler recordHandler3.RecordHandler,
	systemClaims *humanUserLoginClaims.Login,
	companyAdministrator administrator4.Administrator,
	clientAdministrator administrator3.Administrator,
	partyRegistrar registrar.Registrar,
) administrator2.Administrator {
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

func (a *administrator) ValidateGetMyPartyRequest(request *administrator2.GetMyPartyRequest) error {
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

func (a *administrator) GetMyParty(request *administrator2.GetMyPartyRequest) (*administrator2.GetMyPartyResponse, error) {
	if err := a.ValidateGetMyPartyRequest(request); err != nil {
		return nil, err
	}

	response := administrator2.GetMyPartyResponse{}

	switch request.Claims.PartyDetails().PartyType {
	case party.System:
		systemRecordHandlerRetrieveResponse, err := a.systemRecordHandler.Retrieve(&recordHandler3.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Claims.PartyDetails().PartyId,
		})
		if err != nil {
			switch err.(type) {
			case exception3.NotFound:
				return nil, exception.NotFound{}
			default:
				return nil, exception.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.PartyType = party.System
		response.Party = systemRecordHandlerRetrieveResponse.System

	case party.Company:
		companyRecordHandlerRetrieveResponse, err := a.companyRecordHandler.Retrieve(&recordHandler2.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Claims.PartyDetails().PartyId,
		})
		if err != nil {
			switch err.(type) {
			case exception4.NotFound:
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
			case exception2.NotFound:
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

func (a *administrator) ValidateRetrievePartyRequest(request *administrator2.RetrievePartyRequest) error {
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

func (a *administrator) RetrieveParty(request *administrator2.RetrievePartyRequest) (*administrator2.RetrievePartyResponse, error) {
	if err := a.ValidateRetrievePartyRequest(request); err != nil {
		return nil, err
	}
	response := administrator2.RetrievePartyResponse{}
	switch request.PartyType {
	case party.System:
		systemRecordHandlerRetrieveResponse, err := a.systemRecordHandler.Retrieve(&recordHandler3.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Identifier,
		})
		if err != nil {
			switch err.(type) {
			case exception3.NotFound:
				return nil, exception.NotFound{}
			default:
				return nil, exception.PartyRetrieval{Reasons: []string{err.Error()}}
			}
		}
		response.Party = systemRecordHandlerRetrieveResponse.System

	case party.Company:
		companyRecordHandlerRetrieveResponse, err := a.companyRecordHandler.Retrieve(&recordHandler2.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.Identifier,
		})
		if err != nil {
			switch err.(type) {
			case exception4.NotFound:
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
			case exception2.NotFound:
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

func (a *administrator) ValidateCreateAndInviteCompanyRequest(request *administrator2.CreateAndInviteCompanyRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) CreateAndInviteCompany(request *administrator2.CreateAndInviteCompanyRequest) (*administrator2.CreateAndInviteCompanyResponse, error) {
	if err := a.ValidateCreateAndInviteCompanyRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// create company via company administrator
	createResponse, err := a.companyAdministrator.Create(&administrator4.CreateRequest{
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

	return &administrator2.CreateAndInviteCompanyResponse{RegistrationURLToken: inviteResponse.URLToken}, nil
}

func (a *administrator) ValidateCreateAndInviteClientRequest(request *administrator2.CreateAndInviteClientRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) CreateAndInviteClient(request *administrator2.CreateAndInviteClientRequest) (*administrator2.CreateAndInviteClientResponse, error) {
	if err := a.ValidateCreateAndInviteClientRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// create client via client administrator
	createResponse, err := a.clientAdministrator.Create(&administrator3.CreateRequest{
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

	return &administrator2.CreateAndInviteClientResponse{RegistrationURLToken: inviteResponse.URLToken}, nil
}
