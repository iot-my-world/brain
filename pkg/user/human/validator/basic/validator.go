package basic

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/action"
	"github.com/iot-my-world/brain/pkg/party/client"
	clientRecordHandler "github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	companyRecordHandler "github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	companyRecordHandlerException "github.com/iot-my-world/brain/pkg/party/company/recordHandler/exception"
	partyRegistrarAction "github.com/iot-my-world/brain/pkg/party/registrar/action"
	"github.com/iot-my-world/brain/pkg/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/username"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	humanUserAction "github.com/iot-my-world/brain/pkg/user/human/action"
	humanUserRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	humanUserRecordHandlerException "github.com/iot-my-world/brain/pkg/user/human/recordHandler/exception"
	humanUserValidator "github.com/iot-my-world/brain/pkg/user/human/validator"
	humanUserValidatorException "github.com/iot-my-world/brain/pkg/user/human/validator/exception"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
)

type validator struct {
	userRecordHandler    humanUserRecordHandler.RecordHandler
	companyRecordHandler companyRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
	systemClaims         *humanUserLoginClaims.Login
	actionIgnoredReasons map[action.Action]reasonInvalid.IgnoredReasonsInvalid
}

func New(
	userRecordHandler humanUserRecordHandler.RecordHandler,
	companyRecordHandler companyRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
	systemClaims *humanUserLoginClaims.Login,
) humanUserValidator.Validator {

	actionIgnoredReasons := map[action.Action]reasonInvalid.IgnoredReasonsInvalid{
		humanUserAction.Create: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
				"name": {
					reasonInvalid.Blank,
				},
				"surname": {
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

		partyRegistrarAction.InviteCompanyAdminUser: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
				"name": {
					reasonInvalid.Blank,
				},
				"surname": {
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

		partyRegistrarAction.RegisterCompanyAdminUser: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{},
		},

		partyRegistrarAction.InviteCompanyUser: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
				"name": {
					reasonInvalid.Blank,
				},
				"surname": {
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

		partyRegistrarAction.RegisterCompanyUser: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{},
		},

		partyRegistrarAction.InviteClientAdminUser: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
				"name": {
					reasonInvalid.Blank,
				},
				"surname": {
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

		partyRegistrarAction.RegisterClientAdminUser: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{},
		},

		partyRegistrarAction.InviteClientUser: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
				"name": {
					reasonInvalid.Blank,
				},
				"surname": {
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

		partyRegistrarAction.RegisterClientUser: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{},
		},
	}

	return &validator{
		userRecordHandler:    userRecordHandler,
		companyRecordHandler: companyRecordHandler,
		clientRecordHandler:  clientRecordHandler,
		systemClaims:         systemClaims,
		actionIgnoredReasons: actionIgnoredReasons,
	}
}

func (v *validator) ValidateValidateRequest(request *humanUserValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *humanUserValidator.ValidateRequest) (*humanUserValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	userToValidate := &request.User

	if (*userToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*userToValidate).Id,
		})
	}

	if (*userToValidate).EmailAddress == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "emailAddress",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).EmailAddress,
		})
	}

	if (*userToValidate).Name == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Name,
		})
	}

	if (*userToValidate).Surname == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "surname",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Name,
		})
	}

	if (*userToValidate).Username == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "username",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Username,
		})
	}

	if len((*userToValidate).Password) == 0 {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "password",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Password,
		})
	}

	if (*userToValidate).ParentPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).ParentPartyType,
		})
	}

	if (*userToValidate).ParentId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).PartyId,
		})
	}

	if (*userToValidate).PartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "partyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).PartyType,
		})
	}

	if (*userToValidate).PartyId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "partyId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).PartyId,
		})
	}

	if (*userToValidate).Roles == nil {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "roles",
			Type:  reasonInvalid.Nil,
			Help:  "cannot be nil",
			Data:  (*userToValidate).Roles,
		})
	}

	// username and email uniqueness checks
	switch request.Action {

	case partyRegistrarAction.RegisterCompanyAdminUser, partyRegistrarAction.RegisterCompanyUser,
		partyRegistrarAction.RegisterClientAdminUser, partyRegistrarAction.RegisterClientUser:
		// when registering a user the username is scrutinised to ensure that it has not yet been used
		// this is done by checking if the user's username has already been assigned to another user
		if (*userToValidate).Username != "" {
			if userRetrieveResponse, err := v.userRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
				// we use system claims to make sure that all users are visible for this check
				Claims: *v.systemClaims,
				Identifier: username.Identifier{
					Username: (*userToValidate).Username,
				},
			}); err != nil {
				switch err.(type) {
				case humanUserRecordHandlerException.NotFound:
					// this is what we want
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "username",
						Type:  reasonInvalid.Unknown,
						Help:  "retrieve failed",
						Data:  (*userToValidate).Username,
					})
				}
			} else {
				// there was no error, confirm that the username belongs to this user being validated
				if (*userToValidate).Id != userRetrieveResponse.User.Id {
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "username",
						Type:  reasonInvalid.Duplicate,
						Help:  "already exists",
						Data:  (*userToValidate).Username,
					})
				}
			}
		}

	case humanUserAction.Create,
		partyRegistrarAction.InviteCompanyAdminUser, partyRegistrarAction.InviteCompanyUser,
		partyRegistrarAction.InviteClientAdminUser, partyRegistrarAction.InviteClientUser:

		// user cannot have any roles yet for creation
		if len((*userToValidate).Roles) != 0 {
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "roles",
				Type:  reasonInvalid.MustNotBeSet,
				Help:  "can't have roles yet",
				Data:  (*userToValidate).Roles,
			})
		}

		// user cannot be set to registered yet for creation
		// user cannot have any roles yet for creation
		if (*userToValidate).Registered {
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "registered",
				Type:  reasonInvalid.MustNotBeSet,
				Help:  "can't be registered yet",
				Data:  (*userToValidate).Registered,
			})
		}

		// optionally, a username can be provided at this point, it can/will be changed later, but if one
		// is provided now, we check to see if it has been used yet
		if (*userToValidate).Username != "" {
			if _, err := v.userRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
				// we use system claims to make sure that all users are visible for this check
				Claims:     *v.systemClaims,
				Identifier: username.Identifier{Username: (*userToValidate).Username},
			}); err != nil {
				switch err.(type) {
				case humanUserRecordHandlerException.NotFound:
					// this is what we want, user not found so username not taken yet
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "username",
						Type:  reasonInvalid.Unknown,
						Help:  "retrieve failed",
						Data:  (*userToValidate).Username,
					})
				}
			} else {
				// err == nil, i.e. a user was retrieved
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "username",
					Type:  reasonInvalid.Duplicate,
					Help:  "already taken",
					Data:  (*userToValidate).Username,
				})
			}
		}

		// check if the email address is already used
		// is provided now, we check to see if it has been used yet
		if (*userToValidate).EmailAddress != "" {
			if _, err := v.userRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
				// we use system claims to make sure that all users are visible for this check
				Claims:     *v.systemClaims,
				Identifier: emailAddress.Identifier{EmailAddress: (*userToValidate).EmailAddress},
			}); err != nil {
				switch err.(type) {
				case humanUserRecordHandlerException.NotFound:
					// this is what we want, user not found so username not taken yet
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "emailAddress",
						Type:  reasonInvalid.Unknown,
						Help:  "retrieve failed",
						Data:  (*userToValidate).EmailAddress,
					})
				}
			} else {
				// err == nil, i.e. a user was retrieved
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "emailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already taken",
					Data:  (*userToValidate).EmailAddress,
				})
			}
		}

	case humanUserAction.UpdateAllowedFields:
		// username update is allowed, this is to confirm that the username has not been used yet
		// or that it has not changed
		if (*userToValidate).Username != "" {
			if userRetrieveResponse, err := v.userRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
				// we use system claims to make sure that all users are visible for this check
				Claims: *v.systemClaims,
				Identifier: username.Identifier{
					Username: (*userToValidate).Username,
				},
			}); err != nil {
				switch err.(type) {
				case humanUserRecordHandlerException.NotFound:
					// this is what we want
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "username",
						Type:  reasonInvalid.Unknown,
						Help:  "retrieve failed",
						Data:  (*userToValidate).Username,
					})
				}
			} else {
				// there was no error, confirm that the username belongs to this user being validated
				if (*userToValidate).Id != userRetrieveResponse.User.Id {
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "username",
						Type:  reasonInvalid.Duplicate,
						Help:  "already exists",
						Data:  (*userToValidate).Username,
					})
				}
			}
		}
	}

	// entity existence check
	switch request.Action {
	case partyRegistrarAction.RegisterCompanyAdminUser, partyRegistrarAction.RegisterCompanyUser:
		// confirm that the company entity exists
		if _, err := v.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: (*userToValidate).PartyId,
		}); err != nil {
			switch err.(type) {
			case companyRecordHandlerException.NotFound:
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "partyId",
					Type:  reasonInvalid.MustExist,
					Help:  "does not exist",
					Data:  (*userToValidate).PartyId,
				})
			default:
				return nil, humanUserValidatorException.Validate{Reasons: []string{"retrieve company error", err.Error()}}
			}
		}

	case partyRegistrarAction.RegisterClientAdminUser:
		// confirm that the client entity exists
		retrieveResponse, err := v.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: (*userToValidate).PartyId,
		})
		switch err.(type) {
		case companyRecordHandlerException.NotFound:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "partyId",
				Type:  reasonInvalid.MustExist,
				Help:  "does not exist",
				Data:  (*userToValidate).PartyId,
			})
		default:
			return nil, humanUserValidatorException.Validate{Reasons: []string{"retrieve company error", err.Error()}}
		case nil:
			if retrieveResponse.Client.Type == client.Individual {
				// if the client type is individual then the name of the client entity
				// must be the same as the admin user
				if retrieveResponse.Client.Name != (*userToValidate).Name {
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "name",
						Type:  reasonInvalid.Invalid,
						Help:  "must be same as client name",
						Data:  (*userToValidate).Name,
					})
				}
			}
		}

	case partyRegistrarAction.RegisterClientUser:
		// confirm that the client entity exists
		if _, err := v.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: (*userToValidate).PartyId,
		}); err != nil {
			switch err.(type) {
			case companyRecordHandlerException.NotFound:
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "partyId",
					Type:  reasonInvalid.MustExist,
					Help:  "does not exist",
					Data:  (*userToValidate).PartyId,
				})
			default:
				return nil, humanUserValidatorException.Validate{Reasons: []string{"retrieve company error", err.Error()}}
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

	return &humanUserValidator.ValidateResponse{ReasonsInvalid: returnedReasonsInvalid}, nil
}
