package individual

import (
	"github.com/iot-my-world/brain/party"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"gopkg.in/mgo.v2/bson"
)

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

// ContextualiseFilter takes a filter and claims and returns a contextualised filter
func ContextualiseFilter(filter bson.M, claimsToAdd claims.Claims) bson.M {
	if claimsToAdd.PartyDetails().PartyType == party.System {
		// the system party can see everything
		return filter
	}
	// parties other than system can only see their own party
	return bson.M{"id": bson.M{"$eq": claimsToAdd.PartyDetails().PartyId.Id}}
}
