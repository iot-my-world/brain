package wrapped

import (
	"encoding/json"
	"errors"
	brainException "github.com/iot-my-world/brain/internal/exception"
	identifier2 "github.com/iot-my-world/brain/pkg/search/identifier"
	adminEmailAddress2 "github.com/iot-my-world/brain/pkg/search/identifier/adminEmailAddress"
	tk1022 "github.com/iot-my-world/brain/pkg/search/identifier/device/tk102"
	zx3032 "github.com/iot-my-world/brain/pkg/search/identifier/device/zx303"
	emailAddress2 "github.com/iot-my-world/brain/pkg/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/exception"
	id2 "github.com/iot-my-world/brain/pkg/search/identifier/id"
	name2 "github.com/iot-my-world/brain/pkg/search/identifier/name"
	party2 "github.com/iot-my-world/brain/pkg/search/identifier/party"
	username2 "github.com/iot-my-world/brain/pkg/search/identifier/username"
)

type Wrapped struct {
	Type       identifier2.Type       `json:"type"`
	Value      json.RawMessage        `json:"value"`
	Identifier identifier2.Identifier `json:"-"`
}

func Wrap(id identifier2.Identifier) (*Wrapped, error) {
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
	case identifier2.Id:
		var unmarshalledId id2.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier2.Name:
		var unmarshalledId name2.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier2.Username:
		var unmarshalledId username2.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier2.EmailAddress:
		var unmarshalledId emailAddress2.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier2.AdminEmailAddress:
		var unmarshalledId adminEmailAddress2.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier2.DeviceTK102:
		var unmarshalledId tk1022.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier2.DeviceZX303:
		var unmarshalledId zx3032.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier2.Party:
		var unmarshalledId party2.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	default:
		return exception.Invalid{Reasons: []string{"invalid type", string(i.Type)}}
	}

	if i.Identifier == nil {
		return brainException.Unexpected{Reasons: []string{"identifier still nil"}}
	}

	//if err := i.Identifier.IsValid(); err != nil {
	//	return identifierException.Invalid{Reasons: []string{err.Error()}}
	//}

	return nil
}
