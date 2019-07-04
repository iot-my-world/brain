package wrapped

import (
	"encoding/json"
	brainException "github.com/iot-my-world/brain/exception"
	criterion2 "github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/criterion/exact/text"
	"github.com/iot-my-world/brain/pkg/search/criterion/exception"
	text2 "github.com/iot-my-world/brain/pkg/search/criterion/list/text"
	or2 "github.com/iot-my-world/brain/pkg/search/criterion/or"
	text3 "github.com/iot-my-world/brain/pkg/search/criterion/text"
)

type Wrapped struct {
	Type  criterion2.Type `json:"type"`
	Value json.RawMessage `json:"value"`
}

type OrWrapped struct {
	Criteria []Wrapped `json:"criteria"`
}

func Wrap(criterion criterion2.Criterion) (*Wrapped, error) {
	value, err := json.Marshal(criterion)
	if err != nil {
		return nil, exception.Wrapping{Reasons: []string{
			"json marshalling",
			err.Error(),
		}}
	}

	return &Wrapped{
		Type:  criterion.Type(),
		Value: value,
	}, nil
}

func (cw Wrapped) UnWrap() (criterion2.Criterion, error) {
	var result criterion2.Criterion = nil
	switch cw.Type {
	case criterion2.Text:
		var unmarshalledCriterion text3.Criterion
		if err := json.Unmarshal(cw.Value, &unmarshalledCriterion); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledCriterion

	case criterion2.ExactText:
		var unmarshalledCriterion text.Criterion
		if err := json.Unmarshal(cw.Value, &unmarshalledCriterion); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledCriterion

	case criterion2.ListText:
		var unmarshalledCriterion text2.Criterion
		if err := json.Unmarshal(cw.Value, &unmarshalledCriterion); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledCriterion

	case criterion2.Or:
		var wrappedOrCriterion OrWrapped
		var unmarshalledCriterion or2.Criterion
		if err := json.Unmarshal(cw.Value, &wrappedOrCriterion); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		unmarshalledCriterion.Criteria = make([]criterion2.Criterion, 0)
		for wrappedCritIdx := range wrappedOrCriterion.Criteria {
			if crit, err := wrappedOrCriterion.Criteria[wrappedCritIdx].UnWrap(); err != nil {
				return nil, err
			} else {
				unmarshalledCriterion.Criteria = append(unmarshalledCriterion.Criteria, crit)
			}
		}
		result = unmarshalledCriterion

	default:
		return nil, exception.Invalid{Reasons: []string{"invalid type"}}
	}

	if result == nil {
		return nil, brainException.Unexpected{Reasons: []string{"identifier still nil"}}
	}

	return result, nil
}
