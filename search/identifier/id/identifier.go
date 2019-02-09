package id

import (
	"errors"
	"gitlab.com/iotTracker/brain/search/identifier"
)

const TYPE = identifier.Id

type Identifier struct {
	Id string `json:"id"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.Type { return TYPE }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.Id == "" {
		return errors.New("id cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter["id"] = i.Id
	return filter
}
