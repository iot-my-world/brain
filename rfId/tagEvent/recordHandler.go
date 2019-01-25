package tagEvent

import (
	"bitbucket.org/BACKUP/gotimekeeper/rfId/tagEvent"
	"bitbucket.org/gotimekeeper/business"
	"bitbucket.org/gotimekeeper/exoWSC"
)

type RecordHandler interface {
	EmployeeClock(request *EmployeeClockRequest, response *EmployeeClockResponse) error
	RFIDTagEvent(request *RFIDTagEventRequest, response *RFIDTagEventResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Send(message exoWSC.Message) error
}

type EmployeeClockRequest struct {
	TagEvent tagEvent.TagEvent `json:"tagEvent" bson:"tagEvent"`
}

type EmployeeClockResponse struct {
	BusinessDay business.BusinessDay `json:"businessDay"`
	Employee    business.Employee    `json:"employee"`
	ClockEvent  business.ClockEvent  `json:"clockEvent"`
}

type RFIDTagEventRequest struct {
	TagEvent tagEvent.TagEvent `json:"tagEvent" bson:"tagEvent"`
}

type RFIDTagEventResponse struct {}

type RetrieveRequest struct {
	TagId string `json:"tagId" bson:"tagId"`
}

type RetrieveResponse struct {
	TagEvent tagEvent.TagEvent `json:"tagEvent" bson:"tagEvent"`
}