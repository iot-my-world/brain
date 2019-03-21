package mongo

import (
	"fmt"
	"github.com/satori/go.uuid"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/client"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	clientRecordHandlerException "gitlab.com/iotTracker/brain/party/client/recordHandler/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
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
) clientRecordHandler.RecordHandler {

	setupIndices(mongoSession, database, collection)

	return &recordHandler{
		mongoSession: mongoSession,
		database:     database,
		collection:   collection,
	}
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise Client collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	clientCollection := mgoSesh.DB(database).C(collection)

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := clientCollection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}

	// Ensure admin email uniqueness
	adminEmailUnique := mgo.Index{
		Key:    []string{"adminEmailAddress"},
		Unique: true,
	}
	if err := clientCollection.EnsureIndex(adminEmailUnique); err != nil {
		log.Fatal("Could not ensure admin email uniqueness: ", err)
	}

}

func (mrh *recordHandler) ValidateCreateRequest(request *clientRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *recordHandler) Create(request *clientRecordHandler.CreateRequest) (*clientRecordHandler.CreateResponse, error) {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return nil, brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.Client.Id = newId.String()

	if err := clientCollection.Insert(request.Client); err != nil {
		return nil, clientRecordHandlerException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	return &clientRecordHandler.CreateResponse{Client: request.Client}, nil
}

func (mrh *recordHandler) ValidateRetrieveRequest(request *clientRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !client.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for client", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *recordHandler) Retrieve(request *clientRecordHandler.RetrieveRequest) (*clientRecordHandler.RetrieveResponse, error) {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var clientRecord client.Client

	filter := client.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)

	if err := clientCollection.Find(filter).One(&clientRecord); err != nil {
		if err == mgo.ErrNotFound {
			return nil, clientRecordHandlerException.NotFound{}
		} else {
			return nil, brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	return &clientRecordHandler.RetrieveResponse{
		Client: clientRecord,
	}, nil
}

func (mrh *recordHandler) ValidateUpdateRequest(request *clientRecordHandler.UpdateRequest) error {
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

func (mrh *recordHandler) Update(request *clientRecordHandler.UpdateRequest) (*clientRecordHandler.UpdateResponse, error) {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve Client
	retrieveClientResponse, err := mrh.Retrieve(&clientRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		return nil, clientRecordHandlerException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveClientResponse.Client.Id = request.Client.Id // cannot update ever
	retrieveClientResponse.Client.Name = request.Client.Name

	filter := client.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)

	if err := clientCollection.Update(filter, retrieveClientResponse.Client); err != nil {
		return nil, clientRecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return &clientRecordHandler.UpdateResponse{Client: retrieveClientResponse.Client}, nil
}

func (mrh *recordHandler) ValidateDeleteRequest(request *clientRecordHandler.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !client.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for client", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *recordHandler) Delete(request *clientRecordHandler.DeleteRequest) (*clientRecordHandler.DeleteResponse, error) {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return nil, err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	filter := client.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)
	if err := clientCollection.Remove(filter); err != nil {
		return nil, err
	}

	return &clientRecordHandler.DeleteResponse{}, nil
}

func (mrh *recordHandler) ValidateCollectRequest(request *clientRecordHandler.CollectRequest) error {
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

func (mrh *recordHandler) Collect(request *clientRecordHandler.CollectRequest) (*clientRecordHandler.CollectResponse, error) {
	if err := mrh.ValidateCollectRequest(request); err != nil {
		return nil, err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = client.ContextualiseFilter(filter, request.Claims)

	// Get Client Collection
	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()
	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Perform Query
	query := clientCollection.Find(filter)

	response := clientRecordHandler.CollectResponse{}

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
	response.Records = make([]client.Client, 0)
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return nil, err
	}

	return &response, nil
}
