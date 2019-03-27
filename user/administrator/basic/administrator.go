package basic

import (
	"crypto/rsa"
	"fmt"
	"gitlab.com/iotTracker/brain/email/mailer"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/username"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/claims/forgotPassword"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"gitlab.com/iotTracker/brain/security/token"
	"gitlab.com/iotTracker/brain/user"
	userAction "gitlab.com/iotTracker/brain/user/action"
	userAdministrator "gitlab.com/iotTracker/brain/user/administrator"
	userAdministratorException "gitlab.com/iotTracker/brain/user/administrator/exception"
	userRecordHandler "gitlab.com/iotTracker/brain/user/recordHandler"
	userRecordHandlerException "gitlab.com/iotTracker/brain/user/recordHandler/exception"
	userValidator "gitlab.com/iotTracker/brain/user/validator"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type administrator struct {
	userRecordHandler   userRecordHandler.RecordHandler
	userValidator       userValidator.Validator
	mailer              mailer.Mailer
	jwtGenerator        token.JWTGenerator
	mailRedirectBaseUrl string
	systemClaims        *login.Login
}

func New(
	userRecordHandler userRecordHandler.RecordHandler,
	userValidator userValidator.Validator,
	mailer mailer.Mailer,
	rsaPrivateKey *rsa.PrivateKey,
	mailRedirectBaseUrl string,
	systemClaims *login.Login,
) userAdministrator.Administrator {
	return &administrator{
		userRecordHandler:   userRecordHandler,
		userValidator:       userValidator,
		mailer:              mailer,
		jwtGenerator:        token.NewJWTGenerator(rsaPrivateKey),
		mailRedirectBaseUrl: mailRedirectBaseUrl,
		systemClaims:        systemClaims,
	}
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *userAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// claims must be login claims to be able to get user
		if request.Claims.Type() != claims.Login {
			reasonsInvalid = append(reasonsInvalid, "claims must be of type login")
		}

		// user must be valid
		validationResponse, err := a.userValidator.Validate(&userValidator.ValidateRequest{
			Claims: request.Claims,
			User:   request.User,
			Action: userAction.UpdateAllowedFields,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating user: "+err.Error())
		}
		if len(validationResponse.ReasonsInvalid) > 0 {
			for _, reason := range validationResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("user invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *userAdministrator.UpdateAllowedFieldsRequest) (*userAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the user
	userRetrieveResponse, err := a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	})
	if err != nil {
		return nil, userAdministratorException.UserRetrieval{Reasons: []string{err.Error()}}
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
	userUpdateResponse, err := a.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
		User:       userRetrieveResponse.User,
	})
	if err != nil {
		return nil, userAdministratorException.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	return &userAdministrator.UpdateAllowedFieldsResponse{
		User: userUpdateResponse.User,
	}, nil
}

func (a *administrator) ValidateGetMyUserRequest(request *userAdministrator.GetMyUserRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) GetMyUser(request *userAdministrator.GetMyUserRequest) (*userAdministrator.GetMyUserResponse, error) {
	if err := a.ValidateGetMyUserRequest(request); err != nil {
		return nil, err
	}

	// infer the login claims type
	loginClaims, ok := request.Claims.(login.Login)
	if !ok {
		return nil, userAdministratorException.InvalidClaims{Reasons: []string{"cannot assert login claims type"}}
	}

	// retrieve user
	userRetrieveResponse, err := a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	})
	if err != nil {
		return nil, userAdministratorException.UserRetrieval{Reasons: []string{"user retrieval", err.Error()}}
	}
	return &userAdministrator.GetMyUserResponse{User: userRetrieveResponse.User}, nil
}

func (a *administrator) ValidateCreateRequest(request *userAdministrator.CreateRequest) error {
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
			Action: userAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating user: "+err.Error())
		}
		if len(validationResponse.ReasonsInvalid) > 0 {
			for _, reason := range validationResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("user invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *userAdministrator.CreateRequest) (*userAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	// create the user
	createResponse, err := a.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		User: request.User,
	})
	if err != nil {
		return nil, userAdministratorException.UserCreation{Reasons: []string{"user creation", err.Error()}}
	}

	return &userAdministrator.CreateResponse{User: createResponse.User}, nil
}

func (a *administrator) ValidateSetPasswordRequest(request *userAdministrator.SetPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.NewPassword == "" {
		reasonsInvalid = append(reasonsInvalid, "password blank")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "user identifier is nil")
	} else if !user.IsValidIdentifier(request.Identifier) {
		reasonsInvalid = append(reasonsInvalid, "invalid user identifier")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) SetPassword(request *userAdministrator.SetPasswordRequest) (*userAdministrator.SetPasswordResponse, error) {
	if err := a.ValidateSetPasswordRequest(request); err != nil {
		return nil, err
	}

	// Retrieve User
	retrieveUserResponse, err := a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		return nil, userAdministratorException.SetPassword{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, userAdministratorException.SetPassword{Reasons: []string{"hashing password", err.Error()}}
	}

	// update user
	retrieveUserResponse.User.Password = pwdHash

	updateUserResponse, err := a.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		User:       retrieveUserResponse.User,
	})
	if err != nil {
		return nil, userAdministratorException.SetPassword{Reasons: []string{"update user", err.Error()}}
	}

	return &userAdministrator.SetPasswordResponse{
		User: updateUserResponse.User,
	}, nil
}

func (a *administrator) ValidateUpdatePasswordRequest(request *userAdministrator.UpdatePasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// claims must be login claims to be able to get user
		if request.Claims.Type() != claims.Login {
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

func (a *administrator) UpdatePassword(request *userAdministrator.UpdatePasswordRequest) (*userAdministrator.UpdatePasswordResponse, error) {
	if err := a.ValidateUpdatePasswordRequest(request); err != nil {
		return nil, err
	}

	// user identifier taken from claims as you can only update your own password
	loginClaims, ok := request.Claims.(login.Login)
	if !ok {
		return nil, brainException.Unexpected{Reasons: []string{"inferring claims to type login claims"}}
	}

	// Retrieve User
	retrieveUserResponse, err := a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	})
	if err != nil {
		return nil, userAdministratorException.UpdatePassword{Reasons: []string{"retrieving user record", err.Error()}}
	}

	//User record retrieved successfully, check given old password
	if err := bcrypt.CompareHashAndPassword(retrieveUserResponse.User.Password, []byte(request.ExistingPassword)); err != nil {
		//Password Incorrect
		return nil, userAdministratorException.UpdatePassword{Reasons: []string{"given existing password incorrect"}}
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, userAdministratorException.UpdatePassword{Reasons: []string{"hashing password", err.Error()}}
	}

	// update user
	retrieveUserResponse.User.Password = pwdHash

	updateUserResponse, err := a.userRecordHandler.Update(&userRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
		User:       retrieveUserResponse.User,
	})
	if err != nil {
		return nil, userAdministratorException.SetPassword{Reasons: []string{"update user", err.Error()}}
	}

	return &userAdministrator.UpdatePasswordResponse{
		User: updateUserResponse.User,
	}, nil
}

func (a *administrator) ValidateCheckPasswordRequest(request *userAdministrator.CheckPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// claims must be login claims to be able to get user
		if request.Claims.Type() != claims.Login {
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

func (a *administrator) CheckPassword(request *userAdministrator.CheckPasswordRequest) (*userAdministrator.CheckPasswordResponse, error) {
	if err := a.ValidateCheckPasswordRequest(request); err != nil {
		return nil, err
	}

	// user identifier taken from claims as you can only update your own password
	loginClaims, ok := request.Claims.(login.Login)
	if !ok {
		return nil, brainException.Unexpected{Reasons: []string{"inferring claims to type login claims"}}
	}

	// Retrieve User
	retrieveUserResponse, err := a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: loginClaims.UserId,
	})
	if err != nil {
		return nil, userAdministratorException.CheckPassword{Reasons: []string{"retrieving user record", err.Error()}}
	}

	result := true
	//User record retrieved successfully, check given old password
	if err := bcrypt.CompareHashAndPassword(retrieveUserResponse.User.Password, []byte(request.Password)); err != nil {
		//Password Incorrect
		result = false
	}

	return &userAdministrator.CheckPasswordResponse{
		Result: result,
	}, nil
}

func (a *administrator) ValidateForgotPasswordRequest(request *userAdministrator.ForgotPasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UsernameOrEmailAddress == "" {
		reasonsInvalid = append(reasonsInvalid, "UsernameOrEmailAddress blank")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) ForgotPassword(request *userAdministrator.ForgotPasswordRequest) (*userAdministrator.ForgotPasswordResponse, error) {
	if err := a.ValidateForgotPasswordRequest(request); err != nil {
		return nil, err
	}
	var retrieveUserResponse *userRecordHandler.RetrieveResponse
	var err error

	//try and retrieve User record with username
	retrieveUserResponse, err = a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     *a.systemClaims,
		Identifier: username.Identifier{Username: request.UsernameOrEmailAddress},
	})
	switch err.(type) {
	case nil:
		// do nothing, this means that the user could be retrieved

	case userRecordHandlerException.NotFound:
		//try and retrieve User record with email address
		retrieveUserResponse, err = a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
			Claims:     *a.systemClaims,
			Identifier: emailAddress.Identifier{EmailAddress: request.UsernameOrEmailAddress},
		})
		switch err.(type) {
		case nil:
			// do nothing, this means that the user could be retrieved
		case userRecordHandlerException.NotFound:
			return nil, nil
		default:
			// some other retrieval error
			return nil, userAdministratorException.UserRetrieval{Reasons: []string{err.Error()}}
		}
	default:
		// some other retrieval error
		return nil, userAdministratorException.UserRetrieval{Reasons: []string{err.Error()}}
	}

	// User record retrieved successfully
	// generate reset password token for the user
	forgotPasswordToken, err := a.jwtGenerator.GenerateToken(forgotPassword.ForgotPassword{
		UserId:          id.Identifier{Id: retrieveUserResponse.User.Id},
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: retrieveUserResponse.User.ParentPartyType,
		ParentId:        retrieveUserResponse.User.ParentId,
		PartyType:       retrieveUserResponse.User.PartyType,
		PartyId:         retrieveUserResponse.User.PartyId,
	})
	if err != nil {
		return nil, userAdministratorException.TokenGeneration{Reasons: []string{"forgot password", err.Error()}}
	}
	urlToken := fmt.Sprintf("%s/resetPassword?&t=%s", a.mailRedirectBaseUrl, forgotPasswordToken)

	sendMailResponse := mailer.SendResponse{}
	if err := a.mailer.Send(&mailer.SendRequest{
		//From    string
		To: retrieveUserResponse.User.EmailAddress,
		//Cc      string
		Subject: "Password Reset",
		Body:    fmt.Sprintf("Click the link to continue password reset. %s", urlToken),
		//Bcc     []string
	},
		&sendMailResponse); err != nil {
		return nil, err
	}

	return &userAdministrator.ForgotPasswordResponse{}, nil
}
