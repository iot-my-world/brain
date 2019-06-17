package text

import (
	"github.com/go-errors/errors"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Criterion struct {
	Field string `json:"field"`
	Text  string `json:"text"`
}

func (c Criterion) IsValid() error {

	reasonsInvalid := make([]string, 0)

	if c.Text == "" {
		reasonsInvalid = append(reasonsInvalid, "text is blank")
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
	return criterion.ExactText
}

func (c Criterion) ToFilter() map[string]interface{} {
	return bson.M{c.Field: c.Text}
}
