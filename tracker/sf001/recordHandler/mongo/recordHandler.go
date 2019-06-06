package mongo

import (
	brainMongoRecordHandler "gitlab.com/iotTracker/brain/recordHandler/mongo"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/sf001"
	sf001RecordHandler "gitlab.com/iotTracker/brain/tracker/sf001/recordHandler"
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
				Key:    []string{"imei"},
				Unique: true, // Prevent two documents from having the same index key
			},
		},
		sf001.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return sf001RecordHandler.New(
		mongoRecordHandler,
	)
}
