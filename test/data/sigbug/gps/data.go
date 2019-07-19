package gps

import (
	"fmt"
	sigbugGPSReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	"github.com/iot-my-world/brain/pkg/workbook"
	"os"
)

type Data struct {
	Reading     sigbugGPSReading.Reading
	DataMessage string
}

func GetData() (error, map[string][]Data) {
	dataWorkbook, err := workbook.New(
		fmt.Sprintf("%s/src/github.com/iot-my-world/brain/test/data/sigbug/gps/data.data.xlsx", os.Getenv("GOPATH")),
		map[string]int{},
	)
	if err != nil {

	}
}
