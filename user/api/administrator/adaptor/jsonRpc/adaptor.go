package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	apiUser "gitlab.com/iotTracker/brain/user/api"
	apiUserAdministrator "gitlab.com/iotTracker/brain/user/api/administrator"
	"net/http"
)

type Adaptor struct {
	administrator apiUserAdministrator.Administrator
}

func New(administrator apiUserAdministrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	User apiUser.User `json:"apiUser"`
}

type CreateResponse struct {
	User     apiUser.User `json:"apiUser"`
	Password string       `json:"password"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&apiUserAdministrator.CreateRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = createResponse.User
	response.Password = createResponse.Password

	return nil
}

type UpdateAllowedFieldsRequest struct {
	User apiUser.User `json:"apiUser"`
}

type UpdateAllowedFieldsResponse struct {
	User apiUser.User `json:"apiUser"`
}

func (a *Adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.administrator.UpdateAllowedFields(&apiUserAdministrator.UpdateAllowedFieldsRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = updateAllowedFieldsResponse.User

	return nil
}
