package criterion

import (
	"gopkg.in/mgo.v2/bson"
)

func CriteriaToFilter(criteria []Criterion) bson.M {

	// Build filters from criteria
	filter := bson.M{}
	criteriaFilters := make([]bson.M, 0)
	for criterionIdx := range criteria {
		criteriaFilters = append(criteriaFilters, criteria[criterionIdx].ToFilter())
	}

	// and them together
	if len(criteriaFilters) > 0 {
		filter["$and"] = criteriaFilters
	}

	return filter
}
