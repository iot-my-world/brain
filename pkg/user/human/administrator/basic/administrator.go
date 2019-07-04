package basic

import (
	"crypto/rsa"
	"fmt"
	"github.com/iot-my-world/brain/environment"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	emailGenerator "github.com/iot-my-world/brain/pkg/communication/email/generator"
	setPasswordEmail "github.com/iot-my-world/brain/pkg/communication/email/generator/set/password"
	"github.com/iot-my-world/brain/pkg/communication/email/mailer"
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/search/identifier/username"
	"github.com/iot-my-world/brain/pkg/user/human"
	"github.com/iot-my-world/brain/pkg/user/human/action"
	administrator2 "github.com/iot-my-world/brain/pkg/user/human/administrator"
	exception2 "github.com/iot-my-world/brain/pkg/user/human/administrator/exception"
	"github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	"github.com/iot-my-world/brain/pkg/user/human/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/user/human/validator"
	"github.com/iot-my-world/brain/security/claims"
	humanUserLoginClaims "github.com/iot-my-world/brain/security/claims/login/user/human"
	forgotPasswordClaims "github.com/iot-my-world/brain/security/claims/resetPassword"
	"github.com/iot-my-world/brain/security/token"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type administrator struct {
	humanUserRecordHandler    recordHandler.RecordHandler
	userValidator             validator.Validator
	mailer                    mailer.Mailer
	jwtGenerator              token.JWTGenerator
	mailRedirectBaseUrl       string
	systemClaims              *humanUserLoginClaims.Login
	setPasswordEmailGenerator emailGenerator.Generator
	environmentType           environment.Type
}

func New(
	humanUserRecordHandler recordHandler.RecordHandler,
	userValidator validator.Validator,
	mailer mailer.Mailer,
	rsaPrivateKey *rsa.PrivateKey,
	mailRedirectBaseUrl string,
	systemClaims *humanUserLoginClaims.Login,
	setPasswordEmailGenerator emailGenerator.Generator,
	environmentType environment.Type,
) administrator2.Administrator {
	return &administrator{
		humanUserRecordHandler:    humanUserRecordHandler,
		userValidator:             userValidator,
		mailer:                    mailer,
		jwtGenerator:              token.NewJWTGenerator(rsaPrivateKey),
		mailRedirectBaseUrl:       mailRedirectBaseUrl,
		systemClaims:              systemClaims,
		setPasswordEmailGenerator: setPasswordEmailGenerator,
		environmentType:           environmentType,
	}
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *administrator2.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// claims must be login claims to be able to get user
		if request.Claims.Type() != claims.HumanUserLogin {
			reasonsInvalid = append(reasonsInvalid, "claims must be of type login")
		}

		// user must be valid
		validationResponse, err := a.userValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			User:   request.User,
			Action: action.UpdateAllowedFields,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating user: "+err.Error())
		} else {
			if len(validationResponse.ReasonsInvalid) > 0 {
				for _, reason := range validationResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("user invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *administrator2.UpdateAllowedFieldsRequest) (*administrator2.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// retrieve the user
	userRetrieveResponse, err := a.humanUserRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	})
	if err != nil {
		err = exception2.UpdateAllowedFields{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// update allowed fields on the user
	// userRetrieveResponse.user.Id =              request.User.Id
	userRetrieveResponse.User.Name = request.User.Name
	userRetrieveResponse.User.Surname = request.User.Surname
	userRetrieveResponse.User.Username = request.User.Username
	//userRetrieveResponse.User.EmailAddress = request.User.EmailAddress
	//userRetrieveResponse.User.Password = request.User.Password
	//userRetrieveResponse.User.Roles = request.User.Roles
	//userRetrieveResponse.User.ParentPartyType = request.User.ParentPartyType
	//userRetrieveResponse.User.ParentId = request.User.ParentId
	//userRetrieveResponse.User.PartyType = request.User.ParentPartyType
	//userRetrieveResponse.User.PartyId = request.User.PartyId
	//userRetrieveResponse.User.Registered = request.User.Registered

	// update the user
	userUpdateResponse, err := a.humanUserRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
		User:       userRetrieveResponse.User,
	})
	if err != nil {
		err = exception2.UpdateAllowedFields{Reasons: []string{"updating", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.UpdateAllowedFieldsResponse{
		User: userUpdateResponse.User,
	}, nil
}

func (a *administrator) ValidateGetMyUserRequest(request *administrator2.GetMyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) GetMyUser(request *administrator2.GetMyUserRequest) (*administrator2.GetMyUserResponse, error) {
	if err := a.ValidateGetMyUserRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// infer the login claims type
	loginClaims, ok := request.Claims.(humanUserLoginClaims.Login)
	if !ok {
		err := exception2.GetMyUser{Reasons: []string{"cannot assert login claims type"}}
		log.Error(err.Error())
		return nil, err
	}

	// retrieve user
	userRetrieveResponse, err := a.humanUserRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	})
	if err != nil {
		err = exception2.GetMyUser{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.GetMyUserResponse{User: userRetrieveResponse.User}, nil
}

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// confirm that the party details of the user being created matches claims
		// i.e users can only be created by their own party unless the system party
		// is acting
		switch request.Claims.PartyDetails().PartyType {
		case party.System:
			// do nothing, we expect system to know what they are doing
		default:
			if request.User.PartyType != request.Claims.PartyDetails().PartyType {
				reasonsInvalid = append(reasonsInvalid, "user PartyType must be the type of the party doing creation")
			}
			if request.User.PartyId != request.Claims.PartyDetails().PartyId {
				reasonsInvalid = append(reasonsInvalid, "client PartyId must be the id of the party doing creation")
			}
			if request.User.ParentPartyType != request.Claims.PartyDetails().ParentPartyType {
				reasonsInvalid = append(reasonsInvalid, "user ParentPartyType must match that of the party doing creation")
			}
			if request.User.ParentId != request.Claims.PartyDetails().ParentId {
				reasonsInvalid = append(reasonsInvalid, "user ParentId must match that of the party doing creation")
			}
		}

		// user must be valid
		validationResponse, err := a.userValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			User:   request.User,
			Action: action.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating user: "+err.Error())
		} else {
			if len(validationResponse.ReasonsInvalid) > 0 {
				invalidReason := "user invalid: "
				for _, reason := range validationResponse.ReasonsInvalid {
					invalidReason += fmt.Sprintf(" %s - %s - %s,", reason.Field, reason.Type, reason.Help)
				}
				reasonsInvalid = append(reasonsInvalid, invalidReason)
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// create the user
	createResponse, err := a.humanUserRecordHandler.Create(&recordHandler.CreateRequest{
		User: request.User,
	})
	if err != nil {
		err = exception2.Create{Reasons: []string{"user creation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.CreateResponse{User: createResponse.User}, nil
}

func (a *administrator) ValidateSetPasswordRequest(request *administrator2.SetPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.NewPassword == "" {
		reasonsInvalid = append(reasonsInvalid, "password blank")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier is nil")
	} else if !human.IsValidIdentifier(request.Identifier) {
		reasonsInvalid = append(reasonsInvalid, "invalid user identifier")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) SetPassword(request *administrator2.SetPasswordRequest) (*administrator2.SetPasswordResponse, error) {
	if err := a.ValidateSetPasswordRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Retrieve User
	retrieveUserResponse, err := a.humanUserRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		err = exception2.SetPassword{Reasons: []string{"retrieving record", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		err = exception2.SetPassword{Reasons: []string{"hashing password", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// update user
	retrieveUserResponse.User.Password = pwdHash

	if _, err := a.humanUserRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		User:       retrieveUserResponse.User,
	}); err != nil {
		err = exception2.SetPassword{Reasons: []string{"update user", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.SetPasswordResponse{}, nil
}

func (a *administrator) ValidateUpdatePasswordRequest(request *administrator2.UpdatePasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// claims must be login claims to be able to get user
		if request.Claims.Type() != claims.HumanUserLogin {
			reasonsInvalid = append(reasonsInvalid, "claims must be of type login")
		}
	}

	if request.ExistingPassword == "" {
		reasonsInvalid = append(reasonsInvalid, "existing password blank")
	}

	if request.NewPassword == "" {
		reasonsInvalid = append(reasonsInvalid, "new password blank")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdatePassword(request *administrator2.UpdatePasswordRequest) (*administrator2.UpdatePasswordResponse, error) {
	if err := a.ValidateUpdatePasswordRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// user identifier taken from claims as you can only update your own password
	loginClaims, ok := request.Claims.(humanUserLoginClaims.Login)
	if !ok {
		err := brainException.Unexpected{Reasons: []string{"inferring claims to type login claims"}}
		log.Error(err.Error())
		return nil, err
	}

	// Retrieve User
	retrieveUserResponse, err := a.humanUserRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	})
	if err != nil {
		err = exception2.UpdatePassword{Reasons: []string{"retrieving user record", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	//User record retrieved successfully, check given old password
	if err := bcrypt.CompareHashAndPassword(retrieveUserResponse.User.Password, []byte(request.ExistingPassword)); err != nil {
		//Password Incorrect
		err = exception2.UpdatePassword{Reasons: []string{"given existing password incorrect"}}
		log.Error(err.Error())
		return nil, err
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		err = exception2.UpdatePassword{Reasons: []string{"hashing password", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// update user
	retrieveUserResponse.User.Password = pwdHash

	updateUserResponse, err := a.humanUserRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
		User:       retrieveUserResponse.User,
	})
	if err != nil {
		err = exception2.SetPassword{Reasons: []string{"update user", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.UpdatePasswordResponse{
		User: updateUserResponse.User,
	}, nil
}

func (a *administrator) ValidateCheckPasswordRequest(request *administrator2.CheckPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// claims must be login claims to be able to get user
		if request.Claims.Type() != claims.HumanUserLogin {
			reasonsInvalid = append(reasonsInvalid, "claims must be of type login")
		}
	}

	if request.Password == "" {
		reasonsInvalid = append(reasonsInvalid, "password blank")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) CheckPassword(request *administrator2.CheckPasswordRequest) (*administrator2.CheckPasswordResponse, error) {
	if err := a.ValidateCheckPasswordRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// user identifier taken from claims as you can only update your own password
	loginClaims, ok := request.Claims.(humanUserLoginClaims.Login)
	if !ok {
		err := brainException.Unexpected{Reasons: []string{"inferring claims to type login claims"}}
		log.Error(err.Error())
		return nil, err
	}

	// Retrieve User
	retrieveUserResponse, err := a.humanUserRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	})
	if err != nil {
		err = exception2.CheckPassword{Reasons: []string{"retrieving user record", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	result := true
	//User record retrieved successfully, check given old password
	if err := bcrypt.CompareHashAndPassword(retrieveUserResponse.User.Password, []byte(request.Password)); err != nil {
		//Password Incorrect
		result = false
	}

	return &administrator2.CheckPasswordResponse{
		Result: result,
	}, nil
}

func (a *administrator) ValidateForgotPasswordRequest(request *administrator2.ForgotPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UsernameOrEmailAddress == "" {
		reasonsInvalid = append(reasonsInvalid, "UsernameOrEmailAddress blank")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) ForgotPassword(request *administrator2.ForgotPasswordRequest) (*administrator2.ForgotPasswordResponse, error) {
	if err := a.ValidateForgotPasswordRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var retrieveUserResponse *recordHandler.RetrieveResponse
	var err error

	//try and retrieve User record with username
	retrieveUserResponse, err = a.humanUserRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     *a.systemClaims,
		Identifier: username.Identifier{Username: request.UsernameOrEmailAddress},
	})
	switch err.(type) {
	case nil:
		// do nothing, this means that the user could be retrieved

	case exception.NotFound:
		//try and retrieve User record with email address
		retrieveUserResponse, err = a.humanUserRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
			Claims:     *a.systemClaims,
			Identifier: emailAddress.Identifier{EmailAddress: request.UsernameOrEmailAddress},
		})
		if err != nil {
			err = exception2.ForgotPassword{Reasons: []string{"user retrieval", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}

	default:
		// some other retrieval error
		err = exception2.ForgotPassword{Reasons: []string{"user retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// User record retrieved successfully
	// generate reset password token for the user
	forgotPasswordToken, err := a.jwtGenerator.GenerateToken(forgotPasswordClaims.ResetPassword{
		UserId:          id.Identifier{Id: retrieveUserResponse.User.Id},
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: retrieveUserResponse.User.ParentPartyType,
		ParentId:        retrieveUserResponse.User.ParentId,
		PartyType:       retrieveUserResponse.User.PartyType,
		PartyId:         retrieveUserResponse.User.PartyId,
	})
	if err != nil {
		err = exception2.ForgotPassword{Reasons: []string{"token generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}
	urlToken := fmt.Sprintf("%s/resetPassword?&t=%s", a.mailRedirectBaseUrl, forgotPasswordToken)

	generateEmailResponse, err := a.setPasswordEmailGenerator.Generate(&emailGenerator.GenerateRequest{
		Data: setPasswordEmail.Data{
			URLToken: urlToken,
			User:     retrieveUserResponse.User,
		},
	})
	if err != nil {
		err = exception2.ForgotPassword{Reasons: []string{"generating email", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	if a.environmentType == environment.Development {
		// if this is the development environment return response with token
		return &administrator2.ForgotPasswordResponse{URLToken: forgotPasswordToken}, nil
	}

	// otherwise send email and return response without token
	if _, err := a.mailer.Send(&mailer.SendRequest{
		Email: generateEmailResponse.Email,
	}); err != nil {
		err = exception2.ForgotPassword{Reasons: []string{"sending email", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.ForgotPasswordResponse{}, nil
}
