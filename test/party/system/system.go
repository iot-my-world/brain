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
