package validator

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/action"
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	sigbugAction "github.com/iot-my-world/brain/pkg/device/sigbug/action"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	sigbugRecordHandlerException "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/exception"
	sigbugValidator "github.com/iot-my-world/brain/pkg/device/sigbug/validator"
	sigbugValidatorException "github.com/iot-my-world/brain/pkg/device/sigbug/validator/exception"
	"github.com/iot-my-world/brain/pkg/party"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	partyAdministratorException "github.com/iot-my-world/brain/pkg/party/administrator/exception"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
)

type validator struct {
	sigbugRecordHandler  sigbugRecordHandler.RecordHandler
	partyAdministrator   partyAdministrator.Administrator
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
	systemClaims         *humanUserLoginClaims.Login
}

func New(
	sigbugRecordHandler sigbugRecordHandler.RecordHandler,
	partyAdministrator partyAdministrator.Administrator,
	systemClaims *humanUserLoginClaims.Login,
) sigbugValidator.Validator {

	actionIgnoredReasons := map[action.Action]reasonInvalid.IgnoredReasonsInvalid{
		sigbugAction.Create: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
			},
		},
		sigbugAction.UpdateAllowedFields: {
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
		sigbugRecordHandler:  sigbugRecordHandler,
		systemClaims:         systemClaims,
	}
}

func (v *validator) ValidateValidateRequest(request *sigbugValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *sigbugValidator.ValidateRequest) (*sigbugValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	sigbugToValidate := &request.Sigbug

	if (*sigbugToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*sigbugToValidate).Id,
		})
	}

	if (*sigbugToValidate).DeviceId == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "deviceId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*sigbugToValidate).DeviceId,
		})
	}

	// action specific checks
	switch request.Action {
	case sigbugAction.Create:
		if (*sigbugToValidate).DeviceId != "" {
			// if device id is not blank, confirm that it is not a duplicate
			_, err := v.sigbugRecordHandler.Retrieve(&sigbugRecordHandler.RetrieveRequest{
				Claims: v.systemClaims,
				Identifier: sigbug.Identifier{
					DeviceId: (*sigbugToValidate).DeviceId,
				},
			})
			switch err.(type) {
			case nil:
				// this means that there is a duplicate as a retrieval was possible
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "deviceId",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*sigbugToValidate).DeviceId,
				})

			case sigbugRecordHandlerException.NotFound:
				// this is what we want
			default:
				// something else went wrong with retrieval, this is an error
				return nil, sigbugValidatorException.Validate{Reasons: []string{"error retrieving sigbug for duplicate check", err.Error()}}
			}
		}
	}

	// owner id must be set
	if (*sigbugToValidate).OwnerId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "ownerId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*sigbugToValidate).OwnerId,
		})
	}

	// owner party type must be set
	if (*sigbugToValidate).OwnerPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "ownerPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*sigbugToValidate).OwnerPartyType,
		})
	} else {
		// if it is not blank
		// owner party type must be valid. i.e. must be of a valid type and the party must exist
		switch (*sigbugToValidate).OwnerPartyType {
		case party.System, party.Client, party.Company:
			// try and retrieve the owner party if it is not blank
			if (*sigbugToValidate).OwnerId.Id != "" {
				_, err := v.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
					Claims:     request.Claims,
					PartyType:  (*sigbugToValidate).OwnerPartyType,
					Identifier: (*sigbugToValidate).OwnerId,
				})
				if err != nil {
					switch err.(type) {
					case partyAdministratorException.NotFound:
						allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
							Field: "ownerId",
							Type:  reasonInvalid.MustExist,
							Help:  "owner party must exist",
							Data:  (*sigbugToValidate).OwnerId,
						})
					default:
						allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
							Field: "ownerId",
							Type:  reasonInvalid.Unknown,
							Help:  "error retrieving owner party: " + err.Error(),
							Data:  (*sigbugToValidate).OwnerId,
						})
					}
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "ownerPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*sigbugToValidate).OwnerPartyType,
			})
		}
	}

	// although assigned party type can be blank, if it is then the assigned id must also be blank
	if ((*sigbugToValidate).AssignedPartyType == "" && (*sigbugToValidate).AssignedId.Id != "") ||
		((*sigbugToValidate).AssignedId.Id == "" && (*sigbugToValidate).AssignedPartyType != "") {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "assignedPartyType",
			Type:  reasonInvalid.Invalid,
			Help:  "must both be blank or set",
			Data:  (*sigbugToValidate).AssignedPartyType,
		})
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "assignedId",
			Type:  reasonInvalid.Invalid,
			Help:  "must both be blank or set",
			Data:  (*sigbugToValidate).AssignedId,
		})
	} else if (*sigbugToValidate).AssignedPartyType != "" && (*sigbugToValidate).AssignedId.Id != "" {
		// neither are blank
		switch (*sigbugToValidate).AssignedPartyType {
		case party.System, party.Client, party.Company:
			_, err := v.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
				Claims:     request.Claims,
				PartyType:  (*sigbugToValidate).AssignedPartyType,
				Identifier: (*sigbugToValidate).AssignedId,
			})
			if err != nil {
				switch err.(type) {
				case partyAdministratorException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.MustExist,
						Help:  "assigned party must exist",
						Data:  (*sigbugToValidate).AssignedId,
					})
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.Unknown,
						Help:  "error retrieving assigned party: " + err.Error(),
						Data:  (*sigbugToValidate).AssignedId,
					})
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "assignedPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*sigbugToValidate).AssignedPartyType,
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

	return &sigbugValidator.ValidateResponse{
		ReasonsInvalid: returnedReasonsInvalid,
	}, nil
}
