package system

import (
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	"gitlab.com/iotTracker/brain/party"
	partyAdministratorJsonAdaptor "gitlab.com/iotTracker/brain/party/administrator/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/search/identifier/adminEmailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	testData "gitlab.com/iotTracker/brain/test/data"
	systemTest "gitlab.com/iotTracker/brain/test/system"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
	tk102DeviceAdministratorJsonAdaptor "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator/adaptor/jsonRpc"
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
		UsernameOrEmailAddress: systemTest.User.Username,
		Password:               string(systemTest.User.Password),
	}); err != nil {
		suite.Fail("log in error", err.Error())
	}
}

func (suite *System) TestDeviceCreation() {
	pathToDataWorkbook := os.Getenv("GOPATH") + "/src/gitlab.com/iotTracker/brain/test/device/data/deviceData.xlsx"

	var sheetHeaderRowMap = map[string]int{
		"TK102Devices": 1,
	}
	deviceDataWorkBook, err := workbook.New(pathToDataWorkbook, sheetHeaderRowMap)
	if err != nil {
		suite.FailNow("failed to create device data workbook", err.Error())
	}

	// convert sheet to slice of maps
	sheetSliceMap, err := deviceDataWorkBook.SheetAsSliceMap("TK102Devices")
	if err != nil {
		suite.FailNow("failed to get sheet slice map", err.Error())
	}

	// create all of the devices
	for _, rowMap := range sheetSliceMap {
		// create new device
		newDevice := tk102.TK102{
			Id:                "",
			ManufacturerId:    rowMap["ManufacturerId"],
			SimCountryCode:    rowMap["SimCountryCode"],
			SimNumber:         rowMap["SimNumber"],
			OwnerPartyType:    party.Type(rowMap["OwnerPartyType"]),
			OwnerId:           id.Identifier{},
			AssignedPartyType: party.Type(rowMap["AssignedPartyType"]),
			AssignedId:        id.Identifier{},
		}

		// create identifier to retrieve the owner party
		ownerPartyIdentifier, err := wrappedIdentifier.WrapIdentifier(adminEmailAddress.Identifier{AdminEmailAddress: rowMap["Owner Admin Email"]})
		if err != nil {
			suite.FailNow("error wrapping party Identifier", err.Error())
		}

		// try and retrieve the owner party
		retrieveOwnerPartyResponse := partyAdministratorJsonAdaptor.RetrievePartyResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"PartyAdministrator.RetrieveParty",
			partyAdministratorJsonAdaptor.RetrievePartyRequest{
				PartyType:  newDevice.OwnerPartyType,
				Identifier: *ownerPartyIdentifier,
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
		newDevice.OwnerPartyType = unwrappedOwnerParty.Details().PartyType
		newDevice.OwnerId = unwrappedOwnerParty.Details().PartyId

		// if there are assigned party details then retrieve the assigned party and populate for the device
		if newDevice.AssignedPartyType != "" {
			// create identifier to retrieve the assigned party
			assignedPartyIdentifier, err := wrappedIdentifier.WrapIdentifier(adminEmailAddress.Identifier{AdminEmailAddress: rowMap["Assigned Admin Email"]})
			if err != nil {
				suite.FailNow("error wrapping assigned party Identifier", err.Error())
			}

			// try and retrieve the assigned party
			retrieveAssignedPartyResponse := partyAdministratorJsonAdaptor.RetrievePartyResponse{}
			if err := suite.jsonRpcClient.JsonRpcRequest(
				"PartyAdministrator.RetrieveParty",
				partyAdministratorJsonAdaptor.RetrievePartyRequest{
					PartyType:  newDevice.AssignedPartyType,
					Identifier: *assignedPartyIdentifier,
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
			newDevice.AssignedPartyType = unwrappedAssignedParty.Details().PartyType
			newDevice.AssignedId = unwrappedAssignedParty.Details().PartyId
		}

		// create the device
		createDeviceResponse := tk102DeviceAdministratorJsonAdaptor.CreateResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"TK102DeviceAdministrator.Create",
			tk102DeviceAdministratorJsonAdaptor.CreateRequest{
				TK102: newDevice,
			},
			&createDeviceResponse,
		); err != nil {
			suite.FailNow("create device failed", err.Error())
		}
	}
}
