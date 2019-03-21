package system

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	"gitlab.com/iotTracker/brain/search/identifier/device/tk102"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	testData "gitlab.com/iotTracker/brain/test/data"
	systemTest "gitlab.com/iotTracker/brain/test/system"
	tk102DeviceRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/workbook"
	"math/rand"
	"os"
	"time"
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
	rand.Seed(time.Now().Unix())
}

func (suite *System) TestSystemReadingCreation() {
	pathToDeviceDataWorkbook := os.Getenv("GOPATH") + "/src/gitlab.com/iotTracker/brain/test/tracker/device/tk102/data/deviceData.xlsx"
	deviceDataWorkBook, err := workbook.New(pathToDeviceDataWorkbook, map[string]int{
		"TK102Devices": 1,
	})
	if err != nil {
		suite.FailNow("failed to create device data workbook", err.Error())
	}

	pathToReadingDataWorkbook := os.Getenv("GOPATH") + "/src/gitlab.com/iotTracker/brain/test/tracker/reading/data/data.xlsx"
	readingDataWorkBook, err := workbook.New(pathToReadingDataWorkbook, map[string]int{})
	if err != nil {
		suite.FailNow("failed to create device data workbook", err.Error())
	}
	for sheet := range readingDataWorkBook.SheetHeaderRowMap {
		readingDataWorkBook.SheetHeaderRowMap[sheet] = 1
	}

	// convert sheet to slice of maps
	tk102DeviceSheetSliceMap, err := deviceDataWorkBook.SheetAsSliceMap("TK102Devices")
	if err != nil {
		suite.FailNow("failed to get tk102 device sheet slice map", err.Error())
	}

	for _, tk102DeviceRowMap := range tk102DeviceSheetSliceMap {
		// create identifier to retrieve the device
		deviceIdentifier, err := wrappedIdentifier.WrapIdentifier(tk102.Identifier{ManufacturerId: tk102DeviceRowMap["ManufacturerId"]})
		if err != nil {
			suite.FailNow("error wrapping device Identifier", err.Error())
		}

		// try and retrieve the device
		retrieveResponse := tk102DeviceRecordHandlerJsonRpcAdaptor.RetrieveResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"TK102DeviceRecordHandler.Retrieve",
			tk102DeviceRecordHandlerJsonRpcAdaptor.RetrieveRequest{
				Identifier: *deviceIdentifier,
			},
			&retrieveResponse,
		); err != nil {
			suite.FailNow("retrieve device failed", err.Error())
		}

		// readings sheet for device
		sheetName := readingDataWorkBook.GetSheetNames()[rand.Intn(len(readingDataWorkBook.GetSheetNames()))]
		readingSheetSliceMap, err := readingDataWorkBook.SheetAsSliceMap(sheetName)
		if err != nil {
			suite.FailNow("error getting reading sheet as slice map", err.Error())
		}
		for _, readingRow := range readingSheetSliceMap {
			fmt.Println(readingRow)
		}
	}
}
