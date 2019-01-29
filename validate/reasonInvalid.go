package validate

import "gitlab.com/iotTracker/brain/validate/reasonInvalid"

type ReasonInvalid struct {
	Field string             `json:"field"`
	Type  reasonInvalid.Type `json:"type"`
	Help  string             `json:"help"`
	Data  interface{}        `json:"data"`
}