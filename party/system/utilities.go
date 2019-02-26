package system

import (
	"gitlab.com/iotTracker/brain/search/identifier"
)

func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.Id, identifier.Name:
		return true
	default:
		return false
	}
}
