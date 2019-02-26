package criterion

import (
	"gitlab.com/iotTracker/brain/security/claims"
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/iotTracker/brain/party"
)

func CriteriaToFilter(criteria []Criterion, claimsToAdd claims.Claims) bson.M {

	// Build filters from criteria
	filter := bson.M{}
	criteriaFilters := make([]bson.M, 0)
	for criterionIdx := range criteria {
		criteriaFilters = append(criteriaFilters, criteria[criterionIdx].ToFilter())
	}

	// if party is not root then add contextualising filters
	if claimsToAdd.PartyDetails().PartyType != party.System {
		criteriaFilters = append(criteriaFilters, bson.M{"$or": []bson.M{
			{"ownerId.id": bson.M{"$eq": claimsToAdd.PartyDetails().PartyId.Id}},
			{"assignedId.id": bson.M{"$eq": claimsToAdd.PartyDetails().PartyId.Id}},
		}})
	}

	// and them together
	if len(criteriaFilters) > 0 {
		filter["$and"] = criteriaFilters
	}

	return filter
}
