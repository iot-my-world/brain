package main

import (
	"github.com/iot-my-world/brain/internal/log"
	sigbugGPSTestData "github.com/iot-my-world/brain/test/data/sigbug/gps"
	sigbugGPSTestDataGenerator "github.com/iot-my-world/brain/test/data/sigbug/gps/generator"
	"math"
)

const take = 10

func main() {
	gpsDataMap, err := sigbugGPSTestDataGenerator.Generate()
	if err != nil {
		log.Error(err)
		return
	}

	for journeyName := range gpsDataMap {
		journeyData := make([]sigbugGPSTestData.Data, 0)
		if len(gpsDataMap[journeyName]) > take {
			journeyData = gpsDataMap[journeyName]
			continue
		}
		for i := 0; i < take; i++ {
			journeyData = append(
				journeyData,
				gpsDataMap[journeyName][int(math.Ceil(float64(i*len(gpsDataMap[journeyName]))/float64(take)))],
			)
		}
	}
}
