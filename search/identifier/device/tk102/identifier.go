package tk102

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"errors"
)

type Identifier struct {
	ManufacturerId string `json:"manufacturerId"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.Type { return identifier.DeviceTK102 }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.ManufacturerId == "" {
		return errors.New("ManufacturerId cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter["manufacturerId"] = i.ManufacturerId
	return filter
}
