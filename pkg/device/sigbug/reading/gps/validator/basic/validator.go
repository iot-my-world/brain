package validator

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	sigbugGPSReadingAction "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/action"
	sigbugGPSReadingValidator "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/validator"
	sigbugGPSReadingValidatorException "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/validator/exception"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	sigbugRecordHandlerException "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/exception"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
)

type validator struct {
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
	systemClaims         *humanUserLoginClaims.Login
	sigbugRecordHandler  sigbugRecordHandler.RecordHandler
}

func New(
	sigbugRecordHandler sigbugRecordHandler.RecordHandler,
	partyAdministrator partyAdministrator.Administrator,
	systemClaims *humanUserLoginClaims.Login,
) sigbugGPSReadingValidator.Validator {

	actionIgnoredReasons := map[action.Action]reasonInvalid.IgnoredReasonsInvalid{
		sigbugGPSReadingAction.Create: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
			},
		},
	}

	return &validator{
		sigbugRecordHandler:  sigbugRecordHandler,
		actionIgnoredReasons: actionIgnoredReasons,
		systemClaims:         systemClaims,
	}
}

func (v *validator) ValidateValidateRequest(request *sigbugGPSReadingValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *sigbugGPSReadingValidator.ValidateRequest) (*sigbugGPSReadingValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	gpsReadingToValidate := &request.Reading

	if (*gpsReadingToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*gpsReadingToValidate).Id,
		})
	}

	if (*gpsReadingToValidate).DeviceId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "deviceId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*gpsReadingToValidate).DeviceId,
		})
	} else {
		// device must exist
		sigbugRetrieveResponse, err := v.sigbugRecordHandler.Retrieve(&sigbugRecordHandler.RetrieveRequest{
			Claims:     nil,
			Identifier: nil,
		})
		if err != nil {
			switch err.(type) {
			case sigbugRecordHandlerException.NotFound:
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "deviceId",
					Type:  reasonInvalid.MustExist,
					Help:  "associated device must exist",
					Data:  (*gpsReadingToValidate).DeviceId,
				})
			default:
				err = sigbugGPSReadingValidatorException.Validate{
					Reasons: []string{
						"retrieving sigbug",
						err.Error(),
					},
				}
				log.Error(err)
				return nil, err
			}
		} else {
			// device and message must have same owner and assigned party details
			if sigbugRetrieveResponse.Sigbug.OwnerId != (*gpsReadingToValidate).OwnerId {
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "ownerId",
					Type:  reasonInvalid.Invalid,
					Help:  "owner id must be same as associated sigbug",
					Data:  (*gpsReadingToValidate).OwnerId,
				})
			}
			if sigbugRetrieveResponse.Sigbug.OwnerPartyType != (*gpsReadingToValidate).OwnerPartyType {
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "ownerPartyType",
					Type:  reasonInvalid.Invalid,
					Help:  "owner party type must be same as associated sigbug",
					Data:  (*gpsReadingToValidate).OwnerPartyType,
				})
			}
			if sigbugRetrieveResponse.Sigbug.AssignedId != (*gpsReadingToValidate).AssignedId {
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "assignedId",
					Type:  reasonInvalid.Invalid,
					Help:  "assigned id must be same as associated sigbug",
					Data:  (*gpsReadingToValidate).AssignedId,
				})
			}
			if sigbugRetrieveResponse.Sigbug.AssignedPartyType != (*gpsReadingToValidate).AssignedPartyType {
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "assignedPartyType",
					Type:  reasonInvalid.Invalid,
					Help:  "assigned party type must be same as associated sigbug",
					Data:  (*gpsReadingToValidate).AssignedPartyType,
				})
			}
		}
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

	return &sigbugGPSReadingValidator.ValidateResponse{
		ReasonsInvalid: returnedReasonsInvalid,
	}, nil
}
