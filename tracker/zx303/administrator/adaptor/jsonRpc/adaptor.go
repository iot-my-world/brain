package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/tracker/zx303"
	zx303DeviceAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/administrator"
	"net/http"
)

type adaptor struct {
	administrator zx303DeviceAdministrator.Administrator
}

func New(administrator zx303DeviceAdministrator.Administrator) *adaptor {
	return &adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	ZX303 zx303.ZX303 `json:"zx303"`
}

type CreateResponse struct {
	ZX303 zx303.ZX303 `json:"zx303"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&zx303DeviceAdministrator.CreateRequest{
		Claims: claims,
		ZX303:  request.ZX303,
	})
	if err != nil {
		return err
	}

	response.ZX303 = createResponse.ZX303

	return nil
}

type UpdateAllowedFieldsRequest struct {
	ZX303 zx303.ZX303 `json:"zx303"`
}

type UpdateAllowedFieldsResponse struct {
	ZX303 zx303.ZX303 `json:"zx303"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.administrator.UpdateAllowedFields(&zx303DeviceAdministrator.UpdateAllowedFieldsRequest{
		Claims: claims,
		ZX303:  request.ZX303,
	})
	if err != nil {
		return err
	}

	response.ZX303 = updateAllowedFieldsResponse.ZX303

	return nil
}

type HeartbeatRequest struct {
	WrappedZX303Identifier wrappedIdentifier.Wrapped `json:"zx303Identifier"`
}

type HeartbeatResponse struct {
}

func (a *adaptor) Heartbeat(r *http.Request, request *HeartbeatRequest, response *HeartbeatResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	if _, err := a.administrator.Heartbeat(&zx303DeviceAdministrator.HeartbeatRequest{
		Claims:          claims,
		ZX303Identifier: request.WrappedZX303Identifier.Identifier,
	}); err != nil {
		return err
	}

	return nil
}
