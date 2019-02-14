package basic

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"gitlab.com/iotTracker/brain/email/mailer"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	partyRegistrar "gitlab.com/iotTracker/brain/party/registrar"
	registrarException "gitlab.com/iotTracker/brain/party/registrar/exception"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyAdminUser"
	"gitlab.com/iotTracker/brain/security/token"
	"time"
)

type basicRegistrar struct {
	companyRecordHandler companyRecordHandler.RecordHandler
	userRecordHandler    userRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
	mailer               mailer.Mailer
	jwtGenerator         token.JWTGenerator
}

func New(
	companyRecordHandler companyRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
	mailer mailer.Mailer,
	rsaPrivateKey *rsa.PrivateKey,
) *basicRegistrar {
	return &basicRegistrar{
		companyRecordHandler: companyRecordHandler,
		userRecordHandler:    userRecordHandler,
		clientRecordHandler:  clientRecordHandler,
		mailer:               mailer,
		jwtGenerator:         token.NewJWTGenerator(rsaPrivateKey),
	}
}

func (br *basicRegistrar) ValidateInviteCompanyAdminUserRequest(request *partyRegistrar.InviteCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
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
		&companyRetrieveResponse); err != nil {
		return registrarException.UnableToRetrieveParty{Reasons: []string{"company party", err.Error()}}
	}

	// Generate the registration token
	registerCompanyAdminUserClaims := registerCompanyAdminUser.RegisterCompanyAdminUser{
		IssueTime:      time.Now().UTC().Unix(),
		ExpirationTime: time.Now().Add(90 * time.Minute).UTC().Unix(),
		PartyType:      party.Company,
		PartyId:        id.Identifier{Id: companyRetrieveResponse.Company.Id},
		EmailAddress:   companyRetrieveResponse.Company.AdminEmailAddress,
	}

	registrationToken, err := br.jwtGenerator.GenerateToken(registerCompanyAdminUserClaims)
	if err != nil {
		//Unexpected Error!
		return errors.New("log In failed")
	}

	sendMailResponse := mailer.SendResponse{}
	if err := br.mailer.Send(&mailer.SendRequest{
		//From    string
		To: companyRetrieveResponse.Company.AdminEmailAddress,
		//Cc      string
		Subject: "Welcome to SpotNav",
		Body:    fmt.Sprintf("Welcome to Spot Nav. Click the link to continue. http://localhost:3000/register?&t=%s", registrationToken),
		//Bcc     []string
	},
		&sendMailResponse); err != nil {
		return err
	}
	//http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU

	return nil
}

func (br *basicRegistrar) ValidateRegisterCompanyAdminUserRequest(request *partyRegistrar.RegisterCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// user party type and id must be as was in claims otherwise someone is
	// trying to abuse the registration token
	if request.User.PartyType != request.Claims.PartyDetails().PartyType {
		reasonsInvalid = append(reasonsInvalid, "user party type incorrect")
	}
	if request.User.PartyId != request.Claims.PartyDetails().PartyId {
		reasonsInvalid = append(reasonsInvalid, "user party id incorrect")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (br *basicRegistrar) RegisterCompanyAdminUser(request *partyRegistrar.RegisterCompanyAdminUserRequest, response *partyRegistrar.RegisterCompanyAdminUserResponse) error {
	if err := br.ValidateRegisterCompanyAdminUserRequest(request); err != nil {
		return err
	}

	return nil
}
