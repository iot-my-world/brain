package system

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	testData "gitlab.com/iotTracker/brain/test/data"
	systemTest "gitlab.com/iotTracker/brain/test/system"
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

	var sheetFirstRowMap = map[string]int{
		"TK102Devices": 1,
	}
	deviceDataWorkBook, err := workbook.New(pathToDataWorkbook, sheetFirstRowMap)
	if err != nil {
		suite.FailNow("failed to create device data workbook", err.Error())
	}

	sheetSliceMap, err := deviceDataWorkBook.SheetAsSliceMap("TK102Devices")
	if err != nil {
		suite.FailNow("failed to get sheet slice map", err.Error())
	}
	for _, rowMap := range sheetSliceMap {
		fmt.Println(rowMap)
	}
}
