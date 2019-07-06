package public

import (
	"github.com/iot-my-world/brain/test/data"
	publicTestModule "github.com/iot-my-world/brain/test/modules/public"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, publicTestModule.New(
		data.BrainURL,
		[]publicTestModule.CompanyData{},
		[]publicTestModule.ClientData{},
	))
}
