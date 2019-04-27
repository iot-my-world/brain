package task

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/tracker/zx303/task/step"
)

type Type string

type Status string

const Pending Status = "Pending"
const Executing Status = "Executing"
const Finished Status = "Finished"
const Failed Status = "Failed"

type Task struct {
	Id    string      `json:"id" bson:"id"`
	Type  Type        `json:"type" bson:"type"`
	Steps []step.Step `json:"steps" bson:"steps"`
}

func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.Id, identifier.DeviceZX303:
		return true
	default:
		return false
	}
}

func (t *Task) SetId(id string) {
	t.Id = id
}
