package basic

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	"github.com/iot-my-world/brain/pkg/party/client"
	clientAction "github.com/iot-my-world/brain/pkg/party/client/action"
	clientRecordHandler "github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	clientRecordHandlerException "github.com/iot-my-world/brain/pkg/party/client/recordHandler/exception"
	clientValidator "github.com/iot-my-world/brain/pkg/party/client/validator"
	clientValidatorException "github.com/iot-my-world/brain/pkg/party/client/validator/exception"
	partyRegistrarAction "github.com/iot-my-world/brain/pkg/party/registrar/action"
	"github.com/iot-my-world/brain/pkg/search/identifier/adminEmailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	humanUserLogin "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	userRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	userRecordHandlerException "github.com/iot-my-world/brain/pkg/user/human/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
)

type validator struct {
	clientRecordHandler  clientRecordHandler.RecordHandler
	userRecordHandler    userRecordHandler.RecordHandler
	systemClaims         *humanUserLogin.Login
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
}

func New(
	clientRecordHandler clientRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
	systemClaims *humanUserLogin.Login,
) clientValidator.Validator {

	actionIgnoredReasons := map[action.Action]reasonInvalid.IgnoredReasonsInvalid{
		clientAction.Create: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
			},
		},
	}

	return &validator{
		actionIgnoredReasons: actionIgnoredReasons,
		clientRecordHandler:  clientRecordHandler,
		userRecordHandler:    userRecordHandler,
		systemClaims:         systemClaims,
	}
}

func (v *validator) ValidateValidateRequest(request *clientValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (v *validator) Validate(request *clientValidator.ValidateRequest) (*clientValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	clientToValidate := &request.Client

	if (*clientToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*clientToValidate).Id,
		})
	}

	if (*clientToValidate).Name == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).Name,
		})
	} else {
		// check for duplicate
		_, err := v.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Claims: v.systemClaims,
			Identifier: name.Identifier{
				Name: (*clientToValidate).Name,
			},
		})
		switch err.(type) {
		case clientRecordHandlerException.NotFound:
			// this is what we want
		case nil:
			// this means that there is already a backend with this name, i.e. a duplicate
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "name",
				Type:  reasonInvalid.Duplicate,
				Help:  "already exists",
				Data:  (*clientToValidate).Name,
			})
		default:
			err = clientValidatorException.Validate{Reasons: []string{"client retrieval for duplicate name check", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
	}

	if (*clientToValidate).Type == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "type",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).Name,
		})
	} else {
		switch (*clientToValidate).Type {
		case client.Company, client.Individual:
		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "type",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid client type",
				Data:  (*clientToValidate).Type,
			})
		}
	}

	if (*clientToValidate).ParentPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).ParentPartyType,
		})
	}

	if (*clientToValidate).ParentId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).ParentId,
		})
	}

	if (*clientToValidate).AdminEmailAddress == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "adminEmailAddress",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).AdminEmailAddress,
		})
	}

	// Perform additional checks/ignores considering method field
	switch request.Action {
	case clientAction.Create:

		if (*clientToValidate).AdminEmailAddress != "" {

			// Check if there is another client that is already using the same admin email address

			if _, err := v.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
				// system claims as we want to ensure that all clients are visible for this check
				Claims: *v.systemClaims,
				Identifier: adminEmailAddress.Identifier{
					AdminEmailAddress: (*clientToValidate).AdminEmailAddress,
				},
			}); err != nil {
				switch err.(type) {
				case clientRecordHandlerException.NotFound:
					// this is what we want, do nothing
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "adminEmailAddress",
						Type:  reasonInvalid.Unknown,
						Help:  "unknown error",
						Data:  (*clientToValidate).AdminEmailAddress,
					})
				}
			} else {
				// there was no error, this email is already in database
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "adminEmailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*clientToValidate).AdminEmailAddress,
				})
			}

			// Check if there is another user that is already using the same admin email address
			if _, err := v.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
				// system claims as we want to ensure that all clients are visible for this check
				Claims: *v.systemClaims,
				Identifier: emailAddress.Identifier{
					EmailAddress: (*clientToValidate).AdminEmailAddress,
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
						Data:  (*clientToValidate).AdminEmailAddress,
					})
				}
			} else {
				// there was no error, this email is already in database
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "adminEmailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*clientToValidate).AdminEmailAddress,
				})
			}
		}

	case partyRegistrarAction.RegisterClientAdminUser:
		// when registering a client admin user the client entity must exist

	case partyRegistrarAction.RegisterCompanyAdminUser:
		// when registering a company admin user the company entity must exist
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

	return &clientValidator.ValidateResponse{
		ReasonsInvalid: returnedReasonsInvalid,
	}, nil
}
