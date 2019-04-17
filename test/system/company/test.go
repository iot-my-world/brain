package company

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	companyAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/administrator/adaptor/jsonRpc"
	companyRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/recordHandler/adaptor/jsonRpc"
	partyRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyAdminUser"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	companyTestData "gitlab.com/iotTracker/brain/test/company/data"
	testData "gitlab.com/iotTracker/brain/test/data"
	systemTestData "gitlab.com/iotTracker/brain/test/system/data"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"strings"
)

type Company struct {
	suite.Suite
	jsonRpcClient jsonRpcClient.Client
}

func (suite *Company) SetupTest() {
	// create the client
	suite.jsonRpcClient = basicJsonRpcClient.New(testData.BrainURL)

	// log in the client
	if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
		UsernameOrEmailAddress: systemTestData.User.Username,
		Password:               string(systemTestData.User.Password),
	}); err != nil {
		suite.Fail("log in error", err.Error())
	}
}

func (suite *Company) TestSystemCreateCompanies() {
	// confirm that there are no companies in database, should be starting clean
	companyCollectResponse := companyRecordHandlerJsonRpcAdaptor.CollectResponse{}
	if err := suite.jsonRpcClient.JsonRpcRequest(
		"CompanyRecordHandler.Collect",
		companyRecordHandlerJsonRpcAdaptor.CollectRequest{},
		&companyCollectResponse); err != nil {
		suite.Failf("collect companies failed", err.Error())
	}
	if !suite.Equal(0, companyCollectResponse.Total, "company collection should be empty") {
		suite.FailNow("company collection not empty")
	}

	for idx := range companyTestData.EntitiesAndAdminUsersToCreate {
		companyEntity := &(companyTestData.EntitiesAndAdminUsersToCreate[idx].Company)

		// update the new company's details as would be done from the front end
		(*companyEntity).ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
		(*companyEntity).ParentId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

		// create the company
		companyCreateResponse := companyAdministratorJsonRpcAdaptor.CreateResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"CompanyAdministrator.Create",
			companyAdministratorJsonRpcAdaptor.CreateRequest{
				Company: *companyEntity,
			},
			&companyCreateResponse,
		); err != nil {
			suite.FailNow("create company failed", err.Error())
		}

		// update the company
		(*companyEntity).Id = companyCreateResponse.Company.Id
	}
}

func (suite *Company) TestSystemInviteAndRegisterCompanyAdminUsers() {
	for idx := range companyTestData.EntitiesAndAdminUsersToCreate {
		companyEntity := &(companyTestData.EntitiesAndAdminUsersToCreate[idx].Company)
		companyAdminUserEntity := &companyTestData.EntitiesAndAdminUsersToCreate[idx].AdminUser

		// create identifier for the company entity
		companyIdentifier, err := wrappedIdentifier.Wrap(id.Identifier{Id: companyEntity.Id})
		if err != nil {
			suite.FailNow("error wrapping companyIdentifier", err.Error())
		}

		// invite the admin user
		inviteCompanyAdminUserResponse := partyRegistrarJsonRpcAdaptor.InviteCompanyAdminUserResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"PartyRegistrar.InviteCompanyAdminUser",
			partyRegistrarJsonRpcAdaptor.InviteCompanyAdminUserRequest{
				WrappedCompanyIdentifier: *companyIdentifier,
			},
			&inviteCompanyAdminUserResponse,
		); err != nil {
			suite.FailNow("invite company admin user failed", err.Error())
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

		// create a new json rpc client to register the user with
		registerJsonRpcClient := basicJsonRpcClient.New(testData.BrainURL)
		if err := registerJsonRpcClient.SetJWT(jwt); err != nil {
			suite.FailNow("failed to set jwt in registration client", err.Error())
		}

		// register the company admin user
		registerCompanyAdminUserResponse := partyRegistrarJsonRpcAdaptor.RegisterCompanyAdminUserResponse{}
		if err := registerJsonRpcClient.JsonRpcRequest(
			"PartyRegistrar.RegisterCompanyAdminUser",
			partyRegistrarJsonRpcAdaptor.RegisterCompanyAdminUserRequest{
				User: *companyAdminUserEntity,
			},
			&registerCompanyAdminUserResponse,
		); err != nil {
			suite.FailNow("error registering company admin user", err.Error())
		}

		// update the company admin user entity
		(*companyAdminUserEntity).Id = registerCompanyAdminUserResponse.User.Id
		(*companyAdminUserEntity).Roles = registerCompanyAdminUserResponse.User.Roles
	}
}
