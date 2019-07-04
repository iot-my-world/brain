package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	zx3032 "github.com/iot-my-world/brain/pkg/tracker/zx303"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/recordHandler"
	"github.com/iot-my-world/brain/security/claims"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *recordHandler.RecordHandler {
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
		zx3032.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return recordHandler.New(
		mongoRecordHandler,
	)
}
