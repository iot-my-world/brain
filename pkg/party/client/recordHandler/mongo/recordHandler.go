package mongo

import (
	client2 "github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	recordHandler2 "github.com/iot-my-world/brain/pkg/party/client/recordHandler/generic"
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
		client2.IsValidIdentifier,
		client2.ContextualiseFilter,
	)

	return recordHandler2.New(
		mongoRecordHandler,
	)
}
