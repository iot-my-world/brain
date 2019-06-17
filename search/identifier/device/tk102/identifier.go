package tk102

import (
	"errors"
	"github.com/iot-my-world/brain/search/identifier"
	"gopkg.in/mgo.v2/bson"
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

func (i Identifier) ToFilter() bson.M {
	return bson.M{"manufacturerId": i.ManufacturerId}
}
