package mongo

import (
	"gitlab.com/iotTracker/brain/party/company"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	brainMongoRecordHandler "gitlab.com/iotTracker/brain/recordHandler/mongo"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *companyRecordHandler.RecordHandler {
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

	return companyRecordHandler.New(
		mongoRecordHandler,
	)
}
