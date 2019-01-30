package role

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/log"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	globalException "gitlab.com/iotTracker/brain/exception"
	roleException "gitlab.com/iotTracker/brain/security/role/exception"
	"gitlab.com/iotTracker/brain/security"
)

type mongoRecordHandler struct {
	mongoSession         *mgo.Session
	database, collection string
}

func NewMongoRecordHandler(mongoSession *mgo.Session, database, collection string) *mongoRecordHandler {

	setupIndices(mongoSession, database, collection)

	NewMongoRecordHandler := mongoRecordHandler{
		mongoSession: mongoSession,
		database:     database,
		collection:   collection,
	}

	if err := initialRoleSetup(&NewMongoRecordHandler); err != nil {
		log.Fatal("Unable to complete Initial System Role Setup!", err)
	}

	return &NewMongoRecordHandler
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise record collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	roleCollection := mgoSesh.DB(database).C(collection)

	//Ensure index Uniqueness
	uniqueIndex := mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	}

	if err := roleCollection.EnsureIndex(uniqueIndex); err != nil {
		log.Fatal("Could not ensure uniqueness on name index in role collection: ", err)
	}

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := roleCollection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}
}

func (mrh *mongoRecordHandler) Create(request *CreateRequest, response *CreateResponse) error {

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	err := roleCollection.Insert(request.Role)

	if err != nil {
		log.Error("Could not create Role! ", err)
		return err //TODO: Translate Unknown error
	}
	return nil
}

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for role", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Retrieve(request *RetrieveRequest, response *RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var roleRecord security.Role

	if err := userCollection.Find(request.Identifier.ToMap()).One(&roleRecord); err != nil {
		if err == mgo.ErrNotFound {
			return roleException.NotFound{}
		} else {
			return globalException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.Role = roleRecord
	return nil
}

func (mrh *mongoRecordHandler) Update(request *UpdateRequest, response *UpdateResponse) error {

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	err := roleCollection.Update(bson.M{"name": request.Role.Name}, request.Role)
	if err != nil {
		log.Error("Unable to update role!", err)
	}

	return nil
}
