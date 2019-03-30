package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/user"
	userAdministrator "gitlab.com/iotTracker/brain/user/administrator"
	"net/http"
)

type adaptor struct {
	userAdministrator userAdministrator.Administrator
}

func New(
	userAdministrator userAdministrator.Administrator,
) *adaptor {
	return &adaptor{
		userAdministrator: userAdministrator,
	}
}

type GetMyUserRequest struct{}

type GetMyUserResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) GetMyUser(r *http.Request, request *GetMyUserRequest, response *GetMyUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	getMyUserResponse, err := a.userAdministrator.GetMyUser(&userAdministrator.GetMyUserRequest{
		Claims: claims,
	})
	if err != nil {
		return err
	}

	response.User = getMyUserResponse.User

	return nil
}

type UpdateAllowedFieldsRequest struct {
	User user.User `json:"user"`
}

type UpdateAllowedFieldsResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.userAdministrator.UpdateAllowedFields(&userAdministrator.UpdateAllowedFieldsRequest{
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
	User user.User `json:"user"`
}

type CreateResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.userAdministrator.Create(&userAdministrator.CreateRequest{
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
	User user.User `json:"user"`
}

func (a *adaptor) UpdatePassword(r *http.Request, request *UpdatePasswordRequest, response *UpdatePasswordResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updatePasswordResponse, err := a.userAdministrator.UpdatePassword(&userAdministrator.UpdatePasswordRequest{
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

	checkPasswordResponse, err := a.userAdministrator.CheckPassword(&userAdministrator.CheckPasswordRequest{
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
	_, err := a.userAdministrator.ForgotPassword(&userAdministrator.ForgotPasswordRequest{
		UsernameOrEmailAddress: request.UsernameOrEmailAddress,
	})
	if err != nil {
		return err
	}

	return nil
}

type SetPasswordRequest struct {
	Identifier  wrappedIdentifier.Wrapped `json:"identifier"`
	NewPassword string                    `json:"newPassword"`
}

type SetPasswordResponse struct{}

func (a *adaptor) SetPassword(r *http.Request, request *SetPasswordRequest, response *SetPasswordResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	if _, err := a.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      claims,
		Identifier:  id,
		NewPassword: request.NewPassword,
	}); err != nil {
		return err
	}

	return nil
}
