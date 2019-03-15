package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	userAdministrator "gitlab.com/iotTracker/brain/party/user/administrator"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
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
	return nil
}
