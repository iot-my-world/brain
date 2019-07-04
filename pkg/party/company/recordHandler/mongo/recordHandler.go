package mongo

import (
	company2 "github.com/iot-my-world/brain/pkg/party/company"
	"github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	recordHandler2 "github.com/iot-my-world/brain/pkg/party/company/recordHandler/generic"
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
		company2.IsValidIdentifier,
		company2.ContextualiseFilter,
	)

	return recordHandler2.New(
		mongoRecordHandler,
	)
}
