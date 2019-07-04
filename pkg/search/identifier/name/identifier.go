package name

import (
	"errors"
	identifier2 "github.com/iot-my-world/brain/pkg/search/identifier"
	"gopkg.in/mgo.v2/bson"
)

type Identifier struct {
	Name string `json:"name"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier2.Type { return identifier2.Name }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.Name == "" {
		return errors.New("name cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() bson.M {
	return bson.M{"name": i.Name}
}
