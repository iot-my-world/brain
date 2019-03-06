package system

import (
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	companyRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/recordHandler/adaptor/jsonRpc"
	companyTest "gitlab.com/iotTracker/brain/test/party/company"
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

	for _, companyAdminUser := range companyTest.EntitiesAndAdminUsersToCreate {

		// update the new company's details as would be done from the front end
		companyAdminUser.Company.ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
		companyAdminUser.Company.ParentId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

		// create the company
		companyCreateResponse := companyRecordHandlerJsonRpcAdaptor.CreateResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"CompanyRecordHandler.Create",
			companyRecordHandlerJsonRpcAdaptor.CreateRequest{
				Company: companyAdminUser.Company,
			},
			&companyCreateResponse); err != nil {
			suite.FailNow("create company failed", err.Error())
		}

		suite.T().Logf("successfully created company %s", companyAdminUser.Company.Name)
	}
}
