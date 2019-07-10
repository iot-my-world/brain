package sigbug

import (
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client/basic"
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	sigbugAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator"
	sigbugJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator/jsonRpc"
	sigbugGPSReadings "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	sigbugJsonRpcRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/jsonRpc"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	partyAdministratorJsonRpc "github.com/iot-my-world/brain/pkg/party/administrator/jsonRpc"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	authorizationAdministrator "github.com/iot-my-world/brain/pkg/security/authorization/administrator"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	"github.com/stretchr/testify/suite"
)

func New(
	url string,
	user humanUser.User,
	testData []Data,
) *test {
	return &test{
		testData:      testData,
		user:          user,
		jsonRpcClient: basicJsonRpcClient.New(url),
	}
}

type test struct {
	suite.Suite
	jsonRpcClient       jsonRpcClient.Client
	user                humanUser.User
	testData            []Data
	sigbugAdministrator sigbugAdministrator.Administrator
	sigbugRecordHandler sigbugRecordHandler.RecordHandler
	partyAdministrator  partyAdministrator.Administrator
}

type Data struct {
	Device      sigbug.Sigbug
	GPSReadings []sigbugGPSReadings.Reading
}

func (suite *test) SetupTest() {

	// log in the client
	if err := suite.jsonRpcClient.Login(authorizationAdministrator.LoginRequest{
		UsernameOrEmailAddress: suite.user.Username,
		Password:               string(suite.user.Password),
	}); err != nil {
		suite.Fail("log in error", err.Error())
		return
	}

	// set up service provider clients that use jsonRpcClient
	suite.sigbugAdministrator = sigbugJsonRpcAdministrator.New(suite.jsonRpcClient)
	suite.sigbugRecordHandler = sigbugJsonRpcRecordHandler.New(suite.jsonRpcClient)
	suite.partyAdministrator = partyAdministratorJsonRpc.New(suite.jsonRpcClient)
}

func (suite *test) TestSigbug1Create() {
	// create all sigbugs in test data
	for _, data := range suite.testData {
		// retrieve the owner party
		retrieveOwnerPartyResponse, err := suite.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
			PartyType: data.Device.OwnerPartyType,
			Identifier: name.Identifier{
				Name: data.Device.OwnerId.Id,
			},
		})
		if err != nil {
			suite.FailNow("error retrieving owner party", err.Error())
			return
		}

		// set owner party id
		data.Device.OwnerId = retrieveOwnerPartyResponse.Party.Details().PartyId

		// retrieve assigned party if set
		if data.Device.AssignedId.Id != "" {
			// retrieve the assigned party
			retrieveAssignedPartyResponse, err := suite.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
				PartyType: data.Device.AssignedPartyType,
				Identifier: name.Identifier{
					Name: data.Device.AssignedId.Id,
				},
			})
			if err != nil {
				suite.FailNow("error retrieving assigned party", err.Error())
				return
			}

			// set owner party id
			data.Device.OwnerId = retrieveAssignedPartyResponse.Party.Details().PartyId
		}

		// create the device
		if _, err := suite.sigbugAdministrator.Create(&sigbugAdministrator.CreateRequest{
			Sigbug: data.Device,
		}); err != nil {
			suite.FailNow("error creating sigbug device", err.Error())
		}
	}
}
