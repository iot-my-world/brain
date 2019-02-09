package mongo

import (
	"fmt"
	globalException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/client"
	clientException "gitlab.com/iotTracker/brain/party/client/exception"
	"gitlab.com/iotTracker/brain/validate"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoRecordHandler struct {
	mongoSession         *mgo.Session
	database             string
	collection           string
	createIgnoredReasons validate.IgnoredReasonsInvalid
}

func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
) *mongoRecordHandler {

	setupIndices(mongoSession, database, collection)

	createIgnoredReasons := validate.IgnoredReasonsInvalid{
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

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *client.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	// Validate the new client
	clientValidateResponse := client.ValidateResponse{}

	err := mrh.Validate(&client.ValidateRequest{Client: request.Client}, &clientValidateResponse)
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
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Create(request *client.CreateRequest, response *client.CreateResponse) error {
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

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *client.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !client.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for client", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Retrieve(request *client.RetrieveRequest, response *client.RetrieveResponse) error {
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
			return globalException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.Client = clientRecord
	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *client.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Update(request *client.UpdateRequest, response *client.UpdateResponse) error {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve Client
	retrieveClientResponse := client.RetrieveResponse{}
	if err := mrh.Retrieve(&client.RetrieveRequest{Identifier: request.Identifier}, &retrieveClientResponse); err != nil {
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

func (mrh *mongoRecordHandler) ValidateDeleteRequest(request *client.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !client.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for client", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Delete(request *client.DeleteRequest, response *client.DeleteResponse) error {
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

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *client.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Validate(request *client.ValidateRequest, response *client.ValidateResponse) error {
	if err := mrh.ValidateValidateRequest(request); err != nil {
		return err
	}

	reasonsInvalid := make([]validate.ReasonInvalid, 0)
	clientToValidate := &request.Client

	if (*clientToValidate).Id == "" {
		reasonsInvalid = append(reasonsInvalid, validate.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*clientToValidate).Id,
		})
	}

	response.ReasonsInvalid = reasonsInvalid
	return nil
}
