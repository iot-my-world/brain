package businessRole

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
	//Ensure index Uniqueness
	uniqueIndex := mgo.Index{
		Key: []string{"id"},
		Unique: true,
	}
	if err := roleCollection.EnsureIndex(uniqueIndex); err != nil {
		log.Fatal("Could not ensure uniqueness on id in businessRole collection: ", err)
	}
}


func validateServiceReqData(r interface{}) (error) {
	var reasonsInvalid []string

	switch v := r.(type){
	case *CreateRequest:
		if v.BusinessRole.Name == "" {
			reasonsInvalid = append(reasonsInvalid, "Name cannot be blank")
		}
	case *RetrieveAllRequest:
	case *UpdateRequest:
		if v.BusinessRole.Name == "" {
			reasonsInvalid = append(reasonsInvalid, "Name cannot be blank")
		}
		if v.BusinessRole.Id == "" {
			reasonsInvalid = append(reasonsInvalid, "ID cannot be blank")
		}
	case *DeleteRequest:
		if v.BusinessRole.Id == "" {
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

	request.BusinessRole.Id = bson.NewObjectId().Hex()

	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(u.database).C(u.collection)

	err = roleCollection.Insert(request.BusinessRole)

	if err != nil {
		log.Error("Could not create business role! ", err)
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

	var records []BusinessRole

	if err := collection.Find(bson.M{}).All(&records); err != nil {
		log.Warn("unable to retrieve businessRole records!")
		return err
	}

	response.Records = records
	return nil
}

func updateAllowedFields (recordToUpdate *BusinessRole, requestRecord BusinessRole) {
	recordToUpdate.Name = requestRecord.Name
	recordToUpdate.PayRate = requestRecord.PayRate
}

func (u *mongoRecordHandler) Update(request *UpdateRequest, response *UpdateResponse) error {
	if err := validateServiceReqData(request); err != nil {
		return err
	}
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(u.database).C(u.collection)

	//Try and retrieve existing record
	existingRecord := BusinessRole{}
	if err := collection.Find(bson.M{"id": request.BusinessRole.Id}).One(&existingRecord); err != nil {
		return err
	}

	// Update allowed fields
	updateAllowedFields(&existingRecord, request.BusinessRole)

	// Try insert update
	if err := collection.Update(bson.M{"id": request.BusinessRole.Id}, existingRecord); err != nil {
		return err
	}

	response.BusinessRole = existingRecord

	return nil
}

func (u *mongoRecordHandler) Delete(request *DeleteRequest, response *DeleteResponse) error {
	if err := validateServiceReqData(request); err != nil {
		return err
	}

	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(u.database).C(u.collection)

	if err := collection.Remove(bson.M{"id": request.BusinessRole.Id}); err != nil {
		return err
	}

	response.BusinessRole = request.BusinessRole

	return nil
}