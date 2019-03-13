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
	"gopkg.in/mgo.v2"
)

type mongoRecordHandler struct {
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

	return &mongoRecordHandler{
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

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *readingRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (mrh *mongoRecordHandler) Create(request *readingRecordHandler.CreateRequest, response *readingRecordHandler.CreateResponse) error {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	readingCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	newID, err := uuid.NewV4()
	if err != nil {
		return brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.Reading.Id = newID.String()

	if err := readingCollection.Insert(request.Reading); err != nil {
		return err
	}

	response.Reading = request.Reading
	return nil
}

func (mrh *mongoRecordHandler) ValidateCollectRequest(request *readingRecordHandler.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (mrh *mongoRecordHandler) Collect(request *readingRecordHandler.CollectRequest, response *readingRecordHandler.CollectResponse) error {
	if err := mrh.ValidateCollectRequest(request); err != nil {
		return err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = claims.ContextualiseFilter(filter, request.Claims)

	// Get Reading Collection
	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()
	readingCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Perform Query
	query := readingCollection.Find(filter)

	// Apply the count
	if total, err := query.Count(); err == nil {
		response.Total = total
	} else {
		return err
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
		return err
	}

	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *readingRecordHandler.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}
	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}
	readingValidateResponse := readingRecordHandler.ValidateResponse{}
	if err := mrh.Validate(&readingRecordHandler.ValidateRequest{
		Claims:  request.Claims,
		Reading: request.Reading,
	}, &readingValidateResponse); err != nil {
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

func (mrh *mongoRecordHandler) Update(request *readingRecordHandler.UpdateRequest, respose *readingRecordHandler.UpdateResponse) error {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return err
	}

	return nil
}

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *readingRecordHandler.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (mrh *mongoRecordHandler) Validate(request *readingRecordHandler.ValidateRequest, response *readingRecordHandler.ValidateResponse) error {
	if err := mrh.ValidateValidateRequest(request); err != nil {
		return err
	}

	return nil
}
