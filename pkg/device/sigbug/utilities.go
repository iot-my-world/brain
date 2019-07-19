package sigbug

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
)

func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}
	switch id.Type() {
	case identifier.Id, DeviceIdentifier:
		return true
	default:
		return false
	}
}
