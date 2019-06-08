package system

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/authorization/service/adaptor/jsonRpc"
	testData "gitlab.com/iotTracker/brain/test/data"
	systemTestData "gitlab.com/iotTracker/brain/test/system/data"
	"gitlab.com/iotTracker/brain/tracker/sf001"
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
	pathToDataWorkbook := os.Getenv("GOPATH") + "/src/gitlab.com/iotTracker/brain/test/tracker/device/tk102/data/deviceData.xlsx"

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
		newSF001Tracker := sf001.SF001{
			Id:                   "",
			DeviceId:             rowMap["DeviceId"],
			OwnerPartyType:       party.Type(rowMap["OwnerPartyType"]),
			OwnerId:              id.Identifier{},
			AssignedPartyType:    party.Type(rowMap["AssignedPartyType"]),
			AssignedId:           id.Identifier{},
			LastMessageTimestamp: 0,
		}
		fmt.Println("make!", newSF001Tracker)
	}
}
