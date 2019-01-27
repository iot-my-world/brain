package shift

import "gitlab.com/iotTracker/brain/business/businessRole"

type Shift struct {
	Id                 string               `json:"id" bson:"id"`
	StartDateTime      int64                `json:"startDateTime" bson:"startDateTime"`
	EndDateTime        int64                `json:"endDateTime" bson:"endDateTime"`
	BusinessRoleConfig []BusinessRoleConfig `json:"businessRoleConfig" bson:"businessRoleConfig"`
	Register           []RegEntry           `json:"register" bson:"register"`
}

type RegEntry struct {
	EmployeeId    string                    `json:"employeeId" bson:"employeeId"`
	BusinessRole  businessRole.BusinessRole `json:"businessRole" bson:"businessRole"`
	AssignedHours int                       `json:"assignedHours" bson:"assignedHours"`
}

type BusinessRoleConfig struct {
	BusinessRole businessRole.BusinessRole `json:"businessRole" bson:"businessRole"`
	NoRequired   int                       `json:"noRequired" bson:"noRequired"`
}