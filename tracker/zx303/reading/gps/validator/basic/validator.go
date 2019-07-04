package validator

import (
	"github.com/iot-my-world/brain/action"
	brainException "github.com/iot-my-world/brain/exception"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	zx303GPSReadingAction "github.com/iot-my-world/brain/tracker/zx303/reading/gps/action"
	zx303GPSReadingValidator "github.com/iot-my-world/brain/tracker/zx303/reading/gps/validator"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type validator struct {
	partyAdministrator   partyAdministrator.Administrator
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
}

func New(
	partyAdministrator partyAdministrator.Administrator,
) zx303GPSReadingValidator.Validator {

	actionIgnoredReasons := map[action.Action]reasonInvalid.IgnoredReasonsInvalid{
		zx303GPSReadingAction.Create: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
			},
		},
	}

	return &validator{
		partyAdministrator:   partyAdministrator,
		actionIgnoredReasons: actionIgnoredReasons,
	}
}

func (v *validator) ValidateValidateRequest(request *zx303GPSReadingValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *zx303GPSReadingValidator.ValidateRequest) (*zx303GPSReadingValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	zx303GPSReadingToValidate := &request.ZX303GPSReading

	if (*zx303GPSReadingToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*zx303GPSReadingToValidate).Id,
		})
	}

	// Make list of reasons invalid to return
	returnedReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)

	// Add all reasons that cannot be ignored for the given action
	if v.actionIgnoredReasons[request.Action].ReasonsInvalid != nil {
		for _, reason := range allReasonsInvalid {
			if !v.actionIgnoredReasons[request.Action].CanIgnore(reason) {
				returnedReasonsInvalid = append(returnedReasonsInvalid, reason)
			}
		}
	}

	return &zx303GPSReadingValidator.ValidateResponse{
		ReasonsInvalid: returnedReasonsInvalid,
	}, nil
}
