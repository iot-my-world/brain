package validator

import (
	"gitlab.com/iotTracker/brain/action"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	partyAdministrator "gitlab.com/iotTracker/brain/party/administrator"
	partyAdministratorException "gitlab.com/iotTracker/brain/party/administrator/exception"
	"gitlab.com/iotTracker/brain/tracker/device"
	zx303DeviceAction "gitlab.com/iotTracker/brain/tracker/device/zx303/action"
	deviceValidator "gitlab.com/iotTracker/brain/tracker/device/zx303/validator"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
)

type validator struct {
	partyAdministrator   partyAdministrator.Administrator
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
}

func New(
	partyAdministrator partyAdministrator.Administrator,
) deviceValidator.Validator {

	actionIgnoredReasons := map[action.Action]reasonInvalid.IgnoredReasonsInvalid{
		zx303DeviceAction.Create: {
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

func (v *validator) ValidateValidateRequest(request *deviceValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *deviceValidator.ValidateRequest) (*deviceValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	zx303ToValidate := &request.ZX303

	if (*zx303ToValidate).Type != device.ZX303 {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "type",
			Type:  reasonInvalid.Invalid,
			Help:  "must be ZX303",
			Data:  (*zx303ToValidate).Type,
		})
	}

	if (*zx303ToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*zx303ToValidate).Id,
		})
	}

	if (*zx303ToValidate).IMEI == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "imei",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*zx303ToValidate).IMEI,
		})
	}

	if (*zx303ToValidate).SimCountryCode == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "simCountryCode",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*zx303ToValidate).SimCountryCode,
		})
	}

	if (*zx303ToValidate).SimNumber == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "simNumber",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*zx303ToValidate).SimNumber,
		})
	}

	// owner party type must be set, cannot be blank
	if (*zx303ToValidate).OwnerPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "ownerPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*zx303ToValidate).OwnerPartyType,
		})
	} else {
		// if it is not blank
		// owner party type must be valid. i.e. must be of a valid type and the party must exist
		switch (*zx303ToValidate).OwnerPartyType {
		case party.System, party.Client, party.Company:
			_, err := v.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
				Claims:     request.Claims,
				PartyType:  (*zx303ToValidate).OwnerPartyType,
				Identifier: (*zx303ToValidate).OwnerId,
			})
			if err != nil {
				switch err.(type) {
				case partyAdministratorException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "ownerId",
						Type:  reasonInvalid.MustExist,
						Help:  "owner party must exist",
						Data:  (*zx303ToValidate).OwnerId,
					})
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "ownerId",
						Type:  reasonInvalid.Unknown,
						Help:  "error retrieving owner party: " + err.Error(),
						Data:  (*zx303ToValidate).OwnerId,
					})
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "ownerPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*zx303ToValidate).OwnerPartyType,
			})
		}
	}

	// although assigned party type can be blank, if it is then the assigned id must also be blank
	if ((*zx303ToValidate).AssignedPartyType == "" && (*zx303ToValidate).AssignedId.Id != "") ||
		((*zx303ToValidate).AssignedId.Id == "" && (*zx303ToValidate).AssignedPartyType != "") {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "assignedPartyType",
			Type:  reasonInvalid.Invalid,
			Help:  "must both be blank or set",
			Data:  (*zx303ToValidate).AssignedPartyType,
		})
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "assignedId",
			Type:  reasonInvalid.Invalid,
			Help:  "must both be blank or set",
			Data:  (*zx303ToValidate).AssignedId,
		})
	} else if (*zx303ToValidate).AssignedPartyType != "" && (*zx303ToValidate).AssignedId.Id != "" {
		// neither are blank
		switch (*zx303ToValidate).AssignedPartyType {
		case party.System, party.Client, party.Company:
			_, err := v.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
				Claims:     request.Claims,
				PartyType:  (*zx303ToValidate).AssignedPartyType,
				Identifier: (*zx303ToValidate).AssignedId,
			})
			if err != nil {
				switch err.(type) {
				case partyAdministratorException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.MustExist,
						Help:  "assigned party must exist",
						Data:  (*zx303ToValidate).AssignedId,
					})
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.Unknown,
						Help:  "error retrieving assigned party: " + err.Error(),
						Data:  (*zx303ToValidate).AssignedId,
					})
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "assignedPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*zx303ToValidate).AssignedPartyType,
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

	return &deviceValidator.ValidateResponse{
		ReasonsInvalid: returnedReasonsInvalid,
	}, nil
}
