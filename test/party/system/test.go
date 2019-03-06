package system

import (
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	companyRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/recordHandler/adaptor/jsonRpc"
	partyRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	companyTest "gitlab.com/iotTracker/brain/test/party/company"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"encoding/json"
	"strings"
	"gitlab.com/iotTracker/brain/security/claims"
	"fmt"
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

	suite.T().Log("system successfully logged in")
}

func (suite *System) TestCreateCompanies() {
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
		companyData := &companyTest.EntitiesAndAdminUsersToCreate[idx]

		// update the new company's details as would be done from the front end
		(*companyData).Company.ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
		(*companyData).Company.ParentId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

		// create the company
		companyCreateResponse := companyRecordHandlerJsonRpcAdaptor.CreateResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"CompanyRecordHandler.Create",
			companyRecordHandlerJsonRpcAdaptor.CreateRequest{
				Company: (*companyData).Company,
			},
			&companyCreateResponse,
		); err != nil {
			suite.FailNow("create company failed", err.Error())
		}

		// update the company
		(*companyData).Company.Id = companyCreateResponse.Company.Id

		suite.T().Logf("successfully created company %s", (*companyData).Company.Name)

		// create an identifier for the company entity
		companyIdentifier, err := wrappedIdentifier.WrapIdentifier(id.Identifier{
			Id: (*companyData).Company.Id,
		})
		if err != nil {
			suite.FailNow("creating wrapped company identifier failed", err.Error())
		}

		// invite the admin user
		inviteCompanyAdminUserResponse := partyRegistrarJsonRpcAdaptor.InviteCompanyAdminUserResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"PartyRegistrar.InviteCompanyAdminUser",
			partyRegistrarJsonRpcAdaptor.InviteCompanyAdminUserRequest{
				PartyIdentifier: *companyIdentifier,
			},
			&inviteCompanyAdminUserResponse,
		); err != nil {
			suite.FailNow("invite company admin user failed", err.Error())
		}

		// parse the token into register companyAdminUserClaims
		jwt := inviteCompanyAdminUserResponse.URLToken[strings.Index(inviteCompanyAdminUserResponse.URLToken, "&t=")+3:]
		object, err := jose.ParseSigned(jwt)
		if err != nil {
			suite.FailNow("error parsing jwt", err.Error())
		}

		// Access Underlying payload without verification
		fv := reflect.ValueOf(object).Elem().FieldByName("payload")

		wrapped := wrappedClaims.WrappedClaims{}
		if err := json.Unmarshal(fv.Bytes(), &wrapped); err != nil {
			suite.FailNow("error unmarshalling claims", err.Error())
		}

		unwrappedClaims, err := wrapped.Unwrap()
		if err != nil {
			suite.FailNow("error unwrapping claims", err.Error())
		}

		if !suite.Equal(claims.RegisterCompanyAdminUser, unwrappedClaims.Type(), "claims should be "+claims.RegisterCompanyAdminUser) {
			suite.FailNow(fmt.Sprintf("claims are not of type %s", claims.RegisterCompanyAdminUser))
		}

	}
}
