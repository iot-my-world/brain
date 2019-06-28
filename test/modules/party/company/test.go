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
	"github.com/iot-my-world/brain/search/identifier/adminEmailAddress"
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
	testData []Data,
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
	testData             []Data
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

func (suite *test) Test1CreateCompanies() {
	// create all companies in test data
	for _, data := range suite.testData {
		companyEntity := data.Company

		// update the new company's details as would be done from the front end
		companyEntity.ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
		companyEntity.ParentId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

		// create the company
		if _, err := suite.companyAdministrator.Create(&companyAdministrator.CreateRequest{
			Company: companyEntity,
		}); err != nil {
			suite.FailNow("create company failed", err.Error())
			return
		}
	}

	// collect all companies
	companyCollectResponse, err := suite.companyRecordHandler.Collect(&companyRecordHandler.CollectRequest{
		Criteria: make([]criterion.Criterion, 0),
		Query:    query.Query{},
	})
	if err != nil {
		suite.Failf("collect companies failed", err.Error())
		return
	}

	// confirm that each created company can be found
nextCompanyToCreate:
	// for every company that should be created
	for _, companyToCreate := range suite.testData {
		// look for companyToCreate among collected companies
		for _, existingCompany := range companyCollectResponse.Records {
			if companyToCreate.Company.AdminEmailAddress == existingCompany.AdminEmailAddress {
				// update fields set during creation
				companyToCreate.Company.Id = existingCompany.Id
				companyToCreate.Company.ParentPartyType = existingCompany.ParentPartyType
				companyToCreate.Company.ParentId = existingCompany.ParentId
				// assert should be equal
				suite.Equal(companyToCreate.Company, existingCompany, "created company should be equal")
				// if it is found and equal, check for next company to create
				continue nextCompanyToCreate
			}
		}
		// if execution reaches here then companyToCreate was not found among collected companies
	}
}

func (suite *test) Test2CompanyUpdateAllowedFields() {
	for _, data := range suite.testData {

		// retrieve the company by admin email address
		companyRetrieveResponse, err := suite.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Identifier: adminEmailAddress.Identifier{
				AdminEmailAddress: data.Company.AdminEmailAddress,
			},
		})
		if err != nil {
			suite.FailNow("retrieve company entity failed", err.Error())
			return
		}

		// copy retrieved company
		updatedCompanyEntity := companyRetrieveResponse.Company

		// update allowed fields
		updatedCompanyEntity.Name = "Changed Name"

		// perform update
		updateAllowedFieldsResponse, err := suite.companyAdministrator.UpdateAllowedFields(&companyAdministrator.UpdateAllowedFieldsRequest{
			Company: updatedCompanyEntity,
		})
		if err != nil {
			suite.FailNow("company update allowed fields failed", err.Error())
			return
		}

		suite.Equal(
			updatedCompanyEntity,
			updateAllowedFieldsResponse.Company,
			"updated company should equal company in updated response",
		)

		// retrieve the updated entity by id
		updatedCompanyRetrieveResponse, err := suite.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Identifier: id.Identifier{
				Id: updatedCompanyEntity.Id,
			},
		})
		if err != nil {
			suite.FailNow("retrieve updated company entity failed", err.Error())
			return
		}

		suite.Equal(
			updatedCompanyEntity,
			updatedCompanyRetrieveResponse.Company,
			"retrieved company should equal updated company",
		)

		// update company back to original
		updateAllowedFieldsResponse, err = suite.companyAdministrator.UpdateAllowedFields(&companyAdministrator.UpdateAllowedFieldsRequest{
			Company: companyRetrieveResponse.Company,
		})
		if err != nil {
			suite.FailNow("company update allowed fields failed", err.Error())
			return
		}

		suite.Equal(
			companyRetrieveResponse.Company,
			updateAllowedFieldsResponse.Company,
			"updated company should equal company in updated response",
		)
	}
}

func (suite *test) Test3InviteAndRegisterCompanyAdminUsers() {
	for _, data := range suite.testData {
		companyEntity := data.Company
		companyAdminUserEntity := data.AdminUser

		// retrieve the company by admin email address
		companyRetrieveResponse, err := suite.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Identifier: adminEmailAddress.Identifier{
				AdminEmailAddress: companyEntity.AdminEmailAddress,
			},
		})
		if err != nil {
			suite.FailNow("retrieve company entity failed", err.Error())
			return
		}

		// invite the admin user
		inviteCompanyAdminUserResponse, err := suite.partyRegistrar.InviteCompanyAdminUser(&partyRegistrar.InviteCompanyAdminUserRequest{
			CompanyIdentifier: id.Identifier{Id: companyRetrieveResponse.Company.Id},
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
			companyAdminUserEntity.Id = typedClaims.User.Id
			companyAdminUserEntity.EmailAddress = typedClaims.User.EmailAddress
			companyAdminUserEntity.ParentPartyType = typedClaims.User.ParentPartyType
			companyAdminUserEntity.ParentId = typedClaims.User.ParentId
			companyAdminUserEntity.PartyType = typedClaims.User.PartyType
			companyAdminUserEntity.PartyId = typedClaims.User.PartyId
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
		if _, err := suite.partyRegistrar.RegisterCompanyAdminUser(&partyRegistrar.RegisterCompanyAdminUserRequest{
			User: companyAdminUserEntity,
		}); err != nil {
			suite.FailNow("error registering company admin user", err.Error())
			return
		}

		// set token back to logInToken
		if err := suite.jsonRpcClient.SetJWT(logInToken); err != nil {
			suite.FailNow("failed to set json rpc client jwt back to logInToken", err.Error())
		}
	}
}
