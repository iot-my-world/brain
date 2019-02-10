package wrappedCriterion

import (
	"gitlab.com/iotTracker/brain/search/criterion"
	"encoding/json"
)

type WrappedCriterion struct {
	Type  criterion.Type  `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (cw WrappedCriterion) UnWrap() (criterion.Criterion, error) {
	var result criterion.Criterion = nil
	switch cw.Type {
	case criterion.Text:
	default:

	}
}
