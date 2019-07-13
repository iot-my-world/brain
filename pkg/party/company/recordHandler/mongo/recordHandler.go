package mongo

import (
	"github.com/iot-my-world/brain/pkg/party/company"
	"github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	companyGenericRecordHandler "github.com/iot-my-world/brain/pkg/party/company/recordHandler/generic"
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
		[]mgo.Index{
			{
				Key:    []string{"id"},
				Unique: true,
			},
			{
				Key:    []string{"name"},
				Unique: true,
			},
			{
				Key:    []string{"adminEmailAddress"},
				Unique: true,
			},
		},
		company.IsValidIdentifier,
		company.ContextualiseFilter,
	)

	return companyGenericRecordHandler.New(
		mongoRecordHandler,
	)
}
