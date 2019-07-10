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
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	"github.com/iot-my-world/brain/pkg/search/query"
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
			data.Device.AssignedId = retrieveAssignedPartyResponse.Party.Details().PartyId
		}

		// create the device
		if _, err := suite.sigbugAdministrator.Create(&sigbugAdministrator.CreateRequest{
			Sigbug: data.Device,
		}); err != nil {
			suite.FailNow("error creating sigbug device", err.Error())
		}
	}

	// collect all sigbugs
	sigbugCollectResponse, err := suite.sigbugRecordHandler.Collect(&sigbugRecordHandler.CollectRequest{
		Criteria: make([]criterion.Criterion, 0),
		Query:    query.Query{},
	})
	if err != nil {
		suite.Failf("collect sigbugs failed", err.Error())
		return
	}

	// confirm that each created sigbug can be found
nextSigbugToCreate:
	// for every sigbug that should be created
	for _, sigbugToCreate := range suite.testData {
		// look for sigbugToCreate among collected sigbugs
		for _, existingSigbug := range sigbugCollectResponse.Records {
			if sigbugToCreate.Device.DeviceId == existingSigbug.DeviceId {
				// update fields set during creation
				sigbugToCreate.Device.Id = existingSigbug.Id
				// assert should be equal
				suite.Equal(sigbugToCreate.Device, existingSigbug, "created sigbug should be equal")
				// if it is found and equal, check for next sigbug to create
				continue nextSigbugToCreate
			}
		}
		// if execution reaches here then sigbugToCreate was not found among collected sigbugs
	}
}
