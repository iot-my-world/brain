package validator

import (
	"github.com/iot-my-world/brain/action"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/pkg/party"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	partyAdministratorException "github.com/iot-my-world/brain/pkg/party/administrator/exception"
	action2 "github.com/iot-my-world/brain/pkg/user/api/action"
	validator2 "github.com/iot-my-world/brain/pkg/user/api/validator"
	humanUserLoginClaims "github.com/iot-my-world/brain/security/claims/login/user/human"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type validator struct {
	partyAdministrator   partyAdministrator.Administrator
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
}

func New(
	partyAdministrator partyAdministrator.Administrator,
) validator2.Validator {

	actionIgnoredReasons := map[action.Action]reasonInvalid.IgnoredReasonsInvalid{
		action2.Create: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
				"username": {
					reasonInvalid.Blank,
				},
				"password": {
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

func (v *validator) ValidateValidateRequest(request *validator2.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else if _, ok := request.Claims.(humanUserLoginClaims.Login); !ok {
		reasonsInvalid = append(reasonsInvalid, "can only validate against user login claims")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *validator2.ValidateRequest) (*validator2.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	apiUserToValidate := &request.User

	if (*apiUserToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*apiUserToValidate).Id,
		})
	}

	if (*apiUserToValidate).Name == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*apiUserToValidate).Name,
		})
	}

	if (*apiUserToValidate).Description == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "description",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*apiUserToValidate).Description,
		})
	}

	if (*apiUserToValidate).Username == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "username",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*apiUserToValidate).Username,
		})
	}

	if len((*apiUserToValidate).Password) == 0 {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "password",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*apiUserToValidate).Password,
		})
	}

	// party type must be set, cannot be blank
	if (*apiUserToValidate).PartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "partyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*apiUserToValidate).PartyType,
		})
	} else {
		// if it is not blank
		// owner party type must be valid. i.e. must be of a valid type and the party must exist
		switch (*apiUserToValidate).PartyType {
		case party.System, party.Client, party.Company:
			_, err := v.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
				Claims:     request.Claims,
				PartyType:  (*apiUserToValidate).PartyType,
				Identifier: (*apiUserToValidate).PartyId,
			})
			if err != nil {
				switch err.(type) {
				case partyAdministratorException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "partyId",
						Type:  reasonInvalid.MustExist,
						Help:  "party must exist",
						Data:  (*apiUserToValidate).PartyId,
					})
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "partyId",
						Type:  reasonInvalid.Unknown,
						Help:  "error retrieving party: " + err.Error(),
						Data:  (*apiUserToValidate).PartyId,
					})
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "partyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*apiUserToValidate).PartyType,
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

	return &validator2.ValidateResponse{
		ReasonsInvalid: returnedReasonsInvalid,
	}, nil
}
