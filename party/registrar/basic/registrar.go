package basic

import (
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	partyRegistrar "gitlab.com/iotTracker/brain/party/registrar"
	globalException "gitlab.com/iotTracker/brain/exception"
	"fmt"
)

type basicRegistrar struct {
	companyRecordHandler companyRecordHandler.RecordHandler
	userRecordHandler    userRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
}

func New(
	companyRecordHandler companyRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
) *basicRegistrar {
	return &basicRegistrar{
		companyRecordHandler: companyRecordHandler,
		userRecordHandler:    userRecordHandler,
		clientRecordHandler:  clientRecordHandler,
	}
}

func (br *basicRegistrar) ValidateInviteCompanyAdminUserRequest(request *partyRegistrar.InviteCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (br *basicRegistrar) InviteCompanyAdminUser(request *partyRegistrar.InviteCompanyAdminUserRequest, response *partyRegistrar.InviteCompanyAdminUserResponse) error {
	if err := br.ValidateInviteCompanyAdminUserRequest(request); err != nil {
		return err
	}

	fmt.Println("we be here!")

	return nil
}
