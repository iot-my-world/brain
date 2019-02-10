package wrappedIdentifier

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	universalException "gitlab.com/iotTracker/brain/exception"
	searchException "gitlab.com/iotTracker/brain/search/exception"
	"encoding/json"
	"gitlab.com/iotTracker/brain/search/identifier/name"
)

type WrappedIdentifier struct {
	Type  identifier.Type `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (i WrappedIdentifier) UnWrap() (identifier.Identifier, error) {
	var result identifier.Identifier = nil
	switch i.Type {
	case identifier.Id:
		var unmarshalledId name.Identifier
		if err := json.Unmarshal(i.Value, unmarshalledId); err != nil {
			return nil, searchException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId
	case identifier.Name:
		var unmarshalledId name.Identifier
		if err := json.Unmarshal(i.Value, unmarshalledId); err != nil {
			return nil, searchException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId
	case identifier.Username:
		var unmarshalledId name.Identifier
		if err := json.Unmarshal(i.Value, unmarshalledId); err != nil {
			return nil, searchException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId
	case identifier.EmailAddress:
		var unmarshalledId name.Identifier
		if err := json.Unmarshal(i.Value, unmarshalledId); err != nil {
			return nil, searchException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId
	default:
		return nil, searchException.InvalidIdentifier{Reasons: []string{"invalid type"}}
	}

	if result == nil {
		return nil, universalException.Unexpected{Reasons: []string{"identifier still nil"}}
	}

	if err := result.IsValid(); err != nil {
		return nil, searchException.InvalidIdentifier{Reasons: []string{err.Error()}}
	}

	return result, nil
}
