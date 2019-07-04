package basic

import (
	"crypto/rsa"
	"fmt"
	"github.com/iot-my-world/brain/internal/environment"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	emailGenerator "github.com/iot-my-world/brain/pkg/communication/email/generator"
	registrationEmail "github.com/iot-my-world/brain/pkg/communication/email/generator/registration"
	"github.com/iot-my-world/brain/pkg/communication/email/mailer"
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	recordHandler2 "github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	registrar2 "github.com/iot-my-world/brain/pkg/party/registrar"
	"github.com/iot-my-world/brain/pkg/party/registrar/action"
	"github.com/iot-my-world/brain/pkg/party/registrar/exception"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	listText "github.com/iot-my-world/brain/pkg/search/criterion/list/text"
	"github.com/iot-my-world/brain/pkg/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/search/identifier/username"
	humanUserLogin "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	"github.com/iot-my-world/brain/pkg/security/claims/registerClientAdminUser"
	"github.com/iot-my-world/brain/pkg/security/claims/registerClientUser"
	"github.com/iot-my-world/brain/pkg/security/claims/registerCompanyAdminUser"
	"github.com/iot-my-world/brain/pkg/security/claims/registerCompanyUser"
	roleSetup "github.com/iot-my-world/brain/pkg/security/role/setup"
	"github.com/iot-my-world/brain/pkg/security/token"
	userAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator"
	userRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	userRecordHandlerException "github.com/iot-my-world/brain/pkg/user/human/recordHandler/exception"
	userValidator "github.com/iot-my-world/brain/pkg/user/human/validator"
	"time"
)

type registrar struct {
	companyRecordHandler       recordHandler2.RecordHandler
	userRecordHandler          userRecordHandler.RecordHandler
	userValidator              userValidator.Validator
	userAdministrator          userAdministrator.Administrator
	clientRecordHandler        recordHandler.RecordHandler
	mailer                     mailer.Mailer
	jwtGenerator               token.JWTGenerator
	mailRedirectBaseUrl        string
	systemClaims               *humanUserLogin.Login
	registrationEmailGenerator emailGenerator.Generator
	environmentType            environment.Type
}

func New(
	companyRecordHandler recordHandler2.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
	userValidator userValidator.Validator,
	userAdministrator userAdministrator.Administrator,
	clientRecordHandler recordHandler.RecordHandler,
	mailer mailer.Mailer,
	rsaPrivateKey *rsa.PrivateKey,
	mailRedirectBaseUrl string,
	systemClaims *humanUserLogin.Login,
	registrationEmailGenerator emailGenerator.Generator,
	environmentType environment.Type,
) registrar2.Registrar {
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
		environmentType:            environmentType,
	}
}

func (r *registrar) RegisterSystemAdminUser(request *registrar2.RegisterSystemAdminUserRequest) (*registrar2.RegisterSystemAdminUserResponse, error) {

	// check if the system admin user already exists (i.e. has already been registered)
	userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: username.Identifier{Username: request.User.Username},
	})
	switch err.(type) {
	case nil:
		// this means that the user already exists
		return &registrar2.RegisterSystemAdminUserResponse{User: userRetrieveResponse.User}, exception.AlreadyRegistered{}
	case userRecordHandlerException.NotFound:
		// this is fine, we will be creating the user now
	default:
		err = exception.RegisterSystemAdminUser{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// create the user
	userCreateResponse, err := r.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		User: request.User,
	})
	if err != nil {
		err = exception.RegisterSystemAdminUser{Reasons: []string{"user creation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	_, err = r.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: userCreateResponse.User.Id},
		NewPassword: string(request.User.Password),
	})
	if err != nil {
		err = exception.RegisterSystemAdminUser{Reasons: []string{"setting password", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.RegisterSystemAdminUserResponse{User: userCreateResponse.User}, nil
}

func (r *registrar) ValidateInviteCompanyAdminUserRequest(request *registrar2.InviteCompanyAdminUserRequest) error {
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

func (r *registrar) InviteCompanyAdminUser(request *registrar2.InviteCompanyAdminUserRequest) (*registrar2.InviteCompanyAdminUserResponse, error) {
	if err := r.ValidateInviteCompanyAdminUserRequest(request); err != nil {
		return nil, err
	}

	// Retrieve the company party
	companyRetrieveResponse, err := r.companyRecordHandler.Retrieve(&recordHandler2.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.CompanyIdentifier,
	})
	if err != nil {
		err = exception.InviteCompanyAdminUser{Reasons: []string{"company retrieval", err.Error()}}
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
		err = exception.InviteCompanyAdminUser{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		err = exception.AlreadyRegistered{}
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
		err = exception.InviteCompanyAdminUser{Reasons: []string{"token generation", err.Error()}}
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
		err = exception.InviteCompanyAdminUser{Reasons: []string{"email generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	if r.environmentType == environment.Development {
		// if this is the development environment return response with token
		return &registrar2.InviteCompanyAdminUserResponse{URLToken: urlToken}, nil
	}

	// otherwise send email and return response without token
	if _, err := r.mailer.Send(&mailer.SendRequest{
		Email: generateEmailResponse.Email,
	}); err != nil {
		err = exception.InviteCompanyAdminUser{Reasons: []string{"email sending", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.InviteCompanyAdminUserResponse{}, nil
}

func (r *registrar) ValidateRegisterCompanyAdminUserRequest(request *registrar2.RegisterCompanyAdminUserRequest) error {
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
				return exception.AlreadyRegistered{}
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
		Action: action.RegisterCompanyAdminUser,
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

func (r *registrar) RegisterCompanyAdminUser(request *registrar2.RegisterCompanyAdminUserRequest) (*registrar2.RegisterCompanyAdminUserResponse, error) {
	if err := r.ValidateRegisterCompanyAdminUserRequest(request); err != nil {
		log.Error(err.Error())
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
		err = exception.RegisterCompanyAdminUser{Reasons: []string{"user update", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// change the users password
	if _, err := r.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	}); err != nil {
		err = exception.RegisterCompanyAdminUser{Reasons: []string{"user password change", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.RegisterCompanyAdminUserResponse{User: request.User}, nil
}

func (r *registrar) ValidateInviteCompanyUserRequest(request *registrar2.InviteCompanyUserRequest) error {
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

func (r *registrar) InviteCompanyUser(request *registrar2.InviteCompanyUserRequest) (*registrar2.InviteCompanyUserResponse, error) {
	if err := r.ValidateInviteCompanyUserRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// retrieve the user
	userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	})
	if err != nil {
		err = exception.InviteCompanyUser{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		err = exception.AlreadyRegistered{}
		log.Error(err.Error())
		return nil, err
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
		err = exception.InviteCompanyUser{Reasons: []string{"token generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
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
		err = exception.InviteCompanyUser{Reasons: []string{"email generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	if r.environmentType == environment.Development {
		// if this is the development environment return response with token
		return &registrar2.InviteCompanyUserResponse{URLToken: urlToken}, nil
	}

	// otherwise send email and return response without token
	if _, err := r.mailer.Send(&mailer.SendRequest{
		Email: generateEmailResponse.Email,
	}); err != nil {
		err = exception.InviteCompanyUser{Reasons: []string{"email sending", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.InviteCompanyUserResponse{}, nil
}

func (r *registrar) ValidateRegisterCompanyUserRequest(request *registrar2.RegisterCompanyUserRequest) error {
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
				return exception.AlreadyRegistered{}
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
		Action: action.RegisterCompanyUser,
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

func (r *registrar) RegisterCompanyUser(request *registrar2.RegisterCompanyUserRequest) (*registrar2.RegisterCompanyUserResponse, error) {
	if err := r.ValidateRegisterCompanyUserRequest(request); err != nil {
		log.Error(err.Error())
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
		err = exception.RegisterCompanyUser{Reasons: []string{"user update", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// change the users password
	if _, err := r.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	}); err != nil {
		err = exception.RegisterCompanyUser{Reasons: []string{"setting user password", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.RegisterCompanyUserResponse{User: request.User}, nil
}

func (r *registrar) ValidateInviteClientAdminUserRequest(request *registrar2.InviteClientAdminUserRequest) error {
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

func (r *registrar) InviteClientAdminUser(request *registrar2.InviteClientAdminUserRequest) (*registrar2.InviteClientAdminUserResponse, error) {
	if err := r.ValidateInviteClientAdminUserRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// retrieve the client
	clientRetrieveResponse, err := r.clientRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ClientIdentifier,
	})
	if err != nil {
		err = exception.InviteClientAdminUser{Reasons: []string{"client party retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
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
		err = exception.InviteClientAdminUser{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		err = exception.AlreadyRegistered{}
		log.Error(err.Error())
		return nil, err
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
		err = exception.InviteClientAdminUser{Reasons: []string{"token generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
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
		err = exception.InviteClientAdminUser{Reasons: []string{"email generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	if r.environmentType == environment.Development {
		// if this is the development environment return response with token
		return &registrar2.InviteClientAdminUserResponse{URLToken: urlToken}, nil
	}

	// otherwise send email and return response without token
	if _, err := r.mailer.Send(&mailer.SendRequest{
		Email: generateEmailResponse.Email,
	}); err != nil {
		err = exception.InviteClientAdminUser{Reasons: []string{"email sending", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.InviteClientAdminUserResponse{}, nil
}

func (r *registrar) ValidateRegisterClientAdminUserRequest(request *registrar2.RegisterClientAdminUserRequest) error {
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
				return exception.AlreadyRegistered{}
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
		Action: action.RegisterClientAdminUser,
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

func (r *registrar) RegisterClientAdminUser(request *registrar2.RegisterClientAdminUserRequest) (*registrar2.RegisterClientAdminUserResponse, error) {
	if err := r.ValidateRegisterClientAdminUserRequest(request); err != nil {
		log.Error(err.Error())
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
		err = exception.RegisterClientAdminUser{Reasons: []string{"user update", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// change the users password
	if _, err := r.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	}); err != nil {
		err = exception.RegisterClientAdminUser{Reasons: []string{"user password setting", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.RegisterClientAdminUserResponse{User: request.User}, nil
}

func (r *registrar) ValidateInviteClientUserRequest(request *registrar2.InviteClientUserRequest) error {
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

func (r *registrar) InviteClientUser(request *registrar2.InviteClientUserRequest) (*registrar2.InviteClientUserResponse, error) {
	if err := r.ValidateInviteClientUserRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// retrieve the user
	userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	})
	if err != nil {
		err = exception.InviteClientUser{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// if the user is already registered, return an error
	if userRetrieveResponse.User.Registered {
		err = exception.AlreadyRegistered{}
		log.Error(err.Error())
		return nil, err
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
		err = exception.InviteClientUser{Reasons: []string{"token generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
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
		err = exception.InviteClientUser{Reasons: []string{"email generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	if r.environmentType == environment.Development {
		// if this is the development environment return response with token
		return &registrar2.InviteClientUserResponse{URLToken: urlToken}, nil
	}

	// otherwise send email and return response without token
	if _, err := r.mailer.Send(&mailer.SendRequest{
		Email: generateEmailResponse.Email,
	}); err != nil {
		err = exception.InviteClientUser{Reasons: []string{"email sending", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.InviteClientUserResponse{}, nil
}

func (r *registrar) ValidateRegisterClientUserRequest(request *registrar2.RegisterClientUserRequest) error {
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
				return exception.AlreadyRegistered{}
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
		Action: action.RegisterClientUser,
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

func (r *registrar) RegisterClientUser(request *registrar2.RegisterClientUserRequest) (*registrar2.RegisterClientUserResponse, error) {
	if err := r.ValidateRegisterClientUserRequest(request); err != nil {
		log.Error(err.Error())
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
		err = exception.RegisterClientUser{Reasons: []string{"user update", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// change the users password
	if _, err := r.userAdministrator.SetPassword(&userAdministrator.SetPasswordRequest{
		Claims:      request.Claims,
		Identifier:  id.Identifier{Id: request.User.Id},
		NewPassword: string(request.User.Password),
	}); err != nil {
		err = exception.RegisterClientUser{Reasons: []string{"user password setting", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &registrar2.RegisterClientUserResponse{User: request.User}, nil
}

func (r *registrar) ValidateAreAdminsRegisteredRequest(request *registrar2.AreAdminsRegisteredRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *registrar) ValidateInviteUserRequest(request *registrar2.InviteUserRequest) error {
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

func (r *registrar) InviteUser(request *registrar2.InviteUserRequest) (*registrar2.InviteUserResponse, error) {
	if err := r.ValidateInviteUserRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// retrieve the user
	userRetrieveResponse, err := r.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	})
	if err != nil {
		err = exception.InviteUser{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	response := registrar2.InviteUserResponse{}

	// the purpose of this service is to provide a generic way to invite a user from any party, admin user or not
	switch userRetrieveResponse.User.PartyType {
	case party.Company:
		// determine it this is the admin user
		companyRetrieveResponse, err := r.companyRecordHandler.Retrieve(&recordHandler2.RetrieveRequest{
			Claims:     *r.systemClaims,
			Identifier: userRetrieveResponse.User.PartyId,
		})
		if err != nil {
			err = exception.InviteUser{Reasons: []string{"company entity retrieval", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
		if userRetrieveResponse.User.EmailAddress == companyRetrieveResponse.Company.AdminEmailAddress {
			inviteCompanyAdminUserResponse, err := r.InviteCompanyAdminUser(&registrar2.InviteCompanyAdminUserRequest{
				Claims:            request.Claims,
				CompanyIdentifier: userRetrieveResponse.User.PartyId,
			})
			if err != nil {
				err = exception.InviteUser{Reasons: []string{"invite company admin user", err.Error()}}
				log.Error(err.Error())
				return nil, err
			}
			response.URLToken = inviteCompanyAdminUserResponse.URLToken
		} else {
			inviteCompanyUserResponse, err := r.InviteCompanyUser(&registrar2.InviteCompanyUserRequest{
				Claims:         request.Claims,
				UserIdentifier: id.Identifier{Id: userRetrieveResponse.User.Id},
			})
			if err != nil {
				err = exception.InviteUser{Reasons: []string{"invite company user", err.Error()}}
				log.Error(err.Error())
				return nil, err
			}
			response.URLToken = inviteCompanyUserResponse.URLToken
		}

	case party.Client:
		// determine it this is the admin user
		clientRetrieveResponse, err := r.clientRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
			Claims:     *r.systemClaims,
			Identifier: userRetrieveResponse.User.PartyId,
		})
		if err != nil {
			err = exception.InviteUser{Reasons: []string{"retrieving client entity", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
		if userRetrieveResponse.User.EmailAddress == clientRetrieveResponse.Client.AdminEmailAddress {
			inviteClientAdminUserResponse, err := r.InviteClientAdminUser(&registrar2.InviteClientAdminUserRequest{
				Claims:           request.Claims,
				ClientIdentifier: userRetrieveResponse.User.PartyId,
			})
			if err != nil {
				err = exception.InviteUser{Reasons: []string{"inviting client admin user", err.Error()}}
				log.Error(err.Error())
				return nil, err
			}
			response.URLToken = inviteClientAdminUserResponse.URLToken
		} else {
			inviteClientUserResponse, err := r.InviteClientUser(&registrar2.InviteClientUserRequest{
				Claims:         request.Claims,
				UserIdentifier: id.Identifier{Id: userRetrieveResponse.User.Id},
			})
			if err != nil {
				err = exception.InviteUser{Reasons: []string{"inviting client user", err.Error()}}
				log.Error(err.Error())
				return nil, err
			}
			response.URLToken = inviteClientUserResponse.URLToken
		}

	default:
		err = exception.InviteUser{Reasons: []string{"invalid party type", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &response, nil
}

func (r *registrar) AreAdminsRegistered(request *registrar2.AreAdminsRegisteredRequest) (*registrar2.AreAdminsRegisteredResponse, error) {
	if err := r.ValidateAreAdminsRegisteredRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	companyIds := make([]string, 0)
	companyAdminEmails := make([]string, 0)
	clientIds := make([]string, 0)
	clientAdminEmails := make([]string, 0)

	response := registrar2.AreAdminsRegisteredResponse{
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
			err := exception.AreAdminsRegistered{Reasons: []string{"invalid party type", string(partyIdentifier.PartyType)}}
			log.Error(err.Error())
			return nil, err
		}
	}

	// collect companies in request
	companyCollectResponse, err := r.companyRecordHandler.Collect(&recordHandler2.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "id",
				List:  companyIds,
			},
		},
	})
	if err != nil {
		err = exception.AreAdminsRegistered{Reasons: []string{"collecting company parties"}}
		log.Error(err.Error())
		return nil, err
	} else {
		// confirm that for every id received a company was returned
		if len(companyCollectResponse.Records) != len(companyIds) {
			err := exception.AreAdminsRegistered{Reasons: []string{
				"no company records collected different to number of ids given",
				fmt.Sprintf("%d vs %d", len(companyCollectResponse.Records), len(companyIds)),
			}}
			log.Error(err.Error())
			return nil, err
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
		err = exception.AreAdminsRegistered{Reasons: []string{"collecting company admin users"}}
		log.Error(err.Error())
		return nil, err
	} else {
		// confirm that for every admin email a user was returned
		if len(companyAdminUserCollectResponse.Records) != len(companyAdminEmails) {
			err = exception.AreAdminsRegistered{Reasons: []string{
				"no company admin users found different from number of admin emails found",
				fmt.Sprintf("%d vs %d", len(companyAdminUserCollectResponse.Records), len(companyAdminEmails)),
			}}
			log.Error(err.Error())
			return nil, err
		}
	}
	// update result for the company admin users retrieved
	for companyAdminUserIdx := range companyAdminUserCollectResponse.Records {
		response.Result[companyAdminUserCollectResponse.Records[companyAdminUserIdx].PartyId.Id] =
			companyAdminUserCollectResponse.Records[companyAdminUserIdx].Registered
	}

	// collect clients in request
	clientCollectResponse, err := r.clientRecordHandler.Collect(&recordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			listText.Criterion{
				Field: "id",
				List:  clientIds,
			},
		},
	})
	if err != nil {
		err = exception.AreAdminsRegistered{Reasons: []string{"collecting client party entities"}}
		log.Error(err.Error())
		return nil, err
	} else {
		// confirm that for every id received a client was returned
		if len(clientCollectResponse.Records) != len(clientIds) {
			err = exception.AreAdminsRegistered{Reasons: []string{
				"no client records returned different to number of ids given",
				fmt.Sprintf("%d vs %d", len(clientCollectResponse.Records), len(clientIds)),
			}}
			log.Error(err.Error())
			return nil, err
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
		err = exception.AreAdminsRegistered{Reasons: []string{"collecting client admin users"}}
		log.Error(err.Error())
		return nil, err
	} else {
		// confirm that for every admin email user was returned
		if len(clientAdminUserCollectResponse.Records) != len(clientAdminEmails) {
			err = exception.AreAdminsRegistered{Reasons: []string{
				"no client admin users found different from number of admin emails found",
				fmt.Sprintf("%d vs %d", len(clientAdminUserCollectResponse.Records), len(clientAdminEmails)),
			}}
			log.Error(err.Error())
			return nil, err
		}
	}
	// update result for the client admin users retrieved
	for clientAdminUserIdx := range clientAdminUserCollectResponse.Records {
		response.Result[clientAdminUserCollectResponse.Records[clientAdminUserIdx].PartyId.Id] =
			clientAdminUserCollectResponse.Records[clientAdminUserIdx].Registered
	}

	return &response, nil
}
