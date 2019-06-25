package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	partyRegistrar "github.com/iot-my-world/brain/party/registrar"
	partyRegistrarJsonRpcAdaptor "github.com/iot-my-world/brain/party/registrar/adaptor/jsonRpc"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
)

type registrar struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) partyRegistrar.Registrar {
	return &registrar{
		jsonRpcClient: jsonRpcClient,
	}
}

func (r *registrar) RegisterSystemAdminUser(request *partyRegistrar.RegisterSystemAdminUserRequest) (*partyRegistrar.RegisterSystemAdminUserResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *registrar) ValidateInviteCompanyAdminUserRequest(request *partyRegistrar.InviteCompanyAdminUserRequest) error {
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

func (r *registrar) InviteCompanyAdminUser(request *partyRegistrar.InviteCompanyAdminUserRequest) (*partyRegistrar.InviteCompanyAdminUserResponse, error) {
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
	inviteCompanyAdminUserResponse := partyRegistrarJsonRpcAdaptor.InviteCompanyAdminUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		"PartyRegistrar.InviteCompanyAdminUser",
		partyRegistrarJsonRpcAdaptor.InviteCompanyAdminUserRequest{
			WrappedCompanyIdentifier: *companyIdentifier,
		},
		&inviteCompanyAdminUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &partyRegistrar.InviteCompanyAdminUserResponse{URLToken: inviteCompanyAdminUserResponse.URLToken}, nil
}

func (r *registrar) ValidateRegisterCompanyAdminUserRequest(request *partyRegistrar.RegisterCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterCompanyAdminUser(request *partyRegistrar.RegisterCompanyAdminUserRequest) (*partyRegistrar.RegisterCompanyAdminUserResponse, error) {
	if err := r.ValidateRegisterCompanyAdminUserRequest(request); err != nil {
		return nil, err
	}

	registerCompanyAdminUserResponse := partyRegistrarJsonRpcAdaptor.RegisterCompanyAdminUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		"PartyRegistrar.RegisterCompanyAdminUser",
		partyRegistrarJsonRpcAdaptor.RegisterCompanyAdminUserRequest{
			User: request.User,
		},
		&registerCompanyAdminUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &partyRegistrar.RegisterCompanyAdminUserResponse{User: registerCompanyAdminUserResponse.User}, nil
}

func (r *registrar) ValidateInviteCompanyUserRequest(request *partyRegistrar.InviteCompanyUserRequest) error {
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

func (r *registrar) InviteCompanyUser(request *partyRegistrar.InviteCompanyUserRequest) (*partyRegistrar.InviteCompanyUserResponse, error) {
	if err := r.ValidateInviteCompanyUserRequest(request); err != nil {
		return nil, err
	}

	//return &partyRegistrar.InviteCompanyUserResponse{URLToken: urlToken}, nil
	return nil, brainException.NotImplemented{}
}

func (r *registrar) ValidateRegisterCompanyUserRequest(request *partyRegistrar.RegisterCompanyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterCompanyUser(request *partyRegistrar.RegisterCompanyUserRequest) (*partyRegistrar.RegisterCompanyUserResponse, error) {
	if err := r.ValidateRegisterCompanyUserRequest(request); err != nil {
		return nil, err
	}

	registerCompanyUserResponse := partyRegistrarJsonRpcAdaptor.RegisterCompanyUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		"PartyRegistrar.RegisterCompanyUser",
		partyRegistrarJsonRpcAdaptor.RegisterCompanyUserRequest{
			User: request.User,
		},
		&registerCompanyUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &partyRegistrar.RegisterCompanyUserResponse{
		User: registerCompanyUserResponse.User,
	}, nil
}

func (r *registrar) ValidateInviteClientAdminUserRequest(request *partyRegistrar.InviteClientAdminUserRequest) error {
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

func (r *registrar) InviteClientAdminUser(request *partyRegistrar.InviteClientAdminUserRequest) (*partyRegistrar.InviteClientAdminUserResponse, error) {
	if err := r.ValidateInviteClientAdminUserRequest(request); err != nil {
		return nil, err
	}

	//return &partyRegistrar.InviteClientAdminUserResponse{URLToken: urlToken}, nil
	return nil, brainException.NotImplemented{}
}

func (r *registrar) ValidateRegisterClientAdminUserRequest(request *partyRegistrar.RegisterClientAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterClientAdminUser(request *partyRegistrar.RegisterClientAdminUserRequest) (*partyRegistrar.RegisterClientAdminUserResponse, error) {
	if err := r.ValidateRegisterClientAdminUserRequest(request); err != nil {
		return nil, err
	}

	//return &partyRegistrar.RegisterClientAdminUserResponse{User: userChangePasswordResponse.User}, nil
	return nil, brainException.NotImplemented{}
}

func (r *registrar) ValidateInviteClientUserRequest(request *partyRegistrar.InviteClientUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) InviteClientUser(request *partyRegistrar.InviteClientUserRequest) (*partyRegistrar.InviteClientUserResponse, error) {
	if err := r.ValidateInviteClientUserRequest(request); err != nil {
		return nil, err
	}

	//return &partyRegistrar.InviteClientUserResponse{URLToken: urlToken}, nil
	return nil, brainException.NotImplemented{}
}

func (r *registrar) ValidateRegisterClientUserRequest(request *partyRegistrar.RegisterClientUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterClientUser(request *partyRegistrar.RegisterClientUserRequest) (*partyRegistrar.RegisterClientUserResponse, error) {
	if err := r.ValidateRegisterClientUserRequest(request); err != nil {
		return nil, err
	}

	//return &partyRegistrar.RegisterClientUserResponse{User: userChangePasswordResponse.User}, nil
	return nil, brainException.NotImplemented{}
}

func (r *registrar) ValidateInviteUserRequest(request *partyRegistrar.InviteUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) InviteUser(request *partyRegistrar.InviteUserRequest) (*partyRegistrar.InviteUserResponse, error) {
	if err := r.ValidateInviteUserRequest(request); err != nil {
		return nil, err
	}

	id, err := wrappedIdentifier.Wrap(request.UserIdentifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	inviteUserResponse := partyRegistrarJsonRpcAdaptor.InviteUserResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		"PartyRegistrar.InviteUser",
		partyRegistrarJsonRpcAdaptor.InviteUserRequest{
			WrappedUserIdentifier: *id,
		},
		&inviteUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &partyRegistrar.InviteUserResponse{
		URLToken: inviteUserResponse.URLToken,
	}, nil
}

func (r *registrar) ValidateAreAdminsRegisteredRequest(request *partyRegistrar.AreAdminsRegisteredRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) AreAdminsRegistered(request *partyRegistrar.AreAdminsRegisteredRequest) (*partyRegistrar.AreAdminsRegisteredResponse, error) {
	if err := r.ValidateAreAdminsRegisteredRequest(request); err != nil {
		return nil, err
	}

	//return &response, nil
	return nil, brainException.NotImplemented{}
}
