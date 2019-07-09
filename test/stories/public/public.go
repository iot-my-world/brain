package public

import (
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	"github.com/iot-my-world/brain/test/data"
	partyRegistrarAdministratorTestModule "github.com/iot-my-world/brain/test/modules/party/registrarAdministrator"
	humanUserAdministratorTestModule "github.com/iot-my-world/brain/test/modules/user/human/administrator"
	"github.com/stretchr/testify/suite"
)

func New() *test {
	return &test{}
}

type test struct {
	suite.Suite
}

func (t *test) SetupTest() {

}

func (t *test) TestPublic() {
	suite.Run(t.T(), partyRegistrarAdministratorTestModule.New(
		data.BrainURL,
		CompanyTestData,
		ClientTestData,
	))

	for _, companyData := range CompanyTestData {
		suite.Run(t.T(), humanUserAdministratorTestModule.New(
			data.BrainURL,
			[]humanUser.User{companyData.AdminUser},
		))
		suite.Run(t.T(), humanUserAdministratorTestModule.New(
			data.BrainURL,
			companyData.Users,
		))
	}

	for _, clientData := range ClientTestData {
		suite.Run(t.T(), humanUserAdministratorTestModule.New(
			data.BrainURL,
			[]humanUser.User{clientData.AdminUser},
		))
		suite.Run(t.T(), humanUserAdministratorTestModule.New(
			data.BrainURL,
			clientData.Users,
		))
	}
}
