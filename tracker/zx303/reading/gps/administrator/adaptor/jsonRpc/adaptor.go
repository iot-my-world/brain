package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	zx303DeviceAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/administrator"
	zx303GPSReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps"
	"net/http"
)

type Adaptor struct {
	administrator zx303DeviceAdministrator.Administrator
}

func New(administrator zx303DeviceAdministrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	ZX303 zx303GPSReading.Reading `json:"zx303GPSReading"`
}

type CreateResponse struct {
	ZX303 zx303GPSReading.Reading `json:"zx303GPSReading"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
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
	ZX303 zx303GPSReading.Reading `json:"zx303GPSReading"`
}

type UpdateAllowedFieldsResponse struct {
	ZX303 zx303GPSReading.Reading `json:"zx303GPSReading"`
}

func (a *Adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
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
