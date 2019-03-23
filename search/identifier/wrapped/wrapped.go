package wrapped

import (
	"encoding/json"
	"errors"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/identifier/adminEmailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/device/tk102"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	identifierException "gitlab.com/iotTracker/brain/search/identifier/exception"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	"gitlab.com/iotTracker/brain/search/identifier/party"
	"gitlab.com/iotTracker/brain/search/identifier/username"
)

type Wrapped struct {
	Type  identifier.Type `json:"type"`
	Value json.RawMessage `json:"value"`
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

func (i Wrapped) UnWrap() (identifier.Identifier, error) {
	var result identifier.Identifier = nil
	switch i.Type {
	case identifier.Id:
		var unmarshalledId id.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return nil, identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId

	case identifier.Name:
		var unmarshalledId name.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return nil, identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId

	case identifier.Username:
		var unmarshalledId username.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return nil, identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId

	case identifier.EmailAddress:
		var unmarshalledId emailAddress.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return nil, identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId

	case identifier.AdminEmailAddress:
		var unmarshalledId adminEmailAddress.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return nil, identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId

	case identifier.DeviceTK102:
		var unmarshalledId tk102.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return nil, identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId

	case identifier.Party:
		var unmarshalledId party.Identifier
		if err := json.Unmarshal(i.Value, &unmarshalledId); err != nil {
			return nil, identifierException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId

	default:
		return nil, identifierException.Invalid{Reasons: []string{"invalid type", string(i.Type)}}
	}

	if result == nil {
		return nil, brainException.Unexpected{Reasons: []string{"identifier still nil"}}
	}

	if err := result.IsValid(); err != nil {
		return nil, identifierException.Invalid{Reasons: []string{err.Error()}}
	}

	return result, nil
}
