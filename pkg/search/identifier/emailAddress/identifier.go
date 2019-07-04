package emailAddress

import (
	"errors"
	identifier2 "github.com/iot-my-world/brain/pkg/search/identifier"
	"gopkg.in/mgo.v2/bson"
)

type Identifier struct {
	EmailAddress string `json:"emailAddress"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier2.Type { return identifier2.EmailAddress }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.EmailAddress == "" {
		return errors.New("email address cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() bson.M {
	return bson.M{"emailAddress": i.EmailAddress}
}
