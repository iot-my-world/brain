package basic

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"gitlab.com/iotTracker/brain/email/mailer"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	clientRecordHandlerException "gitlab.com/iotTracker/brain/party/client/recordHandler/exception"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	companyRecordHandlerException "gitlab.com/iotTracker/brain/party/company/recordHandler/exception"
	partyRegistrar "gitlab.com/iotTracker/brain/party/registrar"
	registrarException "gitlab.com/iotTracker/brain/party/registrar/exception"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	userRecordHandlerException "gitlab.com/iotTracker/brain/party/user/recordHandler/exception"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/username"
	"gitlab.com/iotTracker/brain/security/claims/registerClientAdminUser"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyAdminUser"
	roleSetup "gitlab.com/iotTracker/brain/security/role/setup"
	"gitlab.com/iotTracker/brain/security/token"
	"time"
	"strings"
	"gitlab.com/iotTracker/brain/search/identifier/adminEmailAddress"
	"gitlab.com/iotTracker/brain/security/claims/login"
)

type basicRegistrar struct {
	companyRecordHandler companyRecordHandler.RecordHandler
	userRecordHandler    userRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
	mailer               mailer.Mailer
	jwtGenerator         token.JWTGenerator
	mailRedirectBaseUrl  string
	systemClaims         login.Login
}

func New(
	companyRecordHandler companyRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
	mailer mailer.Mailer,
	rsaPrivateKey *rsa.PrivateKey,
	mailRedirectBaseUrl string,
	systemClaims login.Login,
) *basicRegistrar {
	return &basicRegistrar{
		companyRecordHandler: companyRecordHandler,
		userRecordHandler:    userRecordHandler,
		clientRecordHandler:  clientRecordHandler,
		mailer:               mailer,
		jwtGenerator:         token.NewJWTGenerator(rsaPrivateKey),
		mailRedirectBaseUrl:  mailRedirectBaseUrl,
		systemClaims:         systemClaims,
	}
}

func (br *basicRegistrar) RegisterSystemAdminUser(request *partyRegistrar.RegisterSystemAdminUserRequest, response *partyRegistrar.RegisterSystemAdminUserResponse) error {

	// check if the system admin user already exists (i.e. has already been registered)
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	err := br.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: username.Identifier{Username: request.User.Username},
	},
		&userRetrieveResponse)
	switch err.(type) {
	case nil:
		// this means that the user already exists
		response.User = userRetrieveResponse.User
		return registrarException.AlreadyRegistered{}
	case userRecordHandlerException.NotFound:
		// this is fine, we will be creating the user now
	default:
		return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
	}

	// create the user
	userCreateResponse := userRecordHandler.CreateResponse{}
	if err := br.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		Claims: request.Claims,
		User:   request.User,
	},
		&userCreateResponse); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userRecordHandler.ChangePasswordResponse{}
	if err := br.userRecordHandler.ChangePassword(&userRecordHandler.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: userCreateResponse.User.Id},
		NewPassword: request.Password,
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	response.User = userCreateResponse.User

	return nil
}

func (br *basicRegistrar) ValidateInviteCompanyAdminUserRequest(request *partyRegistrar.InviteCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

		if request.Claims.PartyDetails().PartyType != party.System {
			// if the user performing the invite is not root then the user's party must be the new users assigned parent party
			if request.User.ParentId.Id != request.Claims.PartyDetails().PartyId.Id {
				reasonsInvalid = append(reasonsInvalid, "parentId must be submitting party's id")
			}

			if request.User.ParentPartyType != request.Claims.PartyDetails().PartyType {
				reasonsInvalid = append(reasonsInvalid, "parentPartyType must be submitting party's type")
			}
		}

		// regardless of who is performing the invite the partyType of the user must be company
		if request.User.PartyType != party.Company {
			reasonsInvalid = append(reasonsInvalid, "user's partyType must be company")
		}

		// validate the new user for the invite admin user method
		userValidateResponse := userRecordHandler.ValidateResponse{}
		err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
			// system claims since we want all users to be visible for the email address check done in validate user
			Claims: br.systemClaims,
			User:   request.User,
			Method: partyRegistrar.InviteAdminUser,
		}, &userValidateResponse)
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "unable to validate newAdminUser")
		} else {
			for _, reason := range userValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
			}
		}

		if request.User.EmailAddress != "" {

			// Check that the admin users email address is correct
			// [1] try and retrieve the company party with the email address
			companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
			if err := br.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
				// system claims since we want all companies to be visible for this retrieval check
				Claims: br.systemClaims,
				Identifier: adminEmailAddress.Identifier{
					AdminEmailAddress: request.User.EmailAddress,
				},
			}, &companyRetrieveResponse); err != nil {
				switch err.(type) {
				case companyRecordHandlerException.NotFound:
					// [2] if no company entity is found this is an issue
					reasonsInvalid = append(reasonsInvalid, "company entity could not be retrieved by the given admin users email address")
				default:
					reasonsInvalid = append(reasonsInvalid, "unable to perform company retrieve to confirm correct email address: "+err.Error())
				}
			} else {
				// [3] if a company was found the id of the company must be the users partyId
				if companyRetrieveResponse.Company.Id != request.User.PartyId.Id {
					// if the id of this other company entity is not the same as the party id of this user then
					reasonsInvalid = append(reasonsInvalid, "emailAddress used as admin email address on another company entity")
				}
			}

			// Check if the users email has already been assigned to another client entity as admin email
			if request.User.EmailAddress != "" {
				if err := br.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
					Claims: br.systemClaims,
					Identifier: adminEmailAddress.Identifier{
						AdminEmailAddress: request.User.EmailAddress,
					},
				},
					&clientRecordHandler.RetrieveResponse{}); err != nil {
					switch err.(type) {
					case clientRecordHandlerException.NotFound:
						// this is what we want, do nothing
					default:
						reasonsInvalid = append(reasonsInvalid, "unable to confirm admin user email address uniqueness")
					}
				} else {
					// there was no error, this email address is already taken by some client entity
					reasonsInvalid = append(reasonsInvalid, "emailAddress used as admin email address on a client entity")
				}
			}
		}
	}

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

	// Generate the registration token for the company admin user to register
	registerCompanyAdminUserClaims := registerCompanyAdminUser.RegisterCompanyAdminUser{
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: request.User.ParentPartyType,
		ParentId:        request.User.ParentId,
		PartyType:       request.User.PartyType,
		PartyId:         request.User.PartyId,
		User:            request.User,
	}
	registrationToken, err := br.jwtGenerator.GenerateToken(registerCompanyAdminUserClaims)
	if err != nil {
		return registrarException.TokenGeneration{Reasons: []string{"inviteCompanyAdminUser", err.Error()}}
	}

	// e.g. //http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
	urlToken := fmt.Sprintf("%s/register?&t=%s", br.mailRedirectBaseUrl, registrationToken)

	sendMailResponse := mailer.SendResponse{}
	if err := br.mailer.Send(&mailer.SendRequest{
		//From    string
		To: request.User.EmailAddress,
		//Cc      string
		Subject: "Welcome to SpotNav",
		Body:    fmt.Sprintf("Welcome to Spot Nav. Click the link to continue. %s", urlToken),
		//Bcc     []string
	},
		&sendMailResponse); err != nil {
		return err
	}

	response.URLToken = urlToken

	return nil
}

func (br *basicRegistrar) ValidateRegisterCompanyAdminUserRequest(request *partyRegistrar.RegisterCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// user party type and id must be as was in claims otherwise someone is
		// trying to abuse the registration token
		if request.User.PartyType != request.Claims.PartyDetails().PartyType {
			reasonsInvalid = append(reasonsInvalid, "user party type incorrect")
		}
		if request.User.PartyId != request.Claims.PartyDetails().PartyId {
			reasonsInvalid = append(reasonsInvalid, "user party id incorrect")
		}
	}

	// email address must be the same as the admin email address on the party entity
	// retrieve party to confirm this
	companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
	if err := br.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Claims.PartyDetails().PartyId,
	},
		&companyRetrieveResponse); err != nil {
		return registrarException.UnableToRetrieveParty{Reasons: []string{"company party", err.Error()}}
	}
	if companyRetrieveResponse.Company.AdminEmailAddress != request.User.EmailAddress {
		reasonsInvalid = append(reasonsInvalid, "user email address incorrect")
	}

	// password field must be blank
	if len(request.User.Password) != 0 {
		reasonsInvalid = append(reasonsInvalid, "user password must be blank")
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

	// give the user the necessary roles
	request.User.Roles = append(request.User.Roles, roleSetup.CompanyAdmin.Name)
	request.User.Roles = append(request.User.Roles, roleSetup.CompanyUser.Name)

	// create the user
	userCreateResponse := userRecordHandler.CreateResponse{}
	if err := br.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		Claims: request.Claims,
		User:   request.User,
	},
		&userCreateResponse); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userRecordHandler.ChangePasswordResponse{}
	if err := br.userRecordHandler.ChangePassword(&userRecordHandler.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: userCreateResponse.User.Id},
		NewPassword: request.Password,
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	response.User = userCreateResponse.User

	return nil
}

func (br *basicRegistrar) ValidateInviteCompanyUserRequest(request *partyRegistrar.InviteCompanyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	// validate the user
	validateUserResponse := userRecordHandler.ValidateResponse{}
	if err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
		Claims: request.Claims,
		User:   request.User,
		Method: partyRegistrar.InviteUser,
	}, &validateUserResponse); err != nil {
		reasonsInvalid = append(reasonsInvalid, "error validating user")
	}
	if len(validateUserResponse.ReasonsInvalid) > 0 {
		reasonsInvalid = append(reasonsInvalid, "user invalid: "+strings.Join(reasonsInvalid, " ;"))
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (br *basicRegistrar) InviteCompanyUser(request *partyRegistrar.InviteCompanyUserRequest, response *partyRegistrar.InviteCompanyUserResponse) error {
	if err := br.ValidateInviteCompanyUserRequest(request); err != nil {
		return err
	}

	// retrieve the Company whose admin will receive invite
	companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
	if err := br.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.User.PartyId,
	},
		&companyRetrieveResponse); err != nil {
		return registrarException.UnableToRetrieveParty{Reasons: []string{"company party", err.Error()}}
	}

	// Generate the registration token
	registerCompanyAdminUserClaims := registerCompanyAdminUser.RegisterCompanyAdminUser{
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: request.Claims.PartyDetails().PartyType,
		ParentId:        request.Claims.PartyDetails().PartyId,
		PartyType:       party.Company,
		PartyId:         id.Identifier{Id: companyRetrieveResponse.Company.Id},
		// User:
	}

	registrationToken, err := br.jwtGenerator.GenerateToken(registerCompanyAdminUserClaims)
	if err != nil {
		//Unexpected Error!
		return errors.New("log In failed")
	}

	// e.g. //http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
	urlToken := fmt.Sprintf("%s/register?&t=%s", br.mailRedirectBaseUrl, registrationToken)

	sendMailResponse := mailer.SendResponse{}
	if err := br.mailer.Send(&mailer.SendRequest{
		//From    string
		To: companyRetrieveResponse.Company.AdminEmailAddress,
		//Cc      string
		Subject: "Welcome to SpotNav",
		Body:    fmt.Sprintf("Welcome to Spot Nav. Click the link to continue. %s", urlToken),
		//Bcc     []string
	},
		&sendMailResponse); err != nil {
		return err
	}

	response.URLToken = urlToken

	return nil
}

func (br *basicRegistrar) ValidateInviteClientAdminUserRequest(request *partyRegistrar.InviteClientAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	// validate the admin user
	validateUserResponse := userRecordHandler.ValidateResponse{}
	if err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
		Claims: request.Claims,
		User:   request.User,
		Method: partyRegistrar.InviteAdminUser,
	}, &validateUserResponse); err != nil {
		reasonsInvalid = append(reasonsInvalid, "error validating user")
	}
	if len(validateUserResponse.ReasonsInvalid) > 0 {
		reasonsInvalid = append(reasonsInvalid, "user invalid: "+strings.Join(reasonsInvalid, " ;"))
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (br *basicRegistrar) InviteClientAdminUser(request *partyRegistrar.InviteClientAdminUserRequest, response *partyRegistrar.InviteClientAdminUserResponse) error {
	if err := br.ValidateInviteClientAdminUserRequest(request); err != nil {
		return err
	}

	// retrieve the Client whose admin will receive invite
	clientRetrieveResponse := clientRecordHandler.RetrieveResponse{}
	if err := br.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.User.PartyId,
	},
		&clientRetrieveResponse); err != nil {
		return registrarException.UnableToRetrieveParty{Reasons: []string{"client party", err.Error()}}
	}

	// Generate the registration token
	registerClientAdminUserClaims := registerClientAdminUser.RegisterClientAdminUser{
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: request.Claims.PartyDetails().PartyType,
		ParentId:        request.Claims.PartyDetails().PartyId,
		PartyType:       party.Client,
		PartyId:         id.Identifier{Id: clientRetrieveResponse.Client.Id},
		EmailAddress:    clientRetrieveResponse.Client.AdminEmailAddress,
	}

	registrationToken, err := br.jwtGenerator.GenerateToken(registerClientAdminUserClaims)
	if err != nil {
		//Unexpected Error!
		return errors.New("log In failed")
	}

	//http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
	urlToken := fmt.Sprintf("%s/register?&t=%s", br.mailRedirectBaseUrl, registrationToken)

	sendMailResponse := mailer.SendResponse{}
	if err := br.mailer.Send(&mailer.SendRequest{
		//From    string
		To: clientRetrieveResponse.Client.AdminEmailAddress,
		//Cc      string
		Subject: "Welcome to SpotNav",
		Body:    fmt.Sprintf("Welcome to Spot Nav. Click the link to continue. %s", urlToken),
		//Bcc     []string
	},
		&sendMailResponse); err != nil {
		return err
	}

	response.URLToken = urlToken

	return nil
}

func (br *basicRegistrar) ValidateRegisterClientAdminUserRequest(request *partyRegistrar.RegisterClientAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// user party type and id must be as was in claims otherwise someone is
		// trying to abuse the registration token
		if request.User.PartyType != request.Claims.PartyDetails().PartyType {
			reasonsInvalid = append(reasonsInvalid, "user party type incorrect")
		}
		if request.User.PartyId != request.Claims.PartyDetails().PartyId {
			reasonsInvalid = append(reasonsInvalid, "user party id incorrect")
		}
	}

	// email address must be the same as the admin email address on the party entity
	// retrieve party to confirm this
	clientRetrieveResponse := clientRecordHandler.RetrieveResponse{}
	if err := br.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Claims.PartyDetails().PartyId,
	},
		&clientRetrieveResponse); err != nil {
		return registrarException.UnableToRetrieveParty{Reasons: []string{"client party", err.Error()}}
	}
	if clientRetrieveResponse.Client.AdminEmailAddress != request.User.EmailAddress {
		reasonsInvalid = append(reasonsInvalid, "user email address incorrect")
	}

	// password field must be blank
	if len(request.User.Password) != 0 {
		reasonsInvalid = append(reasonsInvalid, "user password must be blank")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (br *basicRegistrar) RegisterClientAdminUser(request *partyRegistrar.RegisterClientAdminUserRequest, response *partyRegistrar.RegisterClientAdminUserResponse) error {
	if err := br.ValidateRegisterClientAdminUserRequest(request); err != nil {
		return err
	}

	// give the user the necessary roles
	request.User.Roles = append(request.User.Roles, roleSetup.ClientAdmin.Name)
	request.User.Roles = append(request.User.Roles, roleSetup.ClientUser.Name)

	// create the user
	userCreateResponse := userRecordHandler.CreateResponse{}
	if err := br.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		User: request.User,
	},
		&userCreateResponse); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userRecordHandler.ChangePasswordResponse{}
	if err := br.userRecordHandler.ChangePassword(&userRecordHandler.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: userCreateResponse.User.Id},
		NewPassword: request.Password,
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	response.User = userCreateResponse.User

	return nil
}
