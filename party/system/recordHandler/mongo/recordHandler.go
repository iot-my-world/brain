package mongo

import (
	"fmt"
	"github.com/satori/go.uuid"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	partyRegistrar "gitlab.com/iotTracker/brain/party/registrar"
	"gitlab.com/iotTracker/brain/party/system"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	systemException "gitlab.com/iotTracker/brain/party/system/recordHandler/exception"
	systemSetup "gitlab.com/iotTracker/brain/party/system/setup"
	"gitlab.com/iotTracker/brain/search/criterion"
	loginClaims "gitlab.com/iotTracker/brain/security/claims/login"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gopkg.in/mgo.v2"
)

type recordHandler struct {
	mongoSession         *mgo.Session
	database             string
	collection           string
	createIgnoredReasons reasonInvalid.IgnoredReasonsInvalid
}

func New(
	mongoSession *mgo.Session,
	database,
	collection,
	rootPasswordFileLocation string,
	registrar partyRegistrar.Registrar,
	systemClaims *loginClaims.Login,
) systemRecordHandler.RecordHandler {

	setupIndices(mongoSession, database, collection)

	createIgnoredReasons := reasonInvalid.IgnoredReasonsInvalid{
		ReasonsInvalid: map[string][]reasonInvalid.Type{
			"id": {
				reasonInvalid.Blank,
			},
		},
	}

	newSystemMongoRecordHandler := recordHandler{
		mongoSession:         mongoSession,
		database:             database,
		collection:           collection,
		createIgnoredReasons: createIgnoredReasons,
	}

	if err := systemSetup.InitialSetup(&newSystemMongoRecordHandler, registrar, rootPasswordFileLocation, systemClaims); err != nil {
		log.Fatal("Unable to complete initial system setup!", err.Error())
	}

	return &newSystemMongoRecordHandler
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise System collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	systemCollection := mgoSesh.DB(database).C(collection)

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := systemCollection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}

	// Ensure name uniqueness
	nameUnique := mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	}
	if err := systemCollection.EnsureIndex(nameUnique); err != nil {
		log.Fatal("Could not ensure name uniqueness: ", err)
	}
}

func (r *recordHandler) ValidateCreateRequest(request *systemRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	// Validate the new system
	systemValidateResponse, err := r.Validate(&systemRecordHandler.ValidateRequest{System: request.System})
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate new system")
	} else {
		for _, reason := range systemValidateResponse.ReasonsInvalid {
			if !r.createIgnoredReasons.CanIgnore(reason) {
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

func (r *recordHandler) Create(request *systemRecordHandler.CreateRequest) (*systemRecordHandler.CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	systemCollection := mgoSession.DB(r.database).C(r.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return nil, brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.System.Id = newId.String()

	if err := systemCollection.Insert(request.System); err != nil {
		return nil, systemException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	return &systemRecordHandler.CreateResponse{System: request.System}, nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *systemRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !system.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for system", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *recordHandler) Retrieve(request *systemRecordHandler.RetrieveRequest) (*systemRecordHandler.RetrieveResponse, error) {
	if err := r.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	systemCollection := mgoSession.DB(r.database).C(r.collection)

	var systemRecord system.System

	filter := system.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)
	if err := systemCollection.Find(filter).One(&systemRecord); err != nil {
		if err == mgo.ErrNotFound {
			return nil, systemException.NotFound{}
		} else {
			return nil, brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	return &systemRecordHandler.RetrieveResponse{System: systemRecord}, nil
}

func (r *recordHandler) ValidateUpdateRequest(request *systemRecordHandler.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *recordHandler) Update(request *systemRecordHandler.UpdateRequest) (*systemRecordHandler.UpdateResponse, error) {
	if err := r.ValidateUpdateRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	systemCollection := mgoSession.DB(r.database).C(r.collection)

	// Retrieve System
	retrieveSystemResponse, err := r.Retrieve(&systemRecordHandler.RetrieveRequest{
		Identifier: request.Identifier,
		Claims:     request.Claims,
	})
	if err != nil {
		return nil, systemException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveSystemResponse.System.Id = request.System.Id // cannot update ever
	retrieveSystemResponse.System.Name = request.System.Name
	retrieveSystemResponse.System.AdminEmailAddress = request.System.AdminEmailAddress
	filter := system.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)
	if err := systemCollection.Update(filter, retrieveSystemResponse.System); err != nil {
		return nil, systemException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	return &systemRecordHandler.UpdateResponse{System: retrieveSystemResponse.System}, nil
}

func (r *recordHandler) ValidateDeleteRequest(request *systemRecordHandler.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !system.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for system", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *recordHandler) Delete(request *systemRecordHandler.DeleteRequest) (*systemRecordHandler.DeleteResponse, error) {
	if err := r.ValidateDeleteRequest(request); err != nil {
		return nil, err
	}

	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()

	systemCollection := mgoSession.DB(r.database).C(r.collection)
	filter := system.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)
	if err := systemCollection.Remove(filter); err != nil {
		return nil, err
	}

	return &systemRecordHandler.DeleteResponse{}, nil
}

func (r *recordHandler) ValidateValidateRequest(request *systemRecordHandler.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *recordHandler) Validate(request *systemRecordHandler.ValidateRequest) (*systemRecordHandler.ValidateResponse, error) {
	if err := r.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	systemToValidate := &request.System

	if (*systemToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*systemToValidate).Id,
		})
	}

	if (*systemToValidate).Name == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*systemToValidate).Name,
		})
	}

	if (*systemToValidate).AdminEmailAddress == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "adminEmailAddress",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*systemToValidate).AdminEmailAddress,
		})
	}

	returnedReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)

	// Perform additional checks/ignores considering method field
	switch request.Method {
	case systemRecordHandler.Create:
		// Ignore reasons not applicable for this method
		for _, reason := range allReasonsInvalid {
			if !r.createIgnoredReasons.CanIgnore(reason) {
				returnedReasonsInvalid = append(returnedReasonsInvalid, reason)
			}
		}

	default:
		returnedReasonsInvalid = allReasonsInvalid
	}

	return &systemRecordHandler.ValidateResponse{ReasonsInvalid: returnedReasonsInvalid}, nil
}

func (r *recordHandler) ValidateCollectRequest(request *systemRecordHandler.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (r *recordHandler) Collect(request *systemRecordHandler.CollectRequest) (*systemRecordHandler.CollectResponse, error) {
	if err := r.ValidateCollectRequest(request); err != nil {
		return nil, err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = system.ContextualiseFilter(filter, request.Claims)

	response := systemRecordHandler.CollectResponse{}

	// Get System Collection
	mgoSession := r.mongoSession.Copy()
	defer mgoSession.Close()
	systemCollection := mgoSession.DB(r.database).C(r.collection)

	// Perform Query
	query := systemCollection.Find(filter)

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
	response.Records = make([]system.System, 0)
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return nil, err
	}

	return &response, nil
}
