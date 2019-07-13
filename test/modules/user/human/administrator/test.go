package administrator

import (
	"encoding/json"
	"fmt"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client/basic"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	authorizationAdministrator "github.com/iot-my-world/brain/pkg/security/authorization/administrator"
	"github.com/iot-my-world/brain/pkg/security/claims"
	resetPasswordClaims "github.com/iot-my-world/brain/pkg/security/claims/resetPassword"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	humanUserAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator"
	humanUserJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator/jsonRpc"
	"github.com/stretchr/testify/suite"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"strings"
)

type test struct {
	suite.Suite
	jsonRpcClient          jsonRpcClient.Client
	usersData              []humanUser.User
	humanUserAdministrator humanUserAdministrator.Administrator
}

func New(
	url string,
	usersData []humanUser.User,
) *test {
	return &test{
		jsonRpcClient: basicJsonRpcClient.New(url),
		usersData:     usersData,
	}
}

func (suite *test) SetupTest() {
	// not logging in jsonRpcClient since these tests are done as a public user
	suite.humanUserAdministrator = humanUserJsonRpcAdministrator.New(suite.jsonRpcClient)
}

func (suite *test) TestUserAdministrator1ForgotPassword() {
	for _, user := range suite.usersData {
		forgotPasswordResponse, err := suite.humanUserAdministrator.ForgotPassword(&humanUserAdministrator.ForgotPasswordRequest{
			UsernameOrEmailAddress: user.EmailAddress,
		})
		if err != nil {
			suite.FailNow("error performing forgot password", err.Error())
			return
		}

		// parse the urlToken into a jsonWebToken object
		jwt := forgotPasswordResponse.URLToken[strings.Index(forgotPasswordResponse.URLToken, "&t=")+3:]
		jwtObject, err := jose.ParseSigned(jwt)
		if err != nil {
			suite.FailNow("error parsing jwt", err.Error())
			return
		}

		// Access Underlying jwt payload bytes without verification
		jwtPayload := reflect.ValueOf(jwtObject).Elem().FieldByName("payload")

		// parse the bytes into wrapped claims
		wrapped := wrappedClaims.Wrapped{}
		if err := json.Unmarshal(jwtPayload.Bytes(), &wrapped); err != nil {
			suite.FailNow("error unmarshalling claims", err.Error())
			return
		}

		// unwrap the claims into a claims.Claims interface
		unwrappedClaims, err := wrapped.Unwrap()
		if err != nil {
			suite.FailNow("error unwrapping claims", err.Error())
			return
		}

		// confirm that the claims Type is correct
		if !suite.Equal(claims.ResetPassword, unwrappedClaims.Type(), "claims should be "+claims.ResetPassword) {
			suite.FailNow(fmt.Sprintf("claims are not of type %s", claims.ResetPassword))
		}

		// infer the interface's type and update the client admin user entity with details from them
		var userIdentifier identifier.Identifier
		switch typedClaims := unwrappedClaims.(type) {
		case resetPasswordClaims.ResetPassword:
			userIdentifier = typedClaims.UserId
		default:
			suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterCompanyAdminUser))
			return
		}

		// set reset password token
		if err := suite.jsonRpcClient.SetJWT(jwt); err != nil {
			suite.FailNow("failed to set json rpc client jwt for reset password", err.Error())
			return
		}

		// store the password
		oldPassword := user.Password

		// set the password
		if _, err := suite.humanUserAdministrator.SetPassword(&humanUserAdministrator.SetPasswordRequest{
			Identifier:  userIdentifier,
			NewPassword: "321",
		}); err != nil {
			suite.FailNow("error setting password", err.Error())
			return
		}

		// log out the json rpc client
		suite.jsonRpcClient.Logout()

		// try and log in with the new password
		if err := suite.jsonRpcClient.Login(authorizationAdministrator.LoginRequest{
			UsernameOrEmailAddress: user.Username,
			Password:               "321",
		}); err != nil {
			suite.FailNow("error logging in with new password", err.Error())
			return
		}

		// log out the json rpc client again
		suite.jsonRpcClient.Logout()

		// request a password reset again to set password back to original
		forgotPasswordResponse, err = suite.humanUserAdministrator.ForgotPassword(&humanUserAdministrator.ForgotPasswordRequest{
			UsernameOrEmailAddress: user.EmailAddress,
		})
		if err != nil {
			suite.FailNow("error performing forgot password again", err.Error())
			return
		}

		// parse the urlToken into a jsonWebToken object
		jwt = forgotPasswordResponse.URLToken[strings.Index(forgotPasswordResponse.URLToken, "&t=")+3:]
		jwtObject, err = jose.ParseSigned(jwt)
		if err != nil {
			suite.FailNow("error parsing jwt", err.Error())
			return
		}

		// Access Underlying jwt payload bytes without verification
		jwtPayload = reflect.ValueOf(jwtObject).Elem().FieldByName("payload")

		// parse the bytes into wrapped claims
		if err := json.Unmarshal(jwtPayload.Bytes(), &wrapped); err != nil {
			suite.FailNow("error unmarshalling claims", err.Error())
		}

		// unwrap the claims into a claims.Claims interface
		unwrappedClaims, err = wrapped.Unwrap()
		if err != nil {
			suite.FailNow("error unwrapping claims", err.Error())
			return
		}

		// confirm that the claims Type is correct
		if !suite.Equal(claims.ResetPassword, unwrappedClaims.Type(), "claims should be "+claims.ResetPassword) {
			suite.FailNow(fmt.Sprintf("claims are not of type %s", claims.ResetPassword))
		}

		// infer the interface's type and update the client admin user entity with details from them
		switch typedClaims := unwrappedClaims.(type) {
		case resetPasswordClaims.ResetPassword:
			userIdentifier = typedClaims.UserId
		default:
			suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterCompanyAdminUser))
			return
		}

		// set reset password token
		if err := suite.jsonRpcClient.SetJWT(jwt); err != nil {
			suite.FailNow("failed to set json rpc client jwt for reset password", err.Error())
			return
		}

		// set the password
		if _, err := suite.humanUserAdministrator.SetPassword(&humanUserAdministrator.SetPasswordRequest{
			Identifier:  userIdentifier,
			NewPassword: string(oldPassword),
		}); err != nil {
			suite.FailNow("error setting password", err.Error())
			return
		}

		// log out the json rpc client
		suite.jsonRpcClient.Logout()
	}
}
