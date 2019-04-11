package mongo

import (
	"gitlab.com/iotTracker/brain/party/company"
	deviceRecordHandler "gitlab.com/iotTracker/brain/tracker/device/recordHandler"
	brainMongoRecordHandler "gitlab.com/iotTracker/brain/recordHandler/mongo"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *deviceRecordHandler.RecordHandler {
	mongoRecordHandler := brainMongoRecordHandler.New(
		mongoSession,
		databaseName,
		collectionName,
		[]mgo.Index{{
			Key:    []string{"id"},
			Unique: true,
		}},
		company.IsValidIdentifier,
		company.ContextualiseFilter,
	)

	return deviceRecordHandler.New(
		mongoRecordHandler,
	)
}
