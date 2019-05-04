package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	humanUser "gitlab.com/iotTracker/brain/user/human"
	humanUserAdministrator "gitlab.com/iotTracker/brain/user/human/administrator"
	"net/http"
)

type adaptor struct {
	humanUserAdministrator humanUserAdministrator.Administrator
}

func New(
	humanUserAdministrator humanUserAdministrator.Administrator,
) *adaptor {
	return &adaptor{
		humanUserAdministrator: humanUserAdministrator,
	}
}

type GetMyUserRequest struct{}

type GetMyUserResponse struct {
	User humanUser.User `json:"user"`
}

func (a *adaptor) GetMyUser(r *http.Request, request *GetMyUserRequest, response *GetMyUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	getMyUserResponse, err := a.humanUserAdministrator.GetMyUser(&humanUserAdministrator.GetMyUserRequest{
		Claims: claims,
	})
	if err != nil {
		return err
	}

	response.User = getMyUserResponse.User

	return nil
}

type UpdateAllowedFieldsRequest struct {
	User humanUser.User `json:"user"`
}

type UpdateAllowedFieldsResponse struct {
	User humanUser.User `json:"user"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.humanUserAdministrator.UpdateAllowedFields(&humanUserAdministrator.UpdateAllowedFieldsRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = updateAllowedFieldsResponse.User

	return nil
}

type CreateRequest struct {
	User humanUser.User `json:"user"`
}

type CreateResponse struct {
	User humanUser.User `json:"user"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.humanUserAdministrator.Create(&humanUserAdministrator.CreateRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = createResponse.User

	return nil
}

type UpdatePasswordRequest struct {
	ExistingPassword string `json:"existingPassword"`
	NewPassword      string `json:"newPassword"`
}

type UpdatePasswordResponse struct {
	User humanUser.User `json:"user"`
}

func (a *adaptor) UpdatePassword(r *http.Request, request *UpdatePasswordRequest, response *UpdatePasswordResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updatePasswordResponse, err := a.humanUserAdministrator.UpdatePassword(&humanUserAdministrator.UpdatePasswordRequest{
		Claims:           claims,
		ExistingPassword: request.ExistingPassword,
		NewPassword:      request.NewPassword,
	})
	if err != nil {
		return err
	}

	response.User = updatePasswordResponse.User

	return nil
}

type CheckPasswordRequest struct {
	Password string `json:"password"`
}

type CheckPasswordResponse struct {
	Result bool `json:"result"`
}

func (a *adaptor) CheckPassword(r *http.Request, request *CheckPasswordRequest, response *CheckPasswordResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	checkPasswordResponse, err := a.humanUserAdministrator.CheckPassword(&humanUserAdministrator.CheckPasswordRequest{
		Claims:   claims,
		Password: request.Password,
	})
	if err != nil {
		return err
	}

	response.Result = checkPasswordResponse.Result

	return nil
}

type ForgotPasswordRequest struct {
	UsernameOrEmailAddress string `json:"usernameOrEmailAddress"`
}

type ForgotPasswordResponse struct {
}

func (a *adaptor) ForgotPassword(r *http.Request, request *ForgotPasswordRequest, response *ForgotPasswordResponse) error {
	_, err := a.humanUserAdministrator.ForgotPassword(&humanUserAdministrator.ForgotPasswordRequest{
		UsernameOrEmailAddress: request.UsernameOrEmailAddress,
	})
	if err != nil {
		return err
	}

	return nil
}

type SetPasswordRequest struct {
	WrappedIdentifier wrappedIdentifier.Wrapped `json:"identifier"`
	NewPassword       string                    `json:"newPassword"`
}

type SetPasswordResponse struct{}

func (a *adaptor) SetPassword(r *http.Request, request *SetPasswordRequest, response *SetPasswordResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	if _, err := a.humanUserAdministrator.SetPassword(&humanUserAdministrator.SetPasswordRequest{
		Claims:      claims,
		Identifier:  request.WrappedIdentifier.Identifier,
		NewPassword: request.NewPassword,
	}); err != nil {
		return err
	}

	return nil
}
