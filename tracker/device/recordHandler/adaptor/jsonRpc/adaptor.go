package jsonRpc

import (
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/tracker/device"
	deviceRecordHandler "gitlab.com/iotTracker/brain/tracker/device/recordHandler"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/search/wrappedCriterion"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/log"
)

type adaptor struct {
	RecordHandler deviceRecordHandler.RecordHandler
}

func New(recordHandler deviceRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type CreateRequest struct {
	Device device.Device `json:"device"`
}

type CreateResponse struct {
	Device device.Device `json:"device"`
}

func (s *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createDeviceResponse := deviceRecordHandler.CreateResponse{}

	if err := s.RecordHandler.Create(
		&deviceRecordHandler.CreateRequest{
			Device: request.Device,
			Claims: claims,
		},
		&createDeviceResponse); err != nil {
		return err
	}

	response.Device = createDeviceResponse.Device

	return nil
}

type RetrieveRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type RetrieveResponse struct {
	Device device.Device `json:"device"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrieveDeviceResponse := deviceRecordHandler.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&deviceRecordHandler.RetrieveRequest{
			Identifier: id,
		},
		&retrieveDeviceResponse); err != nil {
		return err
	}

	response.Device = retrieveDeviceResponse.Device

	return nil
}

type UpdateRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
	Device     device.Device                       `json:"device"`
}

type UpdateResponse struct {
	Device device.Device `json:"device"`
}

func (s *adaptor) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	updateDeviceResponse := deviceRecordHandler.UpdateResponse{}
	if err := s.RecordHandler.Update(
		&deviceRecordHandler.UpdateRequest{
			Identifier: id,
		},
		&updateDeviceResponse); err != nil {
		return err
	}

	response.Device = updateDeviceResponse.Device

	return nil
}

type DeleteRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type DeleteResponse struct {
	Device device.Device `json:"device"`
}

func (s *adaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	deleteDeviceResponse := deviceRecordHandler.DeleteResponse{}
	if err := s.RecordHandler.Delete(
		&deviceRecordHandler.DeleteRequest{
			Identifier: id,
		},
		&deleteDeviceResponse); err != nil {
		return err
	}

	response.Device = deleteDeviceResponse.Device

	return nil
}

type ValidateRequest struct {
	Device device.Device `json:"device"`
	Method api.Method    `json:"method"`
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid `json:"reasonsInvalid"`
}

func (s *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {

	validateDeviceResponse := deviceRecordHandler.ValidateResponse{}
	if err := s.RecordHandler.Validate(
		&deviceRecordHandler.ValidateRequest{
			Device: request.Device,
			Method: request.Method,
		},
		&validateDeviceResponse); err != nil {
		return err
	}

	response.ReasonsInvalid = validateDeviceResponse.ReasonsInvalid

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.WrappedCriterion `json:"criteria"`
	Query    query.Query                         `json:"query"`
}

type CollectResponse struct {
	Records []device.Device `json:"records"`
	Total   int             `json:"total"`
}

func (s *adaptor) Collect(r *http.Request, request *CollectRequest, response *CollectResponse) error {
	// unwrap criteria
	criteria := make([]criterion.Criterion, 0)
	for criterionIdx := range request.Criteria {
		if c, err := request.Criteria[criterionIdx].UnWrap(); err == nil {
			criteria = append(criteria, c)
		} else {
			return err
		}
	}

	collectDeviceResponse := deviceRecordHandler.CollectResponse{}
	if err := s.RecordHandler.Collect(&deviceRecordHandler.CollectRequest{
		Criteria: criteria,
		Query:    request.Query,
	},
		&collectDeviceResponse); err != nil {
		return err
	}

	response.Records = collectDeviceResponse.Records
	response.Total = collectDeviceResponse.Total
	return nil
}
