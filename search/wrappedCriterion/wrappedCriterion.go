package wrappedCriterion

import (
	"encoding/json"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	criterionException "gitlab.com/iotTracker/brain/search/criterion/exception"
	"gitlab.com/iotTracker/brain/search/criterion/text"
)

type WrappedCriterion struct {
	Type  criterion.Type  `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (cw WrappedCriterion) UnWrap() (criterion.Criterion, error) {
	var result criterion.Criterion = nil
	switch cw.Type {
	case criterion.Text:
		var unmarshalledCriterion text.Criterion
		if err := json.Unmarshal(cw.Value, &unmarshalledCriterion); err != nil {
			return nil, criterionException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledCriterion
	default:
		return nil, criterionException.Invalid{Reasons: []string{"invalid type"}}
	}

	if result == nil {
		return nil, brainException.Unexpected{Reasons: []string{"identifier still nil"}}
	}

	return result, nil
}
