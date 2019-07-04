package reading

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
)

// IsValidIdentifier determines if a given identifier is valid for a reading entity
func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.Id:
		return true
	default:
		return false
	}
}
