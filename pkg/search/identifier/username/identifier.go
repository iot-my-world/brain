package username

import (
	"errors"
	identifier2 "github.com/iot-my-world/brain/pkg/search/identifier"
	"gopkg.in/mgo.v2/bson"
)

type Identifier struct {
	Username string `json:"username"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier2.Type { return identifier2.Username }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.Username == "" {
		return errors.New("username cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() bson.M {
	return bson.M{"username": i.Username}
}
