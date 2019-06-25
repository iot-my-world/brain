package mongo

import (
	"github.com/iot-my-world/brain/party/company"
	companyRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler"
	companyGenericRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler/generic"
	brainMongoRecordHandler "github.com/iot-my-world/brain/recordHandler/mongo"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) companyRecordHandler.RecordHandler {
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

	return companyGenericRecordHandler.New(
		mongoRecordHandler,
	)
}
