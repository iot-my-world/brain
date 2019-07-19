package wrapped

import (
	"encoding/json"
	"errors"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	adminEmailAddressIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/adminEmailAddress"
	emailAddressIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/exception"
	idIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/id"
	nameIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/name"
	partyIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/party"
	usernameIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/username"
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
		var unmarshalledId idIdentifier.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.Name:
		var unmarshalledId nameIdentifier.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.Username:
		var unmarshalledId usernameIdentifier.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.EmailAddress:
		var unmarshalledId emailAddressIdentifier.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.AdminEmailAddress:
		var unmarshalledId adminEmailAddressIdentifier.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case identifier.Party:
		var unmarshalledId partyIdentifier.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		i.Identifier = unmarshalledId

	case sigbug.DeviceIdentifier:
		var unmarshalledId sigbug.Identifier
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
