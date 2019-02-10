package adminEmailAddress

import (
	"errors"
	"gitlab.com/iotTracker/brain/search/identifier"
)

const TYPE = identifier.EmailAddress

type Identifier struct {
	AdminEmailAddress string `json:"adminEmailAddress"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.Type { return TYPE }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.AdminEmailAddress == "" {
		return errors.New("email address cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter["adminEmailAddress"] = i.AdminEmailAddress
	return filter
}
