package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	zx3032 "github.com/iot-my-world/brain/pkg/tracker/zx303"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/action"
	administrator2 "github.com/iot-my-world/brain/pkg/tracker/zx303/administrator"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/administrator/exception"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/validator"
	"time"
)

type administrator struct {
	zx303DeviceValidator validator.Validator
	zx303RecordHandler   *recordHandler.RecordHandler
}

func New(
	zx303DeviceValidator validator.Validator,
	zx303RecordHandler *recordHandler.RecordHandler,
) administrator2.Administrator {
	return &administrator{
		zx303DeviceValidator: zx303DeviceValidator,
		zx303RecordHandler:   zx303RecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		zx303DeviceValidateResponse, err := a.zx303DeviceValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			ZX303:  request.ZX303,
			Action: action.Create,
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

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.zx303RecordHandler.Create(&recordHandler.CreateRequest{
		ZX303: request.ZX303,
	})
	if err != nil {
		return nil, exception.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &administrator2.CreateResponse{
		ZX303: createResponse.ZX303,
	}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *administrator2.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// device must be valid
		validationResponse, err := a.zx303DeviceValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			Action: action.UpdateAllowedFields,
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

func (a *administrator) UpdateAllowedFields(request *administrator2.UpdateAllowedFieldsRequest) (*administrator2.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the device
	deviceRetrieveResponse, err := a.zx303RecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.ZX303.Id},
	})
	if err != nil {
		return nil, exception.DeviceRetrieval{Reasons: []string{err.Error()}}
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
	_, err = a.zx303RecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.ZX303.Id},
		ZX303:      deviceRetrieveResponse.ZX303,
	})
	if err != nil {
		return nil, exception.DeviceUpdate{Reasons: []string{err.Error()}}
	}

	return &administrator2.UpdateAllowedFieldsResponse{
		ZX303: deviceRetrieveResponse.ZX303,
	}, nil
}

func (a *administrator) ValidateHeartbeatRequest(request *administrator2.HeartbeatRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.ZX303Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "ZX303Identifier is nil")
	} else if !zx3032.IsValidIdentifier(request.ZX303Identifier) {
		reasonsInvalid = append(reasonsInvalid, "ZX303Identifier not valid")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Heartbeat(request *administrator2.HeartbeatRequest) (*administrator2.HeartbeatResponse, error) {
	if err := a.ValidateHeartbeatRequest(request); err != nil {
		return nil, err
	}

	// try and retrieve the device
	retrieveResponse, err := a.zx303RecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303Identifier,
	})
	if err != nil {
		err = exception.Heartbeat{Reasons: []string{"device retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// update the device heartbeat
	retrieveResponse.ZX303.LastHeartbeatTimestamp = time.Now().UTC().Unix()
	if _, err := a.zx303RecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303Identifier,
		ZX303:      retrieveResponse.ZX303,
	}); err != nil {
		err = exception.Heartbeat{Reasons: []string{"device update", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.HeartbeatResponse{}, nil
}
