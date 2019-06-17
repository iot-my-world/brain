package mongo

import (
	brainMongoRecordHandler "gitlab.com/iotTracker/brain/recordHandler/mongo"
	"gitlab.com/iotTracker/brain/security/claims"
	zx303Task "gitlab.com/iotTracker/brain/tracker/zx303/task"
	zx303TaskRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/task/recordHandler"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *zx303TaskRecordHandler.RecordHandler {
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
		zx303Task.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return zx303TaskRecordHandler.New(
		mongoRecordHandler,
	)
}
