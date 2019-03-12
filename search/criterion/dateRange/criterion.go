package dateRange

import (
	"github.com/go-errors/errors"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Criterion struct {
	Field string   `json:"field"`
	List  []string `json:"list"`
}

type DateRange struct {
	Date   int64 `json:"date"`
	Ignore bool  `json:"ignore"`
}

func (c Criterion) IsValid() error {

	reasonsInvalid := make([]string, 0)

	if len(c.List) == 0 {
		reasonsInvalid = append(reasonsInvalid, "list is empty")
	}

	if c.Field == "" {
		reasonsInvalid = append(reasonsInvalid, "field is blank")
	}

	if len(reasonsInvalid) > 0 {
		return errors.New(strings.Join(reasonsInvalid, "; "))
	}

	return nil
}

func (c Criterion) Type() criterion.Type {
	return criterion.DateRange
}

func (c Criterion) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter[c.Field] = bson.M{"$in": c.List}
	return filter
}