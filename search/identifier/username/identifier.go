package username

import (
	"errors"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gopkg.in/mgo.v2/bson"
)

type Identifier struct {
	Username string `json:"username"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier.Type { return identifier.Username }

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
