package system

import (
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	companyRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/recordHandler/adaptor/jsonRpc"
	partyRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	companyTest "gitlab.com/iotTracker/brain/test/party/company"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"encoding/json"
	"strings"
	"gitlab.com/iotTracker/brain/security/claims"
	"fmt"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyAdminUser"
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/party"
)

type System struct {
	suite.Suite
	jsonRpcClient jsonRpcClient.Client
}

func (suite *System) SetupTest() {
	// create the client
	suite.jsonRpcClient = basicJsonRpcClient.New("http://localhost:9010/api")

	// log in the client
	if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
		UsernameOrEmailAddress: User.Username,
		Password:               string(User.Password),
	}); err != nil {
		suite.Fail("log in error", err.Error())
	}
}

func (suite *System) TestSystemCreateCompanies() {
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

	for idx := range companyTest.EntitiesAndAdminUsersToCreate {
		companyEntity := &(companyTest.EntitiesAndAdminUsersToCreate[idx].Company)

		// update the new company's details as would be done from the front end
		(*companyEntity).ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
		(*companyEntity).ParentId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

		// create the company
		companyCreateResponse := companyRecordHandlerJsonRpcAdaptor.CreateResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"CompanyRecordHandler.Create",
			companyRecordHandlerJsonRpcAdaptor.CreateRequest{
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

func (suite *System) TestSystemInviteAndRegisterCompanyAdminUsers() {
	for idx := range companyTest.EntitiesAndAdminUsersToCreate {
		companyEntity := &(companyTest.EntitiesAndAdminUsersToCreate[idx].Company)

		// create the minimal admin user
		adminUser := user.User{
			EmailAddress:    companyEntity.AdminEmailAddress,
			ParentPartyType: suite.jsonRpcClient.Claims().PartyDetails().PartyType,
			ParentId:        suite.jsonRpcClient.Claims().PartyDetails().PartyId,
			PartyType:       party.Company,
			PartyId:         id.Identifier{Id: companyEntity.Id},
		}

		// invite the admin user
		inviteCompanyAdminUserResponse := partyRegistrarJsonRpcAdaptor.InviteCompanyAdminUserResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"PartyRegistrar.InviteCompanyAdminUser",
			partyRegistrarJsonRpcAdaptor.InviteCompanyAdminUserRequest{
				User: adminUser,
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
		wrapped := wrappedClaims.WrappedClaims{}
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
		companyAdminUserEntity := &companyTest.EntitiesAndAdminUsersToCreate[idx].AdminUser
		switch typedClaims := unwrappedClaims.(type) {
		case registerCompanyAdminUser.RegisterCompanyAdminUser:
			(*companyAdminUserEntity).EmailAddress = typedClaims.User.EmailAddress
			(*companyAdminUserEntity).ParentPartyType = typedClaims.User.ParentPartyType
			(*companyAdminUserEntity).ParentId = typedClaims.User.ParentId
			(*companyAdminUserEntity).PartyType = typedClaims.User.PartyType
			(*companyAdminUserEntity).PartyId = typedClaims.User.PartyId
		default:
			suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterCompanyAdminUser))
		}

		// create a new json rpc client to register the user with
		registerJsonRpcClient := basicJsonRpcClient.New("http://localhost:9010/api")
		if err := registerJsonRpcClient.SetJWT(jwt); err != nil {
			suite.FailNow("failed to set jwt in registration client", err.Error())
		}

		// register the company admin user
		registerCompanyAdminUserResponse := partyRegistrarJsonRpcAdaptor.RegisterCompanyAdminUserResponse{}
		password := string(companyAdminUserEntity.Password)
		(*companyAdminUserEntity).Password = []byte{}
		if err := registerJsonRpcClient.JsonRpcRequest(
			"PartyRegistrar.RegisterCompanyAdminUser",
			partyRegistrarJsonRpcAdaptor.RegisterCompanyAdminUserRequest{
				User:     *companyAdminUserEntity,
				Password: password,
			},
			&registerCompanyAdminUserResponse,
		); err != nil {
			suite.FailNow("error registering company admin user", err.Error())
		}

		// update the company admin user entity
		(*companyAdminUserEntity).Id = registerCompanyAdminUserResponse.User.Id
		(*companyAdminUserEntity).Roles = registerCompanyAdminUserResponse.User.Roles
		(*companyAdminUserEntity).Password = []byte(password)
	}
}
