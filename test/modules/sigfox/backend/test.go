package backend

import (
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client/basic"
	jsonRpcServerAuthenticator "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	partyJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/party/administrator/jsonRpc"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/query"
	sigfoxBackend "github.com/iot-my-world/brain/pkg/sigfox/backend"
	sigfoxBackendAdministrator "github.com/iot-my-world/brain/pkg/sigfox/backend/administrator"
	sigfoxBackendJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/sigfox/backend/administrator/jsonRpc"
	sigfoxBackendRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler"
	sigfoxBackendJsonRpcRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler/jsonRpc"
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
	jsonRpcClient              jsonRpcClient.Client
	user                       humanUser.User
	testData                   []Data
	sigfoxBackendAdministrator sigfoxBackendAdministrator.Administrator
	sigfoxBackendRecordHandler sigfoxBackendRecordHandler.RecordHandler
	partyAdministrator         partyAdministrator.Administrator
}

type Data struct {
	Backend sigfoxBackend.Backend
}

func (suite *test) SetupTest() {

	// log in the client
	if err := suite.jsonRpcClient.Login(jsonRpcServerAuthenticator.LoginRequest{
		UsernameOrEmailAddress: suite.user.Username,
		Password:               string(suite.user.Password),
	}); err != nil {
		suite.Fail("log in error", err.Error())
		return
	}

	// set up service provider clients that use jsonRpcClient
	suite.sigfoxBackendAdministrator = sigfoxBackendJsonRpcAdministrator.New(suite.jsonRpcClient)
	suite.sigfoxBackendRecordHandler = sigfoxBackendJsonRpcRecordHandler.New(suite.jsonRpcClient)
	suite.partyAdministrator = partyJsonRpcAdministrator.New(suite.jsonRpcClient)
}

func (suite *test) TestSigfoxBackend1Create() {
	// get logged in party's details
	getMyPartyResponse, err := suite.partyAdministrator.GetMyParty(&partyAdministrator.GetMyPartyRequest{})
	if err != nil {
		suite.FailNow("error getting my party details", err.Error())
		return
	}

	// create all sigfoxBackends in test data
	for _, data := range suite.testData {
		// set owner party details on the backend
		data.Backend.OwnerPartyType = getMyPartyResponse.Party.Details().PartyType
		data.Backend.OwnerId = getMyPartyResponse.Party.Details().PartyId

		// create the device
		if _, err := suite.sigfoxBackendAdministrator.Create(&sigfoxBackendAdministrator.CreateRequest{
			Backend: data.Backend,
		}); err != nil {
			suite.FailNow("error creating sigfox backend", err.Error())
		}
	}

	// collect all sigfoxBackends
	sigfoxBackendCollectResponse, err := suite.sigfoxBackendRecordHandler.Collect(&sigfoxBackendRecordHandler.CollectRequest{
		Criteria: make([]criterion.Criterion, 0),
		Query:    query.Query{},
	})
	if err != nil {
		suite.Failf("collect sigfoxBackends failed", err.Error())
		return
	}

	// confirm that each created sigfoxBackend can be found
nextSigfoxBackendToCreate:
	// for every sigfoxBackend that should be created
	for _, sigfoxBackendToCreate := range suite.testData {
		// look for sigfoxBackendToCreate among collected sigfoxBackends
		for _, existingSigfoxBackend := range sigfoxBackendCollectResponse.Records {
			if sigfoxBackendToCreate.Backend.Name == existingSigfoxBackend.Name {
				// update fields set during creation
				sigfoxBackendToCreate.Backend.Id = existingSigfoxBackend.Id
				sigfoxBackendToCreate.Backend.Token = existingSigfoxBackend.Token
				// assert should be equal
				suite.Equal(sigfoxBackendToCreate.Backend, existingSigfoxBackend, "created sigfoxBackend should be equal")
				// if it is found and equal, check for next sigfoxBackend to create
				continue nextSigfoxBackendToCreate
			}
		}
		// if execution reaches here then sigfoxBackendToCreate was not found among collected sigfoxBackends
	}
}
