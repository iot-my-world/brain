package mongo

import (
	"github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	clientGenericRecordHandler "github.com/iot-my-world/brain/pkg/party/client/recordHandler/generic"
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
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
		[]mgo.Index{{
			Key:    []string{"id"},
			Unique: true,
		}},
		client.IsValidIdentifier,
		client.ContextualiseFilter,
	)

	return clientGenericRecordHandler.New(
		mongoRecordHandler,
	)
}
