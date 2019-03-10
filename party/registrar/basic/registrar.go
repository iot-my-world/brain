package basic

import (
	"crypto/rsa"
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
	"gitlab.com/iotTracker/brain/search/identifier/adminEmailAddress"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyUser"
	"gitlab.com/iotTracker/brain/security/claims/registerClientUser"
	listText "gitlab.com/iotTracker/brain/search/criterion/list/text"
	"gitlab.com/iotTracker/brain/search/criterion"
)

type basicRegistrar struct {
	companyRecordHandler companyRecordHandler.RecordHandler
	userRecordHandler    userRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
	mailer               mailer.Mailer
	jwtGenerator         token.JWTGenerator
	mailRedirectBaseUrl  string
	systemClaims         *login.Login
}

func New(
	companyRecordHandler companyRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
	mailer mailer.Mailer,
	rsaPrivateKey *rsa.PrivateKey,
	mailRedirectBaseUrl string,
	systemClaims *login.Login,
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
		NewPassword: string(request.User.Password),
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	response.User = userCreateResponse.User

	return nil
}

func (br *basicRegistrar) ValidateInviteCompanyAdminUserRequest(request *partyRegistrar.InviteCompanyAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// the user in the invite request must not be registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user cannot be set to registered yet")
	}

	// password field must be blank
	if len(request.User.Password) != 0 {
		reasonsInvalid = append(reasonsInvalid, "user password must be blank")
	}

	// username field must be blank
	if request.User.Username != "" {
		reasonsInvalid = append(reasonsInvalid, "username must be blank")
	}

	// roles must be empty
	if len(request.User.Roles) != 0 {
		reasonsInvalid = append(reasonsInvalid, "user cannot have any roles yet")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

		// at the moment only system is allowed to be the parent of company admin users
		if request.User.ParentId.Id != br.systemClaims.PartyId.Id {
			reasonsInvalid = append(reasonsInvalid, "parentId must be system id")
		}
		if request.User.ParentPartyType != br.systemClaims.PartyType {
			reasonsInvalid = append(reasonsInvalid, "parentPartyType must be system")
		}

		// regardless of who is performing the invite the partyType of the user must be company
		if request.User.PartyType != party.Company {
			reasonsInvalid = append(reasonsInvalid, "user's partyType must be company")
		}

		// validate the new user for the invite company admin user method
		userValidateResponse := userRecordHandler.ValidateResponse{}
		err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
			// system claims since we want all users to be visible for the email address check done in validate user
			Claims: br.systemClaims,
			User:   request.User,
			Method: partyRegistrar.InviteCompanyAdminUser,
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
				Claims: *br.systemClaims,
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
					// system claims since we want all companies to be visible for this retrieval check
					Claims: *br.systemClaims,
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

	// Create the minimal company admin user
	userCreateResponse := userRecordHandler.CreateResponse{}
	if err := br.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		Claims: request.Claims,
		User:   request.User,
	},
		&userCreateResponse); err != nil {
		return err
	}

	// Update the id on the user
	request.User.Id = userCreateResponse.User.Id

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

	// user must not be set to registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user must not yet be registered")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

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
	userValidateResponse := userRecordHandler.ValidateResponse{}
	err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *br.systemClaims,
		User:   request.User,
		Method: partyRegistrar.RegisterCompanyAdminUser,
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

func (br *basicRegistrar) RegisterCompanyAdminUser(request *partyRegistrar.RegisterCompanyAdminUserRequest, response *partyRegistrar.RegisterCompanyAdminUserResponse) error {
	if err := br.ValidateRegisterCompanyAdminUserRequest(request); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userRecordHandler.ChangePasswordResponse{}
	if err := br.userRecordHandler.ChangePassword(&userRecordHandler.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	// retrieve the minimal user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := br.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userRetrieveResponse); err != nil {
		return err
	}

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.CompanyAdmin.Name, roleSetup.CompanyUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	userUpdateResponse := userRecordHandler.UpdateResponse{}
	if err := br.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userUpdateResponse); err != nil {
		return err
	}

	response.User = userUpdateResponse.User

	return nil
}

func (br *basicRegistrar) ValidateInviteCompanyUserRequest(request *partyRegistrar.InviteCompanyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// the user in the invite request must not be registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user cannot be set to registered yet")
	}

	// password field must be blank
	if len(request.User.Password) != 0 {
		reasonsInvalid = append(reasonsInvalid, "user password must be blank")
	}

	// username field must be blank
	if request.User.Username != "" {
		reasonsInvalid = append(reasonsInvalid, "username must be blank")
	}

	// roles must be empty
	if len(request.User.Roles) != 0 {
		reasonsInvalid = append(reasonsInvalid, "user cannot have any roles yet")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

		// unless the user performing the invite is system, the partyDetails of the new user must be the
		// same as the user performing the invite
		if request.Claims.PartyDetails().PartyType != party.System {
			if request.User.ParentPartyType != request.Claims.PartyDetails().ParentPartyType {
				reasonsInvalid = append(reasonsInvalid, "partentPartyType of user must be the same as the user performing invite")
			}
			if request.User.ParentId != request.Claims.PartyDetails().ParentId {
				reasonsInvalid = append(reasonsInvalid, "parentId of user must be the same as the user performing invite")
			}
			if request.User.PartyType != request.Claims.PartyDetails().PartyType {
				reasonsInvalid = append(reasonsInvalid, "partyType of user must be the same as the user performing invite")
			}
			if request.User.PartyId != request.Claims.PartyDetails().PartyId {
				reasonsInvalid = append(reasonsInvalid, "partyId of user must be the same as the user performing invite")
			}
		}

		// regardless of who is performing the invite the partyType of the user must be company
		if request.User.PartyType != party.Company {
			reasonsInvalid = append(reasonsInvalid, "user's partyType must be company")
		}

		// at the moment only system is allowed to be the parent of company users
		if request.User.ParentId.Id != br.systemClaims.PartyId.Id {
			reasonsInvalid = append(reasonsInvalid, "parentId must be system id")
		}
		if request.User.ParentPartyType != br.systemClaims.PartyType {
			reasonsInvalid = append(reasonsInvalid, "parentPartyType must be system")
		}

		// validate the new user for the invite company user method
		userValidateResponse := userRecordHandler.ValidateResponse{}
		err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
			// system claims since we want all users to be visible for the email address check done in validate user
			Claims: br.systemClaims,
			User:   request.User,
			Method: partyRegistrar.InviteCompanyUser,
		}, &userValidateResponse)
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "unable to validate new user")
		} else {
			for _, reason := range userValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
			}
		}

		if request.User.EmailAddress != "" {

			// Check if the users email has already been assigned to a company entity as admin email
			companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
			if err := br.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
				// system claims since we want all companies to be visible for this retrieval check
				Claims: *br.systemClaims,
				Identifier: adminEmailAddress.Identifier{
					AdminEmailAddress: request.User.EmailAddress,
				},
			}, &companyRetrieveResponse); err != nil {
				switch err.(type) {
				case companyRecordHandlerException.NotFound:
					// [2] this is what we want, do nothing
				default:
					reasonsInvalid = append(reasonsInvalid, "unable to perform company retrieve to confirm correct email address: "+err.Error())
				}
			} else {
				// [3] if a company was found, this email address is therefore already being used
				reasonsInvalid = append(reasonsInvalid, "emailAddress used as admin email address on a company entity")
			}

			// Check if the users email has already been assigned to a client entity as admin email
			if request.User.EmailAddress != "" {
				if err := br.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
					// system claims since we want all companies to be visible for this retrieval check
					Claims: *br.systemClaims,
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

func (br *basicRegistrar) InviteCompanyUser(request *partyRegistrar.InviteCompanyUserRequest, response *partyRegistrar.InviteCompanyUserResponse) error {
	if err := br.ValidateInviteCompanyUserRequest(request); err != nil {
		return err
	}
	// Create the minimal company user
	userCreateResponse := userRecordHandler.CreateResponse{}
	if err := br.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		Claims: request.Claims,
		User:   request.User,
	},
		&userCreateResponse); err != nil {
		return err
	}

	// Update the id on the user
	request.User.Id = userCreateResponse.User.Id

	// Generate the registration token for the company user to register
	registerCompanyUserClaims := registerCompanyUser.RegisterCompanyUser{
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: request.User.ParentPartyType,
		ParentId:        request.User.ParentId,
		PartyType:       request.User.PartyType,
		PartyId:         request.User.PartyId,
		User:            request.User,
	}
	registrationToken, err := br.jwtGenerator.GenerateToken(registerCompanyUserClaims)
	if err != nil {
		return registrarException.TokenGeneration{Reasons: []string{"inviteCompanyUser", err.Error()}}
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

func (br *basicRegistrar) ValidateRegisterCompanyUserRequest(request *partyRegistrar.RegisterCompanyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// user must not be set to registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user must not yet be registered")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

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
	userValidateResponse := userRecordHandler.ValidateResponse{}
	err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *br.systemClaims,
		User:   request.User,
		Method: partyRegistrar.RegisterCompanyUser,
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

func (br *basicRegistrar) RegisterCompanyUser(request *partyRegistrar.RegisterCompanyUserRequest, response *partyRegistrar.RegisterCompanyUserResponse) error {
	if err := br.ValidateRegisterCompanyUserRequest(request); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userRecordHandler.ChangePasswordResponse{}
	if err := br.userRecordHandler.ChangePassword(&userRecordHandler.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	// retrieve the minimal user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := br.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userRetrieveResponse); err != nil {
		return err
	}

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.CompanyUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	userUpdateResponse := userRecordHandler.UpdateResponse{}
	if err := br.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userUpdateResponse); err != nil {
		return err
	}

	response.User = userUpdateResponse.User

	return nil
}

func (br *basicRegistrar) ValidateInviteClientAdminUserRequest(request *partyRegistrar.InviteClientAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// the user in the invite request must not be registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user cannot be set to registered yet")
	}

	// password field must be blank
	if len(request.User.Password) != 0 {
		reasonsInvalid = append(reasonsInvalid, "user password must be blank")
	}

	// username field must be blank
	if request.User.Username != "" {
		reasonsInvalid = append(reasonsInvalid, "username must be blank")
	}

	// roles must be empty
	if len(request.User.Roles) != 0 {
		reasonsInvalid = append(reasonsInvalid, "user cannot have any roles yet")
	}

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

		// regardless of who is performing the invite the partyType of the user must be client
		if request.User.PartyType != party.Client {
			reasonsInvalid = append(reasonsInvalid, "user's partyType must be client")
		}

		// validate the new user for the invite admin user method
		userValidateResponse := userRecordHandler.ValidateResponse{}
		err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
			// system claims since we want all users to be visible for the email address check done in validate user
			Claims: *br.systemClaims,
			User:   request.User,
			Method: partyRegistrar.InviteClientAdminUser,
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
			// [1] try and retrieve the client party with the email address
			clientRetrieveResponse := clientRecordHandler.RetrieveResponse{}
			if err := br.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
				// system claims since we want all clients to be visible for this retrieval check
				Claims: *br.systemClaims,
				Identifier: adminEmailAddress.Identifier{
					AdminEmailAddress: request.User.EmailAddress,
				},
			}, &clientRetrieveResponse); err != nil {
				switch err.(type) {
				case clientRecordHandlerException.NotFound:
					// [2] if no client entity is found this is an issue
					reasonsInvalid = append(reasonsInvalid, "client entity could not be retrieved by the given admin users email address")
				default:
					reasonsInvalid = append(reasonsInvalid, "unable to perform client retrieve to confirm correct email address: "+err.Error())
				}
			} else {
				// [3] if a client was found the id of the client must be the users partyId
				if clientRetrieveResponse.Client.Id != request.User.PartyId.Id {
					// if the id of this other client entity is not the same as the party id of this user then
					reasonsInvalid = append(reasonsInvalid, "emailAddress used as admin email address on another client entity")
				}
			}

			// Check if the users email has already been assigned to another company entity as admin email
			if request.User.EmailAddress != "" {
				if err := br.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
					// system claims since we want all companies to be visible for this retrieval check
					Claims: *br.systemClaims,
					Identifier: adminEmailAddress.Identifier{
						AdminEmailAddress: request.User.EmailAddress,
					},
				},
					&companyRecordHandler.RetrieveResponse{}); err != nil {
					switch err.(type) {
					case companyRecordHandlerException.NotFound:
						// this is what we want, do nothing
					default:
						reasonsInvalid = append(reasonsInvalid, "unable to confirm admin user email address uniqueness")
					}
				} else {
					// there was no error, this email address is already taken by some client entity
					reasonsInvalid = append(reasonsInvalid, "emailAddress used as admin email address on a company entity")
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

func (br *basicRegistrar) InviteClientAdminUser(request *partyRegistrar.InviteClientAdminUserRequest, response *partyRegistrar.InviteClientAdminUserResponse) error {
	if err := br.ValidateInviteClientAdminUserRequest(request); err != nil {
		return err
	}

	// Create the minimal client admin user
	userCreateResponse := userRecordHandler.CreateResponse{}
	if err := br.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		Claims: request.Claims,
		User:   request.User,
	},
		&userCreateResponse); err != nil {
		return err
	}

	// update the id on the user
	request.User.Id = userCreateResponse.User.Id

	// Generate the registration token for the client admin user to register
	registerClientAdminUserClaims := registerClientAdminUser.RegisterClientAdminUser{
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: request.User.ParentPartyType,
		ParentId:        request.User.ParentId,
		PartyType:       request.User.PartyType,
		PartyId:         request.User.PartyId,
		User:            request.User,
	}
	registrationToken, err := br.jwtGenerator.GenerateToken(registerClientAdminUserClaims)
	if err != nil {
		//Unexpected Error!
		return registrarException.TokenGeneration{Reasons: []string{"inviteClientAdminUser", err.Error()}}
	}

	//http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
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

func (br *basicRegistrar) ValidateRegisterClientAdminUserRequest(request *partyRegistrar.RegisterClientAdminUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// user must not be set to registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user must not yet be registered")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

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
	userValidateResponse := userRecordHandler.ValidateResponse{}
	err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *br.systemClaims,
		User:   request.User,
		Method: partyRegistrar.RegisterClientAdminUser,
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

func (br *basicRegistrar) RegisterClientAdminUser(request *partyRegistrar.RegisterClientAdminUserRequest, response *partyRegistrar.RegisterClientAdminUserResponse) error {
	if err := br.ValidateRegisterClientAdminUserRequest(request); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userRecordHandler.ChangePasswordResponse{}
	if err := br.userRecordHandler.ChangePassword(&userRecordHandler.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	// retrieve the minimal user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := br.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userRetrieveResponse); err != nil {
		return err
	}

	// give the user the necessary roles
	request.User.Roles = append(request.User.Roles, roleSetup.ClientAdmin.Name)
	request.User.Roles = append(request.User.Roles, roleSetup.ClientUser.Name)

	// set the user to registered
	request.User.Registered = true

	// update the user
	userUpdateResponse := userRecordHandler.UpdateResponse{}
	if err := br.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userUpdateResponse); err != nil {
		return err
	}

	response.User = userUpdateResponse.User

	return nil
}

func (br *basicRegistrar) ValidateInviteClientUserRequest(request *partyRegistrar.InviteClientUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// the user in the invite request must not be registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user cannot be set to registered yet")
	}

	// password field must be blank
	if len(request.User.Password) != 0 {
		reasonsInvalid = append(reasonsInvalid, "user password must be blank")
	}

	// username field must be blank
	if request.User.Username != "" {
		reasonsInvalid = append(reasonsInvalid, "username must be blank")
	}

	// roles must be empty
	if len(request.User.Roles) != 0 {
		reasonsInvalid = append(reasonsInvalid, "user cannot have any roles yet")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

		// unless the user performing the invite is system, the party details of the new user must be the
		// same as the user performing the invite
		if request.Claims.PartyDetails().PartyType != party.System {
			if request.User.ParentPartyType != request.Claims.PartyDetails().ParentPartyType {
				reasonsInvalid = append(reasonsInvalid, "partentPartyType of user must be the same as the user performing invite")
			}
			if request.User.ParentId != request.Claims.PartyDetails().ParentId {
				reasonsInvalid = append(reasonsInvalid, "parentId of user must be the same as the user performing invite")
			}
			if request.User.PartyType != request.Claims.PartyDetails().PartyType {
				reasonsInvalid = append(reasonsInvalid, "partyType of user must be the same as the user performing invite")
			}
			if request.User.PartyId != request.Claims.PartyDetails().PartyId {
				reasonsInvalid = append(reasonsInvalid, "partyId of user must be the same as the user performing invite")
			}
		}

		// regardless of who is performing the invite the partyType of the user must be client
		if request.User.PartyType != party.Client {
			reasonsInvalid = append(reasonsInvalid, "user's partyType must be client")
		}

		// validate the new user for the invite client user method
		userValidateResponse := userRecordHandler.ValidateResponse{}
		err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
			// system claims since we want all users to be visible for the email address check done in validate user
			Claims: br.systemClaims,
			User:   request.User,
			Method: partyRegistrar.InviteClientUser,
		}, &userValidateResponse)
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "unable to validate new user")
		} else {
			for _, reason := range userValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
			}
		}

		if request.User.EmailAddress != "" {

			// Check if the users email has already been assigned to a company entity as admin email
			companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
			if err := br.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
				// system claims since we want all companies to be visible for this retrieval check
				Claims: *br.systemClaims,
				Identifier: adminEmailAddress.Identifier{
					AdminEmailAddress: request.User.EmailAddress,
				},
			}, &companyRetrieveResponse); err != nil {
				switch err.(type) {
				case companyRecordHandlerException.NotFound:
					// [2] this is what we want, do nothing
				default:
					reasonsInvalid = append(reasonsInvalid, "unable to perform company retrieve to confirm correct email address: "+err.Error())
				}
			} else {
				// [3] if a company was found, this email address is therefore already being used
				reasonsInvalid = append(reasonsInvalid, "emailAddress used as admin email address on a company entity")
			}

			// Check if the users email has already been assigned to a client entity as admin email
			if request.User.EmailAddress != "" {
				if err := br.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
					// system claims since we want all companies to be visible for this retrieval check
					Claims: *br.systemClaims,
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

func (br *basicRegistrar) InviteClientUser(request *partyRegistrar.InviteClientUserRequest, response *partyRegistrar.InviteClientUserResponse) error {
	if err := br.ValidateInviteClientUserRequest(request); err != nil {
		return err
	}
	// Create the minimal client user
	userCreateResponse := userRecordHandler.CreateResponse{}
	if err := br.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		Claims: request.Claims,
		User:   request.User,
	},
		&userCreateResponse); err != nil {
		return err
	}

	// Update the id on the user
	request.User.Id = userCreateResponse.User.Id

	// Generate the registration token for the company user to register
	registerClientUserClaims := registerClientUser.RegisterClientUser{
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: request.User.ParentPartyType,
		ParentId:        request.User.ParentId,
		PartyType:       request.User.PartyType,
		PartyId:         request.User.PartyId,
		User:            request.User,
	}
	registrationToken, err := br.jwtGenerator.GenerateToken(registerClientUserClaims)
	if err != nil {
		return registrarException.TokenGeneration{Reasons: []string{"inviteClientUser", err.Error()}}
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

func (br *basicRegistrar) ValidateRegisterClientUserRequest(request *partyRegistrar.RegisterClientUserRequest) error {
	reasonsInvalid := make([]string, 0)

	// user must not be set to registered
	if request.User.Registered {
		reasonsInvalid = append(reasonsInvalid, "user must not yet be registered")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {

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
	userValidateResponse := userRecordHandler.ValidateResponse{}
	err := br.userRecordHandler.Validate(&userRecordHandler.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *br.systemClaims,
		User:   request.User,
		Method: partyRegistrar.RegisterClientUser,
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

func (br *basicRegistrar) RegisterClientUser(request *partyRegistrar.RegisterClientUserRequest, response *partyRegistrar.RegisterClientUserResponse) error {
	if err := br.ValidateRegisterClientUserRequest(request); err != nil {
		return err
	}

	// change the users password
	userChangePasswordResponse := userRecordHandler.ChangePasswordResponse{}
	if err := br.userRecordHandler.ChangePassword(&userRecordHandler.ChangePasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	},
		&userChangePasswordResponse); err != nil {
		return err
	}

	// retrieve the minimal user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := br.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userRetrieveResponse); err != nil {
		return err
	}

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.ClientUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	userUpdateResponse := userRecordHandler.UpdateResponse{}
	if err := br.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	},
		&userUpdateResponse); err != nil {
		return err
	}

	response.User = userUpdateResponse.User

	return nil
}

func (br *basicRegistrar) ValidateAreAdminsRegisteredRequest(request *partyRegistrar.AreAdminsRegisteredRequest) error {
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

func (br *basicRegistrar) AreAdminsRegistered(request *partyRegistrar.AreAdminsRegisteredRequest, response *partyRegistrar.AreAdminsRegisteredResponse) error {
	if err := br.ValidateAreAdminsRegisteredRequest(request); err != nil {
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
			return registrarException.PartyTypeInvalid{Reasons: []string{"areAdminsRegistered", string(partyDetail.PartyType)}}
		}
	}

	// collect companies in request
	companyCollectResponse := companyRecordHandler.CollectResponse{}
	if err := br.companyRecordHandler.Collect(&companyRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "id",
				List:  companyIds,
			},
		},
	}, &companyCollectResponse); err != nil {
		return registrarException.UnableToCollectParties{Reasons: []string{"company", err.Error()}}
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
	if err := br.userRecordHandler.Collect(&userRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "emailAddress",
				List:  companyAdminEmails,
			},
		},
	}, &companyAdminUserCollectResponse); err != nil {
		return registrarException.UnableToCollectParties{Reasons: []string{"companyAdminUsers", err.Error()}}
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
	if err := br.clientRecordHandler.Collect(&clientRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "id",
				List:  clientIds,
			},
		},
	}, &clientCollectResponse); err != nil {
		return registrarException.UnableToCollectParties{Reasons: []string{"client", err.Error()}}
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
	if err := br.userRecordHandler.Collect(&userRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "emailAddress",
				List:  clientAdminEmails,
			},
		},
	}, &clientAdminUserCollectResponse); err != nil {
		return registrarException.UnableToCollectParties{Reasons: []string{"clientAdminUsers", err.Error()}}
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
