package company

import (
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	clientTest "gitlab.com/iotTracker/brain/test/party/client"
	clientRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/client/recordHandler/adaptor/jsonRpc"
	"fmt"
)

type Company struct {
	suite.Suite
	jsonRpcClient jsonRpcClient.Client
}

func (suite *Company) SetupTest() {
	// create the client
	suite.jsonRpcClient = basicJsonRpcClient.New("http://localhost:9010/api")
}

func (suite *Company) TestCompanyCreateClients() {
	// for each test company entity
	for _, companyTestDataEntity := range EntitiesAndAdminUsersToCreate {
		// log in
		if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
			UsernameOrEmailAddress: companyTestDataEntity.AdminUser.Username,
			Password:               string(companyTestDataEntity.AdminUser.Password),
		}); err != nil {
			suite.FailNow(fmt.Sprintf("failed to log in as %s", companyTestDataEntity.AdminUser.Username), err.Error())
		}

		// for each client assigned to be owned by this company
		for idx := range clientTest.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name] {
			clientEntity := &clientTest.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name][idx].Client

			// update the entity
			(*clientEntity).ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
			(*clientEntity).ParentId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

			// create the client
			clientCreateResponse := clientRecordHandlerJsonRpcAdaptor.CreateResponse{}
			if err := suite.jsonRpcClient.JsonRpcRequest(
				"ClientRecordHandler.Create",
				clientRecordHandlerJsonRpcAdaptor.CreateRequest{
					Client: *clientEntity,
				},
				&clientCreateResponse,
			); err != nil {
				suite.FailNow("create client failed", err.Error())
			}

		}

		// log out
		suite.jsonRpcClient.Logout()
	}
}
