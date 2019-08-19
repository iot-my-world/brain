package mongo

import (
	sigbugGPSReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	sigbugGPSReadingRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/recordHandler"
	sigbugGPSReadingGenericRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/recordHandler/generic"
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) sigbugGPSReadingRecordHandler.RecordHandler {
	mongoRecordHandler := brainMongoRecordHandler.New(
		mongoSession,
		databaseName,
		collectionName,
		[]mgo.Index{
			{
				Key:    []string{"id"},
				Unique: true,
			},
		},
		sigbugGPSReading.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return sigbugGPSReadingGenericRecordHandler.New(
		mongoRecordHandler,
	)
}
