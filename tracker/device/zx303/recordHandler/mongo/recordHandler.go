package mongo

import (
	brainMongoRecordHandler "gitlab.com/iotTracker/brain/recordHandler/mongo"
	"gitlab.com/iotTracker/brain/tracker/device/zx303"
	zx303RecordHandler "gitlab.com/iotTracker/brain/tracker/device/zx303/recordHandler"
	"gopkg.in/mgo.v2"
)

func New(
	mongoSession *mgo.Session,
	databaseName string,
	collectionName string,
) *zx303RecordHandler.RecordHandler {
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
				Sparse: true, // Only index documents containing the Key fields
			},
		},
		zx303.IsValidIdentifier,
		zx303.ContextualiseFilter,
	)

	return zx303RecordHandler.New(
		mongoRecordHandler,
	)
}
