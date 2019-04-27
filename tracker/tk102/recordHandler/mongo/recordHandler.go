package mongo

import (
	"fmt"
	"github.com/satori/go.uuid"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/tk102"
	tk102RecordHandler "gitlab.com/iotTracker/brain/tracker/tk102/recordHandler"
	tk102RecordHandlerException "gitlab.com/iotTracker/brain/tracker/tk102/recordHandler/exception"
	"gopkg.in/mgo.v2"
)

type recordHandler struct {
	mongoSession *mgo.Session
	database     string
	collection   string
}

func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
) tk102RecordHandler.RecordHandler {

	setupIndices(mongoSession, database, collection)

	newTK102MongoRecordHandler := recordHandler{
		mongoSession: mongoSession,
		database:     database,
		collection:   collection,
	}

	return &newTK102MongoRecordHandler
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise TK102 collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	tk102Collection := mgoSesh.DB(database).C(collection)

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := tk102Collection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}

	// Ensure admin manufacturerIdUnique uniqueness
	manufacturerIdUnique := mgo.Index{
		Key:    []string{"manufacturerId"},
		Unique: true, // Prevent two documents from having the same index key
	}
	if err := tk102Collection.EnsureIndex(manufacturerIdUnique); err != nil {
		log.Fatal("Could not ensure manufacturerId uniqueness: ", err)
	}

	// Ensure country code + number uniqueness
	countryCodeNumberUnique := mgo.Index{
		Key:    []string{"simCountryCode", "simNumber"},
		Unique: true, // Prevent two documents from having the same index key
	}
	if err := tk102Collection.EnsureIndex(countryCodeNumberUnique); err != nil {
		log.Fatal("Could not ensure sim country code and number combination unique: ", err)
	}
}

func (mrh *recordHandler) ValidateCreateRequest(request *tk102RecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *recordHandler) Create(request *tk102RecordHandler.CreateRequest) (*tk102RecordHandler.CreateResponse, error) {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	tk102Collection := mgoSession.DB(mrh.database).C(mrh.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return nil, brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.TK102.Id = newId.String()

	if err := tk102Collection.Insert(request.TK102); err != nil {
		return nil, tk102RecordHandlerException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	return &tk102RecordHandler.CreateResponse{TK102: request.TK102}, nil
}

func (mrh *recordHandler) ValidateRetrieveRequest(request *tk102RecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !tk102.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for tk102", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *recordHandler) Retrieve(request *tk102RecordHandler.RetrieveRequest) (*tk102RecordHandler.RetrieveResponse, error) {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	tk102Collection := mgoSession.DB(mrh.database).C(mrh.collection)

	var tk102Record tk102.TK102

	filter := claims.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)
	if err := tk102Collection.Find(filter).One(&tk102Record); err != nil {
		if err == mgo.ErrNotFound {
			return nil, tk102RecordHandlerException.NotFound{}
		} else {
			return nil, brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	return &tk102RecordHandler.RetrieveResponse{TK102: tk102Record}, nil
}

func (mrh *recordHandler) ValidateUpdateRequest(request *tk102RecordHandler.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}
	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (mrh *recordHandler) Update(request *tk102RecordHandler.UpdateRequest) (*tk102RecordHandler.UpdateResponse, error) {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	tk102Collection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve TK102
	retrieveTK102Response, err := mrh.Retrieve(&tk102RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		return nil, tk102RecordHandlerException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveTK102Response.TK102.Id = request.TK102.Id // cannot update ever
	retrieveTK102Response.TK102.ManufacturerId = request.TK102.ManufacturerId
	retrieveTK102Response.TK102.SimCountryCode = request.TK102.SimCountryCode
	retrieveTK102Response.TK102.SimNumber = request.TK102.SimNumber
	retrieveTK102Response.TK102.OwnerPartyType = request.TK102.OwnerPartyType
	retrieveTK102Response.TK102.OwnerId = request.TK102.OwnerId
	retrieveTK102Response.TK102.AssignedPartyType = request.TK102.AssignedPartyType
	retrieveTK102Response.TK102.AssignedId = request.TK102.AssignedId

	if err := tk102Collection.Update(request.Identifier.ToFilter(), retrieveTK102Response.TK102); err != nil {
		return nil, tk102RecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return &tk102RecordHandler.UpdateResponse{TK102: retrieveTK102Response.TK102}, nil
}

func (mrh *recordHandler) ValidateDeleteRequest(request *tk102RecordHandler.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !tk102.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for tk102", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *recordHandler) Delete(request *tk102RecordHandler.DeleteRequest) (*tk102RecordHandler.DeleteResponse, error) {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	tk102Collection := mgoSession.DB(mrh.database).C(mrh.collection)

	if err := tk102Collection.Remove(request.Identifier.ToFilter()); err != nil {
		return nil, err
	}

	return &tk102RecordHandler.DeleteResponse{}, nil
}

func (mrh *recordHandler) ValidateCollectRequest(request *tk102RecordHandler.CollectRequest) error {
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

func (mrh *recordHandler) Collect(request *tk102RecordHandler.CollectRequest) (*tk102RecordHandler.CollectResponse, error) {
	if err := mrh.ValidateCollectRequest(request); err != nil {
		return nil, err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = claims.ContextualiseFilter(filter, request.Claims)

	response := tk102RecordHandler.CollectResponse{}

	// Get TK102 Collection
	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()
	tk102Collection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Perform Query
	query := tk102Collection.Find(filter)

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
	response.Records = make([]tk102.TK102, 0)
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return nil, err
	}

	return &response, nil
}
