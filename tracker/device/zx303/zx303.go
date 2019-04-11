package zx303

import (
	"gitlab.com/iotTracker/brain/tracker/device"
)

type ZX303 struct {
	Id   string `json:"id" bson:"id"`
	IMEI string `json:"imei"`
}

func (z ZX303) Type() device.Type {
	return device.ZX303
}
