package mongo

import (
	"github.com/iot-my-world/brain/party/individual"
	sf001RecordHandler "github.com/iot-my-world/brain/party/individual/recordHandler"
	brainMongoRecordHandler "github.com/iot-my-world/brain/recordHandler/mongo"
	"github.com/iot-my-world/brain/security/claims"
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
		individual.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return sf001RecordHandler.New(
		mongoRecordHandler,
	)
}
