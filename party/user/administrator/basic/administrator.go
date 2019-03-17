package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/party/user"
	userAction "gitlab.com/iotTracker/brain/party/user/action"
	userAdministrator "gitlab.com/iotTracker/brain/party/user/administrator"
	userAdministratorException "gitlab.com/iotTracker/brain/party/user/administrator/exception"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	userValidator "gitlab.com/iotTracker/brain/party/user/validator"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"golang.org/x/crypto/bcrypt"
)

type administrator struct {
	userRecordHandler userRecordHandler.RecordHandler
	userValidator     userValidator.Validator
}

func New(
	userRecordHandler userRecordHandler.RecordHandler,
	userValidator userValidator.Validator,
) userAdministrator.Administrator {
	return &administrator{
		userRecordHandler: userRecordHandler,
		userValidator:     userValidator,
	}
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *userAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// claims must be login claims to be able to get user
		if request.Claims.Type() != claims.Login {
			reasonsInvalid = append(reasonsInvalid, "claims must be of type login")
		}

		// user must be valid
		validationResponse := userValidator.ValidateResponse{}
		if err := a.userValidator.Validate(&userValidator.ValidateRequest{
			Claims: request.Claims,
			User:   request.User,
			Action: userAction.UpdateAllowedFields,
		}, &validationResponse); err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating user: "+err.Error())
		}
		if len(validationResponse.ReasonsInvalid) > 0 {
			for _, reason := range validationResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("user invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *userAdministrator.UpdateAllowedFieldsRequest, response *userAdministrator.UpdateAllowedFieldsResponse) error {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return err
	}
	return nil
}

func (a *administrator) ValidateGetMyUserRequest(request *userAdministrator.GetMyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) GetMyUser(request *userAdministrator.GetMyUserRequest, response *userAdministrator.GetMyUserResponse) error {
	if err := a.ValidateGetMyUserRequest(request); err != nil {
		return err
	}

	// infer the login claims type
	loginClaims, ok := request.Claims.(login.Login)
	if !ok {
		return userAdministratorException.InvalidClaims{Reasons: []string{"cannot assert login claims type"}}
	}

	// retrieve user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	}, &userRetrieveResponse); err != nil {
		return userAdministratorException.UserRetrieval{Reasons: []string{"user retrieval", err.Error()}}
	}

	response.User = userRetrieveResponse.User

	return nil
}

func (a *administrator) ValidateCreateRequest(request *userAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// confirm that the party details of the user being created matches claims
		// i.e users can only be created by their own party unless the system party
		// is acting
		switch request.Claims.PartyDetails().PartyType {
		case party.System:
			// do nothing, we expect system to know what they are doing
		default:
			if request.User.PartyType != request.Claims.PartyDetails().PartyType {
				reasonsInvalid = append(reasonsInvalid, "user PartyType must be the type of the party doing creation")
			}
			if request.User.PartyId != request.Claims.PartyDetails().PartyId {
				reasonsInvalid = append(reasonsInvalid, "client PartyId must be the id of the party doing creation")
			}
			if request.User.ParentPartyType != request.Claims.PartyDetails().ParentPartyType {
				reasonsInvalid = append(reasonsInvalid, "user ParentPartyType must match that of the party doing creation")
			}
			if request.User.ParentId != request.Claims.PartyDetails().ParentId {
				reasonsInvalid = append(reasonsInvalid, "user ParentId must match that of the party doing creation")
			}
		}

		// user must be valid
		validationResponse := userValidator.ValidateResponse{}
		if err := a.userValidator.Validate(&userValidator.ValidateRequest{
			Claims: request.Claims,
			User:   request.User,
			Action: userAction.Create,
		}, &validationResponse); err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating user: "+err.Error())
		}
		if len(validationResponse.ReasonsInvalid) > 0 {
			for _, reason := range validationResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("user invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *userAdministrator.CreateRequest, response *userAdministrator.CreateResponse) error {
	if err := a.ValidateCreateRequest(request); err != nil {
		return err
	}

	// create the user
	createResponse := userRecordHandler.CreateResponse{}
	if err := a.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		User: request.User,
	}, &createResponse); err != nil {
		return userAdministratorException.UserCreation{Reasons: []string{"user creation", err.Error()}}
	}

	response.User = createResponse.User

	return nil
}

func (a *administrator) ValidateChangePasswordRequest(request *userAdministrator.ChangePasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.NewPassword == "" {
		reasonsInvalid = append(reasonsInvalid, "password blank")
	}

	if !user.IsValidIdentifier(request.Identifier) {
		reasonsInvalid = append(reasonsInvalid, "invalid user identifier")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) ChangePassword(request *userAdministrator.ChangePasswordRequest, response *userAdministrator.ChangePasswordResponse) error {
	if err := a.ValidateChangePasswordRequest(request); err != nil {
		return err
	}

	// Retrieve User
	retrieveUserResponse := userRecordHandler.RetrieveResponse{}
	if err := a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveUserResponse); err != nil {
		return userAdministratorException.ChangePassword{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return userAdministratorException.ChangePassword{Reasons: []string{"hashing password", err.Error()}}
	}

	// update user
	retrieveUserResponse.User.Password = pwdHash

	updateUserResponse := userRecordHandler.UpdateResponse{}
	if err := a.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		User:       retrieveUserResponse.User,
	}, &updateUserResponse); err != nil {
		return userAdministratorException.ChangePassword{Reasons: []string{"update user", err.Error()}}
	}

	response.User = retrieveUserResponse.User

	return nil
}
