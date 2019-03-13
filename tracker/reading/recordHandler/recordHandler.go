package recordHandler

import (
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/reading"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
)

// RecordHandler handles the reading records
type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Collect(request *CollectRequest, response *CollectResponse) error
	Update(request *UpdateRequest, respose *UpdateResponse) error
	Validate(request *ValidateRequest, response *ValidateResponse) error
}

// CreateRequest is the RecordHandlers's Create request object
type CreateRequest struct {
	Reading reading.Reading
}

// CreateResponse is the RecordHandlers's Create response object
type CreateResponse struct {
	Reading reading.Reading
}

// CollectRequest is the RecordHandlers's Collect request object
type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

// CollectResponse is the RecordHandlers's Collect response object
type CollectResponse struct {
	Records []reading.Reading
	Total   int
}

// UpdateRequest is the RecordHandlers's Update request object
type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Reading    reading.Reading
}

// UpdateResponse is the RecordHandlers's Update response object
type UpdateResponse struct {
	Reading reading.Reading
}

// ValidateRequest is the RecordHandlers's Validate request object
type ValidateRequest struct {
	Claims  claims.Claims
	Reading reading.Reading
	Method  api.Method
}

// ValidateResponse is the RecordHandlers's Validate response object
type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}
