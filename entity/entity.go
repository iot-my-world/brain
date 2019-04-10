package entity

import (
	"gitlab.com/iotTracker/brain/search/identifier"
)

type Entity interface {
	SetId(id string)
	ValidIdentifier(id identifier.Identifier) bool
	//GetBSON() (interface{}, error)
	//SetBSON(raw bson.Raw) error
}
