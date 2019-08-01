package main

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/test/data/sigbug/gps"
)

func main() {
	//gpsData, err := gps.GetData()
	_, err := gps.GetData()
	if err != nil {
		log.Error(err.Error())
		return
	}

}
