package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/recordHandler/mongo"
	apiUser "github.com/iot-my-world/brain/user/api"
	apiUserRecordHandler "github.com/iot-my-world/brain/user/api/recordHandler"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *apiUserRecordHandler.RecordHandler {
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
				Key: []string{"username"},
			},
		},
		apiUser.IsValidIdentifier,
		apiUser.ContextualiseFilter,
	)

	return apiUserRecordHandler.New(
		mongoRecordHandler,
	)
}
