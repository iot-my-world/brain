package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/user/api"
	"github.com/iot-my-world/brain/pkg/user/api/administrator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type Adaptor struct {
	administrator administrator.Administrator
}

func New(administrator administrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	User api.User `json:"apiUser"`
}

type CreateResponse struct {
	User     api.User `json:"apiUser"`
	Password string   `json:"password"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&administrator.CreateRequest{
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
	User api.User `json:"apiUser"`
}

type UpdateAllowedFieldsResponse struct {
	User api.User `json:"apiUser"`
}

func (a *Adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.administrator.UpdateAllowedFields(&administrator.UpdateAllowedFieldsRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = updateAllowedFieldsResponse.User

	return nil
}
