package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	registrar2 "github.com/iot-my-world/brain/pkg/party/registrar"
	"github.com/iot-my-world/brain/pkg/party/registrar/adaptor/jsonRpc"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
)

type registrar struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) registrar2.Registrar {
	return &registrar{
		jsonRpcClient: jsonRpcClient,
	}
}

func (r *registrar) RegisterSystemAdminUser(request *registrar2.RegisterSystemAdminUserRequest) (*registrar2.RegisterSystemAdminUserResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *registrar) ValidateInviteCompanyAdminUserRequest(request *registrar2.InviteCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.CompanyIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "company identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) InviteCompanyAdminUser(request *registrar2.InviteCompanyAdminUserRequest) (*registrar2.InviteCompanyAdminUserResponse, error) {
	if err := r.ValidateInviteCompanyAdminUserRequest(request); err != nil {
		return nil, err
	}

	// create identifier for the company entity
	companyIdentifier, err := wrappedIdentifier.Wrap(request.CompanyIdentifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// invite the admin user
	inviteCompanyAdminUserResponse := jsonRpc.InviteCompanyAdminUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		registrar2.InviteCompanyAdminUserService,
		jsonRpc.InviteCompanyAdminUserRequest{
			WrappedCompanyIdentifier: *companyIdentifier,
		},
		&inviteCompanyAdminUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.InviteCompanyAdminUserResponse{URLToken: inviteCompanyAdminUserResponse.URLToken}, nil
}

func (r *registrar) ValidateRegisterCompanyAdminUserRequest(request *registrar2.RegisterCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterCompanyAdminUser(request *registrar2.RegisterCompanyAdminUserRequest) (*registrar2.RegisterCompanyAdminUserResponse, error) {
	if err := r.ValidateRegisterCompanyAdminUserRequest(request); err != nil {
		return nil, err
	}

	registerCompanyAdminUserResponse := jsonRpc.RegisterCompanyAdminUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		registrar2.RegisterCompanyAdminUserService,
		jsonRpc.RegisterCompanyAdminUserRequest{
			User: request.User,
		},
		&registerCompanyAdminUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.RegisterCompanyAdminUserResponse{User: registerCompanyAdminUserResponse.User}, nil
}

func (r *registrar) ValidateInviteCompanyUserRequest(request *registrar2.InviteCompanyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) InviteCompanyUser(request *registrar2.InviteCompanyUserRequest) (*registrar2.InviteCompanyUserResponse, error) {
	if err := r.ValidateInviteCompanyUserRequest(request); err != nil {
		return nil, err
	}

	return nil, brainException.NotImplemented{}
}

func (r *registrar) ValidateRegisterCompanyUserRequest(request *registrar2.RegisterCompanyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterCompanyUser(request *registrar2.RegisterCompanyUserRequest) (*registrar2.RegisterCompanyUserResponse, error) {
	if err := r.ValidateRegisterCompanyUserRequest(request); err != nil {
		return nil, err
	}

	registerCompanyUserResponse := jsonRpc.RegisterCompanyUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		registrar2.RegisterCompanyUserService,
		jsonRpc.RegisterCompanyUserRequest{
			User: request.User,
		},
		&registerCompanyUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.RegisterCompanyUserResponse{
		User: registerCompanyUserResponse.User,
	}, nil
}

func (r *registrar) ValidateInviteClientAdminUserRequest(request *registrar2.InviteClientAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.ClientIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "clientIdentifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) InviteClientAdminUser(request *registrar2.InviteClientAdminUserRequest) (*registrar2.InviteClientAdminUserResponse, error) {
	if err := r.ValidateInviteClientAdminUserRequest(request); err != nil {
		return nil, err
	}

	// create identifier for the client entity
	clientIdentifier, err := wrappedIdentifier.Wrap(request.ClientIdentifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	inviteClientAdminUserResponse := jsonRpc.InviteClientAdminUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		registrar2.InviteClientAdminUserService,
		jsonRpc.InviteClientAdminUserRequest{
			WrappedClientIdentifier: *clientIdentifier,
		},
		&inviteClientAdminUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.InviteClientAdminUserResponse{
		URLToken: inviteClientAdminUserResponse.URLToken,
	}, nil
}

func (r *registrar) ValidateRegisterClientAdminUserRequest(request *registrar2.RegisterClientAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterClientAdminUser(request *registrar2.RegisterClientAdminUserRequest) (*registrar2.RegisterClientAdminUserResponse, error) {
	if err := r.ValidateRegisterClientAdminUserRequest(request); err != nil {
		return nil, err
	}

	registerClientAdminUserResponse := jsonRpc.RegisterClientAdminUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		registrar2.RegisterClientAdminUserService,
		jsonRpc.RegisterClientAdminUserRequest{
			User: request.User,
		},
		&registerClientAdminUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.RegisterClientAdminUserResponse{
		User: registerClientAdminUserResponse.User,
	}, nil
}

func (r *registrar) ValidateInviteClientUserRequest(request *registrar2.InviteClientUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) InviteClientUser(request *registrar2.InviteClientUserRequest) (*registrar2.InviteClientUserResponse, error) {
	if err := r.ValidateInviteClientUserRequest(request); err != nil {
		return nil, err
	}

	//return &partyRegistrar.InviteClientUserResponse{URLToken: urlToken}, nil
	return nil, brainException.NotImplemented{}
}

func (r *registrar) ValidateRegisterClientUserRequest(request *registrar2.RegisterClientUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterClientUser(request *registrar2.RegisterClientUserRequest) (*registrar2.RegisterClientUserResponse, error) {
	if err := r.ValidateRegisterClientUserRequest(request); err != nil {
		return nil, err
	}

	registerClientUserResponse := jsonRpc.RegisterClientUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		registrar2.RegisterClientUserService,
		jsonRpc.RegisterClientUserRequest{
			User: request.User,
		},
		&registerClientUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.RegisterClientUserResponse{
		User: registerClientUserResponse.User,
	}, nil
}

func (r *registrar) ValidateInviteUserRequest(request *registrar2.InviteUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) InviteUser(request *registrar2.InviteUserRequest) (*registrar2.InviteUserResponse, error) {
	if err := r.ValidateInviteUserRequest(request); err != nil {
		return nil, err
	}

	id, err := wrappedIdentifier.Wrap(request.UserIdentifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	inviteUserResponse := jsonRpc.InviteUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		registrar2.InviteUserService,
		jsonRpc.InviteUserRequest{
			WrappedUserIdentifier: *id,
		},
		&inviteUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.InviteUserResponse{
		URLToken: inviteUserResponse.URLToken,
	}, nil
}

func (r *registrar) ValidateAreAdminsRegisteredRequest(request *registrar2.AreAdminsRegisteredRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) AreAdminsRegistered(request *registrar2.AreAdminsRegisteredRequest) (*registrar2.AreAdminsRegisteredResponse, error) {
	if err := r.ValidateAreAdminsRegisteredRequest(request); err != nil {
		return nil, err
	}

	wrappedPartyIdentifiers := make([]wrappedIdentifier.Wrapped, 0)
	for _, partyIdentifier := range request.PartyIdentifiers {
		id, err := wrappedIdentifier.Wrap(partyIdentifier)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		wrappedPartyIdentifiers = append(wrappedPartyIdentifiers, *id)
	}

	areAdminsRegisteredResponse := jsonRpc.AreAdminsRegisteredResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		registrar2.AreAdminsRegisteredService,
		jsonRpc.AreAdminsRegisteredRequest{
			WrappedPartyIdentifiers: wrappedPartyIdentifiers,
		},
		&areAdminsRegisteredResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.AreAdminsRegisteredResponse{
		Result: areAdminsRegisteredResponse.Result,
	}, nil
}
