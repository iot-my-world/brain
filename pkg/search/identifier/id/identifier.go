package id

import (
	"errors"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"gopkg.in/mgo.v2/bson"
)

type Identifier struct {
	Id string `json:"id"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.Type { return identifier.Id }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.Id == "" {
		return errors.New("id cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() bson.M {
	return bson.M{"id": i.Id}
}
