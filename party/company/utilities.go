package company

import (
	"gitlab.com/iotTracker/brain/search/identifier"
)

func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.Id, identifier.EmailAddress:
		return true
	default:
		return false
	}
}
