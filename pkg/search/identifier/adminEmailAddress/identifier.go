package adminEmailAddress

import (
	"errors"
	identifier2 "github.com/iot-my-world/brain/pkg/search/identifier"
	"gopkg.in/mgo.v2/bson"
)

type Identifier struct {
	AdminEmailAddress string `json:"adminEmailAddress"`
}

// Returns IdentifierType of this Identifier
func (i Identifier) Type() identifier2.Type { return identifier2.AdminEmailAddress }

// Determines and returns the validity of this Identifier
func (i Identifier) IsValid() error {
	if i.AdminEmailAddress == "" {
		return errors.New("email address cannot be blank")
	}
	return nil
}

func (i Identifier) ToFilter() bson.M {
	return bson.M{"adminEmailAddress": i.AdminEmailAddress}
}
