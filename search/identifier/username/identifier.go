package username

import (
	"errors"
	"gitlab.com/iotTracker/brain/search/identifier"
)

const TYPE = identifier.Username

type Identifier struct {
	Username string `json:"username"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.Type { return TYPE }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.Username == "" {
		return errors.New("username cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter["username"] = i.Username
	return filter
}
