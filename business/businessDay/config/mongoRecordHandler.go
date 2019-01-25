package config

import (
	"gopkg.in/mgo.v2"
	"bitbucket.org/gotimekeeper/log"
	"strings"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

type mongoRecordHandler struct{
	mongoSession *mgo.Session
	database, collection string
}

func NewMongoRecordHandler(mongoSession *mgo.Session, database, collection string) *mongoRecordHandler {

	setupRecords(mongoSession, database, collection)

	NewMongoRecordHandler := mongoRecordHandler{
		mongoSession,
		database,
		collection,
	}

	if err := initBusinessDayConfig(&NewMongoRecordHandler); err != nil {
		log.Fatal("Unable to complete initBusinessDayConfig", err)
	}

	return &NewMongoRecordHandler
}

func setupRecords(mongoSession *mgo.Session, database, collectionName string){
	//Initialise record collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	collection := mgoSesh.DB(database).C(collectionName)

	//Ensure index Uniqueness
	uniqueIndex := mgo.Index{
		Key: []string{"id"},
		Unique: true,
	}
	if err := collection.EnsureIndex(uniqueIndex); err != nil {
		log.Fatal("Could not ensure uniqueness on id in businessDayConfig collection: ", err)
	}
}

func initBusinessDayConfig(handler *mongoRecordHandler) error {
	mgoSesh := handler.mongoSession.Copy()
	defer mgoSesh.Close()
	collection := mgoSesh.DB(handler.database).C(handler.collection)

	initConfigRecord := Config{}
	if err := collection.Find(bson.M{}).One(&initConfigRecord); err != nil {
		// Assume this means the record doesn't exist yet, try insert it
		initConfig.Id = bson.NewObjectId().Hex()
		if err := collection.Insert(initConfig); err != nil {
			log.Fatal("Unable to create initial business day config record")
		}
		return nil
	}
	return nil
}


func validateServiceReqData(r interface{}) (error) {
	var reasonsInvalid []string

	//switch v := r.(type){
	//default:
	//case CreateRequest:
	//case UpdateRequest:
	//case RetrieveRequest:
	//	log.Warn("NO CHECK CASE FOR THIS REQUEST!")
	//}

	if len(reasonsInvalid) > 0 {
		return errors.New(strings.Join(reasonsInvalid, ","))
	}
	return nil
}

func (m *mongoRecordHandler) Create(request *CreateRequest, response *CreateResponse) error {

	err := validateServiceReqData(request)
	if err != nil {
		return err
	}

	request.BusinessDayConfig.Id = bson.NewObjectId().Hex()

	mgoSession := m.mongoSession.Copy()
	defer mgoSession.Close()

	collection := mgoSession.DB(m.database).C(m.collection)

	err = collection.Insert(request.BusinessDayConfig)

	if err != nil {
		log.Error("Could not create business day config", err)
		return err //TODO: Translate Unknown error
	}
	
	return nil
}
func (m *mongoRecordHandler) Update(request *UpdateRequest, response *UpdateResponse) error {
	return nil
}
func (m *mongoRecordHandler) Retrieve(request *RetrieveRequest, response *RetrieveResponse) error {
	err := validateServiceReqData(request)
	if err != nil {
		return err
	}

	mgoSession := m.mongoSession.Copy()
	defer mgoSession.Close()

	collection := mgoSession.DB(m.database).C(m.collection)

	latestConfig := Config{}

	if err := collection.Find(bson.M{}).One(&latestConfig); err != nil {
		return err
	}

	response.BusinessDayConfig = latestConfig

	return nil
}