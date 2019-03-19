package system

import (
	"fmt"
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

		// create identifier to retrieve the party
		fmt.Println("Email:", rowMap["Owner Admin Email"])
		partyIdentifier, err := wrappedIdentifier.WrapIdentifier(adminEmailAddress.Identifier{AdminEmailAddress: rowMap["Owner Admin Email"]})
		if err != nil {
			suite.FailNow("error wrapping partyIdentifier", err.Error())
		}

		// try and retrieve the owner
		retrievePartyResponse := partyAdministratorJsonAdaptor.RetrievePartyResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"PartyAdministrator.RetrieveParty",
			partyAdministratorJsonAdaptor.RetrievePartyRequest{
				PartyType:  newDevice.OwnerPartyType,
				Identifier: *partyIdentifier,
			},
			&retrievePartyResponse,
		); err != nil {
			suite.FailNow("retrieve owner party failed", err.Error())
		}

		fmt.Println(retrievePartyResponse.Party)
	}
}
