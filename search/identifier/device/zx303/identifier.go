package zx303

import (
	"errors"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gopkg.in/mgo.v2/bson"
)

type Identifier struct {
	IMEI string `json:"imei"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.Type { return identifier.DeviceZX303 }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.IMEI == "" {
		return errors.New("IMEI cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() bson.M {
	return bson.M{"imei": i.IMEI}
}
