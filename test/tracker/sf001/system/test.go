package system

import (
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	"gitlab.com/iotTracker/brain/party"
	partyAdministratorJsonAdaptor "gitlab.com/iotTracker/brain/party/administrator/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/search/identifier/adminEmailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/authorization/service/adaptor/jsonRpc"
	testData "gitlab.com/iotTracker/brain/test/data"
	systemTestData "gitlab.com/iotTracker/brain/test/system/data"
	"gitlab.com/iotTracker/brain/tracker/sf001"
	sf001TrackerAdministratorJsonAdaptor "gitlab.com/iotTracker/brain/tracker/sf001/administrator/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/workbook"
	"os"
)

type System struct {
	suite.Suite
	jsonRpcClient jsonRpcClient.Client
}

func (suite *System) SetupTest() {
	// create the client
	suite.jsonRpcClient = basicJsonRpcClient.New(testData.BrainURL)

	// log in the client
	if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
		UsernameOrEmailAddress: systemTestData.User.Username,
		Password:               string(systemTestData.User.Password),
	}); err != nil {
		suite.FailNow("log in error", err.Error())
	}
}

func (suite *System) TestSystemDeviceCreation() {
	pathToDataWorkbook := os.Getenv("GOPATH") + "/src/gitlab.com/iotTracker/brain/test/tracker/sf001/data/sf001TrackerTestData.xlsx"

	var sheetHeaderRowMap = map[string]int{
		"SF001Tracker": 1,
	}
	sf001TrackerWorkBook, err := workbook.New(pathToDataWorkbook, sheetHeaderRowMap)
	if err != nil {
		suite.FailNow("failed to create sf001 tracker workbook", err.Error())
	}

	// convert sheet to slice of maps
	sheetSliceMap, err := sf001TrackerWorkBook.SheetAsSliceMap("SF001Tracker")
	if err != nil {
		suite.FailNow("failed to get sheet slice map", err.Error())
	}

	// create all of the sf001 trackers
	for _, rowMap := range sheetSliceMap {
		// new tracker to create
		newSF001Tracker := sf001.SF001{
			Id:                   "",
			DeviceId:             rowMap["DeviceId"],
			OwnerPartyType:       party.Type(rowMap["OwnerPartyType"]),
			OwnerId:              id.Identifier{},
			AssignedPartyType:    party.Type(rowMap["AssignedPartyType"]),
			AssignedId:           id.Identifier{},
			LastMessageTimestamp: 0,
		}

		// create identifier to retrieve the owner party
		ownerPartyIdentifier, err := wrappedIdentifier.Wrap(adminEmailAddress.Identifier{AdminEmailAddress: rowMap["Owner Admin Email"]})
		if err != nil {
			suite.FailNow("error wrapping party Identifier", err.Error())
		}

		// try and retrieve the owner party
		retrieveOwnerPartyResponse := partyAdministratorJsonAdaptor.RetrievePartyResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"PartyAdministrator.RetrieveParty",
			partyAdministratorJsonAdaptor.RetrievePartyRequest{
				PartyType:         newSF001Tracker.OwnerPartyType,
				WrappedIdentifier: *ownerPartyIdentifier,
			},
			&retrieveOwnerPartyResponse,
		); err != nil {
			suite.FailNow("retrieve owner party failed", err.Error())
		}

		// unwrap the owner party from the response
		unwrappedOwnerParty, err := retrieveOwnerPartyResponse.Party.UnWrap()
		if err != nil {
			suite.FailNow("error unwrapping owner party", err.Error())
		}

		// populate the owner details
		newSF001Tracker.OwnerPartyType = unwrappedOwnerParty.Details().PartyType
		newSF001Tracker.OwnerId = unwrappedOwnerParty.Details().PartyId

		// if there are assigned party details then retrieve the assigned party and populate for the device
		if newSF001Tracker.AssignedPartyType != "" {
			// create identifier to retrieve the assigned party
			assignedPartyIdentifier, err := wrappedIdentifier.Wrap(adminEmailAddress.Identifier{AdminEmailAddress: rowMap["Assigned Admin Email"]})
			if err != nil {
				suite.FailNow("error wrapping assigned party Identifier", err.Error())
			}

			// try and retrieve the assigned party
			retrieveAssignedPartyResponse := partyAdministratorJsonAdaptor.RetrievePartyResponse{}
			if err := suite.jsonRpcClient.JsonRpcRequest(
				"PartyAdministrator.RetrieveParty",
				partyAdministratorJsonAdaptor.RetrievePartyRequest{
					PartyType:         newSF001Tracker.AssignedPartyType,
					WrappedIdentifier: *assignedPartyIdentifier,
				},
				&retrieveAssignedPartyResponse,
			); err != nil {
				suite.FailNow("retrieve assigned party failed", err.Error())
			}

			// unwrap the assigned party from the response
			unwrappedAssignedParty, err := retrieveAssignedPartyResponse.Party.UnWrap()
			if err != nil {
				suite.FailNow("error unwrapping assigned party", err.Error())
			}

			// populate the assigned details
			newSF001Tracker.AssignedPartyType = unwrappedAssignedParty.Details().PartyType
			newSF001Tracker.AssignedId = unwrappedAssignedParty.Details().PartyId
		}

		// create the device
		createSF001TrackerResponse := sf001TrackerAdministratorJsonAdaptor.CreateResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"SF001TrackerAdministrator.Create",
			sf001TrackerAdministratorJsonAdaptor.CreateRequest{
				SF001: newSF001Tracker,
			},
			&createSF001TrackerResponse,
		); err != nil {
			suite.FailNow("create sf001 tracker failed", err.Error())
		}
	}
}
