package public

import (
	"encoding/json"
	"fmt"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client/basic"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	partyJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/party/administrator/jsonRpc"
	"github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/party/company"
	partyRegistrar "github.com/iot-my-world/brain/pkg/party/registrar"
	partyJsonRpcRegistrar "github.com/iot-my-world/brain/pkg/party/registrar/jsonRpc"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/claims/registerClientAdminUser"
	"github.com/iot-my-world/brain/pkg/security/claims/registerCompanyAdminUser"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	humanUserAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator"
	humanUserJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator/jsonRpc"
	humanUserRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	humanUserJsonRpcRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler/jsonRpc"
	companyTestModule "github.com/iot-my-world/brain/test/modules/party/company"
	"github.com/stretchr/testify/suite"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"strings"
)

type test struct {
	suite.Suite
	jsonRpcClient          jsonRpcClient.Client
	partyAdministrator     partyAdministrator.Administrator
	companyTestData        []CompanyData
	clientTestData         []ClientData
	partyRegistrar         partyRegistrar.Registrar
	humanUserAdministrator humanUserAdministrator.Administrator
	humanUserRecordHandler humanUserRecordHandler.RecordHandler
}

type CompanyData struct {
	Company   company.Company
	AdminUser humanUser.User
	Users     []humanUser.User
}

type ClientData struct {
	Client    client.Client
	AdminUser humanUser.User
	Users     []humanUser.User
}

func New(
	url string,
	companyTestData []CompanyData,
	clientTestData []ClientData,
) *test {
	return &test{
		jsonRpcClient:   basicJsonRpcClient.New(url),
		companyTestData: companyTestData,
		clientTestData:  clientTestData,
	}
}

func (suite *test) SetupTest() {
	// not logging in jsonRpcClient since these tests are done as a public user
	suite.partyAdministrator = partyJsonRpcAdministrator.New(suite.jsonRpcClient)
	suite.partyRegistrar = partyJsonRpcRegistrar.New(suite.jsonRpcClient)
	suite.humanUserAdministrator = humanUserJsonRpcAdministrator.New(suite.jsonRpcClient)
	suite.humanUserRecordHandler = humanUserJsonRpcRecordHandler.New(suite.jsonRpcClient)
}

func (suite *test) TestPublic1InviteAndRegisterCompanies() {
	for _, companyData := range suite.companyTestData {
		inviteResponse, err := suite.partyAdministrator.CreateAndInviteCompany(&partyAdministrator.CreateAndInviteCompanyRequest{
			Company: companyData.Company,
		})
		if err != nil {
			suite.FailNow(
				"error creating and inviting company",
				err.Error(),
			)
			return
		}

		// parse the urlToken into a jsonWebToken object
		jwt := inviteResponse.RegistrationURLToken[strings.Index(inviteResponse.RegistrationURLToken, "&t=")+3:]
		jwtObject, err := jose.ParseSigned(jwt)
		if err != nil {
			suite.FailNow("error parsing jwt", err.Error())
		}

		// Access Underlying jwt payload bytes without verification
		jwtPayload := reflect.ValueOf(jwtObject).Elem().FieldByName("payload")

		// parse the bytes into wrapped claims
		wrapped := wrappedClaims.Wrapped{}
		if err := json.Unmarshal(jwtPayload.Bytes(), &wrapped); err != nil {
			suite.FailNow("error unmarshalling claims", err.Error())
		}

		// unwrap the claims into a claims.Claims interface
		unwrappedClaims, err := wrapped.Unwrap()
		if err != nil {
			suite.FailNow("error unwrapping claims", err.Error())
			return
		}

		// confirm that the claims Type is correct
		if !suite.Equal(claims.RegisterCompanyAdminUser, unwrappedClaims.Type(), "claims should be "+claims.RegisterCompanyAdminUser) {
			suite.FailNow(fmt.Sprintf("claims are not of type %s", claims.RegisterCompanyAdminUser))
		}

		// infer the interface's type and update the company admin user entity with details from them
		switch typedClaims := unwrappedClaims.(type) {
		case registerCompanyAdminUser.RegisterCompanyAdminUser:
			companyData.AdminUser.Id = typedClaims.User.Id
			companyData.AdminUser.EmailAddress = typedClaims.User.EmailAddress
			companyData.AdminUser.ParentPartyType = typedClaims.User.ParentPartyType
			companyData.AdminUser.ParentId = typedClaims.User.ParentId
			companyData.AdminUser.PartyType = typedClaims.User.PartyType
			companyData.AdminUser.PartyId = typedClaims.User.PartyId
		default:
			suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterCompanyAdminUser))
		}

		// set registration token
		if err := suite.jsonRpcClient.SetJWT(jwt); err != nil {
			suite.FailNow("failed to set json rpc client jwt for registration", err.Error())
		}

		// register the company admin user
		if _, err := suite.partyRegistrar.RegisterCompanyAdminUser(&partyRegistrar.RegisterCompanyAdminUserRequest{
			User: companyData.AdminUser,
		}); err != nil {
			suite.FailNow("error registering company admin user", err.Error())
			return
		}

		// log out the json rpc client
		suite.jsonRpcClient.Logout()
	}
}

func (suite *test) TestPublic2CompanyTests() {
	for _, companyData := range suite.companyTestData {
		companyTests := companyTestModule.New(
			suite.jsonRpcClient.GetURL(),
			companyData.AdminUser,
			[]companyTestModule.Data{
				{
					Company:   companyData.Company,
					AdminUser: companyData.AdminUser,
					Users:     companyData.Users,
				},
			},
		)
		companyTests.Suite = suite.Suite
		companyTests.SetupTest()
		suite.Run("Create Users", companyTests.TestCompany5CreateUsers)
		suite.Run("Invite And Register Users", companyTests.TestCompany6InviteAndRegisterUsers)
		suite.Run("User Login", companyTests.TestCompany7UserLogin)
	}
}

func (suite *test) TestPublic3InviteAndRegisterClients() {
	for _, clientData := range suite.clientTestData {
		inviteResponse, err := suite.partyAdministrator.CreateAndInviteClient(&partyAdministrator.CreateAndInviteClientRequest{
			Client: clientData.Client,
		})
		if err != nil {
			suite.FailNow(
				"error creating and inviting client",
				err.Error(),
			)
			return
		}

		// parse the urlToken into a jsonWebToken object
		jwt := inviteResponse.RegistrationURLToken[strings.Index(inviteResponse.RegistrationURLToken, "&t=")+3:]
		jwtObject, err := jose.ParseSigned(jwt)
		if err != nil {
			suite.FailNow("error parsing jwt", err.Error())
		}

		// Access Underlying jwt payload bytes without verification
		jwtPayload := reflect.ValueOf(jwtObject).Elem().FieldByName("payload")

		// parse the bytes into wrapped claims
		wrapped := wrappedClaims.Wrapped{}
		if err := json.Unmarshal(jwtPayload.Bytes(), &wrapped); err != nil {
			suite.FailNow("error unmarshalling claims", err.Error())
		}

		// unwrap the claims into a claims.Claims interface
		unwrappedClaims, err := wrapped.Unwrap()
		if err != nil {
			suite.FailNow("error unwrapping claims", err.Error())
			return
		}

		// confirm that the claims Type is correct
		if !suite.Equal(claims.RegisterClientAdminUser, unwrappedClaims.Type(), "claims should be "+claims.RegisterClientAdminUser) {
			suite.FailNow(fmt.Sprintf("claims are not of type %s", claims.RegisterClientAdminUser))
		}

		// infer the interface's type and update the client admin user entity with details from them
		switch typedClaims := unwrappedClaims.(type) {
		case registerClientAdminUser.RegisterClientAdminUser:
			clientData.AdminUser.Id = typedClaims.User.Id
			clientData.AdminUser.EmailAddress = typedClaims.User.EmailAddress
			clientData.AdminUser.ParentPartyType = typedClaims.User.ParentPartyType
			clientData.AdminUser.ParentId = typedClaims.User.ParentId
			clientData.AdminUser.PartyType = typedClaims.User.PartyType
			clientData.AdminUser.PartyId = typedClaims.User.PartyId
		default:
			suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterCompanyAdminUser))
		}

		// register the client admin user
		if _, err := suite.partyRegistrar.RegisterClientAdminUser(&partyRegistrar.RegisterClientAdminUserRequest{
			User: clientData.AdminUser,
		}); err != nil {
			suite.FailNow("error registering client admin user", err.Error())
			return
		}

		// log out the json rpc client
		suite.jsonRpcClient.Logout()
	}
}
