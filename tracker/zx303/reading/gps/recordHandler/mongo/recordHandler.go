package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/recordHandler/mongo"
	"github.com/iot-my-world/brain/security/claims"
	zx303GPSReading "github.com/iot-my-world/brain/tracker/zx303/reading/gps"
	zx303GPSReadingRecordHandler "github.com/iot-my-world/brain/tracker/zx303/reading/gps/recordHandler"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *zx303GPSReadingRecordHandler.RecordHandler {
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
		zx303GPSReading.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return zx303GPSReadingRecordHandler.New(
		mongoRecordHandler,
	)
}
