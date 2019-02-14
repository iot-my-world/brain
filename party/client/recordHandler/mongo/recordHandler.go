package mongo

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/client"
	clientException "gitlab.com/iotTracker/brain/party/client/exception"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoRecordHandler struct {
	mongoSession         *mgo.Session
	database             string
	collection           string
	createIgnoredReasons reasonInvalid.IgnoredReasonsInvalid
}

func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
) *mongoRecordHandler {

	setupIndices(mongoSession, database, collection)

	createIgnoredReasons := reasonInvalid.IgnoredReasonsInvalid{
		ReasonsInvalid: map[string][]reasonInvalid.Type{
			"id": {
				reasonInvalid.Blank,
			},
		},
	}

	newClientMongoRecordHandler := mongoRecordHandler{
		mongoSession:         mongoSession,
		database:             database,
		collection:           collection,
		createIgnoredReasons: createIgnoredReasons,
	}

	return &newClientMongoRecordHandler
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
}

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *clientRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	// Validate the new client
	clientValidateResponse := clientRecordHandler.ValidateResponse{}

	err := mrh.Validate(&clientRecordHandler.ValidateRequest{Client: request.Client}, &clientValidateResponse)
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate newClient")
	} else {
		for _, reason := range clientValidateResponse.ReasonsInvalid {
			if !mrh.createIgnoredReasons.CanIgnore(reason) {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Create(request *clientRecordHandler.CreateRequest, response *clientRecordHandler.CreateResponse) error {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	request.Client.Id = bson.NewObjectId().Hex()

	if err := clientCollection.Insert(request.Client); err != nil {
		return clientException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	response.Client = request.Client
	return nil
}

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *clientRecordHandler.RetrieveRequest) error {
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

func (mrh *mongoRecordHandler) Retrieve(request *clientRecordHandler.RetrieveRequest, response *clientRecordHandler.RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var clientRecord client.Client

	if err := clientCollection.Find(request.Identifier.ToFilter()).One(&clientRecord); err != nil {
		if err == mgo.ErrNotFound {
			return clientException.NotFound{}
		} else {
			return brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.Client = clientRecord
	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *clientRecordHandler.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Update(request *clientRecordHandler.UpdateRequest, response *clientRecordHandler.UpdateResponse) error {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve Client
	retrieveClientResponse := clientRecordHandler.RetrieveResponse{}
	if err := mrh.Retrieve(&clientRecordHandler.RetrieveRequest{Identifier: request.Identifier}, &retrieveClientResponse); err != nil {
		return clientException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveClientResponse.Client.Id = request.Client.Id // cannot update ever
	retrieveClientResponse.Client.ParentId = request.Client.ParentId

	if err := clientCollection.Update(request.Identifier.ToFilter(), retrieveClientResponse.Client); err != nil {
		return clientException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	response.Client = retrieveClientResponse.Client

	return nil
}

func (mrh *mongoRecordHandler) ValidateDeleteRequest(request *clientRecordHandler.DeleteRequest) error {
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

func (mrh *mongoRecordHandler) Delete(request *clientRecordHandler.DeleteRequest, response *clientRecordHandler.DeleteResponse) error {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	if err := clientCollection.Remove(request.Identifier.ToFilter()); err != nil {
		return err
	}

	return nil
}

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *clientRecordHandler.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Validate(request *clientRecordHandler.ValidateRequest, response *clientRecordHandler.ValidateResponse) error {
	if err := mrh.ValidateValidateRequest(request); err != nil {
		return err
	}

	reasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	clientToValidate := &request.Client

	if (*clientToValidate).Id == "" {
		reasonsInvalid = append(reasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*clientToValidate).Id,
		})
	}

	response.ReasonsInvalid = reasonsInvalid
	return nil
}
