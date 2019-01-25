package businessDay

import (
	"bitbucket.org/gotimekeeper/business"
	"bitbucket.org/BACKUP/gotimekeeper/rfId/tagEvent"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	GetCurrent(request *GetCurrentRequest, response *GetCurrentResponse) error
	GetBefore(request *GetBeforeRequest, response *GetBeforeResponse) error
	GetAfter(request *GetAfterRequest, response *GetAfterResponse) error
	GetSelected(request *GetSelectedRequest, response *GetSelectedResponse) error
	UpdateShifts(request *UpdateShiftsRequest, response *UpdateShiftsResponse) error
	AssignEmployeesToShift(request *AssignEmployeesToShiftRequest, response *AssignEmployeesToShiftResponse) error
	EmployeeClock(request *EmployeeClockRequest, response *EmployeeClockResponse) error
}

type GetBeforeRequest struct {
	business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type GetBeforeResponse struct {
	business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type GetAfterRequest struct {
	business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type GetAfterResponse struct {
	business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type GetSelectedRequest struct {
	business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type GetSelectedResponse struct {
	business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type CreateRequest struct {
	business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type CreateResponse struct {
	business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type GetCurrentRequest struct {
	ClientDateTime int64 `json:"clientDateTime"`
}

type GetCurrentResponse struct {
	BusinessDay business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type UpdateShiftsRequest struct {
	BusinessDay business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type UpdateShiftsResponse struct {
	BusinessDay business.BusinessDay  `json:"businessDay" bson:"businessDay"`
}

type AssignEmployeesToShiftRequest struct {
	EmployeeIds []string    `json:"employeeIds"`
	ShiftId     string      `json:"shiftId"`
	BusinessDay business.BusinessDay `json:"businessDay"`
}

type AssignEmployeesToShiftResponse struct {
	Succeeded   []string             `json:"succeeded"`
	Failed      []string             `json:"failed"`
	BusinessDay business.BusinessDay `json:"businessDay"`
}

type EmployeeClockRequest struct {
	TagEvent tagEvent.TagEvent `json:"tagEvent"`
	Employee business.Employee `json:"employee"`
}

type EmployeeClockResponse struct {
	BusinessDay business.BusinessDay `json:"businessDay"`
	ClockEvent  business.ClockEvent  `json:"clockEvent"`
}