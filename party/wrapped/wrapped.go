package wrapped

import (
	"encoding/json"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/party/client"
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/party/system"
	wrappedPartyException "gitlab.com/iotTracker/brain/party/wrapped/exception"
)

type Wrapped struct {
	Type  party.Type      `json:"type"`
	Value json.RawMessage `json:"value"`
}

func WrapParty(partyEntity party.Party) (*Wrapped, error) {
	value, err := json.Marshal(partyEntity)
	if err != nil {
		return nil, wrappedPartyException.Wrapping{Reasons: []string{err.Error()}}
	}

	return &Wrapped{
		Type:  partyEntity.Details().PartyType,
		Value: value,
	}, nil
}

func (p Wrapped) UnWrap() (party.Party, error) {
	var result party.Party = nil
	switch p.Type {
	case party.System:
		var unmarshalledParty system.System
		if err := json.Unmarshal(p.Value, &unmarshalledParty); err != nil {
			return nil, wrappedPartyException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledParty

	case party.Company:
		var unmarshalledParty company.Company
		if err := json.Unmarshal(p.Value, &unmarshalledParty); err != nil {
			return nil, wrappedPartyException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledParty

	case party.Client:
		var unmarshalledParty client.Client
		if err := json.Unmarshal(p.Value, &unmarshalledParty); err != nil {
			return nil, wrappedPartyException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledParty

	default:
		return nil, wrappedPartyException.InvalidPartyType{Reasons: []string{"unwrapping party", string(p.Type)}}
	}

	return result, nil
}
