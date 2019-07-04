package or

import (
	"github.com/go-errors/errors"
	criterion2 "github.com/iot-my-world/brain/pkg/search/criterion"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Criterion struct {
	Criteria []criterion2.Criterion `json:"criteria"`
}

func (c Criterion) IsValid() error {

	reasonsInvalid := make([]string, 0)

	if len(c.Criteria) == 0 {
		reasonsInvalid = append(reasonsInvalid, "criteria array has no elements")
	}

	if len(reasonsInvalid) > 0 {
		return errors.New(strings.Join(reasonsInvalid, "; "))
	}

	return nil
}

func (c Criterion) Type() criterion2.Type {
	return criterion2.Or
}

func (c Criterion) ToFilter() map[string]interface{} {
	filter := bson.M{}
	criteriaFilters := make([]bson.M, 0)
	for criterionIdx := range c.Criteria {
		criteriaFilters = append(criteriaFilters, c.Criteria[criterionIdx].ToFilter())
	}
	filter["$or"] = criteriaFilters
	return filter
}
