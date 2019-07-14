package validator

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	"github.com/iot-my-world/brain/pkg/party"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	partyAdministratorException "github.com/iot-my-world/brain/pkg/party/administrator/exception"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	"github.com/iot-my-world/brain/pkg/security/token"
	sigfoxBackendAction "github.com/iot-my-world/brain/pkg/sigfox/backend/action"
	backendRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler"
	backendRecordHandlerException "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler/exception"
	sigfoxBackendValidator "github.com/iot-my-world/brain/pkg/sigfox/backend/validator"
	backendValidatorException "github.com/iot-my-world/brain/pkg/sigfox/backend/validator/exception"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
)

type validator struct {
	jwtValidator         token.JWTValidator
	partyAdministrator   partyAdministrator.Administrator
	backendRecordHandler backendRecordHandler.RecordHandler
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
	systemClaims         *humanUserLoginClaims.Login
}

func New(
	partyAdministrator partyAdministrator.Administrator,
	backendRecordHandler backendRecordHandler.RecordHandler,
	systemClaims *humanUserLoginClaims.Login,
) sigfoxBackendValidator.Validator {

	actionIgnoredReasons := map[action.Action]reasonInvalid.IgnoredReasonsInvalid{
		sigfoxBackendAction.Create: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
				"token": {
					reasonInvalid.Blank,
				},
			},
		},
	}

	return &validator{
		partyAdministrator:   partyAdministrator,
		actionIgnoredReasons: actionIgnoredReasons,
		backendRecordHandler: backendRecordHandler,
		systemClaims:         systemClaims,
	}
}

func (v *validator) ValidateValidateRequest(request *sigfoxBackendValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *sigfoxBackendValidator.ValidateRequest) (*sigfoxBackendValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	backendToValidate := &request.Backend

	if (*backendToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*backendToValidate).Id,
		})
	}

	if (*backendToValidate).OwnerPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "ownerPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*backendToValidate).OwnerPartyType,
		})
	}

	if (*backendToValidate).OwnerId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "ownerId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*backendToValidate).OwnerId,
		})
	}

	// if neither owner party type nor owner id are blank
	if (*backendToValidate).OwnerPartyType != "" && (*backendToValidate).OwnerId.Id != "" {
		// owner party type must be valid. i.e. must be of a valid type and the party must exist
		switch (*backendToValidate).OwnerPartyType {
		case party.System, party.Client, party.Company:
			_, err := v.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
				Claims:     request.Claims,
				PartyType:  (*backendToValidate).OwnerPartyType,
				Identifier: (*backendToValidate).OwnerId,
			})
			if err != nil {
				switch err.(type) {
				case partyAdministratorException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "ownerId",
						Type:  reasonInvalid.MustExist,
						Help:  "owner party must exist",
						Data:  (*backendToValidate).OwnerId,
					})
				default:
					err = backendValidatorException.Validate{Reasons: []string{"retrieving owner party", err.Error()}}
					log.Error(err.Error())
					return nil, err
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "ownerPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*backendToValidate).OwnerPartyType,
			})
		}
	}

	if (*backendToValidate).Name == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*backendToValidate).Name,
		})
	} else {
		// check for duplicate
		_, err := v.backendRecordHandler.Retrieve(&backendRecordHandler.RetrieveRequest{
			Claims: v.systemClaims,
			Identifier: name.Identifier{
				Name: (*backendToValidate).Name,
			},
		})
		switch err.(type) {
		case backendRecordHandlerException.NotFound:
			// this is what we want
		case nil:
			// this means that there is already a backend with this name, i.e. a duplicate
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "name",
				Type:  reasonInvalid.Duplicate,
				Help:  "already exists",
				Data:  (*backendToValidate).Name,
			})
		default:
			err = backendValidatorException.Validate{Reasons: []string{"backend retrieval for duplicate name check", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
	}

	if (*backendToValidate).Token == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "token",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*backendToValidate).Token,
		})
	} else {
		// if token is not blank we check that it is valid
		if (*backendToValidate).OwnerPartyType != "" && (*backendToValidate).OwnerId.Id != "" {
			wrappedJWTClaims, err := v.jwtValidator.ValidateJWT((*backendToValidate).Token)
			if err != nil {
				err = backendValidatorException.Validate{Reasons: []string{"token validation", err.Error()}}
				log.Error(err.Error())
				return nil, err
			}
			unwrappedJWTClaims, err := wrappedJWTClaims.Unwrap()
			if err != nil {
				err = backendValidatorException.Validate{Reasons: []string{"unwrapping claims", err.Error()}}
				log.Error(err.Error())
				return nil, err
			}

			if (*backendToValidate).OwnerPartyType != unwrappedJWTClaims.PartyDetails().PartyType ||
				(*backendToValidate).OwnerId != unwrappedJWTClaims.PartyDetails().PartyId ||
				(*backendToValidate).OwnerPartyType != unwrappedJWTClaims.PartyDetails().ParentPartyType ||
				(*backendToValidate).OwnerId != unwrappedJWTClaims.PartyDetails().ParentId {
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "token",
					Type:  reasonInvalid.Invalid,
					Help:  "party details in claims in token must match that of the backend entity",
					Data:  (*backendToValidate).Token,
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

	return &sigfoxBackendValidator.ValidateResponse{
		ReasonsInvalid: returnedReasonsInvalid,
	}, nil
}
