package wrapped

import (
	"encoding/json"
	"github.com/iot-my-world/brain/pkg/party"
	client2 "github.com/iot-my-world/brain/pkg/party/client"
	company2 "github.com/iot-my-world/brain/pkg/party/company"
	system2 "github.com/iot-my-world/brain/pkg/party/system"
	"github.com/iot-my-world/brain/pkg/party/wrapped/exception"
)

type Wrapped struct {
	Type  party.Type      `json:"type"`
	Value json.RawMessage `json:"value"`
}

func WrapParty(partyEntity party.Party) (*Wrapped, error) {
	value, err := json.Marshal(partyEntity)
	if err != nil {
		return nil, exception.Wrapping{Reasons: []string{err.Error()}}
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
		var unmarshalledParty system2.System
		if err := json.Unmarshal(p.Value, &unmarshalledParty); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledParty

	case party.Company:
		var unmarshalledParty company2.Company
		if err := json.Unmarshal(p.Value, &unmarshalledParty); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledParty

	case party.Client:
		var unmarshalledParty client2.Client
		if err := json.Unmarshal(p.Value, &unmarshalledParty); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledParty

	default:
		return nil, exception.InvalidPartyType{Reasons: []string{"unwrapping party", string(p.Type)}}
	}

	return result, nil
}
