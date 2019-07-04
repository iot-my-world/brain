package jsonRpc

import (
	"fmt"
	"github.com/go-errors/errors"
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/party"
	partyAdministrator "github.com/iot-my-world/brain/party/administrator"
	partyAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/party/administrator/adaptor/jsonRpc"
	"github.com/iot-my-world/brain/party/client"
	"github.com/iot-my-world/brain/party/company"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) partyAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
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
		log.Error(err.Error())
		return nil, err
	}

	getMyPartyResponse := partyAdministratorJsonRpcAdaptor.GetMyPartyResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		partyAdministrator.GetMyPartyService,
		partyAdministratorJsonRpcAdaptor.GetMyPartyRequest{},
		&getMyPartyResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var typedParty party.Party
	var castSuccess bool
	switch getMyPartyResponse.PartyType {
	case party.Client:
		typedParty, castSuccess = getMyPartyResponse.Party.(client.Client)
	case party.Company:
		typedParty, castSuccess = getMyPartyResponse.Party.(company.Company)
	default:
		err := errors.New("invalid party type in get my party response")
		log.Error(err.Error())
		return nil, err
	}

	if !castSuccess {
		err := errors.New("error casting party to particular type")
		log.Error(err.Error())
		return nil, err
	}

	return &partyAdministrator.GetMyPartyResponse{
		Party:     typedParty,
		PartyType: getMyPartyResponse.PartyType,
	}, nil
}

func (a *administrator) ValidateRetrievePartyRequest(request *partyAdministrator.RetrievePartyRequest) error {
	reasonsInvalid := make([]string, 0)

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
		log.Error(err.Error())
		return nil, err
	}

	id, err := wrappedIdentifier.Wrap(request.Identifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	retrievePartyResponse := partyAdministratorJsonRpcAdaptor.RetrievePartyResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		partyAdministrator.RetrievePartyService,
		partyAdministratorJsonRpcAdaptor.RetrievePartyRequest{
			PartyType:         request.PartyType,
			WrappedIdentifier: *id,
		},
		&retrievePartyResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	unwrappedParty, err := retrievePartyResponse.Party.UnWrap()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &partyAdministrator.RetrievePartyResponse{
		Party: unwrappedParty,
	}, nil
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

	createAndInviteCompanyResponse := partyAdministratorJsonRpcAdaptor.CreateAndInviteCompanyResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		partyAdministrator.CreateAndInviteCompanyService,
		partyAdministratorJsonRpcAdaptor.CreateAndInviteCompanyRequest{
			Company: request.Company,
		},
		&createAndInviteCompanyResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &partyAdministrator.CreateAndInviteCompanyResponse{
		RegistrationURLToken: createAndInviteCompanyResponse.RegistrationURLToken,
	}, nil
}

func (a *administrator) ValidateCreateAndInviteCompanyClientRequest(request *partyAdministrator.CreateAndInviteCompanyClientRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) CreateAndInviteCompanyClient(request *partyAdministrator.CreateAndInviteCompanyClientRequest) (*partyAdministrator.CreateAndInviteCompanyClientResponse, error) {
	if err := a.ValidateCreateAndInviteCompanyClientRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	createAndInviteCompanyClientResponse := partyAdministratorJsonRpcAdaptor.CreateAndInviteCompanyClientResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		partyAdministrator.CreateAndInviteCompanyClientService,
		partyAdministratorJsonRpcAdaptor.CreateAndInviteCompanyClientRequest{
			Client: request.Client,
		},
		&createAndInviteCompanyClientResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &partyAdministrator.CreateAndInviteCompanyClientResponse{
		RegistrationURLToken: createAndInviteCompanyClientResponse.RegistrationURLToken,
	}, nil
}

func (a *administrator) ValidateCreateAndInviteIndividualClientRequest(request *partyAdministrator.CreateAndInviteIndividualClientRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) CreateAndInviteIndividualClient(request *partyAdministrator.CreateAndInviteIndividualClientRequest) (*partyAdministrator.CreateAndInviteIndividualClientResponse, error) {
	if err := a.ValidateCreateAndInviteIndividualClientRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	createAndInviteIndividualClientResponse := partyAdministratorJsonRpcAdaptor.CreateAndInviteIndividualClientResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		partyAdministrator.CreateAndInviteIndividualClientService,
		partyAdministratorJsonRpcAdaptor.CreateAndInviteIndividualClientRequest{
			Client: request.Client,
		},
		&createAndInviteIndividualClientResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &partyAdministrator.CreateAndInviteIndividualClientResponse{
		RegistrationURLToken: createAndInviteIndividualClientResponse.RegistrationURLToken,
	}, nil
}
