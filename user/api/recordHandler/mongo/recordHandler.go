package mongo

import (
	brainMongoRecordHandler "gitlab.com/iotTracker/brain/recordHandler/mongo"
	apiUser "gitlab.com/iotTracker/brain/user/api"
	apiUserRecordHandler "gitlab.com/iotTracker/brain/user/api/recordHandler"
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
