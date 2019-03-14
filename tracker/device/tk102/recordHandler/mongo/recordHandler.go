package mongo

import (
	"fmt"
	"github.com/satori/go.uuid"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	clientRecordHandlerException "gitlab.com/iotTracker/brain/party/client/recordHandler/exception"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	companyRecordHandlerException "gitlab.com/iotTracker/brain/party/company/recordHandler/exception"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	systemRecordHandlerException "gitlab.com/iotTracker/brain/party/system/recordHandler/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
	tk102RecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler"
	tk102ExceptionRecordHandlerException "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler/exception"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gopkg.in/mgo.v2"
)

type mongoRecordHandler struct {
	mongoSession         *mgo.Session
	database             string
	collection           string
	createIgnoredReasons reasonInvalid.IgnoredReasonsInvalid
	updateIgnoredReasons reasonInvalid.IgnoredReasonsInvalid
	systemRecordHandler  systemRecordHandler.RecordHandler
	companyRecordHandler companyRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
}

func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
	systemRecordHandler systemRecordHandler.RecordHandler,
	companyRecordHandler companyRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
) *mongoRecordHandler {

	setupIndices(mongoSession, database, collection)

	createIgnoredReasons := reasonInvalid.IgnoredReasonsInvalid{
		ReasonsInvalid: map[string][]reasonInvalid.Type{
			"id": {
				reasonInvalid.Blank,
			},
		},
	}

	updateIgnoredReasons := reasonInvalid.IgnoredReasonsInvalid{
		ReasonsInvalid: map[string][]reasonInvalid.Type{},
	}

	newTK102MongoRecordHandler := mongoRecordHandler{
		mongoSession:         mongoSession,
		database:             database,
		collection:           collection,
		createIgnoredReasons: createIgnoredReasons,
		updateIgnoredReasons: updateIgnoredReasons,
		systemRecordHandler:  systemRecordHandler,
		companyRecordHandler: companyRecordHandler,
		clientRecordHandler:  clientRecordHandler,
	}

	return &newTK102MongoRecordHandler
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise TK102 collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	tk102Collection := mgoSesh.DB(database).C(collection)

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := tk102Collection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}

	// Ensure admin manufacturerIdUnique uniqueness
	manufacturerIdUnique := mgo.Index{
		Key:    []string{"manufacturerId"},
		Unique: true,
	}
	if err := tk102Collection.EnsureIndex(manufacturerIdUnique); err != nil {
		log.Fatal("Could not ensure manufacturerId uniqueness: ", err)
	}

	// Ensure country code + number uniqueness
	countryCodeNumberUnique := mgo.Index{
		Key:    []string{"simCountryCode", "simNumber"},
		Unique: true,
	}
	if err := tk102Collection.EnsureIndex(countryCodeNumberUnique); err != nil {
		log.Fatal("Could not ensure sim country code and number combination unique: ", err)
	}
}

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *tk102RecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "nil claims")
	}

	// Validate the new tk102
	tk102ValidateResponse := tk102RecordHandler.ValidateResponse{}

	if err := mrh.Validate(&tk102RecordHandler.ValidateRequest{
		Claims: request.Claims,
		TK102:  request.TK102,
		Method: tk102RecordHandler.Create},
		&tk102ValidateResponse); err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate newTK102"+err.Error())
	} else {
		for _, reason := range tk102ValidateResponse.ReasonsInvalid {
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

func (mrh *mongoRecordHandler) Create(request *tk102RecordHandler.CreateRequest, response *tk102RecordHandler.CreateResponse) error {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	tk102Collection := mgoSession.DB(mrh.database).C(mrh.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.TK102.Id = newId.String()

	if err := tk102Collection.Insert(request.TK102); err != nil {
		return tk102ExceptionRecordHandlerException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	response.TK102 = request.TK102
	return nil
}

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *tk102RecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !tk102.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for tk102", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Retrieve(request *tk102RecordHandler.RetrieveRequest, response *tk102RecordHandler.RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	tk102Collection := mgoSession.DB(mrh.database).C(mrh.collection)

	var tk102Record tk102.TK102

	filter := claims.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)
	if err := tk102Collection.Find(filter).One(&tk102Record); err != nil {
		if err == mgo.ErrNotFound {
			return tk102ExceptionRecordHandlerException.NotFound{}
		} else {
			return brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.TK102 = tk102Record
	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *tk102RecordHandler.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// validate the device for update
		validateResponse := tk102RecordHandler.ValidateResponse{}
		if err := mrh.Validate(&tk102RecordHandler.ValidateRequest{
			TK102:  request.TK102,
			Claims: request.Claims,
			Method: tk102RecordHandler.Update,
		}, &validateResponse); err != nil {
			reasonsInvalid = append(reasonsInvalid, "validation error: "+err.Error())
		}
		if len(validateResponse.ReasonsInvalid) > 0 {
			for _, reason := range validateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Update(request *tk102RecordHandler.UpdateRequest, response *tk102RecordHandler.UpdateResponse) error {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	tk102Collection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve TK102
	retrieveTK102Response := tk102RecordHandler.RetrieveResponse{}
	if err := mrh.Retrieve(&tk102RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveTK102Response); err != nil {
		return tk102ExceptionRecordHandlerException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveTK102Response.TK102.Id = request.TK102.Id // cannot update ever
	retrieveTK102Response.TK102.ManufacturerId = request.TK102.ManufacturerId
	retrieveTK102Response.TK102.SimCountryCode = request.TK102.SimCountryCode
	retrieveTK102Response.TK102.SimNumber = request.TK102.SimNumber
	retrieveTK102Response.TK102.OwnerPartyType = request.TK102.OwnerPartyType
	retrieveTK102Response.TK102.OwnerId = request.TK102.OwnerId
	retrieveTK102Response.TK102.AssignedPartyType = request.TK102.AssignedPartyType
	retrieveTK102Response.TK102.AssignedId = request.TK102.AssignedId

	if err := tk102Collection.Update(request.Identifier.ToFilter(), retrieveTK102Response.TK102); err != nil {
		return tk102ExceptionRecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	response.TK102 = retrieveTK102Response.TK102

	return nil
}

func (mrh *mongoRecordHandler) ValidateDeleteRequest(request *tk102RecordHandler.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !tk102.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for tk102", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Delete(request *tk102RecordHandler.DeleteRequest, response *tk102RecordHandler.DeleteResponse) error {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	tk102Collection := mgoSession.DB(mrh.database).C(mrh.collection)

	if err := tk102Collection.Remove(request.Identifier.ToFilter()); err != nil {
		return err
	}

	return nil
}

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *tk102RecordHandler.ValidateRequest) error {
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

func (mrh *mongoRecordHandler) Validate(request *tk102RecordHandler.ValidateRequest, response *tk102RecordHandler.ValidateResponse) error {
	if err := mrh.ValidateValidateRequest(request); err != nil {
		return err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	tk102ToValidate := &request.TK102

	if (*tk102ToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*tk102ToValidate).Id,
		})
	}

	if (*tk102ToValidate).ManufacturerId == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "manufacturerId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*tk102ToValidate).ManufacturerId,
		})
	}

	if (*tk102ToValidate).SimCountryCode == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "simCountryCode",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*tk102ToValidate).SimCountryCode,
		})
	}

	if (*tk102ToValidate).SimNumber == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "simNumber",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*tk102ToValidate).SimNumber,
		})
	}

	// owner party type must be set, cannot be blank
	if (*tk102ToValidate).OwnerPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "ownerPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*tk102ToValidate).OwnerPartyType,
		})
	} else {
		// if it is not blank
		// owner party type must be valid. i.e. must be of a valid type and the party must exist
		switch (*tk102ToValidate).OwnerPartyType {
		case party.System:
			// system owner must exist, try and retrieve to confirm
			if err := mrh.systemRecordHandler.Retrieve(&systemRecordHandler.RetrieveRequest{
				Claims:     request.Claims,
				Identifier: (*tk102ToValidate).OwnerId,
			}, &systemRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case systemRecordHandlerException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "ownerId",
						Type:  reasonInvalid.MustExist,
						Help:  "owner party must exist",
						Data:  (*tk102ToValidate).OwnerId,
					})
				default:
					return brainException.Unexpected{Reasons: []string{"error retrieving system", err.Error()}}
				}
			}

		case party.Company:
			// company owner must exist, try and retrieve to confirm
			if err := mrh.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
				Claims:     request.Claims,
				Identifier: (*tk102ToValidate).OwnerId,
			}, &companyRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case companyRecordHandlerException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "ownerId",
						Type:  reasonInvalid.MustExist,
						Help:  "owner party must exist",
						Data:  (*tk102ToValidate).OwnerId,
					})
				default:
					return brainException.Unexpected{Reasons: []string{"error retrieving company", err.Error()}}
				}
			}

		case party.Client:
			// client owner must exist, try and retrieve to confirm
			if err := mrh.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
				Claims:     request.Claims,
				Identifier: (*tk102ToValidate).OwnerId,
			}, &clientRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case clientRecordHandlerException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "ownerId",
						Type:  reasonInvalid.MustExist,
						Help:  "owner party must exist",
						Data:  (*tk102ToValidate).OwnerId,
					})
				default:
					return brainException.Unexpected{Reasons: []string{"error retrieving client", err.Error()}}
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "ownerPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*tk102ToValidate).OwnerPartyType,
			})
		}
	}

	// although assigned party type can be blank, if it is then the assigned id must also be blank
	if ((*tk102ToValidate).AssignedPartyType == "" && (*tk102ToValidate).AssignedId.Id != "") ||
		((*tk102ToValidate).AssignedId.Id == "" && (*tk102ToValidate).AssignedPartyType != "") {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "assignedPartyType",
			Type:  reasonInvalid.Invalid,
			Help:  "must both be blank or set",
			Data:  (*tk102ToValidate).AssignedPartyType,
		})
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "assignedId",
			Type:  reasonInvalid.Invalid,
			Help:  "must both be blank or set",
			Data:  (*tk102ToValidate).AssignedId,
		})
	} else if (*tk102ToValidate).AssignedPartyType != "" && (*tk102ToValidate).AssignedId.Id != "" {
		// neither are blank
		switch (*tk102ToValidate).AssignedPartyType {
		case party.System:
			// system assigned must exist, try and retrieve to confirm
			if err := mrh.systemRecordHandler.Retrieve(&systemRecordHandler.RetrieveRequest{
				Claims:     request.Claims,
				Identifier: (*tk102ToValidate).AssignedId,
			},
				&systemRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case systemRecordHandlerException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.MustExist,
						Help:  "assigned party must exist",
						Data:  (*tk102ToValidate).AssignedId,
					})
				default:
					return brainException.Unexpected{Reasons: []string{"error retrieving system", err.Error()}}
				}
			}

		case party.Company:
			// company assigned must exist, try and retrieve to confirm
			if err := mrh.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
				Claims:     request.Claims,
				Identifier: (*tk102ToValidate).AssignedId,
			},
				&companyRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case companyRecordHandlerException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.MustExist,
						Help:  "assigned party must exist",
						Data:  (*tk102ToValidate).AssignedId,
					})
				default:
					return brainException.Unexpected{Reasons: []string{"error retrieving company", err.Error()}}
				}
			}

		case party.Client:
			// client assigned must exist, try and retrieve to confirm
			if err := mrh.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
				Claims:     request.Claims,
				Identifier: (*tk102ToValidate).AssignedId,
			},
				&clientRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case clientRecordHandlerException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.MustExist,
						Help:  "assigned party must exist",
						Data:  (*tk102ToValidate).AssignedId,
					})
				default:
					return brainException.Unexpected{Reasons: []string{"error retrieving client", err.Error()}}
				}
			}

		default:
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "ownerPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "must be a valid type",
				Data:  (*tk102ToValidate).OwnerPartyType,
			})
		}
	}

	returnedReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)

	// Perform additional checks/ignores considering method field
	switch request.Method {
	case tk102RecordHandler.Create:

		// Ignore reasons not applicable for this method
		for _, reason := range allReasonsInvalid {
			if !mrh.createIgnoredReasons.CanIgnore(reason) {
				returnedReasonsInvalid = append(returnedReasonsInvalid, reason)
			}
		}

	case tk102RecordHandler.Update:

		// Ignore reasons not applicable for this method
		for _, reason := range allReasonsInvalid {
			if !mrh.updateIgnoredReasons.CanIgnore(reason) {
				returnedReasonsInvalid = append(returnedReasonsInvalid, reason)
			}
		}

	default:
		returnedReasonsInvalid = allReasonsInvalid
	}

	response.ReasonsInvalid = returnedReasonsInvalid
	return nil
}

func (mrh *mongoRecordHandler) ValidateCollectRequest(request *tk102RecordHandler.CollectRequest) error {
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

func (mrh *mongoRecordHandler) Collect(request *tk102RecordHandler.CollectRequest, response *tk102RecordHandler.CollectResponse) error {
	if err := mrh.ValidateCollectRequest(request); err != nil {
		return err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = claims.ContextualiseFilter(filter, request.Claims)

	// Get TK102 Collection
	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()
	tk102Collection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Perform Query
	query := tk102Collection.Find(filter)

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
	response.Records = make([]tk102.TK102, 0)
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return err
	}

	return nil
}
