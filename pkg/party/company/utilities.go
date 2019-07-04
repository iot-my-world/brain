package company

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"gopkg.in/mgo.v2/bson"
)

// IsValidIdentifier determines if a given identifier is valid for a company entity
func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.Id, identifier.AdminEmailAddress:
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
	// parties other than system can only see
	return bson.M{"$and": []bson.M{
		filter,
		{"$or": []bson.M{
			// their own company party
			{"id": bson.M{"$eq": claimsToAdd.PartyDetails().PartyId.Id}},
			// OR a company party who is their parent
			{"id": bson.M{"$eq": claimsToAdd.PartyDetails().ParentId.Id}},
		}},
	}}
}
