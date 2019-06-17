package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/recordHandler/mongo"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/tracker/sf001"
	sf001RecordHandler "github.com/iot-my-world/brain/tracker/sf001/recordHandler"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *sf001RecordHandler.RecordHandler {
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
				Key:    []string{"deviceId"},
				Unique: true,
			},
		},
		sf001.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return sf001RecordHandler.New(
		mongoRecordHandler,
	)
}
