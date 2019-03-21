package wrapped

import (
	"encoding/json"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	criterionException "gitlab.com/iotTracker/brain/search/criterion/exception"
	listText "gitlab.com/iotTracker/brain/search/criterion/list/text"
	"gitlab.com/iotTracker/brain/search/criterion/or"
	"gitlab.com/iotTracker/brain/search/criterion/text"
)

type Wrapped struct {
	Type  criterion.Type  `json:"type"`
	Value json.RawMessage `json:"value"`
}

type OrWrapped struct {
	Criteria []Wrapped `json:"criteria"`
}

func (cw Wrapped) UnWrap() (criterion.Criterion, error) {
	var result criterion.Criterion = nil
	switch cw.Type {
	case criterion.Text:
		var unmarshalledCriterion text.Criterion
		if err := json.Unmarshal(cw.Value, &unmarshalledCriterion); err != nil {
			return nil, criterionException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledCriterion

	case criterion.ListText:
		var unmarshalledCriterion listText.Criterion
		if err := json.Unmarshal(cw.Value, &unmarshalledCriterion); err != nil {
			return nil, criterionException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledCriterion

	case criterion.Or:
		var wrappedOrCriterion OrWrapped
		var unmarshalledCriterion or.Criterion
		if err := json.Unmarshal(cw.Value, &wrappedOrCriterion); err != nil {
			return nil, criterionException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		unmarshalledCriterion.Criteria = make([]criterion.Criterion, 0)
		for wrappedCritIdx := range wrappedOrCriterion.Criteria {
			if crit, err := wrappedOrCriterion.Criteria[wrappedCritIdx].UnWrap(); err != nil {
				return nil, err
			} else {
				unmarshalledCriterion.Criteria = append(unmarshalledCriterion.Criteria, crit)
			}
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
