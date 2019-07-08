package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	administrator2 "github.com/iot-my-world/brain/pkg/user/human/administrator"
	"github.com/iot-my-world/brain/pkg/user/human/administrator/adaptor/jsonRpc"
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

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *administrator2.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *administrator2.UpdateAllowedFieldsRequest) (*administrator2.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	updateAllowedFieldsResponse := jsonRpc.UpdateAllowedFieldsResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.UpdateAllowedFieldsService,
		jsonRpc.UpdateAllowedFieldsRequest{
			User: request.User,
		},
		&updateAllowedFieldsResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.UpdateAllowedFieldsResponse{
		User: updateAllowedFieldsResponse.User,
	}, nil
}

func (a *administrator) ValidateGetMyUserRequest(request *administrator2.GetMyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) GetMyUser(request *administrator2.GetMyUserRequest) (*administrator2.GetMyUserResponse, error) {
	if err := a.ValidateGetMyUserRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	getMyUserResponse := jsonRpc.GetMyUserResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.GetMyUserService,
		jsonRpc.GetMyUserRequest{},
		&getMyUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.GetMyUserResponse{User: getMyUserResponse.User}, nil
}

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	createResponse := jsonRpc.CreateResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.CreateService,
		jsonRpc.CreateRequest{
			User: request.User,
		},
		&createResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.CreateResponse{User: createResponse.User}, nil
}

func (a *administrator) ValidateSetPasswordRequest(request *administrator2.SetPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) SetPassword(request *administrator2.SetPasswordRequest) (*administrator2.SetPasswordResponse, error) {
	if err := a.ValidateSetPasswordRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	id, err := wrappedIdentifier.Wrap(request.Identifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	setPasswordResponse := jsonRpc.SetPasswordResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.SetPasswordService,
		jsonRpc.SetPasswordRequest{
			WrappedIdentifier: *id,
			NewPassword:       request.NewPassword,
		},
		&setPasswordResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.SetPasswordResponse{}, nil
}

func (a *administrator) ValidateUpdatePasswordRequest(request *administrator2.UpdatePasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdatePassword(request *administrator2.UpdatePasswordRequest) (*administrator2.UpdatePasswordResponse, error) {
	if err := a.ValidateUpdatePasswordRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	updatePasswordResponse := jsonRpc.UpdatePasswordResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.UpdatePasswordService,
		jsonRpc.UpdatePasswordRequest{
			ExistingPassword: request.ExistingPassword,
			NewPassword:      request.NewPassword,
		},
		&updatePasswordResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.UpdatePasswordResponse{
		User: updatePasswordResponse.User,
	}, nil
}

func (a *administrator) ValidateCheckPasswordRequest(request *administrator2.CheckPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) CheckPassword(request *administrator2.CheckPasswordRequest) (*administrator2.CheckPasswordResponse, error) {
	if err := a.ValidateCheckPasswordRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	updatePasswordResponse := jsonRpc.CheckPasswordResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.CheckPasswordService,
		jsonRpc.CheckPasswordRequest{
			Password: request.Password,
		},
		&updatePasswordResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.CheckPasswordResponse{
		Result: updatePasswordResponse.Result,
	}, nil
}

func (a *administrator) ValidateForgotPasswordRequest(request *administrator2.ForgotPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) ForgotPassword(request *administrator2.ForgotPasswordRequest) (*administrator2.ForgotPasswordResponse, error) {
	if err := a.ValidateForgotPasswordRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	forgotPasswordResponse := jsonRpc.ForgotPasswordResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.ForgotPasswordService,
		jsonRpc.ForgotPasswordRequest{
			UsernameOrEmailAddress: request.UsernameOrEmailAddress,
		},
		&forgotPasswordResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.ForgotPasswordResponse{URLToken: forgotPasswordResponse.URLToken}, nil
}
