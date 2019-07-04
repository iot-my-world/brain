package mongo

import (
	brainMongoRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler/mongo"
	humanUser "github.com/iot-my-world/brain/user/human"
	humanUserRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler"
	humanUserGenericRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler/generic"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) humanUserRecordHandler.RecordHandler {
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
				Key:    []string{"emailAddress"},
				Unique: true,
			},
		},
		humanUser.IsValidIdentifier,
		humanUser.ContextualiseFilter,
	)

	return humanUserGenericRecordHandler.New(
		mongoRecordHandler,
	)
}
