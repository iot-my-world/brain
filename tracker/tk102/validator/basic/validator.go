package basic

import (
	"github.com/iot-my-world/brain/action"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/pkg/party"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	partyAdministratorException "github.com/iot-my-world/brain/pkg/party/administrator/exception"
	tk102DeviceAction "github.com/iot-my-world/brain/tracker/tk102/action"
	tk102DeviceValidator "github.com/iot-my-world/brain/tracker/tk102/validator"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type validator struct {
	partyAdministrator   partyAdministrator.Administrator
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
}

func New(
	partyAdministrator partyAdministrator.Administrator,
) tk102DeviceValidator.Validator {

	actionIgnoredReasons := map[action.Action]reasonInvalid.IgnoredReasonsInvalid{
		tk102DeviceAction.Create: {
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

func (v *validator) ValidateValidateRequest(request *tk102DeviceValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *tk102DeviceValidator.ValidateRequest) (*tk102DeviceValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	tk102ToValidate := &request.TK102

	if (*tk102ToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*tk102ToValidate).Id,
		})
	}

	if (*tk102ToValidate).ManufacturerId == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "manufacturerId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*tk102ToValidate).ManufacturerId,
		})
	}

	if (*tk102ToValidate).SimCountryCode == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "simCountryCode",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*tk102ToValidate).SimCountryCode,
		})
	}

	if (*tk102ToValidate).SimNumber == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "simNumber",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*tk102ToValidate).SimNumber,
		})
	}

	// owner party type must be set, cannot be blank
	if (*tk102ToValidate).OwnerPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "ownerPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*tk102ToValidate).OwnerPartyType,
		})
	} else {
		// if it is not blank
		// owner party type must be valid. i.e. must be of a valid type and the party must exist
		switch (*tk102ToValidate).OwnerPartyType {
		case party.System, party.Client, party.Company:
			_, err := v.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
				Claims:     request.Claims,
				PartyType:  (*tk102ToValidate).OwnerPartyType,
				Identifier: (*tk102ToValidate).OwnerId,
			})
			if err != nil {
				switch err.(type) {
				case partyAdministratorException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "ownerId",
						Type:  reasonInvalid.MustExist,
						Help:  "owner party must exist",
						Data:  (*tk102ToValidate).OwnerId,
					})
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "ownerId",
						Type:  reasonInvalid.Unknown,
						Help:  "error retrieving owner party: " + err.Error(),
						Data:  (*tk102ToValidate).OwnerId,
					})
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "ownerPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*tk102ToValidate).OwnerPartyType,
			})
		}
	}

	// although assigned party type can be blank, if it is then the assigned id must also be blank
	if ((*tk102ToValidate).AssignedPartyType == "" && (*tk102ToValidate).AssignedId.Id != "") ||
		((*tk102ToValidate).AssignedId.Id == "" && (*tk102ToValidate).AssignedPartyType != "") {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "assignedPartyType",
			Type:  reasonInvalid.Invalid,
			Help:  "must both be blank or set",
			Data:  (*tk102ToValidate).AssignedPartyType,
		})
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "assignedId",
			Type:  reasonInvalid.Invalid,
			Help:  "must both be blank or set",
			Data:  (*tk102ToValidate).AssignedId,
		})
	} else if (*tk102ToValidate).AssignedPartyType != "" && (*tk102ToValidate).AssignedId.Id != "" {
		// neither are blank
		switch (*tk102ToValidate).AssignedPartyType {
		case party.System, party.Client, party.Company:
			_, err := v.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
				Claims:     request.Claims,
				PartyType:  (*tk102ToValidate).AssignedPartyType,
				Identifier: (*tk102ToValidate).AssignedId,
			})
			if err != nil {
				switch err.(type) {
				case partyAdministratorException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.MustExist,
						Help:  "assigned party must exist",
						Data:  (*tk102ToValidate).AssignedId,
					})
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.Unknown,
						Help:  "error retrieving assigned party: " + err.Error(),
						Data:  (*tk102ToValidate).AssignedId,
					})
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "assignedPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*tk102ToValidate).AssignedPartyType,
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

	return &tk102DeviceValidator.ValidateResponse{ReasonsInvalid: returnedReasonsInvalid}, nil
}
