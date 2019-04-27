package task

import "gitlab.com/iotTracker/brain/tracker/zx303/task/step"

type Type string

type Status string

const Pending Status = "Pending"
const Executing Status = "Executing"
const Finished Status = "Finished"
const Failed Status = "Failed"

type Task struct {
	Type  Type        `json:"type" bson:"type"`
	Steps []step.Step `json:"steps" bson:"steps"`
}
