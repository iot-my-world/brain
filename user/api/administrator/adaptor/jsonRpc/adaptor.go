package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/tracker/device/apiUser"
	apiUserDeviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/apiUser/administrator"
	"net/http"
)

type Adaptor struct {
	administrator apiUserDeviceAdministrator.Administrator
}

func New(administrator apiUserDeviceAdministrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	User apiUser.User `json:"apiUser"`
}

type CreateResponse struct {
	User apiUser.User `json:"apiUser"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&apiUserDeviceAdministrator.CreateRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = createResponse.User

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

	updateAllowedFieldsResponse, err := a.administrator.UpdateAllowedFields(&apiUserDeviceAdministrator.UpdateAllowedFieldsRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = updateAllowedFieldsResponse.User

	return nil
}
