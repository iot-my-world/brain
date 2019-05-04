package tk102

import (
	"gitlab.com/iotTracker/brain/search/identifier"
)

func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.Id, identifier.DeviceTK102:
		return true
	default:
		return false
	}
}
