package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/sigfox/backend"
	backendRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler"
	backendGenericRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler/generic"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) backendRecordHandler.RecordHandler {
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
		},
		backend.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return backendGenericRecordHandler.New(
		mongoRecordHandler,
	)
}
