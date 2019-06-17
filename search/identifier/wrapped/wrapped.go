package wrapped

import (
	"encoding/json"
	"errors"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/identifier/adminEmailAddress"
	"github.com/iot-my-world/brain/search/identifier/device/tk102"
	"github.com/iot-my-world/brain/search/identifier/device/zx303"
	"github.com/iot-my-world/brain/search/identifier/emailAddress"
	identifierException "github.com/iot-my-world/brain/search/identifier/exception"
	"github.com/iot-my-world/brain/search/identifier/id"
	"github.com/iot-my-world/brain/search/identifier/name"
	"github.com/iot-my-world/brain/search/identifier/party"
	"github.com/iot-my-world/brain/search/identifier/username"
)

type Wrapped struct {
	Type       identifier.Type       `json:"type"`
	Value      json.RawMessage       `json:"value"`
	Identifier identifier.Identifier `json:"-"`
}

func Wrap(id identifier.Identifier) (*Wrapped, error) {
	value, err := json.Marshal(id)
	if err != nil {
		return nil, errors.New("wrapping error " + err.Error())
	}

	return &Wrapped{
		Type:  id.Type(),
		Value: value,
	}, nil
}

func (i *Wrapped) UnmarshalJSON(data []byte) error {
	type Alias Wrapped
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch aux.Type {
	case identifier.Id:
		var unmarshalledId id.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.Name:
		var unmarshalledId name.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.Username:
		var unmarshalledId username.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.EmailAddress:
		var unmarshalledId emailAddress.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.AdminEmailAddress:
		var unmarshalledId adminEmailAddress.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.DeviceTK102:
		var unmarshalledId tk102.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.DeviceZX303:
		var unmarshalledId zx303.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.Party:
		var unmarshalledId party.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	default:
		return identifierException.Invalid{Reasons: []string{"invalid type", string(i.Type)}}
	}

	if i.Identifier == nil {
		return brainException.Unexpected{Reasons: []string{"identifier still nil"}}
	}

	//if err := i.Identifier.IsValid(); err != nil {
	//	return identifierException.Invalid{Reasons: []string{err.Error()}}
	//}

	return nil
}
