package date

import (
	"gitlab.com/iotTracker/brain/search/criterion"
	"gopkg.in/mgo.v2/bson"
)

const Type = criterion.DateRange

type Criterion struct {
	Field     string     `json:"field"`
	StartDate RangeValue `json:"startDate"`
	EndDate   RangeValue `json:"endDate"`
}

type RangeValue struct {
	Date      int64 `json:"date"`
	Inclusive bool  `json:"inclusive"`
	Ignore    bool  `json:"ignore"`
}

func (c Criterion) IsValid() error {
	return nil
}

func (c Criterion) Type() criterion.Type {
	return Type
}

func (c Criterion) ToFilter() map[string]interface{} {

	filter := make([]bson.M, 0)

	if !c.StartDate.Ignore {
		if c.StartDate.Inclusive {
			filter = append(
				filter,
				bson.M{c.Field: bson.M{"$gte": c.StartDate.Date}},
			)
		} else {
			filter = append(
				filter,
				bson.M{c.Field: bson.M{"$gt": c.StartDate.Date}},
			)
		}
	}

	return filter
}
