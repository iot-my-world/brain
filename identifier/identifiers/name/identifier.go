package name

import (
	"gitlab.com/iotTracker/brain/identifier"
	"errors"
	"encoding/json"
	"gitlab.com/iotTracker/brain/log"
)

const TYPE = identifier.NAME
type Identifier string

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.IdentifierType { return TYPE }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i == "" {
		return errors.New("id cannot be blank")
	}
	return nil
}

// Struct which shall be used for marshalling and unmarshalling
// Note that this is done to store important type information
// for when the identifier is persisted.
type id struct{
	Id string `json:"id"`
}

func (i Identifier) Marshall() identifier.MarshalledIdentifier {
	// Attempt to marshall this identifer
	data, err := json.Marshal(id{Id: string(i)})
	if err != nil {
		log.Error("Error While marshalling " + TYPE + "identifier: ", err)
	}


}