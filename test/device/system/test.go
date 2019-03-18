package system

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	testData "gitlab.com/iotTracker/brain/test/data"
	systemTest "gitlab.com/iotTracker/brain/test/system"
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

func ColumnHeaderMap(xlsxFile *excelize.File, sheet string) (map[string]string, error) {
	columnHeaderMap := make(map[string]string)

	rows := xlsxFile.GetRows(sheet)
	if len(rows) == 0 {
		return nil, errors.New("no rows in sheet")
	}
	for colIdx, colCell := range rows[0] {
		columnHeaderMap[colCell] = excelize.ToAlphaString(colIdx)
	}

	return columnHeaderMap, nil
}

func (suite *System) TestDeviceCreation() {
	pathToDataWorkbook := os.Getenv("GOPATH") + "/src/gitlab.com/iotTracker/brain/test/device/data/deviceData.xlsx"

	deviceDataWorkBook, err := excelize.OpenFile(pathToDataWorkbook)
	if err != nil {
		suite.FailNow("failed to open device data workbook", err.Error())
	}
	columnHeaderMap, err := ColumnHeaderMap(deviceDataWorkBook, "DevicesToCreate")
	if err != nil {
		suite.FailNow("failed to get header map of workbook", err.Error())
	}

	fmt.Println(columnHeaderMap)
}
