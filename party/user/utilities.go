package user

import (
	"gitlab.com/iotTracker/brain/search/identifier"
)

func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.EmailAddress, identifier.Id, identifier.Username:
		return true
	default:
		return false
	}
}
