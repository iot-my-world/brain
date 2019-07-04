package mongo

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	reading2 "github.com/iot-my-world/brain/pkg/tracker/tk102/reading"
	recordHandler2 "github.com/iot-my-world/brain/pkg/tracker/tk102/reading/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/tk102/reading/recordHandler/exception"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
)

type recordHandler struct {
	mongoSession *mgo.Session
	database     string
	collection   string
}

// New reading mongo record handler
func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
) recordHandler2.RecordHandler {

	setupIndices(mongoSession, database, collection)

	return &recordHandler{
		mongoSession: mongoSession,
		database:     database,
		collection:   collection,
	}
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise Reading collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	readingCollection := mgoSesh.DB(database).C(collection)

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := readingCollection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}

}

func (r *recordHandler) ValidateCreateRequest(request *recordHandler2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Create(request *recordHandler2.CreateRequest) (*recordHandler2.CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	readingCollection := mgoSession.DB(r.database).C(r.collection)

	newID, err := uuid.NewV4()
	if err != nil {
		return nil, brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.Reading.Id = newID.String()

	if err := readingCollection.Insert(request.Reading); err != nil {
		return nil, err
	}

	return &recordHandler2.CreateResponse{Reading: request.Reading}, nil
}

func (r *recordHandler) ValidateCreateBulkRequest(request *recordHandler2.CreateBulkRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *recordHandler) CreateBulk(request *recordHandler2.CreateBulkRequest) (*recordHandler2.CreateBulkResponse, error) {
	if err := r.ValidateCreateBulkRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	readingCollection := mgoSession.DB(r.database).C(r.collection)
	readingBulkOperation := readingCollection.Bulk()

	// docs interface to hold readings
	var readings []interface{}

	// generate uuid for all readings
	for readingIdx := range request.Readings {
		newID, err := uuid.NewV4()
		if err != nil {
			return nil, brainException.UUIDGeneration{Reasons: []string{err.Error()}}
		}
		request.Readings[readingIdx].Id = newID.String()
		readings = append(readings, request.Readings[readingIdx])
	}

	// queue insert operations
	readingBulkOperation.Insert(readings...)
	if _, err := readingBulkOperation.Run(); err != nil {
		return nil, exception.CreateBulk{Reason: err.Error()}
	}

	return &recordHandler2.CreateBulkResponse{
		Readings: request.Readings,
	}, nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *recordHandler2.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !reading2.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for reading", request.Identifier.Type()))
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

	readingCollection := mgoSession.DB(r.database).C(r.collection)

	var readingRecord reading2.Reading

	filter := request.Identifier.ToFilter()
	filter = claims.ContextualiseFilter(filter, request.Claims)

	if err := readingCollection.Find(filter).One(&readingRecord); err != nil {
		if err == mgo.ErrNotFound {
			return nil, exception.NotFound{}
		}
		return nil, brainException.Unexpected{Reasons: []string{err.Error()}}
	}

	return &recordHandler2.RetrieveResponse{Reading: readingRecord}, nil
}

func (r *recordHandler) ValidateCollectRequest(request *recordHandler2.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Collect(request *recordHandler2.CollectRequest) (*recordHandler2.CollectResponse, error) {
	if err := r.ValidateCollectRequest(request); err != nil {
		return nil, err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = claims.ContextualiseFilter(filter, request.Claims)

	response := recordHandler2.CollectResponse{}

	// Get Reading Collection
	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()
	readingCollection := mgoSession.DB(r.database).C(r.collection)

	// Perform Query
	query := readingCollection.Find(filter)

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
	response.Records = make([]reading2.Reading, 0)
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *recordHandler) ValidateUpdateRequest(request *recordHandler2.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else if !reading2.IsValidIdentifier(request.Identifier) {
		reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for reading", request.Identifier.Type()))
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Update(request *recordHandler2.UpdateRequest) (*recordHandler2.UpdateResponse, error) {
	if err := r.ValidateUpdateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	readingCollection := mgoSession.DB(r.database).C(r.collection)

	// Retrieve reading
	retrieveReadingResponse, err := r.Retrieve(&recordHandler2.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		return nil, exception.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields
	// retrieveReadingResponse.Reading.Id = request.Reading.Id // cannot update ever
	// retrieveReadingResponse.DeviceId = request.Reading.DeviceId // cannot update ever
	// retrieveReadingResponse.Reading.DeviceType = request.Reading.DeviceType // cannot update ever
	retrieveReadingResponse.Reading.OwnerPartyType = request.Reading.OwnerPartyType
	retrieveReadingResponse.Reading.OwnerId = request.Reading.OwnerId
	retrieveReadingResponse.Reading.AssignedPartyType = request.Reading.AssignedPartyType
	retrieveReadingResponse.Reading.AssignedId = request.Reading.AssignedId
	// retrieveReadingResponse.Reading.Raw = request.Reading.Raw // cannot update ever
	// retrieveReadingResponse.Reading.TimeStamp = request.Reading.TimeStamp // cannot update ever
	// retrieveReadingResponse.Reading.Latitude = request.Reading.Latitude // cannot update ever
	// retrieveReadingResponse.Reading.Longitude = request.Reading.Longitude // cannot update ever

	filter := request.Identifier.ToFilter()
	filter = claims.ContextualiseFilter(filter, request.Claims)
	if err := readingCollection.Update(filter, retrieveReadingResponse.Reading); err != nil {
		return nil, exception.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return &recordHandler2.UpdateResponse{Reading: retrieveReadingResponse.Reading}, nil
}
