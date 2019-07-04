package text

import (
	"github.com/go-errors/errors"
	criterion2 "github.com/iot-my-world/brain/pkg/search/criterion"
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

func (c Criterion) Type() criterion2.Type {
	return criterion2.ExactText
}

func (c Criterion) ToFilter() map[string]interface{} {
	return bson.M{c.Field: c.Text}
}
