package ship

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

	return &NewMongoRecordHandler
}

func setupRecords(mongoSession *mgo.Session, database, collection string){
	//Initialise record collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	roleCollection := mgoSesh.DB(database).C(collection)

	//Ensure index Uniqueness
	uniqueIndex := mgo.Index{
		Key: []string{"id"},
		Unique: true,
	}
	if err := roleCollection.EnsureIndex(uniqueIndex); err != nil {
		log.Fatal("Could not ensure uniqueness on id in ship collection: ", err)
	}
}


func validateServiceReqData(r interface{}) (error) {
	var reasonsInvalid []string

	switch v := r.(type){
	case *CreateRequest:
		if v.Ship.Name == "" {
			reasonsInvalid = append(reasonsInvalid, "Name cannot be blank")
		}
		if v.Ship.Berth == "" {
			reasonsInvalid = append(reasonsInvalid, "Berth cannot be blank")
		}
		// TODO: ensure outDateTime > inDateTime
	case *RetrieveAllRequest:
	case *UpdateRequest:
		if v.Ship.Id == "" {
			reasonsInvalid = append(reasonsInvalid, "ID cannot be blank")
		}
	case *DeleteRequest:
		if v.Ship.Id == "" {
			reasonsInvalid = append(reasonsInvalid, "ID cannot be blank")
		}
	default:
		log.Warn("NO CHECK CASE FOR THIS REQUEST!")
	}

	if len(reasonsInvalid) > 0 {
		return errors.New(strings.Join(reasonsInvalid, ","))
	}
	return nil
}

func (u *mongoRecordHandler) Create(request *CreateRequest, response *CreateResponse) error {

	err := validateServiceReqData(request)
	if err != nil {
		return err
	}

	request.Ship.Id = bson.NewObjectId().Hex()

	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(u.database).C(u.collection)

	err = roleCollection.Insert(request.Ship)

	if err != nil {
		return err //TODO: Translate Unknown error
	}

	return nil
}

func (u *mongoRecordHandler) RetrieveAll(request *RetrieveAllRequest, response *RetrieveAllResponse) error {
	if err := validateServiceReqData(request); err != nil {
		return err
	}
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	collection := mgoSession.DB(u.database).C(u.collection)

	var records []Ship

	if err := collection.Find(bson.M{"deleted": false}).All(&records); err != nil {
		return err
	}

	response.Records = records
	return nil
}

func updateAllowedFields (recordToUpdate *Ship, requestRecord Ship) error {
	if recordToUpdate.Deleted {
		return errors.New("cannot update deleted record")
	}

	recordToUpdate.Name = requestRecord.Name
	recordToUpdate.Berth = requestRecord.Berth
	recordToUpdate.InDateTime = requestRecord.InDateTime
	recordToUpdate.OutDateTime = requestRecord.OutDateTime

	return nil
}

func (u *mongoRecordHandler) Update(request *UpdateRequest, response *UpdateResponse) error {
	if err := validateServiceReqData(request); err != nil {
		return err
	}
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(u.database).C(u.collection)

	//Try and retrieve existing record
	existingRecord := Ship{}
	if err := collection.Find(bson.M{"id": request.Ship.Id}).One(&existingRecord); err != nil {
		return err
	}

	// Update allowed fields
	if err := updateAllowedFields(&existingRecord, request.Ship); err != nil {
		return err
	}

	// Try insert update
	if err := collection.Update(bson.M{"id": request.Ship.Id}, existingRecord); err != nil {
		return err
	}

	response.Ship = existingRecord

	return nil
}

func (u *mongoRecordHandler) Delete(request *DeleteRequest, response *DeleteResponse) error {
	if err := validateServiceReqData(request); err != nil {
		return err
	}

	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(u.database).C(u.collection)

	//Try and retrieve existing record
	existingRecord := Ship{}
	if err := collection.Find(bson.M{"id": request.Ship.Id}).One(&existingRecord); err != nil {
		return err
	}

	// Mark Ship Record as Deleted
	existingRecord.Deleted = true

	// Try make Update
	if err := collection.Update(bson.M{"id": request.Ship.Id}, existingRecord); err != nil {
		return err
	}

	response.Ship = existingRecord

	return nil
}