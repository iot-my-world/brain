package role

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search"
)

func IsValidIdentifier(id search.Identifier) bool {
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