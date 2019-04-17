package wrapped

import (
	"encoding/json"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/tracker/device"
	deviceException "gitlab.com/iotTracker/brain/tracker/device/exception"
	"gitlab.com/iotTracker/brain/tracker/device/zx303"
)

type Wrapped struct {
	Type   device.Type     `json:"type"`
	Value  json.RawMessage `json:"value"`
	Device device.Device   `json:"-"`
}

func Wrap(device device.Device) (*Wrapped, error) {
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

func (d *Wrapped) UnmarshalJSON(data []byte) error {
	type Alias Wrapped
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch aux.Type {
	case device.ZX303:
		var unmarshalledDevice zx303.ZX303
		if err := json.Unmarshal(aux.Value, &unmarshalledDevice); err != nil {
			return deviceException.UnWrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		d.Device = unmarshalledDevice
	default:
		return deviceException.UnWrapping{Reasons: []string{"invalid type", string(d.Type)}}
	}

	if d.Device == nil {
		return brainException.Unexpected{Reasons: []string{"device still nil"}}
	}

	return nil
}
