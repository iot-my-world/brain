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

func ColumnHeaderMap(xlsxFile *excelize.File, sheet string, topRowIdx int) (map[string]string, error) {
	columnHeaderMap := make(map[string]string)
	rows := xlsxFile.GetRows(sheet)
	fmt.Println("norows", len(rows))
	if len(rows)-1 < topRowIdx {
		return nil, errors.New("not enough rows in sheet")
	}
	for colIdx, colCell := range rows[topRowIdx] {
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
	columnHeaderMap, err := ColumnHeaderMap(deviceDataWorkBook, "TK102Devices", 1)
	if err != nil {
		suite.FailNow("failed to get header map of workbook", err.Error())
	}

	fmt.Println(columnHeaderMap)
}
