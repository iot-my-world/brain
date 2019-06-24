package mongo

import (
	"github.com/iot-my-world/brain/party/individual"
	individualRecordHandler "github.com/iot-my-world/brain/party/individual/recordHandler"
	brainMongoRecordHandler "github.com/iot-my-world/brain/recordHandler/mongo"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *individualRecordHandler.RecordHandler {
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
		individual.IsValidIdentifier,
		individual.ContextualiseFilter,
	)

	return individualRecordHandler.New(
		mongoRecordHandler,
	)
}
