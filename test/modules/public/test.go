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
	"github.com/iot-my-world/brain/pkg/security/claims/registerCompanyAdminUser"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	"github.com/stretchr/testify/suite"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"strings"
)

type test struct {
	suite.Suite
	jsonRpcClient      jsonRpcClient.Client
	partyAdministrator partyAdministrator.Administrator
	companyTestData    []CompanyData
	clientTestData     []ClientData
	partyRegistrar     partyRegistrar.Registrar
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
