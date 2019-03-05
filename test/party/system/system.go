package system

import (
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	companyRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/recordHandler/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/party/company"
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
		UsernameOrEmailAddress: "root",
		Password:               "12345",
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
		suite.Fail("company collection is not empty")
	}

	companyCreateRequest := companyRecordHandlerJsonRpcAdaptor.CreateRequest{
		Company: company.Company{
			Name:              "Monteagle Logistics",
			AdminEmailAddress: "brbitzbussy@gmail.com",
			ParentPartyType:   suite.jsonRpcClient.Claims().PartyDetails().PartyType,
			ParentId:          suite.jsonRpcClient.Claims().PartyDetails().PartyId,
		},
	}
	companyCreateResponse := companyRecordHandlerJsonRpcAdaptor.CreateResponse{}

	if err := suite.jsonRpcClient.JsonRpcRequest(
		"CompanyRecordHandler.Create",
		companyCreateRequest,
		&companyCreateResponse); err != nil {
		suite.Failf("create company failed", err.Error())
	}
}
