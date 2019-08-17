package validator

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/action"
	sigbugGPSReadingAction "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/action"
	sigbugGPSReadingValidator "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/validator"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	"github.com/iot-my-world/brain/pkg/party"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	partyAdministratorException "github.com/iot-my-world/brain/pkg/party/administrator/exception"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
)

type validator struct {
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
	partyAdministrator   partyAdministrator.Administrator
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
		partyAdministrator:   partyAdministrator,
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
	}

	// owner id must be set
	if (*gpsReadingToValidate).OwnerId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "ownerId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*gpsReadingToValidate).OwnerId,
		})
	}

	// owner party type must be set
	if (*gpsReadingToValidate).OwnerPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "ownerPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*gpsReadingToValidate).OwnerPartyType,
		})
	} else {
		// if it is not blank
		// owner party type must be valid. i.e. must be of a valid type and the party must exist
		switch (*gpsReadingToValidate).OwnerPartyType {
		case party.System, party.Client, party.Company:
			// try and retrieve the owner party if it is not blank
			if (*gpsReadingToValidate).OwnerId.Id != "" {
				_, err := v.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
					Claims:     request.Claims,
					PartyType:  (*gpsReadingToValidate).OwnerPartyType,
					Identifier: (*gpsReadingToValidate).OwnerId,
				})
				if err != nil {
					switch err.(type) {
					case partyAdministratorException.NotFound:
						allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
							Field: "ownerId",
							Type:  reasonInvalid.MustExist,
							Help:  "owner party must exist",
							Data:  (*gpsReadingToValidate).OwnerId,
						})
					default:
						allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
							Field: "ownerId",
							Type:  reasonInvalid.Unknown,
							Help:  "error retrieving owner party: " + err.Error(),
							Data:  (*gpsReadingToValidate).OwnerId,
						})
					}
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "ownerPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*gpsReadingToValidate).OwnerPartyType,
			})
		}
	}

	// although assigned party type can be blank, if it is then the assigned id must also be blank
	if ((*gpsReadingToValidate).AssignedPartyType == "" && (*gpsReadingToValidate).AssignedId.Id != "") ||
		((*gpsReadingToValidate).AssignedId.Id == "" && (*gpsReadingToValidate).AssignedPartyType != "") {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "assignedPartyType",
			Type:  reasonInvalid.Invalid,
			Help:  "must both be blank or set",
			Data:  (*gpsReadingToValidate).AssignedPartyType,
		})
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "assignedId",
			Type:  reasonInvalid.Invalid,
			Help:  "must both be blank or set",
			Data:  (*gpsReadingToValidate).AssignedId,
		})
	} else if (*gpsReadingToValidate).AssignedPartyType != "" && (*gpsReadingToValidate).AssignedId.Id != "" {
		// neither are blank
		switch (*gpsReadingToValidate).AssignedPartyType {
		case party.System, party.Client, party.Company:
			_, err := v.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
				Claims:     request.Claims,
				PartyType:  (*gpsReadingToValidate).AssignedPartyType,
				Identifier: (*gpsReadingToValidate).AssignedId,
			})
			if err != nil {
				switch err.(type) {
				case partyAdministratorException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.MustExist,
						Help:  "assigned party must exist",
						Data:  (*gpsReadingToValidate).AssignedId,
					})
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.Unknown,
						Help:  "error retrieving assigned party: " + err.Error(),
						Data:  (*gpsReadingToValidate).AssignedId,
					})
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "assignedPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*gpsReadingToValidate).AssignedPartyType,
			})
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
