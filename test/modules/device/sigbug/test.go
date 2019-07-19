package sigbug

import (
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client/basic"
	jsonRpcServerAuthenticator "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator"
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
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	"github.com/stretchr/testify/suite"
)

func New(
	humanUserURL string,
	sigfoxBackendAPIUserURL string,
	user humanUser.User,
	testData []Data,
	sigfoxBackendJWT string,
) *test {
	return &test{
		testData:                          testData,
		user:                              user,
		humanUserJsonRpcClient:            basicJsonRpcClient.New(humanUserURL),
		sigfoxBackendAPIUserJsonRpcClient: basicJsonRpcClient.New(sigfoxBackendAPIUserURL),
		sigfoxBackendAPIUserURL:           sigfoxBackendAPIUserURL,
		sigfoxBackendJWT:                  sigfoxBackendJWT,
	}
}

type test struct {
	suite.Suite
	humanUserJsonRpcClient            jsonRpcClient.Client
	sigfoxBackendAPIUserJsonRpcClient jsonRpcClient.Client
	user                              humanUser.User
	testData                          []Data
	sigbugAdministrator               sigbugAdministrator.Administrator
	sigbugRecordHandler               sigbugRecordHandler.RecordHandler
	partyAdministrator                partyAdministrator.Administrator
	sigfoxBackendAPIUserURL           string
	sigfoxBackendJWT                  string
}

type Data struct {
	Device      sigbug.Sigbug
	GPSReadings []sigbugGPSReadings.Reading
}

func (suite *test) SetupTest() {

	// log in the human user client
	if err := suite.humanUserJsonRpcClient.Login(jsonRpcServerAuthenticator.LoginRequest{
		UsernameOrEmailAddress: suite.user.Username,
		Password:               string(suite.user.Password),
	}); err != nil {
		suite.FailNow("log in error", err.Error())
		return
	}

	// set token for api user json rpc client
	if err := suite.sigfoxBackendAPIUserJsonRpcClient.SetJWT(suite.sigfoxBackendJWT); err != nil {
		suite.FailNow("error setting api user json rpc client", err.Error())
		return
	}

	// set up service provider clients that use jsonRpcClient
	suite.sigbugAdministrator = sigbugJsonRpcAdministrator.New(suite.humanUserJsonRpcClient)
	suite.sigbugRecordHandler = sigbugJsonRpcRecordHandler.New(suite.humanUserJsonRpcClient)
	suite.partyAdministrator = partyAdministratorJsonRpc.New(suite.humanUserJsonRpcClient)
}

func (suite *test) TestSigbug1Create() {
	// create all sigbugs in test data
	for idx := range suite.testData {
		// retrieve the owner party
		retrieveOwnerPartyResponse, err := suite.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
			PartyType: suite.testData[idx].Device.OwnerPartyType,
			Identifier: name.Identifier{
				Name: suite.testData[idx].Device.OwnerId.Id,
			},
		})
		if err != nil {
			suite.FailNow("error retrieving owner party", err.Error())
			return
		}

		// set owner party id
		suite.testData[idx].Device.OwnerId = retrieveOwnerPartyResponse.Party.Details().PartyId

		// retrieve assigned party if set
		if suite.testData[idx].Device.AssignedId.Id != "" {
			// retrieve the assigned party
			retrieveAssignedPartyResponse, err := suite.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
				PartyType: suite.testData[idx].Device.AssignedPartyType,
				Identifier: name.Identifier{
					Name: suite.testData[idx].Device.AssignedId.Id,
				},
			})
			if err != nil {
				suite.FailNow("error retrieving assigned party", err.Error())
				return
			}

			// set owner party id
			suite.testData[idx].Device.AssignedId = retrieveAssignedPartyResponse.Party.Details().PartyId
		}

		// create the device
		if _, err := suite.sigbugAdministrator.Create(&sigbugAdministrator.CreateRequest{
			Sigbug: suite.testData[idx].Device,
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
