package system

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	"gitlab.com/iotTracker/brain/search/identifier/device/tk102"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	testData "gitlab.com/iotTracker/brain/test/data"
	systemTestData "gitlab.com/iotTracker/brain/test/system/data"
	"gitlab.com/iotTracker/brain/tracker/device"
	tk102DeviceRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/tracker/reading"
	readingAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/reading/administrator/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/workbook"
	"math/rand"
	"os"
	"strconv"
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
		UsernameOrEmailAddress: systemTestData.User.Username,
		Password:               string(systemTestData.User.Password),
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
	noDevices := len(tk102DeviceSheetSliceMap)
	for tk102DeviceIdx, tk102DeviceRowMap := range tk102DeviceSheetSliceMap {
		// create identifier to retrieve the device
		deviceIdentifier, err := wrappedIdentifier.Wrap(tk102.Identifier{ManufacturerId: tk102DeviceRowMap["ManufacturerId"]})
		if err != nil {
			suite.FailNow("error wrapping device Identifier", err.Error())
		}

		// try and retrieve the device
		retrieveTK102DeviceResponse := tk102DeviceRecordHandlerJsonRpcAdaptor.RetrieveResponse{}
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"TK102DeviceRecordHandler.Retrieve",
			tk102DeviceRecordHandlerJsonRpcAdaptor.RetrieveRequest{
				Identifier: *deviceIdentifier,
			},
			&retrieveTK102DeviceResponse,
		); err != nil {
			suite.FailNow("retrieve device failed", err.Error())
		}

		// readings sheet for device
		sheetName := readingDataWorkBook.GetSheetNames()[rand.Intn(len(readingDataWorkBook.GetSheetNames()))]
		readingSheetSliceMap, err := readingDataWorkBook.SheetAsSliceMap(sheetName)
		if err != nil {
			suite.FailNow("error getting reading sheet as slice map", err.Error())
		}

		// prepare a batch of readings
		readingsToCreate := make([]reading.Reading, 0)
		for _, readingRow := range readingSheetSliceMap {
			lat, err := strconv.ParseFloat(readingRow["Lat"], 32)
			if err != nil {
				suite.FailNow("error parsing lat value to float", err.Error())
			}
			lon, err := strconv.ParseFloat(readingRow["Lon"], 32)
			if err != nil {
				suite.FailNow("error parsing lon value to float", err.Error())
			}
			timeStamp, err := strconv.ParseInt(readingRow["stamp"], 10, 64)
			if err != nil {
				suite.FailNow("error parsing stamp value to float", err.Error())
			}
			readingsToCreate = append(readingsToCreate, reading.Reading{
				//Id:                "",
				DeviceId:          id.Identifier{Id: retrieveTK102DeviceResponse.TK102.Id},
				DeviceType:        device.TK102,
				OwnerPartyType:    retrieveTK102DeviceResponse.TK102.OwnerPartyType,
				OwnerId:           retrieveTK102DeviceResponse.TK102.OwnerId,
				AssignedPartyType: retrieveTK102DeviceResponse.TK102.AssignedPartyType,
				AssignedId:        retrieveTK102DeviceResponse.TK102.AssignedId,
				Raw:               "_dummy_data_",
				TimeStamp:         timeStamp,
				Latitude:          float32(lat),
				Longitude:         float32(lon),
			})
		}

		// try and create the readings in bulk
		if err := suite.jsonRpcClient.JsonRpcRequest(
			"ReadingAdministrator.CreateBulk",
			readingAdministratorJsonRpcAdaptor.CreateBulkRequest{
				Readings: readingsToCreate,
			},
			&readingAdministratorJsonRpcAdaptor.CreateBulkResponse{},
		); err != nil {
			suite.FailNow("creating bulk readings failed", err.Error())
		}

		fmt.Printf("Device %d/%d\n", tk102DeviceIdx+1, noDevices)
	}
}
