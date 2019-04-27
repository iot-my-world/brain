package step

type Type string

type Status string

const Pending Status = "Pending"
const Executing Status = "Executing"
const Finished Status = "Finished"
const Failed Status = "Failed"

type Step struct {
	Type   Type   `json:"type" bson:"type"`
	Status Status `json:"status" bson:"status"`
}
