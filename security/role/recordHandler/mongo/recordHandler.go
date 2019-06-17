package mongo

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/security/role"
	roleRecordHandler "github.com/iot-my-world/brain/security/role/recordHandler"
	roleRecordHandlerException "github.com/iot-my-world/brain/security/role/recordHandler/exception"
	roleSetup "github.com/iot-my-world/brain/security/role/setup"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
)

type recordHandler struct {
	mongoSession         *mgo.Session
	database, collection string
}

func New(mongoSession *mgo.Session, database, collection string) roleRecordHandler.RecordHandler {

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

func (r *recordHandler) Create(request *roleRecordHandler.CreateRequest) (*roleRecordHandler.CreateResponse, error) {

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(r.database).C(r.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return nil, brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.Role.Id = newId.String()

	if err := roleCollection.Insert(request.Role); err != nil {
		log.Error("Could not create Role! ", err)
		return nil, err
	}
	return &roleRecordHandler.CreateResponse{}, nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *roleRecordHandler.RetrieveRequest) error {
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
	}
	return nil
}

func (r *recordHandler) Retrieve(request *roleRecordHandler.RetrieveRequest) (*roleRecordHandler.RetrieveResponse, error) {
	if err := r.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(r.database).C(r.collection)

	var roleRecord role.Role

	if err := userCollection.Find(request.Identifier.ToFilter()).One(&roleRecord); err != nil {
		if err == mgo.ErrNotFound {
			return nil, roleRecordHandlerException.NotFound{}
		} else {
			return nil, brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	return &roleRecordHandler.RetrieveResponse{Role: roleRecord}, nil
}

func (r *recordHandler) Update(request *roleRecordHandler.UpdateRequest) (*roleRecordHandler.UpdateResponse, error) {

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(r.database).C(r.collection)

	// Retrieve role
	retrieveRoleResponse, err := r.Retrieve(&roleRecordHandler.RetrieveRequest{
		Identifier: request.Identifier,
	})
	if err != nil {
		return nil, roleRecordHandlerException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields
	// retrieveRoleResponse.Role.Id = request.Role.Id // cannot update ever
	// retrieveRoleResponse.Role.Name = request.Role.Name cannot update ever
	retrieveRoleResponse.Role.ViewPermissions = request.Role.ViewPermissions
	retrieveRoleResponse.Role.APIPermissions = request.Role.APIPermissions

	if err := roleCollection.Update(request.Identifier.ToFilter(), retrieveRoleResponse.Role); err != nil {
		return nil, roleRecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return &roleRecordHandler.UpdateResponse{}, nil
}
