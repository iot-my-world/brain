package zx303

import (
	"errors"
	identifier2 "github.com/iot-my-world/brain/pkg/search/identifier"
	"gopkg.in/mgo.v2/bson"
)

type Identifier struct {
	IMEI string `json:"imei"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier2.Type { return identifier2.DeviceZX303 }

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
