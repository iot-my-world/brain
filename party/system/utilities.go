package system

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/party"
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

func ContextualiseFilter(filter bson.M, claimsToAdd claims.Claims) bson.M {
	if claimsToAdd.PartyDetails().PartyType == party.System {
		// the system party can see everything
		return filter
	} else {
		// parties other than system can only see
		return bson.M{"$and": []bson.M{
			filter,
			{"$or": []bson.M{
				// their own system party
				{"id": bson.M{"$eq": claimsToAdd.PartyDetails().PartyId.Id}},
				// OR a system party who is their parent
				{"id": bson.M{"$eq": claimsToAdd.PartyDetails().ParentId.Id}},
			}},
		}}
	}
}
