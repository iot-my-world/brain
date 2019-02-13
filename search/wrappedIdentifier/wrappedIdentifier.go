package wrappedIdentifier

import (
	"encoding/json"
	universalException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/identifier/adminEmailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	identifierException "gitlab.com/iotTracker/brain/search/identifier/exception"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	"gitlab.com/iotTracker/brain/search/identifier/username"
)

type WrappedIdentifier struct {
	Type  identifier.Type `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (i WrappedIdentifier) UnWrap() (identifier.Identifier, error) {
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

	default:
		return nil, identifierException.Invalid{Reasons: []string{"invalid type"}}
	}

	if result == nil {
		return nil, universalException.Unexpected{Reasons: []string{"identifier still nil"}}
	}

	if err := result.IsValid(); err != nil {
		return nil, identifierException.Invalid{Reasons: []string{err.Error()}}
	}

	return result, nil
}
