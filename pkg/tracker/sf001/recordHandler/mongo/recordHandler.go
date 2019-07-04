package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	"github.com/iot-my-world/brain/pkg/security/claims"
	sf0012 "github.com/iot-my-world/brain/pkg/tracker/sf001"
	recordHandler2 "github.com/iot-my-world/brain/pkg/tracker/sf001/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/sf001/recordHandler/generic"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) recordHandler2.RecordHandler {
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
				Key:    []string{"deviceId"},
				Unique: true,
			},
		},
		sf0012.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return recordHandler.New(
		mongoRecordHandler,
	)
}
