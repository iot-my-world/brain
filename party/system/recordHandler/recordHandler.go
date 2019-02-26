package recordHandler

import (
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gitlab.com/iotTracker/brain/party/system"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
	Validate(request *ValidateRequest, response *ValidateResponse) error
}

const Create api.Method = "Create"
const Retrieve api.Method = "Retrieve"
const Update api.Method = "Update"
const Delete api.Method = "Delete"
const Validate api.Method = "Validate"

type ValidateRequest struct {
	System system.System
	Method api.Method
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

type CreateRequest struct {
	System system.System
}

type CreateResponse struct {
	System system.System
}

type DeleteRequest struct {
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	System system.System
}

type UpdateRequest struct {
	Identifier identifier.Identifier
	System     system.System
}

type UpdateResponse struct {
	System system.System
}

type RetrieveRequest struct {
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	System system.System
}
