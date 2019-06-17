package tk102

import (
	"github.com/iot-my-world/brain/search/identifier"
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
