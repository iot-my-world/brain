package sigbug

import (
	"errors"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"gopkg.in/mgo.v2/bson"
)

const DeviceIdentifier identifier.Type = "SigbugDeviceIdentifier"

type Identifier struct {
	DeviceId string `json:"deviceId"`
}

func (i Identifier) IsValid() error {
	if i.DeviceId == "" {
		return errors.New("id cannot be blank")
	}
	return nil
}

func (i Identifier) Type() identifier.Type {
	return DeviceIdentifier
}

func (i Identifier) ToFilter() bson.M {
	return bson.M{"deviceId": i.DeviceId}
}
