package public

import (
	"github.com/iot-my-world/brain/test/data"
	partyRegistrarAdministratorTestModule "github.com/iot-my-world/brain/test/modules/party/registrarAdministrator"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, partyRegistrarAdministratorTestModule.New(
		data.BrainURL,
		CompanyTestData,
		ClientTestData,
	))
}
