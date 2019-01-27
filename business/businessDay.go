package business

import (
	"gitlab.com/iotTracker/brain/business/shift"
	"gitlab.com/iotTracker/brain/rfId"
)

type BusinessDay struct {
	Id            string        `json:"id" bson:"id"`
	StartDateTime int64         `json:"startDateTime" bson:"startDateTime"`
	EndDateTime   int64         `json:"endDateTime" bson:"endDateTime"`
	Shifts        []shift.Shift `json:"shifts" json:"shifts"`
	ClockedIn     []string      `json:"clockedIn"` // Employee id Strings
	ClockedOut    []string      `json:"clockedOut"`
	ClockHistory  []ClockEvent  `json:"clockHistory"`
}

type ClockEvent struct {
	TagEvent   rfId.TagEvent  `json:"tagEvent" bson:"tagEvent"`
	EmployeeId string         `json:"employeeId" bson:"employeeId"`
	Direction  clockDirection `json:"direction" bson:"direction"`
}

type clockDirection string
// Comment
const CLOCK_IN clockDirection = "IN"
const CLOCK_OUT clockDirection = "OUT"