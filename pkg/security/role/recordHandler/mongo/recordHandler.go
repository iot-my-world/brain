package mongo

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	role2 "github.com/iot-my-world/brain/pkg/security/role"
	recordHandler2 "github.com/iot-my-world/brain/pkg/security/role/recordHandler"
	"github.com/iot-my-world/brain/pkg/security/role/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/security/role/setup"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
)

type recordHandler struct {
	mongoSession         *mgo.Session
	database, collection string
}

func New(mongoSession *mgo.Session, database, collection string) recordHandler2.RecordHandler {

	setupIndices(mongoSession, database, collection)

	NewMongoRecordHandler := recordHandler{
		mongoSession: mongoSession,
		database:     database,
		collection:   collection,
	}

	if err := setup.InitialSetup(&NewMongoRecordHandler); err != nil {
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

func (r *recordHandler) Create(request *recordHandler2.CreateRequest) (*recordHandler2.CreateResponse, error) {

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
	return &recordHandler2.CreateResponse{}, nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *recordHandler2.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !role2.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for role", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Retrieve(request *recordHandler2.RetrieveRequest) (*recordHandler2.RetrieveResponse, error) {
	if err := r.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(r.database).C(r.collection)

	var roleRecord role2.Role

	if err := userCollection.Find(request.Identifier.ToFilter()).One(&roleRecord); err != nil {
		if err == mgo.ErrNotFound {
			return nil, exception.NotFound{}
		} else {
			return nil, brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	return &recordHandler2.RetrieveResponse{Role: roleRecord}, nil
}

func (r *recordHandler) Update(request *recordHandler2.UpdateRequest) (*recordHandler2.UpdateResponse, error) {

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(r.database).C(r.collection)

	// Retrieve role
	retrieveRoleResponse, err := r.Retrieve(&recordHandler2.RetrieveRequest{
		Identifier: request.Identifier,
	})
	if err != nil {
		return nil, exception.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields
	// retrieveRoleResponse.Role.Id = request.Role.Id // cannot update ever
	// retrieveRoleResponse.Role.Name = request.Role.Name cannot update ever
	retrieveRoleResponse.Role.ViewPermissions = request.Role.ViewPermissions
	retrieveRoleResponse.Role.APIPermissions = request.Role.APIPermissions

	if err := roleCollection.Update(request.Identifier.ToFilter(), retrieveRoleResponse.Role); err != nil {
		return nil, exception.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return &recordHandler2.UpdateResponse{}, nil
}
