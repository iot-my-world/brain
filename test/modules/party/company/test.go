package company

import (
	"encoding/json"
	"fmt"
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client/basic"
	"github.com/iot-my-world/brain/party/company"
	companyAdministrator "github.com/iot-my-world/brain/party/company/administrator"
	companyJsonRpcAdministrator "github.com/iot-my-world/brain/party/company/administrator/jsonRpc"
	companyRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler"
	companyJsonRpcRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler/jsonRpc"
	partyRegistrar "github.com/iot-my-world/brain/party/registrar"
	partyJsonRpcRegistrar "github.com/iot-my-world/brain/party/registrar/jsonRpc"
	"github.com/iot-my-world/brain/search/criterion"
	"github.com/iot-my-world/brain/search/identifier/id"
	"github.com/iot-my-world/brain/search/query"
	authJsonRpcAdaptor "github.com/iot-my-world/brain/security/authorization/service/adaptor/jsonRpc"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/claims/registerCompanyAdminUser"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	humanUser "github.com/iot-my-world/brain/user/human"
	"github.com/stretchr/testify/suite"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"strings"
)

func New(
	url string,
	user humanUser.User,
	testData *[]Data,
) *test {
	return &test{
		testData:      testData,
		user:          user,
		jsonRpcClient: basicJsonRpcClient.New(url),
	}
}

type test struct {
	suite.Suite
	jsonRpcClient        jsonRpcClient.Client
	companyRecordHandler companyRecordHandler.RecordHandler
	companyAdministrator companyAdministrator.Administrator
	partyRegistrar       partyRegistrar.Registrar
	user                 humanUser.User
	testData             *[]Data
}

type Data struct {
	Company   company.Company
	AdminUser humanUser.User
	Users     []humanUser.User
}

func (suite *test) SetupTest() {

	// log in the client
	if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
		UsernameOrEmailAddress: suite.user.Username,
		Password:               string(suite.user.Password),
	}); err != nil {
		suite.Fail("log in error", err.Error())
		return
	}

	// set up service provider clients that use jsonRpcClient
	suite.companyRecordHandler = companyJsonRpcRecordHandler.New(suite.jsonRpcClient)
	suite.companyAdministrator = companyJsonRpcAdministrator.New(suite.jsonRpcClient)
	suite.partyRegistrar = partyJsonRpcRegistrar.New(suite.jsonRpcClient)
}

func (suite *test) TestCreateCompanies() {
	// confirm that there are no companies in database, should be starting clean
	companyCollectResponse, err := suite.companyRecordHandler.Collect(&companyRecordHandler.CollectRequest{
		Criteria: make([]criterion.Criterion, 0),
		Query:    query.Query{},
	})
	if err != nil {
		suite.Failf("collect companies failed", err.Error())
		return
	}
	if !suite.Equal(0, companyCollectResponse.Total, "company collection should be empty") {
		suite.FailNow("company collection not empty")
	}

	for idx := range *suite.testData {
		companyEntity := &((*suite.testData)[idx].Company)

		// update the new company's details as would be done from the front end
		(*companyEntity).ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
		(*companyEntity).ParentId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

		// create the company
		companyCreateResponse, err := suite.companyAdministrator.Create(&companyAdministrator.CreateRequest{
			Company: *companyEntity,
		})
		if err != nil {
			suite.FailNow("create company failed", err.Error())
			return
		}

		// update the company
		(*companyEntity).Id = companyCreateResponse.Company.Id
	}
}

func (suite *test) TestInviteAndRegisterCompanyAdminUsers() {
	for idx := range *suite.testData {
		companyEntity := &((*suite.testData)[idx].Company)
		companyAdminUserEntity := &((*suite.testData)[idx].AdminUser)

		// invite the admin user
		inviteCompanyAdminUserResponse, err := suite.partyRegistrar.InviteCompanyAdminUser(&partyRegistrar.InviteCompanyAdminUserRequest{
			CompanyIdentifier: id.Identifier{Id: companyEntity.Id},
		})
		if err != nil {
			suite.FailNow("invite company admin user failed", err.Error())
			return
		}

		// parse the urlToken into a jsonWebToken object
		jwt := inviteCompanyAdminUserResponse.URLToken[strings.Index(inviteCompanyAdminUserResponse.URLToken, "&t=")+3:]
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

		// infer the interfaces type and update the company admin user entity with details from them
		switch typedClaims := unwrappedClaims.(type) {
		case registerCompanyAdminUser.RegisterCompanyAdminUser:
			(*companyAdminUserEntity).Id = typedClaims.User.Id
			(*companyAdminUserEntity).EmailAddress = typedClaims.User.EmailAddress
			(*companyAdminUserEntity).ParentPartyType = typedClaims.User.ParentPartyType
			(*companyAdminUserEntity).ParentId = typedClaims.User.ParentId
			(*companyAdminUserEntity).PartyType = typedClaims.User.PartyType
			(*companyAdminUserEntity).PartyId = typedClaims.User.PartyId
		default:
			suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterCompanyAdminUser))
		}

		// store login token
		logInToken := suite.jsonRpcClient.GetJWT()
		// change token to registration token
		if err := suite.jsonRpcClient.SetJWT(jwt); err != nil {
			suite.FailNow("failed to set json rpc client jwt for registration", err.Error())
		}

		// register the company admin user
		registerCompanyAdminUserResponse, err := suite.partyRegistrar.RegisterCompanyAdminUser(&partyRegistrar.RegisterCompanyAdminUserRequest{
			User: *companyAdminUserEntity,
		})
		if err != nil {
			suite.FailNow("error registering company admin user", err.Error())
			return
		}

		// set token back to logInToken
		if err := suite.jsonRpcClient.SetJWT(logInToken); err != nil {
			suite.FailNow("failed to set json rpc client jwt back to logInToken", err.Error())
		}

		// update the company admin user entity
		(*companyAdminUserEntity).Id = registerCompanyAdminUserResponse.User.Id
		(*companyAdminUserEntity).Roles = registerCompanyAdminUserResponse.User.Roles
	}
}
