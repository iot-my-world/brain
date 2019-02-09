package name

import (
	"errors"
	"gitlab.com/iotTracker/brain/search/identifier"
)

const TYPE = identifier.Name

type Identifier struct {
	Name string `json:"name"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.Type { return TYPE }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.Name == "" {
		return errors.New("name cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter["name"] = i.Name
	return filter
}
