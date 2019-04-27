package mongo

import (
	brainMongoRecordHandler "gitlab.com/iotTracker/brain/recordHandler/mongo"
	"gitlab.com/iotTracker/brain/security/claims"
	zx303GPSReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps"
	zx303GPSReadigRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps/recordHandler"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *zx303GPSReadigRecordHandler.RecordHandler {
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

	return zx303GPSReadigRecordHandler.New(
		mongoRecordHandler,
	)
}
