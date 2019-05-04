package mongo

import (
	brainMongoRecordHandler "gitlab.com/iotTracker/brain/recordHandler/mongo"
	"gitlab.com/iotTracker/brain/security/claims"
	zx303StatusReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/status"
	zx303StatusReadingRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/recordHandler"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *zx303StatusReadingRecordHandler.RecordHandler {
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
		zx303StatusReading.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return zx303StatusReadingRecordHandler.New(
		mongoRecordHandler,
	)
}
