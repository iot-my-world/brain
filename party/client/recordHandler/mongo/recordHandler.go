package mongo

import (
	"gitlab.com/iotTracker/brain/party/client"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	brainMongoRecordHandler "gitlab.com/iotTracker/brain/recordHandler/mongo"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *clientRecordHandler.RecordHandler {
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

	return clientRecordHandler.New(
		mongoRecordHandler,
	)
}
