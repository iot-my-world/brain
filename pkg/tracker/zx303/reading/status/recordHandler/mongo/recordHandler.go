package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/recordHandler"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *recordHandler.RecordHandler {
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
		status.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return recordHandler.New(
		mongoRecordHandler,
	)
}
