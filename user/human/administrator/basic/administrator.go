package basic

import (
	"crypto/rsa"
	"fmt"
	emailGenerator "github.com/iot-my-world/brain/communication/email/generator"
	setPasswordEmail "github.com/iot-my-world/brain/communication/email/generator/set/password"
	"github.com/iot-my-world/brain/communication/email/mailer"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/party"
	"github.com/iot-my-world/brain/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/search/identifier/id"
	"github.com/iot-my-world/brain/search/identifier/username"
	"github.com/iot-my-world/brain/security/claims"
	humanUserLoginClaims "github.com/iot-my-world/brain/security/claims/login/user/human"
	forgotPasswordClaims "github.com/iot-my-world/brain/security/claims/resetPassword"
	"github.com/iot-my-world/brain/security/token"
	humanUser "github.com/iot-my-world/brain/user/human"
	humanUserAction "github.com/iot-my-world/brain/user/human/action"
	humanUserAdministrator "github.com/iot-my-world/brain/user/human/administrator"
	humanUserAdministratorException "github.com/iot-my-world/brain/user/human/administrator/exception"
	humanUserRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler"
	humanUserRecordHandlerException "github.com/iot-my-world/brain/user/human/recordHandler/exception"
	userValidator "github.com/iot-my-world/brain/user/human/validator"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type administrator struct {
	humanUserRecordHandler    humanUserRecordHandler.RecordHandler
	userValidator             userValidator.Validator
	mailer                    mailer.Mailer
	jwtGenerator              token.JWTGenerator
	mailRedirectBaseUrl       string
	systemClaims              *humanUserLoginClaims.Login
	setPasswordEmailGenerator emailGenerator.Generator
}

func New(
	humanUserRecordHandler humanUserRecordHandler.RecordHandler,
	userValidator userValidator.Validator,
	mailer mailer.Mailer,
	rsaPrivateKey *rsa.PrivateKey,
	mailRedirectBaseUrl string,
	systemClaims *humanUserLoginClaims.Login,
	setPasswordEmailGenerator emailGenerator.Generator,
) humanUserAdministrator.Administrator {
	return &administrator{
		humanUserRecordHandler:    humanUserRecordHandler,
		userValidator:             userValidator,
		mailer:                    mailer,
		jwtGenerator:              token.NewJWTGenerator(rsaPrivateKey),
		mailRedirectBaseUrl:       mailRedirectBaseUrl,
		systemClaims:              systemClaims,
		setPasswordEmailGenerator: setPasswordEmailGenerator,
	}
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *humanUserAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// claims must be login claims to be able to get user
		if request.Claims.Type() != claims.HumanUserLogin {
			reasonsInvalid = append(reasonsInvalid, "claims must be of type login")
		}

		// user must be valid
		validationResponse, err := a.userValidator.Validate(&userValidator.ValidateRequest{
			Claims: request.Claims,
			User:   request.User,
			Action: humanUserAction.UpdateAllowedFields,
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

func (a *administrator) UpdateAllowedFields(request *humanUserAdministrator.UpdateAllowedFieldsRequest) (*humanUserAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the user
	userRetrieveResponse, err := a.humanUserRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	})
	if err != nil {
		return nil, humanUserAdministratorException.UserRetrieval{Reasons: []string{err.Error()}}
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
	userUpdateResponse, err := a.humanUserRecordHandler.Update(&humanUserRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
		User:       userRetrieveResponse.User,
	})
	if err != nil {
		return nil, humanUserAdministratorException.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	return &humanUserAdministrator.UpdateAllowedFieldsResponse{
		User: userUpdateResponse.User,
	}, nil
}

func (a *administrator) ValidateGetMyUserRequest(request *humanUserAdministrator.GetMyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) GetMyUser(request *humanUserAdministrator.GetMyUserRequest) (*humanUserAdministrator.GetMyUserResponse, error) {
	if err := a.ValidateGetMyUserRequest(request); err != nil {
		return nil, err
	}

	// infer the login claims type
	loginClaims, ok := request.Claims.(humanUserLoginClaims.Login)
	if !ok {
		return nil, humanUserAdministratorException.InvalidClaims{Reasons: []string{"cannot assert login claims type"}}
	}

	// retrieve user
	userRetrieveResponse, err := a.humanUserRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	})
	if err != nil {
		return nil, humanUserAdministratorException.UserRetrieval{Reasons: []string{"user retrieval", err.Error()}}
	}
	return &humanUserAdministrator.GetMyUserResponse{User: userRetrieveResponse.User}, nil
}

func (a *administrator) ValidateCreateRequest(request *humanUserAdministrator.CreateRequest) error {
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
		validationResponse, err := a.userValidator.Validate(&userValidator.ValidateRequest{
			Claims: request.Claims,
			User:   request.User,
			Action: humanUserAction.Create,
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

func (a *administrator) Create(request *humanUserAdministrator.CreateRequest) (*humanUserAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	// create the user
	createResponse, err := a.humanUserRecordHandler.Create(&humanUserRecordHandler.CreateRequest{
		User: request.User,
	})
	if err != nil {
		return nil, humanUserAdministratorException.UserCreation{Reasons: []string{"user creation", err.Error()}}
	}

	return &humanUserAdministrator.CreateResponse{User: createResponse.User}, nil
}

func (a *administrator) ValidateSetPasswordRequest(request *humanUserAdministrator.SetPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.NewPassword == "" {
		reasonsInvalid = append(reasonsInvalid, "password blank")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier is nil")
	} else if !humanUser.IsValidIdentifier(request.Identifier) {
		reasonsInvalid = append(reasonsInvalid, "invalid user identifier")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) SetPassword(request *humanUserAdministrator.SetPasswordRequest) (*humanUserAdministrator.SetPasswordResponse, error) {
	if err := a.ValidateSetPasswordRequest(request); err != nil {
		return nil, err
	}

	// Retrieve User
	retrieveUserResponse, err := a.humanUserRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		return nil, humanUserAdministratorException.SetPassword{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, humanUserAdministratorException.SetPassword{Reasons: []string{"hashing password", err.Error()}}
	}

	// update user
	retrieveUserResponse.User.Password = pwdHash

	if _, err := a.humanUserRecordHandler.Update(&humanUserRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		User:       retrieveUserResponse.User,
	}); err != nil {
		return nil, humanUserAdministratorException.SetPassword{Reasons: []string{"update user", err.Error()}}
	}

	return &humanUserAdministrator.SetPasswordResponse{}, nil
}

func (a *administrator) ValidateUpdatePasswordRequest(request *humanUserAdministrator.UpdatePasswordRequest) error {
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

func (a *administrator) UpdatePassword(request *humanUserAdministrator.UpdatePasswordRequest) (*humanUserAdministrator.UpdatePasswordResponse, error) {
	if err := a.ValidateUpdatePasswordRequest(request); err != nil {
		return nil, err
	}

	// user identifier taken from claims as you can only update your own password
	loginClaims, ok := request.Claims.(humanUserLoginClaims.Login)
	if !ok {
		return nil, brainException.Unexpected{Reasons: []string{"inferring claims to type login claims"}}
	}

	// Retrieve User
	retrieveUserResponse, err := a.humanUserRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	})
	if err != nil {
		return nil, humanUserAdministratorException.UpdatePassword{Reasons: []string{"retrieving user record", err.Error()}}
	}

	//User record retrieved successfully, check given old password
	if err := bcrypt.CompareHashAndPassword(retrieveUserResponse.User.Password, []byte(request.ExistingPassword)); err != nil {
		//Password Incorrect
		return nil, humanUserAdministratorException.UpdatePassword{Reasons: []string{"given existing password incorrect"}}
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, humanUserAdministratorException.UpdatePassword{Reasons: []string{"hashing password", err.Error()}}
	}

	// update user
	retrieveUserResponse.User.Password = pwdHash

	updateUserResponse, err := a.humanUserRecordHandler.Update(&humanUserRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
		User:       retrieveUserResponse.User,
	})
	if err != nil {
		return nil, humanUserAdministratorException.SetPassword{Reasons: []string{"update user", err.Error()}}
	}

	return &humanUserAdministrator.UpdatePasswordResponse{
		User: updateUserResponse.User,
	}, nil
}

func (a *administrator) ValidateCheckPasswordRequest(request *humanUserAdministrator.CheckPasswordRequest) error {
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

func (a *administrator) CheckPassword(request *humanUserAdministrator.CheckPasswordRequest) (*humanUserAdministrator.CheckPasswordResponse, error) {
	if err := a.ValidateCheckPasswordRequest(request); err != nil {
		return nil, err
	}

	// user identifier taken from claims as you can only update your own password
	loginClaims, ok := request.Claims.(humanUserLoginClaims.Login)
	if !ok {
		return nil, brainException.Unexpected{Reasons: []string{"inferring claims to type login claims"}}
	}

	// Retrieve User
	retrieveUserResponse, err := a.humanUserRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	})
	if err != nil {
		return nil, humanUserAdministratorException.CheckPassword{Reasons: []string{"retrieving user record", err.Error()}}
	}

	result := true
	//User record retrieved successfully, check given old password
	if err := bcrypt.CompareHashAndPassword(retrieveUserResponse.User.Password, []byte(request.Password)); err != nil {
		//Password Incorrect
		result = false
	}

	return &humanUserAdministrator.CheckPasswordResponse{
		Result: result,
	}, nil
}

func (a *administrator) ValidateForgotPasswordRequest(request *humanUserAdministrator.ForgotPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UsernameOrEmailAddress == "" {
		reasonsInvalid = append(reasonsInvalid, "UsernameOrEmailAddress blank")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) ForgotPassword(request *humanUserAdministrator.ForgotPasswordRequest) (*humanUserAdministrator.ForgotPasswordResponse, error) {
	if err := a.ValidateForgotPasswordRequest(request); err != nil {
		return nil, err
	}
	var retrieveUserResponse *humanUserRecordHandler.RetrieveResponse
	var err error

	//try and retrieve User record with username
	retrieveUserResponse, err = a.humanUserRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
		Claims:     *a.systemClaims,
		Identifier: username.Identifier{Username: request.UsernameOrEmailAddress},
	})
	switch err.(type) {
	case nil:
		// do nothing, this means that the user could be retrieved

	case humanUserRecordHandlerException.NotFound:
		//try and retrieve User record with email address
		retrieveUserResponse, err = a.humanUserRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
			Claims:     *a.systemClaims,
			Identifier: emailAddress.Identifier{EmailAddress: request.UsernameOrEmailAddress},
		})
		switch err.(type) {
		case nil:
			// do nothing, this means that the user could be retrieved
		case humanUserRecordHandlerException.NotFound:
			return nil, nil
		default:
			// some other retrieval error
			return nil, humanUserAdministratorException.UserRetrieval{Reasons: []string{err.Error()}}
		}
	default:
		// some other retrieval error
		return nil, humanUserAdministratorException.UserRetrieval{Reasons: []string{err.Error()}}
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
		return nil, humanUserAdministratorException.TokenGeneration{Reasons: []string{"forgot password", err.Error()}}
	}
	urlToken := fmt.Sprintf("%s/resetPassword?&t=%s", a.mailRedirectBaseUrl, forgotPasswordToken)

	generateEmailResponse, err := a.setPasswordEmailGenerator.Generate(&emailGenerator.GenerateRequest{
		Data: setPasswordEmail.Data{
			URLToken: urlToken,
			User:     retrieveUserResponse.User,
		},
	})
	if err != nil {
		return nil, humanUserAdministratorException.EmailGeneration{Reasons: []string{"set password", err.Error()}}
	}

	if _, err := a.mailer.Send(&mailer.SendRequest{
		Email: generateEmailResponse.Email,
	}); err != nil {
		return nil, err
	}

	return &humanUserAdministrator.ForgotPasswordResponse{}, nil
}
