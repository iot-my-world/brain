package zx303

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/device"
	"gopkg.in/mgo.v2/bson"
)

func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.Id, identifier.DeviceZX303:
		return true
	default:
		return false
	}
}

func ContextualiseFilter(filter bson.M, claimsToAdd claims.Claims) bson.M {
	contextualFilter := claims.ContextualiseFilter(filter, claimsToAdd)

	return bson.M{
		"$and": []bson.M{
			contextualFilter,
			{"type": bson.M{"$eq": device.ZX303}},
		},
	}

	return claims.ContextualiseFilter(filter, claimsToAdd)
}
