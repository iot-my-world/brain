package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	"github.com/iot-my-world/brain/pkg/user/api"
	"github.com/iot-my-world/brain/pkg/user/api/recordHandler"
	recordHandler2 "github.com/iot-my-world/brain/pkg/user/api/recordHandler/generic"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) recordHandler.RecordHandler {
	mongoRecordHandler := brainMongoRecordHandler.New(
		mongoSession,
		databaseName,
		collectionName,
		[]mgo.Index{
			{
				Key:    []string{"id"},
				Unique: true,
			},
			{
				Key: []string{"username"},
			},
		},
		api.IsValidIdentifier,
		api.ContextualiseFilter,
	)

	return recordHandler2.New(
		mongoRecordHandler,
	)
}
