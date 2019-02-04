package search

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"encoding/json"
	"gitlab.com/iotTracker/brain/search/identifiers/name"
	searchException "gitlab.com/iotTracker/brain/search/exception"
	universalException "gitlab.com/iotTracker/brain/exception"
)

type IdentifierWrapper struct {
	Type       identifier.Type `json:"type"`
	Identifier json.RawMessage `json:"identifier"`
}

func (i IdentifierWrapper) UnWrap() (Identifier, error) {
	var result Identifier = nil
	switch i.Type {
	case identifier.Id:
		var unmarshalledId name.Identifier
		if err := json.Unmarshal(i.Identifier, unmarshalledId); err != nil {
			return nil, searchException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId
	case identifier.Name:
		var unmarshalledId name.Identifier
		if err := json.Unmarshal(i.Identifier, unmarshalledId); err != nil {
			return nil, searchException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId
	case identifier.Username:
		var unmarshalledId name.Identifier
		if err := json.Unmarshal(i.Identifier, unmarshalledId); err != nil {
			return nil, searchException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId
	case identifier.EmailAddress:
		var unmarshalledId name.Identifier
		if err := json.Unmarshal(i.Identifier, unmarshalledId); err != nil {
			return nil, searchException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledId
	default:
		return nil, searchException.Unwrapping{Reasons: []string{"invalid type"}}
	}

	if result == nil {
		return nil, universalException.Unexpected{Reasons: []string{"identifier still nil"}}
	}

	if err := result.IsValid(); err != nil {
		return nil, searchException.InvalidIdentifier{Reasons: []string{err.Error()}}
	}

	return result, nil
}

type Identifier interface {
	IsValid() error                   // Returns the validity of the Identifier
	Type() identifier.Type            // Returns the IdentifierType of the Identifier
	ToFilter() map[string]interface{} // Returns a map to use to query the databases
	//Marshal() MarshalledIdentifier // Returns the Identifier in Marshalled form
}
