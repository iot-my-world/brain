package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	"github.com/iot-my-world/brain/security/claims"
	zx303StatusReading "github.com/iot-my-world/brain/tracker/zx303/reading/status"
	zx303StatusReadingRecordHandler "github.com/iot-my-world/brain/tracker/zx303/reading/status/recordHandler"
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
