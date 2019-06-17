package client

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
	case identifier.Id, identifier.EmailAddress, identifier.AdminEmailAddress:
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
				// their own client party
				{"id": bson.M{"$eq": claimsToAdd.PartyDetails().PartyId.Id}},
				// OR a client party who they are the parent of
				{"parentId.id": bson.M{"$eq": claimsToAdd.PartyDetails().PartyId.Id}},
			}},
		}}
	}
}
