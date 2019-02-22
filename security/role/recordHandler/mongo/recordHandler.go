package mongo

import (
	"fmt"
	"github.com/satori/go.uuid"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/security/role"
	roleException "gitlab.com/iotTracker/brain/security/role/exception"
	roleRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler"
	roleSetup "gitlab.com/iotTracker/brain/security/role/setup"
	"gopkg.in/mgo.v2"
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

	if err := roleSetup.InitialSetup(&NewMongoRecordHandler); err != nil {
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

func (mrh *recordHandler) Create(request *roleRecordHandler.CreateRequest, response *roleRecordHandler.CreateResponse) error {

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.Role.Id = newId.String()

	if err := roleCollection.Insert(request.Role); err != nil {
		log.Error("Could not create Role! ", err)
		return err //TODO: Translate Unknown error
	}
	return nil
}

func (mrh *recordHandler) ValidateRetrieveRequest(request *roleRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !role.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for role", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *recordHandler) Retrieve(request *roleRecordHandler.RetrieveRequest, response *roleRecordHandler.RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var roleRecord role.Role

	if err := userCollection.Find(request.Identifier.ToFilter()).One(&roleRecord); err != nil {
		if err == mgo.ErrNotFound {
			return roleException.NotFound{}
		} else {
			return brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.Role = roleRecord
	return nil
}

func (mrh *recordHandler) Update(request *roleRecordHandler.UpdateRequest, response *roleRecordHandler.UpdateResponse) error {

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve role
	retrieveRoleResponse := roleRecordHandler.RetrieveResponse{}
	if err := mrh.Retrieve(&roleRecordHandler.RetrieveRequest{Identifier: request.Identifier}, &retrieveRoleResponse); err != nil {
		return roleException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields
	// retrieveRoleResponse.Role.Id = request.Role.Id // cannot update ever
	// retrieveRoleResponse.Role.Name = request.Role.Name cannot update ever
	retrieveRoleResponse.Role.ViewPermissions = request.Role.ViewPermissions
	retrieveRoleResponse.Role.APIPermissions = request.Role.APIPermissions

	if err := roleCollection.Update(request.Identifier.ToFilter(), retrieveRoleResponse.Role); err != nil {
		return roleException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return nil
}
