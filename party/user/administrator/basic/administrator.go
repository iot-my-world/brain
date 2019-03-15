package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	userAdministrator "gitlab.com/iotTracker/brain/party/user/administrator"
	userAdministratorException "gitlab.com/iotTracker/brain/party/user/administrator/exception"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/claims/login"
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
