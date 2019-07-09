package system

import (
	"github.com/iot-my-world/brain/test/data"
	clientTestModule "github.com/iot-my-world/brain/test/modules/party/client"
	companyTestModule "github.com/iot-my-world/brain/test/modules/party/company"
	"github.com/iot-my-world/brain/test/stories/client"
	"github.com/iot-my-world/brain/test/stories/company"
	"github.com/stretchr/testify/suite"
)

func New() *test {
	return &test{}
}

type test struct {
	suite.Suite
	companyTestModule *companyTestModule.Test
}

func (t *test) SetupTest() {
	companyTestData := make([]companyTestModule.Data, 0)
	for _, companyData := range company.TestData {
		companyTestData = append(companyTestData, companyData.CompanyTestData)
	}
	companyTestSuite, err := companyTestModule.New(
		t.Suite,
		data.BrainURL,
		User,
		companyTestData,
	)
	if err != nil {
		t.FailNow("error creating company test module", err.Error())
		return
	}
	t.companyTestModule = companyTestSuite
}

func (t *test) TestSystem() {
	// perform system company tests
	t.Run("Company Create", t.companyTestModule.TestCompany1Create)
	t.Run("Company Update Allowed Fields", t.companyTestModule.TestCompany2UpdateAllowedFields)
	t.Run("Company Delete", t.companyTestModule.TestCompany3Delete)
	t.Run("Company Invite and Register Admin", t.companyTestModule.TestCompany4InviteAndRegisterAdmin)
	t.Run("Company Create Users", t.companyTestModule.TestCompany5CreateUsers)
	t.Run("Company Invite and Register Users", t.companyTestModule.TestCompany6InviteAndRegisterUsers)
	t.Run("Company Company User Login", t.companyTestModule.TestCompany7UserLogin)

	// perform system client tests
	clientData, found := client.TestData["root"]
	if !found {
		t.FailNow("root client data not found")
		return
	}

	clientTestData := make([]clientTestModule.Data, 0)
	for _, clientData := range clientData {
		clientTestData = append(clientTestData, clientData.ClientTestData)
	}
	suite.Run(t.T(), clientTestModule.New(
		data.BrainURL,
		User,
		clientTestData,
	))
}
