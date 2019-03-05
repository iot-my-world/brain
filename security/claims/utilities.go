package claims

import (
	"gitlab.com/iotTracker/brain/party"
	"gopkg.in/mgo.v2/bson"
)

func ContextualiseFilter(filter bson.M, claimsToAdd Claims) bson.M {

	if claimsToAdd.PartyDetails().PartyType == party.System {
		// if the party is system, then no contextual filter is added
		return filter
	} else {
		// otherwise we return a contextual filter
		return bson.M{"$and": []bson.M{
			filter,
			{"$or": []bson.M{
				{"ownerId.id": bson.M{"$eq": claimsToAdd.PartyDetails().PartyId.Id}},
				{"assignedId.id": bson.M{"$eq": claimsToAdd.PartyDetails().PartyId.Id}},
			}},
		}}
	}
}
