package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/user"
	userAdministrator "gitlab.com/iotTracker/brain/party/user/administrator"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
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

	getMyUserResponse := userAdministrator.GetMyUserResponse{}
	if err := a.userAdministrator.GetMyUser(&userAdministrator.GetMyUserRequest{
		Claims: claims,
	}, &getMyUserResponse); err != nil {
		return err
	}

	response.User = getMyUserResponse.User

	return nil
}
