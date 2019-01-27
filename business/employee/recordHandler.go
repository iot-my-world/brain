package employee

import (
	"gitlab.com/iotTracker/brain/business"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	RetrieveAll(request *RetrieveAllRequest, response *RetrieveAllResponse) error
	RetrieveByShiftAssignment(request *RetrieveByShiftAssignmentRequest, response *RetrieveByShiftAssignmentResponse) error
	RetrieveByTagID(request *RetrieveByTagIDRequest, response *RetrieveByTagIDResponse) error
}

type CreateRequest struct {
	Employee business.Employee `json:"employee"`
}

type CreateResponse struct {
	Employee business.Employee `json:"employee"`
}

type UpdateRequest struct {
	Employee business.Employee `json:"employee"`
}

type UpdateResponse struct {
	Employee business.Employee `json:"employee"`
}

type RetrieveAllRequest struct {

}

type RetrieveAllResponse struct {
	Records []business.Employee `json:"records"`
}

type RetrieveByShiftAssignmentRequest struct {
	BusinessDay business.BusinessDay `json:"businessDay"`
}

type RetrieveByShiftAssignmentResponse struct {
	ShiftGroups [][]business.Employee `json:"shiftGroups"`
	Unassigned  []business.Employee   `json:"unassigned"`
}

type RetrieveByTagIDRequest struct {
	TagID string `json:"tagID"`
}

type RetrieveByTagIDResponse struct {
	Employee business.Employee `json:"employee"`
}