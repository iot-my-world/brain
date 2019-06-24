package basic

import (
	"crypto/rsa"
	"fmt"
	emailGenerator "github.com/iot-my-world/brain/communication/email/generator"
	registrationEmail "github.com/iot-my-world/brain/communication/email/generator/registration"
	"github.com/iot-my-world/brain/communication/email/mailer"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/party"
	clientRecordHandler "github.com/iot-my-world/brain/party/client/recordHandler"
	companyRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler"
	partyRegistrar "github.com/iot-my-world/brain/party/registrar"
	partyRegistrarAction "github.com/iot-my-world/brain/party/registrar/action"
	partyRegistrarException "github.com/iot-my-world/brain/party/registrar/exception"
	"github.com/iot-my-world/brain/search/criterion"
	listText "github.com/iot-my-world/brain/search/criterion/list/text"
	"github.com/iot-my-world/brain/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/search/identifier/id"
	"github.com/iot-my-world/brain/search/identifier/username"
	humanUserLogin "github.com/iot-my-world/brain/security/claims/login/user/human"
	"github.com/iot-my-world/brain/security/claims/registerClientAdminUser"
	"github.com/iot-my-world/brain/security/claims/registerClientUser"
	"github.com/iot-my-world/brain/security/claims/registerCompanyAdminUser"
	"github.com/iot-my-world/brain/security/claims/registerCompanyUser"
	roleSetup "github.com/iot-my-world/brain/security/role/setup"
	"github.com/iot-my-world/brain/security/token"
	userAdministrator "github.com/iot-my-world/brain/user/human/administrator"
	userRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler"
	userRecordHandlerException "github.com/iot-my-world/brain/user/human/recordHandler/exception"
	userValidator "github.com/iot-my-world/brain/user/human/validator"
	"time"
)

type registrar struct {
	companyRecordHandler       companyRecordHandler.RecordHandler
	userRecordHandler          userRecordHandler.RecordHandler
	userValidator              userValidator.Validator
	userAdministrator          userAdministrator.Administrator
	clientRecordHandler        *clientRecordHandler.RecordHandler
	mailer                     mailer.Mailer
	jwtGenerator               token.JWTGenerator
	mailRedirectBaseUrl        string
	systemClaims               *humanUserLogin.Login
	registrationEmailGenerator emailGenerator.Generator
}

func New(
	companyRecordHandler companyRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
	userValidator userValidator.Validator,
	userAdministrator userAdministrator.Administrator,
	clientRecordHandler *clientRecordHandler.RecordHandler,
	mailer mailer.Mailer,
	rsaPrivateKey *rsa.PrivateKey,
	mailRedirectBaseUrl string,
	systemClaims *humanUserLogin.Login,
	registrationEmailGenerator emailGenerator.Generator,
) partyRegistrar.Registrar {
	return &registrar{
		companyRecordHandler:       companyRecordHandler,
		userRecordHandler:          userRecordHandler,
		userValidator:              userValidator,
		userAdministrator:          userAdministrator,
		clientRecordHandler:        clientRecordHandler,
		mailer:                     mailer,
		jwtGenerator:               token.NewJWTGenerator(rsaPrivateKey),
		mailRedirectBaseUrl:        mailRedirectBaseUrl,
		systemClaims:               systemClaims,
		registrationEmailGenerator: registrationEmailGenerator,
	}
}

func (r *registrar) RegisterSystemAdminUser(request *partyRegistrar.RegisterSystemAdminUserRequest) (*partyRegistrar.RegisterSystemAdminUserResponse, error) {

	// check if the system admin user already exists (i.e. has already been registered)
	userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: username.Identifier{Username: request.User.Username},
	})
	switch err.(type) {
	case nil:
		// this means that the user already exists
		return &partyRegistrar.RegisterSystemAdminUserResponse{User: userRetrieveResponse.User}, partyRegistrarException.AlreadyRegistered{}
	case userRecordHandlerException.NotFound:
		// this is fine, we will be creating the user now
	default:
		err = partyRegistrarException.RegisterSystemAdminUser{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// create the user
	userCreateResponse, err := r.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		User: request.User,
	})
	if err != nil {
		err = partyRegistrarException.RegisterSystemAdminUser{Reasons: []string{"user creation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	_, err = r.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: userCreateResponse.User.Id},
		NewPassword: string(request.User.Password),
	})
	if err != nil {
		err = partyRegistrarException.RegisterSystemAdminUser{Reasons: []string{"setting password", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &partyRegistrar.RegisterSystemAdminUserResponse{User: userCreateResponse.User}, nil
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

func (r *registrar) InviteCompanyAdminUser(request *partyRegistrar.InviteCompanyAdminUserRequest) (*partyRegistrar.InviteCompanyAdminUserResponse, error) {
	if err := r.ValidateInviteCompanyAdminUserRequest(request); err != nil {
		return nil, err
	}

	// Retrieve the company party
	companyRetrieveResponse, err := r.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.CompanyIdentifier,
	})
	if err != nil {
		err = partyRegistrarException.InviteCompanyAdminUser{Reasons: []string{"company retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// Retrieve the minimal company admin user which was created on company creation
	userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims: request.Claims,
		Identifier: emailAddress.Identifier{
			EmailAddress: companyRetrieveResponse.Company.AdminEmailAddress,
		},
	})
	if err != nil {
		err = partyRegistrarException.InviteCompanyAdminUser{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		err = partyRegistrarException.AlreadyRegistered{}
		log.Error(err.Error())
		return nil, err
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
		err = partyRegistrarException.InviteCompanyAdminUser{Reasons: []string{"token generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	urlToken := fmt.Sprintf("%s/register?&t=%s", r.mailRedirectBaseUrl, registrationToken)

	userRetrieveResponse.User.Name = fmt.Sprintf("%s Administrator", companyRetrieveResponse.Company.Name)

	generateEmailResponse, err := r.registrationEmailGenerator.Generate(&emailGenerator.GenerateRequest{
		Data: registrationEmail.Data{
			URLToken: urlToken,
			User:     userRetrieveResponse.User,
		},
	})
	if err != nil {
		err = partyRegistrarException.InviteCompanyAdminUser{Reasons: []string{"email generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	if _, err := r.mailer.Send(&mailer.SendRequest{
		Email: generateEmailResponse.Email,
	}); err != nil {
		err = partyRegistrarException.InviteCompanyAdminUser{Reasons: []string{"email sending", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &partyRegistrar.InviteCompanyAdminUserResponse{URLToken: urlToken}, nil
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
		userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: id.Identifier{Id: request.User.Id},
		})
		if err == nil {
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
	userValidateResponse, err := r.userValidator.Validate(&userValidator.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *r.systemClaims,
		User:   request.User,
		Action: partyRegistrarAction.RegisterCompanyAdminUser,
	})
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate newAdminUser")
	} else {
		for _, reason := range userValidateResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterCompanyAdminUser(request *partyRegistrar.RegisterCompanyAdminUserRequest) (*partyRegistrar.RegisterCompanyAdminUserResponse, error) {
	if err := r.ValidateRegisterCompanyAdminUserRequest(request); err != nil {
		return nil, err
	}

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.CompanyAdmin.Name, roleSetup.CompanyUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	_, err := r.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	})
	if err != nil {
		return nil, err
	}

	// change the users password
	userChangePasswordResponse, err := r.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	})
	if err != nil {
		return nil, err
	}

	return &partyRegistrar.RegisterCompanyAdminUserResponse{User: userChangePasswordResponse.User}, nil
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

func (r *registrar) InviteCompanyUser(request *partyRegistrar.InviteCompanyUserRequest) (*partyRegistrar.InviteCompanyUserResponse, error) {
	if err := r.ValidateInviteCompanyUserRequest(request); err != nil {
		return nil, err
	}

	// retrieve the user
	userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	})
	if err != nil {
		return nil, partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"user retrieval", err.Error()}}
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		return nil, partyRegistrarException.AlreadyRegistered{}
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
		return nil, partyRegistrarException.TokenGeneration{Reasons: []string{"inviteCompanyUser", err.Error()}}
	}

	// e.g. //http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
	urlToken := fmt.Sprintf("%s/register?&t=%s", r.mailRedirectBaseUrl, registrationToken)

	generateEmailResponse, err := r.registrationEmailGenerator.Generate(&emailGenerator.GenerateRequest{
		Data: registrationEmail.Data{
			URLToken: urlToken,
			User:     userRetrieveResponse.User,
		},
	})
	if err != nil {
		return nil, partyRegistrarException.EmailGeneration{Reasons: []string{"invite company admin user", err.Error()}}
	}

	if _, err := r.mailer.Send(&mailer.SendRequest{
		Email: generateEmailResponse.Email,
	}); err != nil {
		return nil, err
	}

	return &partyRegistrar.InviteCompanyUserResponse{URLToken: urlToken}, nil
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
		userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: id.Identifier{Id: request.User.Id},
		})
		if err == nil {
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
	userValidateResponse, err := r.userValidator.Validate(&userValidator.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *r.systemClaims,
		User:   request.User,
		Action: partyRegistrarAction.RegisterCompanyUser,
	})
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate new user")
	} else {
		for _, reason := range userValidateResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterCompanyUser(request *partyRegistrar.RegisterCompanyUserRequest) (*partyRegistrar.RegisterCompanyUserResponse, error) {
	if err := r.ValidateRegisterCompanyUserRequest(request); err != nil {
		return nil, err
	}

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.CompanyUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	_, err := r.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	})
	if err != nil {
		return nil, err
	}

	// change the users password
	userChangePasswordResponse, err := r.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	})
	if err != nil {
		return nil, err
	}

	return &partyRegistrar.RegisterCompanyUserResponse{User: userChangePasswordResponse.User}, nil
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

func (r *registrar) InviteClientAdminUser(request *partyRegistrar.InviteClientAdminUserRequest) (*partyRegistrar.InviteClientAdminUserResponse, error) {
	if err := r.ValidateInviteClientAdminUserRequest(request); err != nil {
		return nil, err
	}

	// retrieve the client
	clientRetrieveResponse, err := r.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ClientIdentifier,
	})
	if err != nil {
		return nil, partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"client", err.Error()}}
	}

	// retrieve the minimal client admin user
	userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		// we use system claims as users can typically only be retrieved by a user of the same party
		Claims: *r.systemClaims,
		Identifier: emailAddress.Identifier{
			EmailAddress: clientRetrieveResponse.Client.AdminEmailAddress,
		},
	})
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"user retrieval", err.Error()}}
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		return nil, partyRegistrarException.AlreadyRegistered{}
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
		return nil, partyRegistrarException.TokenGeneration{Reasons: []string{"inviteClientAdminUser", err.Error()}}
	}

	//http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
	urlToken := fmt.Sprintf("%s/register?&t=%s", r.mailRedirectBaseUrl, registrationToken)

	userRetrieveResponse.User.Name = fmt.Sprintf("%s Administrator", clientRetrieveResponse.Client.Name)

	generateEmailResponse, err := r.registrationEmailGenerator.Generate(&emailGenerator.GenerateRequest{
		Data: registrationEmail.Data{
			URLToken: urlToken,
			User:     userRetrieveResponse.User,
		},
	})
	if err != nil {
		return nil, partyRegistrarException.EmailGeneration{Reasons: []string{"invite company admin user", err.Error()}}
	}

	if _, err := r.mailer.Send(&mailer.SendRequest{
		Email: generateEmailResponse.Email,
	}); err != nil {
		return nil, err
	}

	return &partyRegistrar.InviteClientAdminUserResponse{URLToken: urlToken}, nil
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
		userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: id.Identifier{Id: request.User.Id},
		})
		if err == nil {
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
	userValidateResponse, err := r.userValidator.Validate(&userValidator.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *r.systemClaims,
		User:   request.User,
		Action: partyRegistrarAction.RegisterClientAdminUser,
	})
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate newAdminUser")
	} else {
		for _, reason := range userValidateResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterClientAdminUser(request *partyRegistrar.RegisterClientAdminUserRequest) (*partyRegistrar.RegisterClientAdminUserResponse, error) {
	if err := r.ValidateRegisterClientAdminUserRequest(request); err != nil {
		return nil, err
	}

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.ClientAdmin.Name, roleSetup.ClientUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	_, err := r.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	})
	if err != nil {
		return nil, err
	}

	// change the users password
	userChangePasswordResponse, err := r.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	})
	if err != nil {
		return nil, err
	}

	return &partyRegistrar.RegisterClientAdminUserResponse{User: userChangePasswordResponse.User}, nil
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
	}
	return nil
}

func (r *registrar) InviteClientUser(request *partyRegistrar.InviteClientUserRequest) (*partyRegistrar.InviteClientUserResponse, error) {
	if err := r.ValidateInviteClientUserRequest(request); err != nil {
		return nil, err
	}

	// retrieve the user
	userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	})
	if err != nil {
		return nil, partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"user retrieval", err.Error()}}
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		return nil, partyRegistrarException.AlreadyRegistered{}
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
		return nil, partyRegistrarException.TokenGeneration{Reasons: []string{"inviteClientUser", err.Error()}}
	}

	// e.g. //http://localhost:3000/register?&t=eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJ0eXBlIjoiUmVnaXN0cmF0aW9uIiwiZXhwIjoxNTUwMDM0NjYxLCJpYXQiOjE1NDk5NDgyNjIsImNvbnRleHQiOnsibmFtZSI6IkJvYidzIE93biBNYW4iLCJwYXJ0eUNvZGUiOiJCT0IiLCJwYXJ0eVR5cGUiOiJJTkRJVklEVUFMIn19.CrqxhOs_NSk1buXQyEykyCsPtNQCoWWFkxQ_HphgjSc2idchlov8SdlpdjYxtqaRv7zpDrPwKHaeR4inbcf0Xat1vasqXEPqgE5WzSWtt-GbXi5iUEc-pg79yx0zQ8riIeSkho84BRZbh252ePuOXBK1Yqa4MG9O2xblDOsfQgDVa-9Ha6XZvxHbNOFYKchiKfsclaZ_osQn9Ll6p8GAw9wqCStWp_kRSJM81RUc8rFIfxNgBwqoab_r6QhFHLT9jm90eU3RrVkGv_bB4hRcwhwE_0ksRL9lXRCIKs5ctuZkcYtPvhdKMRCaXPlV-Bm6sgx4qpS-nzmOmc0bNCrOZlP0JUAHdKSBHmw9mSw5QRLkVTPgAuAm9qOj5PjU95DiFLY1q9X0pyRL2uG7xiE8F-Q_g_5q0vXLZkvgwcEpc604ZGgMsH3Sw5mCl0aKsF6c7eiKjTCBkSv46hDqED4cP4KBrxhEgNN_oKrYPqjElZ0xrFe7P3fAyt1jh3SqgaYoZQB4ORJ76CByLhTRAtTmX2SnVQJhMwgtZu9kPXtpKTfdyAUZcd4eUmfLpJ1VXCzvFlIXQW9rN1TgsE2eMqSbmOtgwHQqQD52M-CW8w7CLBfWG7-GQ68GUA42IErMVKlL9mp22LbOkzvpiFEOx5V0cXyVzndPDKNPZ278gwablyU
	urlToken := fmt.Sprintf("%s/register?&t=%s", r.mailRedirectBaseUrl, registrationToken)

	generateEmailResponse, err := r.registrationEmailGenerator.Generate(&emailGenerator.GenerateRequest{
		Data: registrationEmail.Data{
			URLToken: urlToken,
			User:     userRetrieveResponse.User,
		},
	})
	if err != nil {
		return nil, partyRegistrarException.EmailGeneration{Reasons: []string{"invite company admin user", err.Error()}}
	}

	if _, err := r.mailer.Send(&mailer.SendRequest{
		Email: generateEmailResponse.Email,
	}); err != nil {
		return nil, err
	}

	return &partyRegistrar.InviteClientUserResponse{URLToken: urlToken}, nil
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
		userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: id.Identifier{Id: request.User.Id},
		})
		if err == nil {
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
	userValidateResponse, err := r.userValidator.Validate(&userValidator.ValidateRequest{
		// system claims since we want all users to be visible for the email address check done in validate user
		Claims: *r.systemClaims,
		User:   request.User,
		Action: partyRegistrarAction.RegisterClientUser,
	})
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate new user")
	} else {
		for _, reason := range userValidateResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) RegisterClientUser(request *partyRegistrar.RegisterClientUserRequest) (*partyRegistrar.RegisterClientUserResponse, error) {
	if err := r.ValidateRegisterClientUserRequest(request); err != nil {
		return nil, err
	}

	// retrieve the minimal user

	// give the user the necessary roles
	request.User.Roles = []string{roleSetup.ClientUser.Name}

	// set the user to registered
	request.User.Registered = true

	// update the user
	_, err := r.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		User:       request.User,
		Identifier: id.Identifier{Id: request.User.Id},
	})
	if err != nil {
		return nil, err
	}

	// change the users password
	userChangePasswordResponse, err := r.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	})
	if err != nil {
		return nil, err
	}

	return &partyRegistrar.RegisterClientUserResponse{User: userChangePasswordResponse.User}, nil
}

func (r *registrar) ValidateAreAdminsRegisteredRequest(request *partyRegistrar.AreAdminsRegisteredRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
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

func (r *registrar) InviteUser(request *partyRegistrar.InviteUserRequest) (*partyRegistrar.InviteUserResponse, error) {
	if err := r.ValidateInviteUserRequest(request); err != nil {
		return nil, err
	}

	// retrieve the user
	userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	})
	if err != nil {
		return nil, partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"user retrieval", err.Error()}}
	}

	response := partyRegistrar.InviteUserResponse{}

	// the purpose of this service is to provide a generic way to invite a user from any party, admin user or not
	switch userRetrieveResponse.User.PartyType {
	case party.Company:
		// determine it this is the admin user
		companyRetrieveResponse, err := r.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Claims:     *r.systemClaims,
			Identifier: userRetrieveResponse.User.PartyId,
		})
		if err != nil {
			return nil, partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"company", err.Error()}}
		}
		if userRetrieveResponse.User.EmailAddress == companyRetrieveResponse.Company.AdminEmailAddress {
			inviteCompanyAdminUserResponse, err := r.InviteCompanyAdminUser(&partyRegistrar.InviteCompanyAdminUserRequest{
				Claims:            request.Claims,
				CompanyIdentifier: userRetrieveResponse.User.PartyId,
			})
			if err != nil {
				return nil, err
			}
			response.URLToken = inviteCompanyAdminUserResponse.URLToken
		} else {
			inviteCompanyUserResponse, err := r.InviteCompanyUser(&partyRegistrar.InviteCompanyUserRequest{
				Claims:         request.Claims,
				UserIdentifier: id.Identifier{Id: userRetrieveResponse.User.Id},
			})
			if err != nil {
				return nil, err
			}
			response.URLToken = inviteCompanyUserResponse.URLToken
		}

	case party.Client:
		// determine it this is the admin user
		clientRetrieveResponse, err := r.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Claims:     *r.systemClaims,
			Identifier: userRetrieveResponse.User.PartyId,
		})
		if err != nil {
			return nil, partyRegistrarException.UnableToRetrieveParty{Reasons: []string{"company", err.Error()}}
		}
		if userRetrieveResponse.User.EmailAddress == clientRetrieveResponse.Client.AdminEmailAddress {
			inviteClientAdminUserResponse, err := r.InviteClientAdminUser(&partyRegistrar.InviteClientAdminUserRequest{
				Claims:           request.Claims,
				ClientIdentifier: userRetrieveResponse.User.PartyId,
			})
			if err != nil {
				return nil, err
			}
			response.URLToken = inviteClientAdminUserResponse.URLToken
		} else {
			inviteClientUserResponse, err := r.InviteClientUser(&partyRegistrar.InviteClientUserRequest{
				Claims:         request.Claims,
				UserIdentifier: id.Identifier{Id: userRetrieveResponse.User.Id},
			})
			if err != nil {
				return nil, err
			}
			response.URLToken = inviteClientUserResponse.URLToken
		}

	default:
		return nil, partyRegistrarException.PartyTypeInvalid{Reasons: []string{string(userRetrieveResponse.User.PartyType)}}

	}

	return &response, nil
}

func (r *registrar) AreAdminsRegistered(request *partyRegistrar.AreAdminsRegisteredRequest) (*partyRegistrar.AreAdminsRegisteredResponse, error) {
	if err := r.ValidateAreAdminsRegisteredRequest(request); err != nil {
		return nil, err
	}

	companyIds := make([]string, 0)
	companyAdminEmails := make([]string, 0)
	clientIds := make([]string, 0)
	clientAdminEmails := make([]string, 0)

	response := partyRegistrar.AreAdminsRegisteredResponse{
		Result: make(map[string]bool),
	}

	// compose id lists for exact list criteria
	for _, partyIdentifier := range request.PartyIdentifiers {
		switch partyIdentifier.PartyType {
		case party.System:
			response.Result[partyIdentifier.PartyIdIdentifier.Id] = true
		case party.Company:
			companyIds = append(companyIds, partyIdentifier.PartyIdIdentifier.Id)
		case party.Client:
			clientIds = append(clientIds, partyIdentifier.PartyIdIdentifier.Id)
		default:
			return nil, partyRegistrarException.PartyTypeInvalid{Reasons: []string{"areAdminsRegistered", string(partyIdentifier.PartyType)}}
		}
	}

	// collect companies in request
	companyCollectResponse, err := r.companyRecordHandler.Collect(&companyRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "id",
				List:  companyIds,
			},
		},
	})
	if err != nil {
		return nil, partyRegistrarException.UnableToCollectParties{Reasons: []string{"company", err.Error()}}
	} else {
		// confirm that for every id received a company was returned
		if len(companyCollectResponse.Records) != len(companyIds) {
			return nil, brainException.Unexpected{Reasons: []string{
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
	companyAdminUserCollectResponse, err := r.userRecordHandler.Collect(&userRecordHandler.CollectRequest{
		// use system claims as usually users can only be retrieved by a user of the same party
		Claims: *r.systemClaims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "emailAddress",
				List:  companyAdminEmails,
			},
		},
	})
	if err != nil {
		return nil, partyRegistrarException.UnableToCollectParties{Reasons: []string{"companyAdminUsers", err.Error()}}
	} else {
		// confirm that for every admin email a user was returned
		if len(companyAdminUserCollectResponse.Records) != len(companyAdminEmails) {
			return nil, brainException.Unexpected{Reasons: []string{
				"no company admin users found different from number of admin emails found",
				fmt.Sprintf("%d vs %d", len(companyAdminUserCollectResponse.Records), len(companyAdminEmails)),
			}}
		}
	}
	// update result for the company admin users retrieved
	for companyAdminUserIdx := range companyAdminUserCollectResponse.Records {
		response.Result[companyAdminUserCollectResponse.Records[companyAdminUserIdx].PartyId.Id] =
			companyAdminUserCollectResponse.Records[companyAdminUserIdx].Registered
	}

	// collect clients in request
	clientCollectResponse, err := r.clientRecordHandler.Collect(&clientRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "id",
				List:  clientIds,
			},
		},
	})
	if err != nil {
		return nil, partyRegistrarException.UnableToCollectParties{Reasons: []string{"client", err.Error()}}
	} else {
		// confirm that for every id received a client was returned
		if len(clientCollectResponse.Records) != len(clientIds) {
			return nil, brainException.Unexpected{Reasons: []string{
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
	clientAdminUserCollectResponse, err := r.userRecordHandler.Collect(&userRecordHandler.CollectRequest{
		Claims: *r.systemClaims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "emailAddress",
				List:  clientAdminEmails,
			},
		},
	})
	if err != nil {
		return nil, partyRegistrarException.UnableToCollectParties{Reasons: []string{"clientAdminUsers", err.Error()}}
	} else {
		// confirm that for every admin email user was returned
		if len(clientAdminUserCollectResponse.Records) != len(clientAdminEmails) {
			return nil, brainException.Unexpected{Reasons: []string{
				"no client admin users found different from number of admin emails found",
				fmt.Sprintf("%d vs %d", len(clientAdminUserCollectResponse.Records), len(clientAdminEmails)),
			}}
		}
	}
	// update result for the client admin users retrieved
	for clientAdminUserIdx := range clientAdminUserCollectResponse.Records {
		response.Result[clientAdminUserCollectResponse.Records[clientAdminUserIdx].PartyId.Id] =
			clientAdminUserCollectResponse.Records[clientAdminUserIdx].Registered
	}

	return &response, nil
}
