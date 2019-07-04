package mongo

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	recordHandler2 "github.com/iot-my-world/brain/pkg/recordHandler"
	"github.com/iot-my-world/brain/pkg/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type recordHandler struct {
	mongoSession        *mgo.Session
	database            string
	collection          string
	validIdentifier     func(id identifier.Identifier) bool
	contextualiseFilter func(filter bson.M, claimsToAdd claims.Claims) bson.M
}

// New mongo record handler
func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
	uniqueIndexes []mgo.Index,
	validIdentifier func(id identifier.Identifier) bool,
	contextualiseFilter func(filter bson.M, claimsToAdd claims.Claims) bson.M,
) recordHandler2.RecordHandler {

	if contextualiseFilter == nil {
		contextualiseFilter = claims.ContextualiseFilter
	}

	setupIndices(mongoSession, database, collection, uniqueIndexes)
	newRecordHandler := recordHandler{
		mongoSession:        mongoSession,
		database:            database,
		collection:          collection,
		validIdentifier:     validIdentifier,
		contextualiseFilter: contextualiseFilter,
	}

	return &newRecordHandler
}

func setupIndices(mongoSession *mgo.Session, database, collectionName string, uniqueIndexes []mgo.Index) {

	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	collection := mgoSesh.DB(database).C(collectionName)

	for _, uI := range uniqueIndexes {
		if err := collection.EnsureIndex(uI); err != nil {
			log.Fatal(fmt.Sprintf("Could not ensure uniqueness in %s collection: %s", collectionName, err.Error()))
		}
	}
}

func (r *recordHandler) ValidateCreateRequest(request *recordHandler2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Create(request *recordHandler2.CreateRequest, response *recordHandler2.CreateResponse) error {
	if err := r.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	collection := mgoSession.DB(r.database).C(r.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}

	request.Entity.SetId(newId.String())

	if err := collection.Insert(request.Entity); err != nil {
		return exception.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	response.Entity = request.Entity

	return nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *recordHandler2.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !r.validIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for %s entity type", request.Identifier.Type(), r.collection))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Retrieve(request *recordHandler2.RetrieveRequest, response *recordHandler2.RetrieveResponse) error {
	if err := r.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	collection := mgoSession.DB(r.database).C(r.collection)

	filter := request.Identifier.ToFilter()
	filter = r.contextualiseFilter(filter, request.Claims)

	if err := collection.Find(filter).One(response.Entity); err != nil {
		if err == mgo.ErrNotFound {
			return exception.NotFound{}
		}
		return brainException.Unexpected{Reasons: []string{err.Error()}}
	}

	return nil
}

func (r *recordHandler) ValidateUpdateRequest(request *recordHandler2.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else if !r.validIdentifier(request.Identifier) {
		reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for %s entity", request.Identifier.Type(), r.collection))
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Update(request *recordHandler2.UpdateRequest, response *recordHandler2.UpdateResponse) error {
	if err := r.ValidateUpdateRequest(request); err != nil {
		return err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	collection := mgoSession.DB(r.database).C(r.collection)

	filter := request.Identifier.ToFilter()
	filter = r.contextualiseFilter(filter, request.Claims)

	if err := collection.Update(filter, request.Entity); err != nil {
		return exception.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return nil
}

func (r *recordHandler) ValidateDeleteRequest(request *recordHandler2.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !r.validIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for %s entity", request.Identifier.Type(), r.collection))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Delete(request *recordHandler2.DeleteRequest, response *recordHandler2.DeleteResponse) error {
	if err := r.ValidateDeleteRequest(request); err != nil {
		return err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	collection := mgoSession.DB(r.database).C(r.collection)

	filter := request.Identifier.ToFilter()
	filter = r.contextualiseFilter(filter, request.Claims)

	if err := collection.Remove(filter); err != nil {
		return err
	}

	return nil
}

func (r *recordHandler) ValidateCollectRequest(request *recordHandler2.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Criteria == nil {
		reasonsInvalid = append(reasonsInvalid, "criteria is nil")
	} else {
		for _, c := range request.Criteria {
			if c == nil {
				reasonsInvalid = append(reasonsInvalid, "a criterion is nil")
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Collect(request *recordHandler2.CollectRequest, response *recordHandler2.CollectResponse) error {
	if err := r.ValidateCollectRequest(request); err != nil {
		return err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = r.contextualiseFilter(filter, request.Claims)

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(r.database).C(r.collection)

	// Perform Query
	query := collection.Find(filter)

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
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(response.Records); err != nil {
		return err
	}

	return nil
}
