package mongo

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/genRecordHandler"
	genRecordHandlerException "gitlab.com/iotTracker/brain/genRecordHandler/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gopkg.in/mgo.v2"
)

type MongoRecordHandler struct {
	mongoSession  *mgo.Session
	database      string
	collection    string
	uniqueIndexes []mgo.Index
}

// New mongo record handler
func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
	uniqueIndexes []mgo.Index,
) genRecordHandler.RecordHandler {

	setupIndices(mongoSession, database, collection, uniqueIndexes)
	newCompanyMongoRecordHandler := MongoRecordHandler{
		mongoSession:  mongoSession,
		database:      database,
		collection:    collection,
		uniqueIndexes: uniqueIndexes,
	}

	return &newCompanyMongoRecordHandler
}

func (mrh *MongoRecordHandler) Start() {
	setupIndices(mrh.mongoSession, mrh.database, mrh.collection, mrh.uniqueIndexes)
}

func setupIndices(mongoSession *mgo.Session, database, collection string, uniqueIndexes []mgo.Index) {
	//Initialise Company collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	companyCollection := mgoSesh.DB(database).C(collection)

	for _, uI := range uniqueIndexes {
		if err := companyCollection.EnsureIndex(uI); err != nil {
			log.Fatal("Could not ensure uniqueness: ", err)
		}
	}
}

func (mrh *MongoRecordHandler) ValidateCreateRequest(request *genRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (mrh *MongoRecordHandler) GCreate(request *genRecordHandler.CreateRequest) (*genRecordHandler.CreateResponse, error) {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	if err := request.Entity.SetId(); err != nil {
		return nil, brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}

	if err := companyCollection.Insert(request.Entity); err != nil {
		return nil, genRecordHandlerException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	return &genRecordHandler.CreateResponse{Entity: request.Entity}, nil
}

func (mrh *MongoRecordHandler) ValidateRetrieveRequest(request *genRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !company.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for company", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (mrh *MongoRecordHandler) GRetrieve(request *genRecordHandler.RetrieveRequest) (*genRecordHandler.RetrieveResponse, error) {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var entityRecord genRecordHandler.GenEntity

	filter := request.Identifier.ToFilter()
	filter = company.ContextualiseFilter(filter, request.Claims)

	if err := companyCollection.Find(filter).One(&entityRecord); err != nil {
		if err == mgo.ErrNotFound {
			return nil, genRecordHandlerException.NotFound{}
		}
		return nil, brainException.Unexpected{Reasons: []string{err.Error()}}
	}

	return &genRecordHandler.RetrieveResponse{Entity: entityRecord}, nil
}

func (mrh *MongoRecordHandler) ValidateUpdateRequest(request *genRecordHandler.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else if !company.IsValidIdentifier(request.Identifier) {
		reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for company", request.Identifier.Type()))
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (mrh *MongoRecordHandler) GUpdate(request *genRecordHandler.UpdateRequest) (*genRecordHandler.UpdateResponse, error) {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	//// Retrieve Company
	//retrieveCompanyResponse, err := mrh.Retrieve(&genRecordHandler.RetrieveRequest{
	//	Claims:     request.Claims,
	//	Identifier: request.Identifier,
	//})
	//if err != nil {
	//	return nil, genRecordHandlerException.Update{Reasons: []string{"retrieving record", err.Error()}}
	//}

	// Update fields:
	// retrieveCompanyResponse.Company.Id = request.Company.Id // cannot update ever
	//retrieveCompanyResponse.Company.Name = request.Company.Name //Update

	filter := request.Identifier.ToFilter()
	filter = company.ContextualiseFilter(filter, request.Claims)
	if err := companyCollection.Update(filter, request.Entity); err != nil {
		return nil, genRecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return &genRecordHandler.UpdateResponse{Entity: request.Entity}, nil
}

func (mrh *MongoRecordHandler) ValidateDeleteRequest(request *genRecordHandler.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !company.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for company", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (mrh *MongoRecordHandler) GDelete(request *genRecordHandler.DeleteRequest) (*genRecordHandler.DeleteResponse, error) {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	filter := request.Identifier.ToFilter()
	filter = company.ContextualiseFilter(filter, request.Claims)
	if err := companyCollection.Remove(filter); err != nil {
		return nil, err
	}

	return &genRecordHandler.DeleteResponse{}, nil
}

func (mrh *MongoRecordHandler) ValidateCollectRequest(request *genRecordHandler.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (mrh *MongoRecordHandler) GCollect(request *genRecordHandler.CollectRequest) (*genRecordHandler.CollectResponse, error) {
	if err := mrh.ValidateCollectRequest(request); err != nil {
		return nil, err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = company.ContextualiseFilter(filter, request.Claims)

	response := genRecordHandler.CollectResponse{}

	// Get Company Collection
	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()
	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Perform Query
	query := companyCollection.Find(filter)

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
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return nil, err
	}

	return &response, nil
}
