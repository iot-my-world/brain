package role

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/log"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	globalException "gitlab.com/iotTracker/brain/exception"
	roleException "gitlab.com/iotTracker/brain/security/role/exception"
	"gitlab.com/iotTracker/brain/security"
	"gitlab.com/iotTracker/brain/security/role"
)

type recordHandler struct {
	mongoSession         *mgo.Session
	database, collection string
}

func New(mongoSession *mgo.Session, database, collection string) *recordHandler {

	setupIndices(mongoSession, database, collection)

	NewMongoRecordHandler := recordHandler{
		mongoSession: mongoSession,
		database:     database,
		collection:   collection,
	}

	if err := role.InitialSetup(&NewMongoRecordHandler); err != nil {
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

func (mrh *recordHandler) Create(request *role.CreateRequest, response *role.CreateResponse) error {

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	request.Role.Id = bson.NewObjectId().Hex()

	err := roleCollection.Insert(request.Role)

	if err != nil {
		log.Error("Could not create Role! ", err)
		return err //TODO: Translate Unknown error
	}
	return nil
}

func (mrh *recordHandler) ValidateRetrieveRequest(request *role.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !role.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for role", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *recordHandler) Retrieve(request *role.RetrieveRequest, response *role.RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var roleRecord security.Role

	if err := userCollection.Find(request.Identifier.ToFilter()).One(&roleRecord); err != nil {
		if err == mgo.ErrNotFound {
			return roleException.NotFound{}
		} else {
			return globalException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.Role = roleRecord
	return nil
}

func (mrh *recordHandler) Update(request *role.UpdateRequest, response *role.UpdateResponse) error {

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	err := roleCollection.Update(bson.M{"name": request.Role.Name}, request.Role)
	if err != nil {
		log.Error("Unable to update role!", err)
	}

	return nil
}