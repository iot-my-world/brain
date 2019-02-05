package username

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"errors"
)

const TYPE = identifier.Username

type Identifier string

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.Type { return TYPE }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i == "" {
		return errors.New("username cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter["username"] = i
	return filter
}
