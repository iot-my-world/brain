package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/tracker/zx303"
	zx303DeviceAction "gitlab.com/iotTracker/brain/tracker/zx303/action"
	zx303DeviceAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/administrator"
	zx303DeviceAdministratorException "gitlab.com/iotTracker/brain/tracker/zx303/administrator/exception"
	zx303RecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/recordHandler"
	zx303DeviceValidator "gitlab.com/iotTracker/brain/tracker/zx303/validator"
	"time"
)

type administrator struct {
	zx303DeviceValidator zx303DeviceValidator.Validator
	zx303RecordHandler   *zx303RecordHandler.RecordHandler
}

func New(
	zx303DeviceValidator zx303DeviceValidator.Validator,
	zx303RecordHandler *zx303RecordHandler.RecordHandler,
) zx303DeviceAdministrator.Administrator {
	return &administrator{
		zx303DeviceValidator: zx303DeviceValidator,
		zx303RecordHandler:   zx303RecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *zx303DeviceAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		zx303DeviceValidateResponse, err := a.zx303DeviceValidator.Validate(&zx303DeviceValidator.ValidateRequest{
			Claims: request.Claims,
			ZX303:  request.ZX303,
			Action: zx303DeviceAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating zx303 device: "+err.Error())
		}
		if len(zx303DeviceValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range zx303DeviceValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("zx303 device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *zx303DeviceAdministrator.CreateRequest) (*zx303DeviceAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.zx303RecordHandler.Create(&zx303RecordHandler.CreateRequest{
		ZX303: request.ZX303,
	})
	if err != nil {
		return nil, zx303DeviceAdministratorException.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &zx303DeviceAdministrator.CreateResponse{
		ZX303: createResponse.ZX303,
	}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *zx303DeviceAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// device must be valid
		validationResponse, err := a.zx303DeviceValidator.Validate(&zx303DeviceValidator.ValidateRequest{
			Claims: request.Claims,
			Action: zx303DeviceAction.UpdateAllowedFields,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating device: "+err.Error())
		}
		if len(validationResponse.ReasonsInvalid) > 0 {
			for _, reason := range validationResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("zx303 device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *zx303DeviceAdministrator.UpdateAllowedFieldsRequest) (*zx303DeviceAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the device
	deviceRetrieveResponse, err := a.zx303RecordHandler.Retrieve(&zx303RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.ZX303.Id},
	})
	if err != nil {
		return nil, zx303DeviceAdministratorException.DeviceRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the device
	//deviceRetrieveResponse.ZX303.Type = request.ZX303.Type
	//deviceRetrieveResponse.ZX303.Id = request.ZX303.Id
	deviceRetrieveResponse.ZX303.IMEI = request.ZX303.IMEI
	deviceRetrieveResponse.ZX303.SimCountryCode = request.ZX303.SimCountryCode
	deviceRetrieveResponse.ZX303.SimNumber = request.ZX303.SimNumber
	//deviceRetrieveResponse.ZX303.OwnerPartyType = request.ZX303.OwnerPartyType
	//deviceRetrieveResponse.ZX303.OwnerId = request.ZX303.OwnerId
	//deviceRetrieveResponse.ZX303.AssignedPartyType = request.ZX303.AssignedPartyType
	//deviceRetrieveResponse.ZX303.AssignedId = request.ZX303.AssignedId

	// update the device
	_, err = a.zx303RecordHandler.Update(&zx303RecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.ZX303.Id},
		ZX303:      deviceRetrieveResponse.ZX303,
	})
	if err != nil {
		return nil, zx303DeviceAdministratorException.DeviceUpdate{Reasons: []string{err.Error()}}
	}

	return &zx303DeviceAdministrator.UpdateAllowedFieldsResponse{
		ZX303: deviceRetrieveResponse.ZX303,
	}, nil
}

func (a *administrator) ValidateHeartbeatRequest(request *zx303DeviceAdministrator.HeartbeatRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.ZX303Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "ZX303Identifier is nil")
	} else if !zx303.IsValidIdentifier(request.ZX303Identifier) {
		reasonsInvalid = append(reasonsInvalid, "ZX303Identifier not valid")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Heartbeat(request *zx303DeviceAdministrator.HeartbeatRequest) (*zx303DeviceAdministrator.HeartbeatResponse, error) {
	if err := a.ValidateHeartbeatRequest(request); err != nil {
		return nil, err
	}

	// try and retrieve the device
	retrieveResponse, err := a.zx303RecordHandler.Retrieve(&zx303RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303Identifier,
	})
	if err != nil {
		err = zx303DeviceAdministratorException.Heartbeat{Reasons: []string{"device retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// update the device heartbeat
	retrieveResponse.ZX303.LastHeartbeatTimestamp = time.Now().UTC().Unix()
	if _, err := a.zx303RecordHandler.Update(&zx303RecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303Identifier,
		ZX303:      retrieveResponse.ZX303,
	}); err != nil {
		err = zx303DeviceAdministratorException.Heartbeat{Reasons: []string{"device update", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &zx303DeviceAdministrator.HeartbeatResponse{}, nil
}
