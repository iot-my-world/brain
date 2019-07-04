package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	"github.com/iot-my-world/brain/security/claims"
	zx303Task "github.com/iot-my-world/brain/tracker/zx303/task"
	zx303TaskRecordHandler "github.com/iot-my-world/brain/tracker/zx303/task/recordHandler"
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
