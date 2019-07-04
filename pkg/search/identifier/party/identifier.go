package party

import (
	"github.com/go-errors/errors"
	"github.com/iot-my-world/brain/pkg/party"
	identifier2 "github.com/iot-my-world/brain/pkg/search/identifier"
	id2 "github.com/iot-my-world/brain/pkg/search/identifier/id"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Identifier struct {
	PartyType         party.Type     `json:"partyType"`
	PartyIdIdentifier id2.Identifier `json:"partyIdIdentifier"`
}

func (i Identifier) Type() identifier2.Type {
	return identifier2.Party
}

func (i Identifier) IsValid() error {
	reasons := make([]string, 0)
	switch i.PartyType {
	case party.System, party.Client, party.Company:
		// do nothing
	default:
		reasons = append(reasons, "invalid party type: "+string(i.PartyType))
	}

	if err := i.PartyIdIdentifier.IsValid(); err != nil {
		reasons = append(reasons, "partyIdIdentifier invalid: "+err.Error())
	}

	if len(reasons) > 0 {
		return errors.New(strings.Join(reasons, "; "))
	}

	return nil
}

func (i Identifier) ToFilter() bson.M {
	return i.PartyIdIdentifier.ToFilter()
}
