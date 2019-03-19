package basic

import (
	"crypto/rsa"
	"fmt"
	"gitlab.com/iotTracker/brain/email/mailer"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	partyRegistrar "gitlab.com/iotTracker/brain/party/registrar"
	partyRegistrarAction "gitlab.com/iotTracker/brain/party/registrar/action"
	partyRegistrarException "gitlab.com/iotTracker/brain/party/registrar/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	listText "gitlab.com/iotTracker/brain/search/criterion/list/text"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/username"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"gitlab.com/iotTracker/brain/security/claims/registerClientAdminUser"
	"gitlab.com/iotTracker/brain/security/claims/registerClientUser"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyAdminUser"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyUser"
	roleSetup "gitlab.com/iotTracker/brain/security/role/setup"
	"gitlab.com/iotTracker/brain/security/token"
	userAdministrator "gitlab.com/iotTracker/brain/user/administrator"
	userRecordHandler "gitlab.com/iotTracker/brain/user/recordHandler"
	userRecordHandlerException "gitlab.com/iotTracker/brain/user/recordHandler/exception"
	userValidator "gitlab.com/iotTracker/brain/user/validator"
	"time"
)

type registrar struct {
	companyRecordHandler companyRecordHandler.RecordHandler
	userRecordHandler    userRecordHandler.RecordHandler
	userValidator        userValidator.Validator
	userAdministrator    userAdministrator.Administrator
	clientRecordHandler  clientRecordHandler.RecordHandler
	mailer               mailer.Mailer
	jwtGenerator         token.JWTGenerator
	mailRedirectBaseUrl  string
	systemClaims         *login.Login
}

func New(
	companyRecordHandler companyRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
	userValidator userValidator.Validator,
	userAdministrator userAdministrator.Administrator,
	clientRecordHandler clientRecordHandler.RecordHandler,
	mailer mailer.Mailer,
	rsaPrivateKey *rsa.PrivateKey,
	mailRedirectBaseUrl string,
	systemClaims *login.Login,
) *registrar {
	return &registrar{
		companyRecordHandler: companyRecordHandler,
		userRecordHandler:    userRecordHandler,
		userValidator:        userValidator,
		userAdministrator:    userAdministrator,
		clientRecordHandler:  clientRecordHandler,
		mailer:               mailer,
		jwtGenerator:         token.NewJWTGenerator(rsaPrivateKey),
		mailRedirectBaseUrl:  mailRedirectBaseUrl,
		systemClaims:         systemClaims,
	}
}

func (r *registrar) RegisterSystemAdminUser(request *partyRegistrar.RegisterSystemAdminUserRequest, response *partyRegistrar.RegisterSystemAdminUserResponse) error {

	// check if the system admin user already exists (i.e. has already been registered)
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: username.Identifier{Username: request.User.Username},
	},
		&userRetrieveResponse)
	switch err.(type) {
	case nil:
		// this means that the user already exists
		response.User = userRetrieveResponse.User
		return partyRegistrarException.AlreadyRegistered{}
	case userRecordHandlerException.NotFound:
		// this is fine, we will be creating the user now
	default:
		return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
	}

	// create the user
	userCreateResponse := userRecordHandler.CreateResponse{}
	if err := r.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		User: request.User,
	},
		&userCreateResponse); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userAdministrator.ChangePasswordResponse{}
	if err := r.userAdministrator.ChangePassword(&userAdministrator.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: userCreateResponse.User.Id},
		NewPassword: string(request.User.Password),
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	response.User = userCreateResponse.User

	return nil
}

func (r *registrar) ValidateInviteCompanyAdminUserRequest(request *partyRegistrar.InviteCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.CompanyIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "company identifier is nil")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) InviteCompanyAdminUser(request *partyRegistrar.InviteCompanyAdminUserRequest, response *partyRegistrar.InviteCompanyAdminUserResponse) error {
	if err := r.ValidateInviteCompanyAdminUserRequest(request); err != nil {
		return err
	}

	// Retrieve the company party
	companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
	if err := r.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.CompanyIdentifier,
	}, &companyRetrieveResponse); err != nil {
		return partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"company", err.Error()}}
	}

	// Retrieve the minimal company admin user which was created on company creation
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims: request.Claims,
		Identifier: emailAddress.Identifier{
			EmailAddress: companyRetrieveResponse.Company.AdminEmailAddress,
		},
	}, &userRetrieveResponse); err != nil {
		return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		return partyRegistrarException.AlreadyRegistered{}
	}

	// Generate the registration token for the company admin user to register
	registerCompanyAdminUserClaims := registerCompanyAdminUser.RegisterCompanyAdminUser{
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: userRetrieveResponse.User.ParentPartyType,
		ParentId:        userRetrieveResponse.User.ParentId,
		PartyType:       userRetrieveResponse.User.PartyType,
		PartyId:         userRetrieveResponse.User.PartyId,
		User:            userRetrieveResponse.User,
	}
	registrationToken, err := r.jwtGenerator.GenerateToken(registerCompanyAdminUserClaims)
	if err != nil {
		return partyRegistrarException.TokenGeneration{Reasons: []string{"inviteCompanyAdminUser", err.Error()}}
	}

	// e.g. //http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
	urlToken := fmt.Sprintf("%s/register?&t=%s", r.mailRedirectBaseUrl, registrationToken)

	sendMailResponse := mailer.SendResponse{}
	if err := r.mailer.Send(&mailer.SendRequest{
		//From    string
		To: userRetrieveResponse.User.EmailAddress,
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

func (r *registrar) ValidateRegisterCompanyAdminUserRequest(request *partyRegistrar.RegisterCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// user must not be set to registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user must not yet be registered")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

		// try and retrieve a user with this id to see if they have already been invited
		userRetrieveResponse := userRecordHandler.RetrieveResponse{}
		if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: id.Identifier{Id: request.User.Id},
		}, &userRetrieveResponse); err == nil {
			// user should exist but should not yet be registered
			if userRetrieveResponse.User.Registered {
				return partyRegistrarException.AlreadyRegistered{}
			}
		} else {
			return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
		}

		switch typedClaims := request.Claims.(type) {
		default:
			reasonsInvalid = append(reasonsInvalid, "cannot infer correct type from claims")

		case registerCompanyAdminUser.RegisterCompanyAdminUser:
			// confirm that all fields that were set on the user when the claims were generated have not been changed
			if request.User.Id != typedClaims.User.Id {
				reasonsInvalid = append(reasonsInvalid, "id has changed")
			}
			if request.User.EmailAddress != typedClaims.User.EmailAddress {
				reasonsInvalid = append(reasonsInvalid, "email address has changed")
			}
			if request.User.ParentPartyType != typedClaims.User.ParentPartyType {
				reasonsInvalid = append(reasonsInvalid, "parent party type has changed")
			}
			if request.User.ParentId != typedClaims.User.ParentId {
				reasonsInvalid = append(reasonsInvalid, "parent id has changed")
			}
			if request.User.PartyType != typedClaims.User.PartyType {
				reasonsInvalid = append(reasonsInvalid, "party type has changed")
			}
			if request.User.PartyId != typedClaims.User.PartyId {
				reasonsInvalid = append(reasonsInvalid, "party id has changed")
			}
			if len(request.User.Roles) != len(typedClaims.User.Roles) {
				reasonsInvalid = append(reasonsInvalid, "no of roles has changed")
			} else {
				// no of roles the same, compare roles
				for _, requestUserRole := range request.User.Roles {
					for roleIdx, claimsUserRole := range typedClaims.User.Roles {
						if claimsUserRole == requestUserRole {
							break
						}
						if roleIdx == len(typedClaims.User.Roles)-1 {
							reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("could not find role %s in user in claims", requestUserRole))
						}
					}
				}
			}
		}
	}

	// validate the user for the registration process
	userValidateResponse := userValidator.ValidateResponse{}
	err := r.userValidator.Validate(&userValidator.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *r.systemClaims,
		User:   request.User,
		Action: partyRegistrarAction.RegisterCompanyAdminUser,
	}, &userValidateResponse)
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate newAdminUser")
	} else {
		for _, reason := range userValidateResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) RegisterCompanyAdminUser(request *partyRegistrar.RegisterCompanyAdminUserRequest, response *partyRegistrar.RegisterCompanyAdminUserResponse) error {
	if err := r.ValidateRegisterCompanyAdminUserRequest(request); err != nil {
		return err
	}

	// retrieve the minimal user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	}, &userRetrieveResponse); err != nil {
		return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
	}

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.CompanyAdmin.Name, roleSetup.CompanyUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	userUpdateResponse := userRecordHandler.UpdateResponse{}
	if err := r.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userUpdateResponse); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userAdministrator.ChangePasswordResponse{}
	if err := r.userAdministrator.ChangePassword(&userAdministrator.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	response.User = userUpdateResponse.User

	return nil
}

func (r *registrar) ValidateInviteCompanyUserRequest(request *partyRegistrar.InviteCompanyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}
	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) InviteCompanyUser(request *partyRegistrar.InviteCompanyUserRequest, response *partyRegistrar.InviteCompanyUserResponse) error {
	if err := r.ValidateInviteCompanyUserRequest(request); err != nil {
		return err
	}

	// retrieve the user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	}, &userRetrieveResponse); err != nil {
		return partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"user retrieval", err.Error()}}
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		return partyRegistrarException.AlreadyRegistered{}
	}

	// Generate the registration token for the company user to register
	registerCompanyUserClaims := registerCompanyUser.RegisterCompanyUser{
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: userRetrieveResponse.User.ParentPartyType,
		ParentId:        userRetrieveResponse.User.ParentId,
		PartyType:       userRetrieveResponse.User.PartyType,
		PartyId:         userRetrieveResponse.User.PartyId,
		User:            userRetrieveResponse.User,
	}
	registrationToken, err := r.jwtGenerator.GenerateToken(registerCompanyUserClaims)
	if err != nil {
		return partyRegistrarException.TokenGeneration{Reasons: []string{"inviteCompanyUser", err.Error()}}
	}

	// e.g. //http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
	urlToken := fmt.Sprintf("%s/register?&t=%s", r.mailRedirectBaseUrl, registrationToken)

	sendMailResponse := mailer.SendResponse{}
	if err := r.mailer.Send(&mailer.SendRequest{
		//From    string
		To: userRetrieveResponse.User.EmailAddress,
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

func (r *registrar) ValidateRegisterCompanyUserRequest(request *partyRegistrar.RegisterCompanyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// user must not be set to registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user must not yet be registered")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

		// try and retrieve a user with this id to see if they have already been invited
		userRetrieveResponse := userRecordHandler.RetrieveResponse{}
		if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: id.Identifier{Id: request.User.Id},
		}, &userRetrieveResponse); err == nil {
			// user should exist but should not yet be registered
			if userRetrieveResponse.User.Registered {
				return partyRegistrarException.AlreadyRegistered{}
			}
		} else {
			return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
		}

		switch typedClaims := request.Claims.(type) {
		default:
			reasonsInvalid = append(reasonsInvalid, "cannot infer correct type from claims")

		case registerCompanyUser.RegisterCompanyUser:
			// confirm that all fields that were set on the user when the claims were generated have not been changed
			if request.User.Id != typedClaims.User.Id {
				reasonsInvalid = append(reasonsInvalid, "id has changed")
			}
			if request.User.EmailAddress != typedClaims.User.EmailAddress {
				reasonsInvalid = append(reasonsInvalid, "email address has changed")
			}
			if request.User.ParentPartyType != typedClaims.User.ParentPartyType {
				reasonsInvalid = append(reasonsInvalid, "parent party type has changed")
			}
			if request.User.ParentId != typedClaims.User.ParentId {
				reasonsInvalid = append(reasonsInvalid, "parent id has changed")
			}
			if request.User.PartyType != typedClaims.User.PartyType {
				reasonsInvalid = append(reasonsInvalid, "party type has changed")
			}
			if request.User.PartyId != typedClaims.User.PartyId {
				reasonsInvalid = append(reasonsInvalid, "party id has changed")
			}
			if len(request.User.Roles) != len(typedClaims.User.Roles) {
				reasonsInvalid = append(reasonsInvalid, "no of roles has changed")
			} else {
				// no of roles the same, compare roles
				for _, requestUserRole := range request.User.Roles {
					for roleIdx, claimsUserRole := range typedClaims.User.Roles {
						if claimsUserRole == requestUserRole {
							break
						}
						if roleIdx == len(typedClaims.User.Roles)-1 {
							reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("could not find role %s in user in claims", requestUserRole))
						}
					}
				}
			}
		}
	}

	// validate the user for the registration process
	userValidateResponse := userValidator.ValidateResponse{}
	err := r.userValidator.Validate(&userValidator.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *r.systemClaims,
		User:   request.User,
		Action: partyRegistrarAction.RegisterCompanyUser,
	}, &userValidateResponse)
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate new user")
	} else {
		for _, reason := range userValidateResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) RegisterCompanyUser(request *partyRegistrar.RegisterCompanyUserRequest, response *partyRegistrar.RegisterCompanyUserResponse) error {
	if err := r.ValidateRegisterCompanyUserRequest(request); err != nil {
		return err
	}

	// retrieve the minimal user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	}, &userRetrieveResponse); err != nil {
		return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
	}

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.CompanyUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	userUpdateResponse := userRecordHandler.UpdateResponse{}
	if err := r.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userUpdateResponse); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userAdministrator.ChangePasswordResponse{}
	if err := r.userAdministrator.ChangePassword(&userAdministrator.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	response.User = userUpdateResponse.User

	return nil
}

func (r *registrar) ValidateInviteClientAdminUserRequest(request *partyRegistrar.InviteClientAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.ClientIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "clientIdentifier is nil")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) InviteClientAdminUser(request *partyRegistrar.InviteClientAdminUserRequest, response *partyRegistrar.InviteClientAdminUserResponse) error {
	if err := r.ValidateInviteClientAdminUserRequest(request); err != nil {
		return err
	}

	// retrieve the client
	clientRetrieveResponse := clientRecordHandler.RetrieveResponse{}
	if err := r.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ClientIdentifier,
	}, &clientRetrieveResponse); err != nil {
		return partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"client", err.Error()}}
	}

	// retrieve the minimal client admin user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims: request.Claims,
		Identifier: emailAddress.Identifier{
			EmailAddress: clientRetrieveResponse.Client.AdminEmailAddress,
		},
	}, &userRetrieveResponse); err != nil {
		return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		return partyRegistrarException.AlreadyRegistered{}
	}

	// Generate the registration token for the client admin user to register
	registerClientAdminUserClaims := registerClientAdminUser.RegisterClientAdminUser{
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: userRetrieveResponse.User.ParentPartyType,
		ParentId:        userRetrieveResponse.User.ParentId,
		PartyType:       userRetrieveResponse.User.PartyType,
		PartyId:         userRetrieveResponse.User.PartyId,
		User:            userRetrieveResponse.User,
	}
	registrationToken, err := r.jwtGenerator.GenerateToken(registerClientAdminUserClaims)
	if err != nil {
		return partyRegistrarException.TokenGeneration{Reasons: []string{"inviteClientAdminUser", err.Error()}}
	}

	//http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
	urlToken := fmt.Sprintf("%s/register?&t=%s", r.mailRedirectBaseUrl, registrationToken)

	sendMailResponse := mailer.SendResponse{}
	if err := r.mailer.Send(&mailer.SendRequest{
		//From    string
		To: userRetrieveResponse.User.EmailAddress,
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

func (r *registrar) ValidateRegisterClientAdminUserRequest(request *partyRegistrar.RegisterClientAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// user must not be set to registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user must not yet be registered")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

		// try and retrieve a user with this id to see if they have already been invited
		userRetrieveResponse := userRecordHandler.RetrieveResponse{}
		if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: id.Identifier{Id: request.User.Id},
		}, &userRetrieveResponse); err == nil {
			// user should exist but should not yet be registered
			if userRetrieveResponse.User.Registered {
				return partyRegistrarException.AlreadyRegistered{}
			}
		} else {
			return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
		}

		switch typedClaims := request.Claims.(type) {
		default:
			reasonsInvalid = append(reasonsInvalid, "cannot infer correct type from claims")

		case registerClientAdminUser.RegisterClientAdminUser:
			// confirm that all fields that were set on the user when the claims were generated have not been changed
			if request.User.Id != typedClaims.User.Id {
				reasonsInvalid = append(reasonsInvalid, "id has changed")
			}
			if request.User.EmailAddress != typedClaims.User.EmailAddress {
				reasonsInvalid = append(reasonsInvalid, "email address has changed")
			}
			if request.User.ParentPartyType != typedClaims.User.ParentPartyType {
				reasonsInvalid = append(reasonsInvalid, "parent party type has changed")
			}
			if request.User.ParentId != typedClaims.User.ParentId {
				reasonsInvalid = append(reasonsInvalid, "parent id has changed")
			}
			if request.User.PartyType != typedClaims.User.PartyType {
				reasonsInvalid = append(reasonsInvalid, "party type has changed")
			}
			if request.User.PartyId != typedClaims.User.PartyId {
				reasonsInvalid = append(reasonsInvalid, "party id has changed")
			}
			if len(request.User.Roles) != len(typedClaims.User.Roles) {
				reasonsInvalid = append(reasonsInvalid, "no of roles has changed")
			} else {
				// no of roles the same, compare roles
				for _, requestUserRole := range request.User.Roles {
					for roleIdx, claimsUserRole := range typedClaims.User.Roles {
						if claimsUserRole == requestUserRole {
							break
						}
						if roleIdx == len(typedClaims.User.Roles)-1 {
							reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("could not find role %s in user in claims", requestUserRole))
						}
					}
				}
			}
		}
	}

	// validate the user for the registration process
	userValidateResponse := userValidator.ValidateResponse{}
	err := r.userValidator.Validate(&userValidator.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *r.systemClaims,
		User:   request.User,
		Action: partyRegistrarAction.RegisterClientAdminUser,
	}, &userValidateResponse)
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate newAdminUser")
	} else {
		for _, reason := range userValidateResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) RegisterClientAdminUser(request *partyRegistrar.RegisterClientAdminUserRequest, response *partyRegistrar.RegisterClientAdminUserResponse) error {
	if err := r.ValidateRegisterClientAdminUserRequest(request); err != nil {
		return err
	}

	// retrieve the minimal user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	}, &userRetrieveResponse); err != nil {
		return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
	}

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.ClientAdmin.Name, roleSetup.ClientUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	userUpdateResponse := userRecordHandler.UpdateResponse{}
	if err := r.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userUpdateResponse); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userAdministrator.ChangePasswordResponse{}
	if err := r.userAdministrator.ChangePassword(&userAdministrator.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	response.User = userUpdateResponse.User

	return nil
}

func (r *registrar) ValidateInviteClientUserRequest(request *partyRegistrar.InviteClientUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) InviteClientUser(request *partyRegistrar.InviteClientUserRequest, response *partyRegistrar.InviteClientUserResponse) error {
	if err := r.ValidateInviteClientUserRequest(request); err != nil {
		return err
	}

	// retrieve the user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	}, &userRetrieveResponse); err != nil {
		return partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"user retrieval", err.Error()}}
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		return partyRegistrarException.AlreadyRegistered{}
	}

	// Generate the registration token for the client user to register
	registerClientUserClaims := registerClientUser.RegisterClientUser{
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: userRetrieveResponse.User.ParentPartyType,
		ParentId:        userRetrieveResponse.User.ParentId,
		PartyType:       userRetrieveResponse.User.PartyType,
		PartyId:         userRetrieveResponse.User.PartyId,
		User:            userRetrieveResponse.User,
	}
	registrationToken, err := r.jwtGenerator.GenerateToken(registerClientUserClaims)
	if err != nil {
		return partyRegistrarException.TokenGeneration{Reasons: []string{"inviteClientUser", err.Error()}}
	}

	// e.g. //http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
	urlToken := fmt.Sprintf("%s/register?&t=%s", r.mailRedirectBaseUrl, registrationToken)

	sendMailResponse := mailer.SendResponse{}
	if err := r.mailer.Send(&mailer.SendRequest{
		//From    string
		To: userRetrieveResponse.User.EmailAddress,
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

func (r *registrar) ValidateRegisterClientUserRequest(request *partyRegistrar.RegisterClientUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// user must not be set to registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user must not yet be registered")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

		// try and retrieve a user with this id to see if they have already been invited
		userRetrieveResponse := userRecordHandler.RetrieveResponse{}
		if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: id.Identifier{Id: request.User.Id},
		}, &userRetrieveResponse); err == nil {
			// user should exist but should not yet be registered
			if userRetrieveResponse.User.Registered {
				return partyRegistrarException.AlreadyRegistered{}
			}
		} else {
			return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
		}

		switch typedClaims := request.Claims.(type) {
		default:
			reasonsInvalid = append(reasonsInvalid, "cannot infer correct type from claims")

		case registerClientUser.RegisterClientUser:
			// confirm that all fields that were set on the user when the claims were generated have not been changed
			if request.User.Id != typedClaims.User.Id {
				reasonsInvalid = append(reasonsInvalid, "id has changed")
			}
			if request.User.EmailAddress != typedClaims.User.EmailAddress {
				reasonsInvalid = append(reasonsInvalid, "email address has changed")
			}
			if request.User.ParentPartyType != typedClaims.User.ParentPartyType {
				reasonsInvalid = append(reasonsInvalid, "parent party type has changed")
			}
			if request.User.ParentId != typedClaims.User.ParentId {
				reasonsInvalid = append(reasonsInvalid, "parent id has changed")
			}
			if request.User.PartyType != typedClaims.User.PartyType {
				reasonsInvalid = append(reasonsInvalid, "party type has changed")
			}
			if request.User.PartyId != typedClaims.User.PartyId {
				reasonsInvalid = append(reasonsInvalid, "party id has changed")
			}
			if len(request.User.Roles) != len(typedClaims.User.Roles) {
				reasonsInvalid = append(reasonsInvalid, "no of roles has changed")
			} else {
				// no of roles the same, compare roles
				for _, requestUserRole := range request.User.Roles {
					for roleIdx, claimsUserRole := range typedClaims.User.Roles {
						if claimsUserRole == requestUserRole {
							break
						}
						if roleIdx == len(typedClaims.User.Roles)-1 {
							reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("could not find role %s in user in claims", requestUserRole))
						}
					}
				}
			}
		}
	}

	// validate the user for the registration process
	userValidateResponse := userValidator.ValidateResponse{}
	err := r.userValidator.Validate(&userValidator.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *r.systemClaims,
		User:   request.User,
		Action: partyRegistrarAction.RegisterClientUser,
	}, &userValidateResponse)
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate new user")
	} else {
		for _, reason := range userValidateResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *registrar) RegisterClientUser(request *partyRegistrar.RegisterClientUserRequest, response *partyRegistrar.RegisterClientUserResponse) error {
	if err := r.ValidateRegisterClientUserRequest(request); err != nil {
		return err
	}

	// retrieve the minimal user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userRetrieveResponse); err != nil {
		return brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
	}

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.ClientUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	userUpdateResponse := userRecordHandler.UpdateResponse{}
	if err := r.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userUpdateResponse); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userAdministrator.ChangePasswordResponse{}
	if err := r.userAdministrator.ChangePassword(&userAdministrator.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	response.User = userUpdateResponse.User

	return nil
}

func (r *registrar) ValidateAreAdminsRegisteredRequest(request *partyRegistrar.AreAdminsRegisteredRequest) error {
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

func (r *registrar) ValidateInviteUserRequest(request *partyRegistrar.InviteUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) InviteUser(request *partyRegistrar.InviteUserRequest, response *partyRegistrar.InviteUserResponse) error {
	if err := r.ValidateInviteUserRequest(request); err != nil {
		return err
	}

	// retrieve the user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	}, &userRetrieveResponse); err != nil {
		return partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"user retrieval", err.Error()}}
	}

	// the purpose of this service is to provide a generic way to invite a user from any party, admin user or not
	switch userRetrieveResponse.User.PartyType {
	case party.Company:
		// determine it this is the admin user
		companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
		if err := r.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Claims:     *r.systemClaims,
			Identifier: userRetrieveResponse.User.PartyId,
		}, &companyRetrieveResponse); err != nil {
			return partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"company", err.Error()}}
		}
		if userRetrieveResponse.User.EmailAddress == companyRetrieveResponse.Company.AdminEmailAddress {
			inviteCompanyAdminUserResponse := partyRegistrar.InviteCompanyAdminUserResponse{}
			if err := r.InviteCompanyAdminUser(&partyRegistrar.InviteCompanyAdminUserRequest{
				Claims:            request.Claims,
				CompanyIdentifier: userRetrieveResponse.User.PartyId,
			}, &inviteCompanyAdminUserResponse); err != nil {
				return err
			}
			response.URLToken = inviteCompanyAdminUserResponse.URLToken
		} else {
			inviteCompanyUserResponse := partyRegistrar.InviteCompanyUserResponse{}
			if err := r.InviteCompanyUser(&partyRegistrar.InviteCompanyUserRequest{
				Claims:         request.Claims,
				UserIdentifier: id.Identifier{Id: userRetrieveResponse.User.Id},
			}, &inviteCompanyUserResponse); err != nil {
				return err
			}
			response.URLToken = inviteCompanyUserResponse.URLToken
		}

	case party.Client:
		// determine it this is the admin user
		clientRetrieveResponse := clientRecordHandler.RetrieveResponse{}
		if err := r.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Claims:     *r.systemClaims,
			Identifier: userRetrieveResponse.User.PartyId,
		}, &clientRetrieveResponse); err != nil {
			return partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"company", err.Error()}}
		}
		if userRetrieveResponse.User.EmailAddress == clientRetrieveResponse.Client.AdminEmailAddress {
			inviteClientAdminUserResponse := partyRegistrar.InviteClientAdminUserResponse{}
			if err := r.InviteClientAdminUser(&partyRegistrar.InviteClientAdminUserRequest{
				Claims:           request.Claims,
				ClientIdentifier: userRetrieveResponse.User.PartyId,
			}, &inviteClientAdminUserResponse); err != nil {
				return err
			}
			response.URLToken = inviteClientAdminUserResponse.URLToken
		} else {
			inviteClientUserResponse := partyRegistrar.InviteClientUserResponse{}
			if err := r.InviteClientUser(&partyRegistrar.InviteClientUserRequest{
				Claims:         request.Claims,
				UserIdentifier: id.Identifier{Id: userRetrieveResponse.User.Id},
			}, &inviteClientUserResponse); err != nil {
				return err
			}
			response.URLToken = inviteClientUserResponse.URLToken
		}

	default:
		return partyRegistrarException.PartyTypeInvalid{Reasons: []string{string(userRetrieveResponse.User.PartyType)}}

	}

	return nil
}

func (r *registrar) AreAdminsRegistered(request *partyRegistrar.AreAdminsRegisteredRequest, response *partyRegistrar.AreAdminsRegisteredResponse) error {
	if err := r.ValidateAreAdminsRegisteredRequest(request); err != nil {
		return err
	}

	result := make(map[string]bool)
	companyIds := make([]string, 0)
	companyAdminEmails := make([]string, 0)
	clientIds := make([]string, 0)
	clientAdminEmails := make([]string, 0)

	// compose id lists for exact list criteria
	for _, partyDetail := range request.PartyDetails {
		switch partyDetail.PartyType {
		case party.System:
			result[partyDetail.PartyId.Id] = true
		case party.Company:
			companyIds = append(companyIds, partyDetail.PartyId.Id)
		case party.Client:
			clientIds = append(clientIds, partyDetail.PartyId.Id)
		default:
			return partyRegistrarException.PartyTypeInvalid{Reasons: []string{"areAdminsRegistered", string(partyDetail.PartyType)}}
		}
	}

	// collect companies in request
	companyCollectResponse := companyRecordHandler.CollectResponse{}
	if err := r.companyRecordHandler.Collect(&companyRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "id",
				List:  companyIds,
			},
		},
	}, &companyCollectResponse); err != nil {
		return partyRegistrarException.UnableToCollectParties{Reasons: []string{"company", err.Error()}}
	} else {
		// confirm that for every id received a company was returned
		if len(companyCollectResponse.Records) != len(companyIds) {
			return brainException.Unexpected{Reasons: []string{
				"no company records returned different to number of ids given",
				fmt.Sprintf("%d vs %d", len(companyCollectResponse.Records), len(companyIds)),
			}}
		}
	}
	// compose list of admin email addresses
	for companyIdx := range companyCollectResponse.Records {
		companyAdminEmails = append(companyAdminEmails, companyCollectResponse.Records[companyIdx].AdminEmailAddress)
	}
	// collect users with these admin email addresses
	companyAdminUserCollectResponse := userRecordHandler.CollectResponse{}
	if err := r.userRecordHandler.Collect(&userRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "emailAddress",
				List:  companyAdminEmails,
			},
		},
	}, &companyAdminUserCollectResponse); err != nil {
		return partyRegistrarException.UnableToCollectParties{Reasons: []string{"companyAdminUsers", err.Error()}}
	} else {
		// confirm that for every admin email a user was returned
		if len(companyAdminUserCollectResponse.Records) != len(companyAdminEmails) {
			return brainException.Unexpected{Reasons: []string{
				"no company admin users found different from number of admin emails found",
				fmt.Sprintf("%d vs %d", len(companyAdminUserCollectResponse.Records), len(companyAdminEmails)),
			}}
		}
	}
	// update result for the company admin users retrieved
	for companyAdminUserIdx := range companyAdminUserCollectResponse.Records {
		result[companyAdminUserCollectResponse.Records[companyAdminUserIdx].PartyId.Id] =
			companyAdminUserCollectResponse.Records[companyAdminUserIdx].Registered
	}

	// collect clients in request
	clientCollectResponse := clientRecordHandler.CollectResponse{}
	if err := r.clientRecordHandler.Collect(&clientRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "id",
				List:  clientIds,
			},
		},
	}, &clientCollectResponse); err != nil {
		return partyRegistrarException.UnableToCollectParties{Reasons: []string{"client", err.Error()}}
	} else {
		// confirm that for every id received a client was returned
		if len(clientCollectResponse.Records) != len(clientIds) {
			return brainException.Unexpected{Reasons: []string{
				"no client records returned different to number of ids given",
				fmt.Sprintf("%d vs %d", len(clientCollectResponse.Records), len(clientIds)),
			}}
		}
	}
	// compose list of admin email addresses
	for clientIdx := range clientCollectResponse.Records {
		clientAdminEmails = append(clientAdminEmails, clientCollectResponse.Records[clientIdx].AdminEmailAddress)
	}
	// collect users with these admin email addresses
	clientAdminUserCollectResponse := userRecordHandler.CollectResponse{}
	if err := r.userRecordHandler.Collect(&userRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "emailAddress",
				List:  clientAdminEmails,
			},
		},
	}, &clientAdminUserCollectResponse); err != nil {
		return partyRegistrarException.UnableToCollectParties{Reasons: []string{"clientAdminUsers", err.Error()}}
	} else {
		// confirm that for every admin email user was returned
		if len(clientAdminUserCollectResponse.Records) != len(clientAdminEmails) {
			return brainException.Unexpected{Reasons: []string{
				"no client admin users found different from number of admin emails found",
				fmt.Sprintf("%d vs %d", len(clientAdminUserCollectResponse.Records), len(clientAdminEmails)),
			}}
		}
	}
	// update result for the client admin users retrieved
	for clientAdminUserIdx := range clientAdminUserCollectResponse.Records {
		result[clientAdminUserCollectResponse.Records[clientAdminUserIdx].PartyId.Id] =
			clientAdminUserCollectResponse.Records[clientAdminUserIdx].Registered
	}

	response.Result = result
	return nil
}
