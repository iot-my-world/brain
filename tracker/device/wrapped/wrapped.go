package wrapped

import (
	"encoding/json"
	brainException "gitlab.com/iotTracker/brain/exception"
	trackerDevice "gitlab.com/iotTracker/brain/tracker/device"
	deviceException "gitlab.com/iotTracker/brain/tracker/device/exception"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
)

type Wrapped struct {
	Type  trackerDevice.Type `json:"type"`
	Value json.RawMessage    `json:"value"`
}

func Wrap(device trackerDevice.Device) (*Wrapped, error) {
	if device == nil {
		return nil, deviceException.Wrapping{Reasons: []string{"device is nil"}}
	}

	value, err := json.Marshal(device)
	if err != nil {
		return nil, deviceException.Wrapping{Reasons: []string{"marshaling", err.Error()}}
	}

	return &Wrapped{
		Type:  device.Type(),
		Value: value,
	}, nil
}

func (d Wrapped) UnWrap() (trackerDevice.Device, error) {
	var result trackerDevice.Device = nil
	switch d.Type {
	case trackerDevice.TK102:
		var unmarshalledDevice tk102.TK102
		if err := json.Unmarshal(d.Value, &unmarshalledDevice); err != nil {
			return nil, deviceException.UnWrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledDevice
	default:
		return nil, deviceException.UnWrapping{Reasons: []string{"invalid type", string(d.Type)}}
	}

	if result == nil {
		return nil, brainException.Unexpected{Reasons: []string{"device still nil"}}
	}

	return result, nil
}
