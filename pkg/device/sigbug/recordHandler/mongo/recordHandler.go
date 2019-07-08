package mongo

import (
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	sigbugGenericRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/generic"
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) sigbugRecordHandler.RecordHandler {
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
		sigbug.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return sigbugGenericRecordHandler.New(
		mongoRecordHandler,
	)
}
