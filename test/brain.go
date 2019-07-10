package test

import (
	clientTest "github.com/iot-my-world/brain/test/stories/client"
	companyTest "github.com/iot-my-world/brain/test/stories/company"
	publicTest "github.com/iot-my-world/brain/test/stories/public"
	systemTest "github.com/iot-my-world/brain/test/stories/system"
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

func (t *test) TestBrain() {
	suite.Run(t.T(), systemTest.New())
	suite.Run(t.T(), companyTest.New())
	suite.Run(t.T(), clientTest.New())
	suite.Run(t.T(), publicTest.New())
}
