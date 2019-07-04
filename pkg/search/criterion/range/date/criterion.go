package date

import (
	criterion2 "github.com/iot-my-world/brain/pkg/search/criterion"
	"gopkg.in/mgo.v2/bson"
)

const Type = criterion2.DateRange

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

func (c Criterion) Type() criterion2.Type {
	return Type
}

func (c Criterion) ToFilter() map[string]interface{} {

	startDateOperator := "$gt"
	if !c.StartDate.Ignore {
		if c.StartDate.Inclusive {
			startDateOperator = "$gte"
		}
	}

	endDateOperator := "$lt"
	if !c.EndDate.Ignore {
		if c.EndDate.Inclusive {
			endDateOperator = "$lte"
		}
	}

	if !c.StartDate.Ignore && c.EndDate.Ignore {
		// only consider start date
		return bson.M{c.Field: bson.M{startDateOperator: c.StartDate.Date}}
	} else if c.StartDate.Ignore && !c.EndDate.Ignore {
		// only consider end date
		return bson.M{c.Field: bson.M{endDateOperator: c.EndDate.Date}}
	} else if !(c.StartDate.Ignore || c.EndDate.Ignore) {
		// consider both start and end dates
		return bson.M{c.Field: bson.M{
			startDateOperator: c.StartDate.Date,
			endDateOperator:   c.EndDate.Date,
		}}
	}

	// consider neither
	return bson.M{c.Field: bson.M{}}
}
