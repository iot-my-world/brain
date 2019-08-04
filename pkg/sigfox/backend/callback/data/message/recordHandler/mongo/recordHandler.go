package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	sigfoxBackendDataCallbackMessageRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler"
	sigfoxBackendDataCallbackMessageGenericRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler/generic"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) sigfoxBackendDataCallbackMessageRecordHandler.RecordHandler {
	mongoRecordHandler := brainMongoRecordHandler.New(
		mongoSession,
		databaseName,
		collectionName,
		[]mgo.Index{
			{
				Key:    []string{"id"},
				Unique: true,
			},
		},
		message.IsValidIdentifier,
		claims.ContextualiseFilter,
	)

	return sigfoxBackendDataCallbackMessageGenericRecordHandler.New(
		mongoRecordHandler,
	)
}
