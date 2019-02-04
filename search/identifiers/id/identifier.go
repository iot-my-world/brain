package id

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"errors"
)

const TYPE = identifier.Id
type Identifier string

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.Type { return TYPE }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i == "" {
		return errors.New("id cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter["id"] = i
	return filter
}

//// Struct which shall be used for marshalling and unmarshalling
//// Note that this is done to store important type information
//// for when the identifier is persisted.
//type id struct{
//	Id string `json:"id"`
//}
//
//func (i Identifier) Marshall() identifier.MarshalledIdentifier {
//	// Attempt to marshall this identifer
//	data, err := json.Marshal(id{Id: string(i)})
//	if err != nil {
//		log.Error("Error While marshalling " + string(TYPE) + "identifier: ", err)
//	}
//
//
//}