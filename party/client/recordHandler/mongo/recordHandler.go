package mongo

import (
	"fmt"
	"github.com/satori/go.uuid"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/party/client"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	clientRecordHandlerException "gitlab.com/iotTracker/brain/party/client/recordHandler/exception"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	userRecordHandlerException "gitlab.com/iotTracker/brain/party/user/recordHandler/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier/adminEmailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type mongoRecordHandler struct {
	mongoSession         *mgo.Session
	database             string
	collection           string
	createIgnoredReasons reasonInvalid.IgnoredReasonsInvalid
	userRecordHandler    userRecordHandler.RecordHandler
}

func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
	userRecordHandler userRecordHandler.RecordHandler,
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
		userRecordHandler:    userRecordHandler,
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

	// Ensure admin email uniqueness
	adminEmailUnique := mgo.Index{
		Key:    []string{"adminEmailAddress"},
		Unique: true,
	}
	if err := clientCollection.EnsureIndex(adminEmailUnique); err != nil {
		log.Fatal("Could not ensure admin email uniqueness: ", err)
	}

}

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *clientRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	// Validate the new client
	clientValidateResponse := clientRecordHandler.ValidateResponse{}

	if err := mrh.Validate(&clientRecordHandler.ValidateRequest{
		Claims: request.Claims,
		Client: request.Client,
		Method: clientRecordHandler.Create},
		&clientValidateResponse); err != nil {
		reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("unable to validate newClient: %s", err.Error()))
	} else {
		for _, reason := range clientValidateResponse.ReasonsInvalid {
			if !mrh.createIgnoredReasons.CanIgnore(reason) {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
			}
		}
	}

	// except for root, who is free to assign parents to other parties, the new clients parent party type and
	// parent id must match claims
	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		switch request.Claims.PartyDetails().PartyType {
		case party.System:
			// do nothing, we expect system to not make a mistake
		default:
			if request.Client.ParentPartyType != request.Claims.PartyDetails().PartyType {
				reasonsInvalid = append(reasonsInvalid, "client ParentPartyType must be the type of the party doing creation")
			}
			if request.Client.ParentId != request.Claims.PartyDetails().PartyId {
				reasonsInvalid = append(reasonsInvalid, "client ParentId must be the id of the party doing creation")
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

	newId, err := uuid.NewV4()
	if err != nil {
		return brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.Client.Id = newId.String()

	if err := clientCollection.Insert(request.Client); err != nil {
		return clientRecordHandlerException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	// create minimal admin user for the client
	if err := mrh.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		Claims: request.Claims,
		User: user.User{
			EmailAddress:    request.Client.AdminEmailAddress,
			ParentPartyType: request.Client.ParentPartyType,
			ParentId:        request.Client.ParentId,
			PartyType:       party.Client,
			PartyId:         id.Identifier{Id: request.Client.Id},
		},
	}, &userRecordHandler.CreateResponse{}); err != nil {
		return clientRecordHandlerException.Create{Reasons: []string{"creating admin user", err.Error()}}
	}

	response.Client = request.Client
	return nil
}

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *clientRecordHandler.RetrieveRequest) error {
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

func (mrh *mongoRecordHandler) Retrieve(request *clientRecordHandler.RetrieveRequest, response *clientRecordHandler.RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var clientRecord client.Client

	filter := client.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)

	if err := clientCollection.Find(filter).One(&clientRecord); err != nil {
		if err == mgo.ErrNotFound {
			return clientRecordHandlerException.NotFound{}
		} else {
			return brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.Client = clientRecord
	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *clientRecordHandler.UpdateRequest) error {
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

func (mrh *mongoRecordHandler) Update(request *clientRecordHandler.UpdateRequest, response *clientRecordHandler.UpdateResponse) error {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve Client
	retrieveClientResponse := clientRecordHandler.RetrieveResponse{}
	if err := mrh.Retrieve(&clientRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveClientResponse); err != nil {
		return clientRecordHandlerException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveClientResponse.Client.Id = request.Client.Id // cannot update ever
	retrieveClientResponse.Client.Name = request.Client.Name

	filter := client.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)

	if err := clientCollection.Update(filter, retrieveClientResponse.Client); err != nil {
		return clientRecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
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

	filter := client.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)
	if err := clientCollection.Remove(filter); err != nil {
		return err
	}

	return nil
}

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *clientRecordHandler.ValidateRequest) error {
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

func (mrh *mongoRecordHandler) Validate(request *clientRecordHandler.ValidateRequest, response *clientRecordHandler.ValidateResponse) error {
	if err := mrh.ValidateValidateRequest(request); err != nil {
		return err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	clientToValidate := &request.Client

	if (*clientToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*clientToValidate).Id,
		})
	}

	if (*clientToValidate).Name == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).Name,
		})
	}

	if (*clientToValidate).ParentPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).ParentPartyType,
		})
	}

	if (*clientToValidate).ParentId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).ParentId,
		})
	}

	if (*clientToValidate).AdminEmailAddress == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "adminEmailAddress",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).AdminEmailAddress,
		})
	}

	returnedReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)

	// Perform additional checks/ignores considering method field
	switch request.Method {
	case clientRecordHandler.Create:

		if (*clientToValidate).AdminEmailAddress != "" {

			// Check if there is another client that is already using the same admin email address
			if err := mrh.Retrieve(&clientRecordHandler.RetrieveRequest{
				Claims: request.Claims,
				Identifier: adminEmailAddress.Identifier{
					AdminEmailAddress: (*clientToValidate).AdminEmailAddress,
				},
			},
				&clientRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case clientRecordHandlerException.NotFound:
					// this is what we want, do nothing
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "adminEmailAddress",
						Type:  reasonInvalid.Unknown,
						Help:  "unknown error",
						Data:  (*clientToValidate).AdminEmailAddress,
					})
				}
			} else {
				// there was no error, this email is already in database
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "adminEmailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*clientToValidate).AdminEmailAddress,
				})
			}

			// Check if there is another user that is already using the same admin email address
			if err := mrh.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
				Claims: request.Claims,
				Identifier: emailAddress.Identifier{
					EmailAddress: (*clientToValidate).AdminEmailAddress,
				},
			},
				&userRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case userRecordHandlerException.NotFound:
					// this is what we want, do nothing
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "adminEmailAddress",
						Type:  reasonInvalid.Unknown,
						Help:  "unknown error",
						Data:  (*clientToValidate).AdminEmailAddress,
					})
				}
			} else {
				// there was no error, this email is already in database
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "adminEmailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*clientToValidate).AdminEmailAddress,
				})
			}
		}

		// Check if there is another user that is already using the same admin email address

		// Ignore reasons not applicable for this method
		for _, reason := range allReasonsInvalid {
			if !mrh.createIgnoredReasons.CanIgnore(reason) {
				returnedReasonsInvalid = append(returnedReasonsInvalid, reason)
			}
		}
	default:
		returnedReasonsInvalid = allReasonsInvalid
	}

	response.ReasonsInvalid = returnedReasonsInvalid
	return nil
}

func (mrh *mongoRecordHandler) ValidateCollectRequest(request *clientRecordHandler.CollectRequest) error {
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

func (mrh *mongoRecordHandler) Collect(request *clientRecordHandler.CollectRequest, response *clientRecordHandler.CollectResponse) error {
	if err := mrh.ValidateCollectRequest(request); err != nil {
		return err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = client.ContextualiseFilter(filter, request.Claims)

	// Get Client Collection
	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()
	clientCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Perform Query
	query := clientCollection.Find(filter)

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
	response.Records = make([]client.Client, 0)
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return err
	}

	return nil
}
