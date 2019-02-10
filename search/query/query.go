package query

import (
	"strings"
	"gitlab.com/iotTracker/brain/log"
)

type Query struct {
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Order  []string `json:"order"`
	SortBy []string `json:"sortBy"`
}

func (q Query) ToMongoSortFormat() []string {
	if len(q.Order) != len(q.SortBy) {
		log.Error("query and sortBy are not the same length")
		return make([]string, 0)
	}
	var sortOrder []string
	for i, field := range q.SortBy {
		if len(q.Order) > i && strings.ToLower(q.Order[i]) == "desc" {
			sortOrder = append(sortOrder, "-"+field)
		} else {
			sortOrder = append(sortOrder, field)
		}
	}
	return sortOrder
}
