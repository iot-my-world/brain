package mongo

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/search/criterion"
	humanUser "github.com/iot-my-world/brain/user/human"
	userRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler"
	userRecordHandlerException "github.com/iot-my-world/brain/user/human/recordHandler/exception"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
)

type recordHandler struct {
	mongoSession *mgo.Session
	database     string
	collection   string
}

func New(
	mongoSession *mgo.Session,
	database,
	collection string,
) *recordHandler {

	setupIndices(mongoSession, database, collection)

	newUserMongoRecordHandler := recordHandler{
		mongoSession: mongoSession,
		database:     database,
		collection:   collection,
	}

	return &newUserMongoRecordHandler
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise User collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	userCollection := mgoSesh.DB(database).C(collection)

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := userCollection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}

	//Ensure emailAddress Uniqueness
	emailAddressUnique := mgo.Index{
		Key:    []string{"emailAddress"},
		Unique: true,
	}
	if err := userCollection.EnsureIndex(emailAddressUnique); err != nil {
		log.Fatal("Could not ensure email address uniqueness: ", err)
	}
}

func (r *recordHandler) ValidateCreateRequest(request *userRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *recordHandler) Create(request *userRecordHandler.CreateRequest) (*userRecordHandler.CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(r.database).C(r.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return nil, brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.User.Id = newId.String()

	if err := userCollection.Insert(request.User); err != nil {
		return nil, userRecordHandlerException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	return &userRecordHandler.CreateResponse{User: request.User}, nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *userRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !humanUser.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *recordHandler) Retrieve(request *userRecordHandler.RetrieveRequest) (*userRecordHandler.RetrieveResponse, error) {
	if err := r.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(r.database).C(r.collection)

	var userRecord humanUser.User

	filter := humanUser.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)
	if err := userCollection.Find(filter).One(&userRecord); err != nil {
		if err == mgo.ErrNotFound {
			return nil, userRecordHandlerException.NotFound{}
		} else {
			return nil, brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	return &userRecordHandler.RetrieveResponse{User: userRecord}, nil
}

func (r *recordHandler) ValidateUpdateRequest(request *userRecordHandler.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *recordHandler) Update(request *userRecordHandler.UpdateRequest) (*userRecordHandler.UpdateResponse, error) {
	if err := r.ValidateUpdateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(r.database).C(r.collection)

	// Retrieve User
	retrieveUserResponse, err := r.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		return nil, userRecordHandlerException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields
	// Note that all fields are updated here, higher level services which utilise this service
	// control how these updates are handled
	// retrieveUserResponse.User.Id = request.User.Id // cannot update ever
	retrieveUserResponse.User.Name = request.User.Name
	retrieveUserResponse.User.Surname = request.User.Surname
	retrieveUserResponse.User.Username = request.User.Username
	retrieveUserResponse.User.EmailAddress = request.User.EmailAddress
	retrieveUserResponse.User.Password = request.User.Password
	retrieveUserResponse.User.Roles = request.User.Roles
	retrieveUserResponse.User.ParentPartyType = request.User.ParentPartyType
	retrieveUserResponse.User.ParentId = request.User.ParentId
	retrieveUserResponse.User.PartyType = request.User.PartyType
	retrieveUserResponse.User.PartyId = request.User.PartyId
	retrieveUserResponse.User.Registered = request.User.Registered

	if err := userCollection.Update(request.Identifier.ToFilter(), retrieveUserResponse.User); err != nil {
		return nil, userRecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return &userRecordHandler.UpdateResponse{User: retrieveUserResponse.User}, nil
}

func (r *recordHandler) ValidateDeleteRequest(request *userRecordHandler.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !humanUser.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *recordHandler) Delete(request *userRecordHandler.DeleteRequest) (*userRecordHandler.DeleteResponse, error) {
	if err := r.ValidateDeleteRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(r.database).C(r.collection)

	filter := humanUser.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)
	if err := userCollection.Remove(filter); err != nil {
		return nil, err
	}

	return &userRecordHandler.DeleteResponse{}, nil
}

func (r *recordHandler) ValidateCollectRequest(request *userRecordHandler.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *recordHandler) Collect(request *userRecordHandler.CollectRequest) (*userRecordHandler.CollectResponse, error) {
	if err := r.ValidateCollectRequest(request); err != nil {
		return nil, err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = humanUser.ContextualiseFilter(filter, request.Claims)

	response := userRecordHandler.CollectResponse{}

	// Get User Collection
	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()
	userCollection := mgoSession.DB(r.database).C(r.collection)

	// Perform Query
	query := userCollection.Find(filter)

	// Apply the count
	if total, err := query.Count(); err == nil {
		response.Total = total
	} else {
		return nil, err
	}

	// Apply limit if applicable
	if request.Query.Limit > 0 {
		query.Limit(request.Query.Limit)
	}

	// Determine the Sort Order
	mongoSortOrder := request.Query.ToMongoSortFormat()

	// Populate records
	response.Records = make([]humanUser.User, 0)
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return nil, err
	}

	return &response, nil
}
