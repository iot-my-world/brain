package text

import (
	"github.com/go-errors/errors"
	"github.com/iot-my-world/brain/search/criterion"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

const Type = criterion.Text

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
	return Type
}

func (c Criterion) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter[c.Field] = bson.RegEx{Pattern: ".*" + c.Text + ".*", Options: "i"}
	return filter
}
