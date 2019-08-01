package generator

import (
	"fmt"
	"github.com/go-errors/errors"
	sigbugGPSReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	"github.com/iot-my-world/brain/pkg/workbook"
	sigbugGPSTestData "github.com/iot-my-world/brain/test/data/sigbug/gps"
	"os"
	"strconv"
)

func Generate() (map[string][]sigbugGPSTestData.Data, error) {
	dataWorkbook, err := workbook.New(
		fmt.Sprintf("%s/src/github.com/iot-my-world/brain/test/data/sigbug/gps/data.xlsx", os.Getenv("GOPATH")),
		map[string]int{
			"dbnRbay":      0,
			"dbnJhb":       0,
			"dbnBloem":     0,
			"cptKmbrlyJhb": 0,
			"dbnCpt":       0,
			"homeToWork":   0,
		},
	)
	if err != nil {
		return nil, err
	}

	dataToReturn := make(map[string][]sigbugGPSTestData.Data)
	for _, sheetName := range dataWorkbook.GetSheetNames() {
		// create a slice entry at this sheet name if one has not yet been made
		if _, found := dataToReturn[sheetName]; !found {
			dataToReturn[sheetName] = make([]sigbugGPSTestData.Data, 0)
		}

		// get rows for this sheet
		rows, err := dataWorkbook.SheetAsSliceMap(sheetName)
		if err != nil {
			return nil, err
		}

		// populate data
		for rowIdx := range rows {
			dataMessage, found := rows[rowIdx]["messageData"]
			if !found {
				return nil, errors.New("message data entry not found in row")
			}
			latString, found := rows[rowIdx]["Lat"]
			if !found {
				return nil, errors.New("lat entry not found in row")
			}
			lat, err := strconv.ParseFloat(latString, 32)
			if err != nil {
				return nil, errors.New("error parsing lat: " + err.Error())
			}
			lonString, found := rows[rowIdx]["Lon"]
			if !found {
				return nil, errors.New("lon entry not found in row")
			}
			lon, err := strconv.ParseFloat(lonString, 32)
			if err != nil {
				return nil, errors.New("error parsing lon: " + err.Error())
			}

			dataToReturn[sheetName] = append(
				dataToReturn[sheetName],
				sigbugGPSTestData.Data{
					Reading: sigbugGPSReading.Reading{
						Latitude:  float32(lat),
						Longitude: float32(lon),
					},
					DataMessage: dataMessage,
				},
			)
		}
	}

	return dataToReturn, nil
}
