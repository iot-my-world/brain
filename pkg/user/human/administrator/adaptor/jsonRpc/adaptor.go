package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	"github.com/iot-my-world/brain/pkg/user/human"
	"github.com/iot-my-world/brain/pkg/user/human/administrator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	humanUserAdministrator administrator.Administrator
}

func New(
	humanUserAdministrator administrator.Administrator,
) *adaptor {
	return &adaptor{
		humanUserAdministrator: humanUserAdministrator,
	}
}

type GetMyUserRequest struct{}

type GetMyUserResponse struct {
	User human.User `json:"user"`
}

func (a *adaptor) GetMyUser(r *http.Request, request *GetMyUserRequest, response *GetMyUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	getMyUserResponse, err := a.humanUserAdministrator.GetMyUser(&administrator.GetMyUserRequest{
		Claims: claims,
	})
	if err != nil {
		return err
	}

	response.User = getMyUserResponse.User

	return nil
}

type UpdateAllowedFieldsRequest struct {
	User human.User `json:"user"`
}

type UpdateAllowedFieldsResponse struct {
	User human.User `json:"user"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.humanUserAdministrator.UpdateAllowedFields(&administrator.UpdateAllowedFieldsRequest{
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
	User human.User `json:"user"`
}

type CreateResponse struct {
	User human.User `json:"user"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.humanUserAdministrator.Create(&administrator.CreateRequest{
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
	User human.User `json:"user"`
}

func (a *adaptor) UpdatePassword(r *http.Request, request *UpdatePasswordRequest, response *UpdatePasswordResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updatePasswordResponse, err := a.humanUserAdministrator.UpdatePassword(&administrator.UpdatePasswordRequest{
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

	checkPasswordResponse, err := a.humanUserAdministrator.CheckPassword(&administrator.CheckPasswordRequest{
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
	URLToken string `json:"urlToken"`
}

func (a *adaptor) ForgotPassword(r *http.Request, request *ForgotPasswordRequest, response *ForgotPasswordResponse) error {
	forgotPasswordResponse, err := a.humanUserAdministrator.ForgotPassword(&administrator.ForgotPasswordRequest{
		UsernameOrEmailAddress: request.UsernameOrEmailAddress,
	})
	if err != nil {
		return err
	}

	response.URLToken = forgotPasswordResponse.URLToken

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

	if _, err := a.humanUserAdministrator.SetPassword(&administrator.SetPasswordRequest{
		Claims:      claims,
		Identifier:  request.WrappedIdentifier.Identifier,
		NewPassword: request.NewPassword,
	}); err != nil {
		return err
	}

	return nil
}
