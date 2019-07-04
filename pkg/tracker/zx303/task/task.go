package task

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task/exception"
	step2 "github.com/iot-my-world/brain/pkg/tracker/zx303/task/step"
)

type Type string

type Status string

const Pending Status = "Pending"
const Executing Status = "Executing"
const Finished Status = "Finished"
const Failed Status = "Failed"

type Task struct {
	Id       string        `json:"id" bson:"id"`
	DeviceId id.Identifier `json:"deviceId" bson:"deviceId"`
	Type     Type          `json:"type" bson:"type"`
	Status   Status        `json:"status" bson:"status"`
	Steps    []step2.Step  `json:"steps" bson:"steps"`
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

func (t *Task) ExecutingStep() (*step2.Step, int, error) {
	if t.Status == Finished || t.Status == Failed {
		return nil, -1, exception.NextStep{Reasons: []string{"task in Finished or Failed status", string(t.Status)}}
	}

	for stepIdx := range t.Steps {
		switch t.Steps[stepIdx].Status {
		case step2.Pending:
			return nil, -1, exception.ExecutingStep{Reasons: []string{"pending step found before an executing step"}}

		case step2.Executing:
			return &t.Steps[stepIdx], stepIdx, nil

		case step2.Finished:
			// we continue, the next step should be executing and return this function

		case step2.Failed:
			// if the task has a failed step we can't determine what should be the executing step
			return nil, -1, exception.ExecutingStep{Reasons: []string{"task has a failed step"}}

		default:
			return nil, -1, exception.ExecutingStep{Reasons: []string{"task step has invalid status", string(t.Steps[stepIdx].Status)}}
		}
	}

	return nil, -1, exception.ExecutingStep{Reasons: []string{"task complete"}}
}

func (t *Task) PendingStep() (*step2.Step, int, error) {
	if t.Status == Finished || t.Status == Failed {
		return nil, -1, exception.NextStep{Reasons: []string{"task in Finished or Failed status", string(t.Status)}}
	}

	for stepIdx := range t.Steps {
		switch t.Steps[stepIdx].Status {
		case step2.Pending:
			return &t.Steps[stepIdx], stepIdx, nil

		case step2.Executing:
			return nil, -1, exception.NextStep{Reasons: []string{"task has an executing step"}}

		case step2.Finished:
			// we continue, the next step should be pending and return this function

		case step2.Failed:
			// if the task has a failed step we can't determine what should be the pending step
			return nil, -1, exception.NextStep{Reasons: []string{"task has a failed step"}}

		default:
			return nil, -1, exception.NextStep{Reasons: []string{"task step has invalid status", string(t.Steps[stepIdx].Status)}}
		}
	}

	return nil, -1, exception.NextStep{Reasons: []string{"task complete"}}
}
