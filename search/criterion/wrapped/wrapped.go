package wrapped

import (
	"encoding/json"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/search/criterion"
	exactText "github.com/iot-my-world/brain/search/criterion/exact/text"
	criterionException "github.com/iot-my-world/brain/search/criterion/exception"
	listText "github.com/iot-my-world/brain/search/criterion/list/text"
	"github.com/iot-my-world/brain/search/criterion/or"
	"github.com/iot-my-world/brain/search/criterion/text"
)

type Wrapped struct {
	Type  criterion.Type  `json:"type"`
	Value json.RawMessage `json:"value"`
}

type OrWrapped struct {
	Criteria []Wrapped `json:"criteria"`
}

func Wrap(criterion criterion.Criterion) (*Wrapped, error) {
	value, err := json.Marshal(criterion)
	if err != nil {
		return nil, criterionException.Wrapping{Reasons: []string{
			"json marshalling",
			err.Error(),
		}}
	}

	return &Wrapped{
		Type:  criterion.Type(),
		Value: value,
	}, nil
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

	case criterion.ExactText:
		var unmarshalledCriterion exactText.Criterion
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
