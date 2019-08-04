package gps

import (
	sigbugGPSReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
)

type Data struct {
	Reading     sigbugGPSReading.Reading
	DataMessage string
}
