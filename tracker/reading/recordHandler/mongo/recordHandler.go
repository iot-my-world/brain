package mongo

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/log"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gopkg.in/mgo.v2/bson"
)

type mongoRecordHandler struct {
	mongoSession *mgo.Session
	database     string
	collection   string
}

func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
) *mongoRecordHandler {

	setupIndices(mongoSession, database, collection)

	return &mongoRecordHandler{
		mongoSession: mongoSession,
		database:     database,
		collection:   collection,
	}
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise Company collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	readingCollection := mgoSesh.DB(database).C(collection)

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := readingCollection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}

}

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *readingRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Create(request *readingRecordHandler.CreateRequest, response *readingRecordHandler.CreateResponse) error {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	readingCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	request.Reading.Id = bson.NewObjectId().Hex()

	if err := readingCollection.Insert(request.Reading); err != nil {
		return err
	}

	response.Reading = request.Reading
	return nil
}