package public

import (
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client/basic"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	partyJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/party/administrator/jsonRpc"
	"github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/party/company"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	"github.com/stretchr/testify/suite"
)

type test struct {
	suite.Suite
	jsonRpcClient      jsonRpcClient.Client
	partyAdministrator partyAdministrator.Administrator
	companyTestData    []CompanyData
	clientTestData     []ClientData
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
}
