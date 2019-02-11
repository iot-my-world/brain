package basic

import (
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	partyRegistrar "gitlab.com/iotTracker/brain/party/registrar"
	globalException "gitlab.com/iotTracker/brain/exception"
	registrarException "gitlab.com/iotTracker/brain/party/registrar/exception"
	"gitlab.com/iotTracker/brain/email/mailer"
)

type basicRegistrar struct {
	companyRecordHandler companyRecordHandler.RecordHandler
	userRecordHandler    userRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
	mailer               mailer.Mailer
}

func New(
	companyRecordHandler companyRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
	mailer mailer.Mailer,
) *basicRegistrar {
	return &basicRegistrar{
		companyRecordHandler: companyRecordHandler,
		userRecordHandler:    userRecordHandler,
		clientRecordHandler:  clientRecordHandler,
		mailer:               mailer,
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

	// retrieve the Company whose admin will receive invite
	companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
	if err := br.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
		Identifier: request.PartyIdentifier,
	},
		&companyRetrieveResponse);
		err != nil {
		return registrarException.UnableToRetrieveParty{Reasons: []string{"company party", err.Error()}}
	}

	sendMailResponse := mailer.SendResponse{}
	if err := br.mailer.Send(&mailer.SendRequest{
		//From    string
		To: companyRetrieveResponse.Company.AdminEmailAddress,
		//Cc      string
		Subject: "Welcome to SpotNav",
		Body:    "Welcome to Spot Nav. Click the link to continue.",
		//Bcc     []string
	},
		&sendMailResponse);
		err != nil {
		return err
	}

	return nil
}
