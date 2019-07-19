package public

import (
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	"github.com/iot-my-world/brain/test/data/environment"
	partyRegistrarAdministratorTestModule "github.com/iot-my-world/brain/test/modules/party/registrarAdministrator"
	humanUserAdministratorTestModule "github.com/iot-my-world/brain/test/modules/user/human/administrator"
	"github.com/iot-my-world/brain/test/stories/public/data"
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
		environment.BrainHumanUserURL,
		data.CompanyTestData,
		data.ClientTestData,
	))

	for _, companyData := range data.CompanyTestData {
		suite.Run(t.T(), humanUserAdministratorTestModule.New(
			environment.BrainHumanUserURL,
			[]humanUser.User{companyData.AdminUser},
		))
		suite.Run(t.T(), humanUserAdministratorTestModule.New(
			environment.BrainHumanUserURL,
			companyData.Users,
		))
	}

	for _, clientData := range data.ClientTestData {
		suite.Run(t.T(), humanUserAdministratorTestModule.New(
			environment.BrainHumanUserURL,
			[]humanUser.User{clientData.AdminUser},
		))
		suite.Run(t.T(), humanUserAdministratorTestModule.New(
			environment.BrainHumanUserURL,
			clientData.Users,
		))
	}
}
