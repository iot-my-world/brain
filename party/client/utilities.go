package client

import (
	"gitlab.com/iotTracker/brain/search/identifier"
)

func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.Id, identifier.EmailAddress, identifier.AdminEmailAddress:
		return true
	default:
		return false
	}
}
