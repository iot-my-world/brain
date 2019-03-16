package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party/user"
	userAdministrator "gitlab.com/iotTracker/brain/party/user/administrator"
	userAdministratorException "gitlab.com/iotTracker/brain/party/user/administrator/exception"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"golang.org/x/crypto/bcrypt"
)

type basicAdministrator struct {
	userRecordHandler userRecordHandler.RecordHandler
}

func New(
	userRecordHandler userRecordHandler.RecordHandler,
) userAdministrator.Administrator {
	return &basicAdministrator{
		userRecordHandler: userRecordHandler,
	}
}

func (ba *basicAdministrator) ValidateUpdateAllowedFieldsRequest(request *userAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// claims must be login claims to be able to get user
		if request.Claims.Type() != claims.Login {
			reasonsInvalid = append(reasonsInvalid, "claims must be of type login")
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (ba *basicAdministrator) UpdateAllowedFields(request *userAdministrator.UpdateAllowedFieldsRequest, response *userAdministrator.UpdateAllowedFieldsResponse) error {
	if err := ba.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return err
	}
	return nil
}

func (ba *basicAdministrator) ValidateGetMyUserRequest(request *userAdministrator.GetMyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (ba *basicAdministrator) GetMyUser(request *userAdministrator.GetMyUserRequest, response *userAdministrator.GetMyUserResponse) error {
	if err := ba.ValidateGetMyUserRequest(request); err != nil {
		return err
	}

	// infer the login claims type
	loginClaims, ok := request.Claims.(login.Login)
	if !ok {
		return userAdministratorException.InvalidClaims{Reasons: []string{"cannot assert login claims type"}}
	}

	// retrieve user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := ba.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	}, &userRetrieveResponse); err != nil {
		return userAdministratorException.UserRetrieval{Reasons: []string{"user retrieval", err.Error()}}
	}

	response.User = userRetrieveResponse.User

	return nil
}

func (ba *basicAdministrator) ValidateChangePasswordRequest(request *userAdministrator.ChangePasswordRequest) error {
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

func (ba *basicAdministrator) ChangePassword(request *userAdministrator.ChangePasswordRequest, response *userAdministrator.ChangePasswordResponse) error {
	if err := ba.ValidateChangePasswordRequest(request); err != nil {
		return err
	}

	// Retrieve User
	retrieveUserResponse := userRecordHandler.RetrieveResponse{}
	if err := ba.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
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
	if err := ba.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Identifier: request.Identifier,
		User:       retrieveUserResponse.User,
	}, &updateUserResponse); err != nil {
		return userAdministratorException.ChangePassword{Reasons: []string{"update user", err.Error()}}
	}

	response.User = retrieveUserResponse.User

	return nil
}
