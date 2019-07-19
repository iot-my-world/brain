package basic

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	companyAction "github.com/iot-my-world/brain/pkg/party/company/action"
	"github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	companyRecordHandler "github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	"github.com/iot-my-world/brain/pkg/party/company/recordHandler/exception"
	companyRecordHandlerException "github.com/iot-my-world/brain/pkg/party/company/recordHandler/exception"
	companyValidator "github.com/iot-my-world/brain/pkg/party/company/validator"
	companyValidatorException "github.com/iot-my-world/brain/pkg/party/company/validator/exception"
	"github.com/iot-my-world/brain/pkg/search/identifier/adminEmailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	userRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	userRecordHandlerException "github.com/iot-my-world/brain/pkg/user/human/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
)

type validator struct {
	companyRecordHandler recordHandler.RecordHandler
	userRecordHandler    userRecordHandler.RecordHandler
	systemClaims         *humanUserLoginClaims.Login
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
}

// New mongo record handler
func New(
	companyRecordHandler recordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
	systemClaims *humanUserLoginClaims.Login,
) companyValidator.Validator {

	actionIgnoredReasons := map[action.Action]reasonInvalid.IgnoredReasonsInvalid{
		companyAction.Create: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
			},
		},
	}

	return &validator{
		actionIgnoredReasons: actionIgnoredReasons,
		userRecordHandler:    userRecordHandler,
		companyRecordHandler: companyRecordHandler,
		systemClaims:         systemClaims,
	}
}

func (v *validator) ValidateValidateRequest(request *companyValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *companyValidator.ValidateRequest) (*companyValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	companyToValidate := &request.Company

	if (*companyToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*companyToValidate).Id,
		})
	}

	if (*companyToValidate).Name == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*companyToValidate).Name,
		})
	} else {
		// check for duplicate
		_, err := v.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Claims: v.systemClaims,
			Identifier: name.Identifier{
				Name: (*companyToValidate).Name,
			},
		})
		switch err.(type) {
		case companyRecordHandlerException.NotFound:
			// this is what we want
		case nil:
			// this means that there is already a backend with this name, i.e. a duplicate
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "name",
				Type:  reasonInvalid.Duplicate,
				Help:  "already exists",
				Data:  (*companyToValidate).Name,
			})
		default:
			err = companyValidatorException.Validate{Reasons: []string{"company retrieval for duplicate name check", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
	}

	if (*companyToValidate).AdminEmailAddress == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "adminEmailAddress",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*companyToValidate).AdminEmailAddress,
		})
	}

	if (*companyToValidate).ParentPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*companyToValidate).ParentPartyType,
		})
	}

	if (*companyToValidate).ParentId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*companyToValidate).ParentId,
		})
	}

	// Perform additional checks/ignores considering method field
	switch request.Action {
	case companyAction.Create:

		// Check if there is another client that is already using the same admin email address
		if (*companyToValidate).AdminEmailAddress != "" {
			if _, err := v.companyRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
				// system claims as we want to ensure that all companies are visible for this check
				Claims: *v.systemClaims,
				Identifier: adminEmailAddress.Identifier{
					AdminEmailAddress: (*companyToValidate).AdminEmailAddress,
				},
			}); err != nil {
				switch err.(type) {
				case exception.NotFound:
					// this is what we want, do nothing
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "adminEmailAddress",
						Type:  reasonInvalid.Unknown,
						Help:  "unknown error",
						Data:  (*companyToValidate).AdminEmailAddress,
					})
				}
			} else {
				// there was no error, this email is already in database
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "adminEmailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*companyToValidate).AdminEmailAddress,
				})
			}

			// check if there any users with this email address
			if _, err := v.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
				// system claims as we want to ensure that all companies are visible for this check
				Claims: request.Claims,
				Identifier: emailAddress.Identifier{
					EmailAddress: (*companyToValidate).AdminEmailAddress,
				},
			}); err != nil {
				switch err.(type) {
				case userRecordHandlerException.NotFound:
					// this is what we want, do nothing
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "adminEmailAddress",
						Type:  reasonInvalid.Unknown,
						Help:  "unknown error",
						Data:  (*companyToValidate).AdminEmailAddress,
					})
				}
			} else {
				// there was no error, this email is already in database
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "adminEmailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*companyToValidate).AdminEmailAddress,
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

	return &companyValidator.ValidateResponse{
		ReasonsInvalid: returnedReasonsInvalid,
	}, nil
}
