package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	humanUserAdministrator "github.com/iot-my-world/brain/user/human/administrator"
	userAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/user/human/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) humanUserAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *humanUserAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *humanUserAdministrator.UpdateAllowedFieldsRequest) (*humanUserAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	updateAllowedFieldsResponse := userAdministratorJsonRpcAdaptor.UpdateAllowedFieldsResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"UserAdministrator.Create",
		userAdministratorJsonRpcAdaptor.UpdateAllowedFieldsRequest{
			User: request.User,
		},
		&updateAllowedFieldsResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &humanUserAdministrator.UpdateAllowedFieldsResponse{
		User: updateAllowedFieldsResponse.User,
	}, nil
}

func (a *administrator) ValidateGetMyUserRequest(request *humanUserAdministrator.GetMyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) GetMyUser(request *humanUserAdministrator.GetMyUserRequest) (*humanUserAdministrator.GetMyUserResponse, error) {
	if err := a.ValidateGetMyUserRequest(request); err != nil {
		return nil, err
	}

	getMyUserResponse := userAdministratorJsonRpcAdaptor.GetMyUserResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"UserAdministrator.GetMyUser",
		userAdministratorJsonRpcAdaptor.GetMyUserRequest{},
		&getMyUserResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &humanUserAdministrator.GetMyUserResponse{User: getMyUserResponse.User}, nil
}

func (a *administrator) ValidateCreateRequest(request *humanUserAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *humanUserAdministrator.CreateRequest) (*humanUserAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse := userAdministratorJsonRpcAdaptor.CreateResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"UserAdministrator.Create",
		userAdministratorJsonRpcAdaptor.CreateRequest{},
		&createResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &humanUserAdministrator.CreateResponse{User: createResponse.User}, nil
}

func (a *administrator) ValidateSetPasswordRequest(request *humanUserAdministrator.SetPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) SetPassword(request *humanUserAdministrator.SetPasswordRequest) (*humanUserAdministrator.SetPasswordResponse, error) {
	if err := a.ValidateSetPasswordRequest(request); err != nil {
		return nil, err
	}

	setPasswordResponse := userAdministratorJsonRpcAdaptor.SetPasswordResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"UserAdministrator.SetPassword",
		userAdministratorJsonRpcAdaptor.SetPasswordRequest{},
		&setPasswordResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &humanUserAdministrator.SetPasswordResponse{}, nil
}

func (a *administrator) ValidateUpdatePasswordRequest(request *humanUserAdministrator.UpdatePasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdatePassword(request *humanUserAdministrator.UpdatePasswordRequest) (*humanUserAdministrator.UpdatePasswordResponse, error) {
	if err := a.ValidateUpdatePasswordRequest(request); err != nil {
		return nil, err
	}

	updatePasswordResponse := userAdministratorJsonRpcAdaptor.UpdatePasswordResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"UserAdministrator.UpdatePassword",
		userAdministratorJsonRpcAdaptor.UpdatePasswordRequest{
			ExistingPassword: request.ExistingPassword,
			NewPassword:      request.NewPassword,
		},
		&updatePasswordResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &humanUserAdministrator.UpdatePasswordResponse{
		User: updatePasswordResponse.User,
	}, nil
}

func (a *administrator) ValidateCheckPasswordRequest(request *humanUserAdministrator.CheckPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) CheckPassword(request *humanUserAdministrator.CheckPasswordRequest) (*humanUserAdministrator.CheckPasswordResponse, error) {
	if err := a.ValidateCheckPasswordRequest(request); err != nil {
		return nil, err
	}

	updatePasswordResponse := userAdministratorJsonRpcAdaptor.CheckPasswordResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"UserAdministrator.CheckPassword",
		userAdministratorJsonRpcAdaptor.CheckPasswordRequest{
			Password: request.Password,
		},
		&updatePasswordResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &humanUserAdministrator.CheckPasswordResponse{
		Result: updatePasswordResponse.Result,
	}, nil
}

func (a *administrator) ValidateForgotPasswordRequest(request *humanUserAdministrator.ForgotPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) ForgotPassword(request *humanUserAdministrator.ForgotPasswordRequest) (*humanUserAdministrator.ForgotPasswordResponse, error) {
	if err := a.ValidateForgotPasswordRequest(request); err != nil {
		return nil, err
	}

	forgotPasswordResponse := userAdministratorJsonRpcAdaptor.ForgotPasswordResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"UserAdministrator.ForgotPassword",
		userAdministratorJsonRpcAdaptor.ForgotPasswordRequest{
			UsernameOrEmailAddress: request.UsernameOrEmailAddress,
		},
		&forgotPasswordResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &humanUserAdministrator.ForgotPasswordResponse{}, nil
}
