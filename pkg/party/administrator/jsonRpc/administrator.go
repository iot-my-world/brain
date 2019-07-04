package jsonRpc

import (
	"fmt"
	"github.com/go-errors/errors"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	"github.com/iot-my-world/brain/pkg/party"
	administrator2 "github.com/iot-my-world/brain/pkg/party/administrator"
	"github.com/iot-my-world/brain/pkg/party/administrator/adaptor/jsonRpc"
	client2 "github.com/iot-my-world/brain/pkg/party/client"
	company2 "github.com/iot-my-world/brain/pkg/party/company"
	partyAdministrator "github.com/iot-my-world/brain/pkg/partyarty/administrator"
	partyAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/partyarty/administrator/adaptor/jsonRpc"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) administrator2.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
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
		log.Error(err.Error())
		return nil, err
	}

	getMyPartyResponse := jsonRpc.GetMyPartyResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.GetMyPartyService,
		jsonRpc.GetMyPartyRequest{},
		&getMyPartyResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var typedParty party.Party
	var castSuccess bool
	switch getMyPartyResponse.PartyType {
	case party.Client:
		typedParty, castSuccess = getMyPartyResponse.Party.(client2.Client)
	case party.Company:
		typedParty, castSuccess = getMyPartyResponse.Party.(company2.Company)
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

	return &administrator2.GetMyPartyResponse{
		Party:     typedParty,
		PartyType: getMyPartyResponse.PartyType,
	}, nil
}

func (a *administrator) ValidateRetrievePartyRequest(request *administrator2.RetrievePartyRequest) error {
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

func (a *administrator) RetrieveParty(request *administrator2.RetrievePartyRequest) (*administrator2.RetrievePartyResponse, error) {
	if err := a.ValidateRetrievePartyRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	id, err := wrappedIdentifier.Wrap(request.Identifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	retrievePartyResponse := jsonRpc.RetrievePartyResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.RetrievePartyService,
		jsonRpc.RetrievePartyRequest{
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

	return &administrator2.RetrievePartyResponse{
		Party: unwrappedParty,
	}, nil
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

	createAndInviteCompanyResponse := jsonRpc.CreateAndInviteCompanyResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.CreateAndInviteCompanyService,
		jsonRpc.CreateAndInviteCompanyRequest{
			Company: request.Company,
		},
		&createAndInviteCompanyResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.CreateAndInviteCompanyResponse{
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
