package mongo

import (
	"fmt"
	"github.com/satori/go.uuid"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/reading"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
	readingRecordHandlerException "gitlab.com/iotTracker/brain/tracker/reading/recordHandler/exception"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
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
) readingRecordHandler.RecordHandler {

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

func (r *recordHandler) ValidateCreateRequest(request *readingRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		readingValidateResponse, err := r.Validate(&readingRecordHandler.ValidateRequest{
			Claims:  request.Claims,
			Reading: request.Reading,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating reading: "+err.Error())
		}
		if len(readingValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range readingValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("invalid reading: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Create(request *readingRecordHandler.CreateRequest) (*readingRecordHandler.CreateResponse, error) {
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

	return &readingRecordHandler.CreateResponse{Reading: request.Reading}, nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *readingRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !reading.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for reading", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Retrieve(request *readingRecordHandler.RetrieveRequest) (*readingRecordHandler.RetrieveResponse, error) {
	if err := r.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	readingCollection := mgoSession.DB(r.database).C(r.collection)

	var readingRecord reading.Reading

	filter := request.Identifier.ToFilter()
	filter = claims.ContextualiseFilter(filter, request.Claims)

	if err := readingCollection.Find(filter).One(&readingRecord); err != nil {
		if err == mgo.ErrNotFound {
			return nil, readingRecordHandlerException.NotFound{}
		}
		return nil, brainException.Unexpected{Reasons: []string{err.Error()}}
	}

	return &readingRecordHandler.RetrieveResponse{Reading: readingRecord}, nil
}

func (r *recordHandler) ValidateCollectRequest(request *readingRecordHandler.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Collect(request *readingRecordHandler.CollectRequest) (*readingRecordHandler.CollectResponse, error) {
	if err := r.ValidateCollectRequest(request); err != nil {
		return nil, err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = claims.ContextualiseFilter(filter, request.Claims)

	response := readingRecordHandler.CollectResponse{}

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
	response.Records = make([]reading.Reading, 0)
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *recordHandler) ValidateUpdateRequest(request *readingRecordHandler.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else if !reading.IsValidIdentifier(request.Identifier) {
		reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for reading", request.Identifier.Type()))
	}

	readingValidateResponse, err := r.Validate(&readingRecordHandler.ValidateRequest{
		Claims:  request.Claims,
		Reading: request.Reading,
	})
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "error validating reading: "+err.Error())
	}
	if len(readingValidateResponse.ReasonsInvalid) > 0 {
		for _, reason := range readingValidateResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("invalid reading: %s - %s - %s", reason.Field, reason.Type, reason.Help))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Update(request *readingRecordHandler.UpdateRequest) (*readingRecordHandler.UpdateResponse, error) {
	if err := r.ValidateUpdateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	readingCollection := mgoSession.DB(r.database).C(r.collection)

	// Retrieve reading
	retrieveReadingResponse, err := r.Retrieve(&readingRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		return nil, readingRecordHandlerException.Update{Reasons: []string{"retrieving record", err.Error()}}
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
		return nil, readingRecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return &readingRecordHandler.UpdateResponse{Reading: retrieveReadingResponse.Reading}, nil
}

func (r *recordHandler) ValidateValidateRequest(request *readingRecordHandler.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *recordHandler) Validate(request *readingRecordHandler.ValidateRequest) (*readingRecordHandler.ValidateResponse, error) {
	if err := r.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	reasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)

	return &readingRecordHandler.ValidateResponse{
		ReasonsInvalid: reasonsInvalid,
	}, nil
}
