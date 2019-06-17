package query

import (
	"github.com/iot-my-world/brain/log"
)

type Query struct {
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Order  []SortOrder `json:"order"`
	SortBy []string    `json:"sortBy"`
}

type SortOrder string

const SortOrderAscending SortOrder = "asc"
const SortOrderDescending SortOrder = "desc"

func (q Query) ToMongoSortFormat() []string {
	if len(q.Order) != len(q.SortBy) {
		log.Error("query and sortBy are not the same length")
		return make([]string, 0)
	}
	var sortOrder []string
	for i, field := range q.SortBy {
		if len(q.Order) > i && q.Order[i] == SortOrderDescending {
			sortOrder = append(sortOrder, "-"+field)
		} else {
			sortOrder = append(sortOrder, field)
		}
	}
	return sortOrder
}
