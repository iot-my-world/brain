package mongo

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/tracker/device"
	deviceException "gitlab.com/iotTracker/brain/tracker/device/exception"
	deviceRecordHandler "gitlab.com/iotTracker/brain/tracker/device/recordHandler"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/satori/go.uuid"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	companyException "gitlab.com/iotTracker/brain/party/company/exception"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	clientException "gitlab.com/iotTracker/brain/party/client/exception"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type mongoRecordHandler struct {
	mongoSession         *mgo.Session
	database             string
	collection           string
	createIgnoredReasons reasonInvalid.IgnoredReasonsInvalid
	companyRecordHandler companyRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
}

func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
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

	newDeviceMongoRecordHandler := mongoRecordHandler{
		mongoSession:         mongoSession,
		database:             database,
		collection:           collection,
		createIgnoredReasons: createIgnoredReasons,
		companyRecordHandler: companyRecordHandler,
		clientRecordHandler:  clientRecordHandler,
	}

	return &newDeviceMongoRecordHandler
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise Device collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	deviceCollection := mgoSesh.DB(database).C(collection)

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := deviceCollection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}

	// Ensure admin imei uniqueness
	imeiUnique := mgo.Index{
		Key:    []string{"imei"},
		Unique: true,
	}
	if err := deviceCollection.EnsureIndex(imeiUnique); err != nil {
		log.Fatal("Could not ensure imei: ", err)
	}

	// Ensure country code + number uniqueness
	countryCodeNumberUnique := mgo.Index{
		Key:    []string{"simCountryCode", "simNumber"},
		Unique: true,
	}
	if err := deviceCollection.EnsureIndex(countryCodeNumberUnique); err != nil {
		log.Fatal("Could not ensure sim country code and number combination unique: ", err)
	}
}

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *deviceRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "nil claims")
	}

	// Validate the new device
	deviceValidateResponse := deviceRecordHandler.ValidateResponse{}

	if err := mrh.Validate(&deviceRecordHandler.ValidateRequest{
		Device: request.Device,
		Method: deviceRecordHandler.Create},
		&deviceValidateResponse); err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate newDevice")
	} else {
		for _, reason := range deviceValidateResponse.ReasonsInvalid {
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

func (mrh *mongoRecordHandler) Create(request *deviceRecordHandler.CreateRequest, response *deviceRecordHandler.CreateResponse) error {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	deviceCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.Device.Id = newId.String()

	if err := deviceCollection.Insert(request.Device); err != nil {
		return deviceException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	response.Device = request.Device
	return nil
}

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *deviceRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !device.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for device", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Retrieve(request *deviceRecordHandler.RetrieveRequest, response *deviceRecordHandler.RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	deviceCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var deviceRecord device.Device

	filter := request.Identifier.ToFilter()
	if err := deviceCollection.Find(filter).One(&deviceRecord); err != nil {
		if err == mgo.ErrNotFound {
			return deviceException.NotFound{}
		} else {
			return brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.Device = deviceRecord
	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *deviceRecordHandler.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Update(request *deviceRecordHandler.UpdateRequest, response *deviceRecordHandler.UpdateResponse) error {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	deviceCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve Device
	retrieveDeviceResponse := deviceRecordHandler.RetrieveResponse{}
	if err := mrh.Retrieve(&deviceRecordHandler.RetrieveRequest{Identifier: request.Identifier}, &retrieveDeviceResponse); err != nil {
		return deviceException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveDeviceResponse.Device.Id = request.Device.Id // cannot update ever

	if err := deviceCollection.Update(request.Identifier.ToFilter(), retrieveDeviceResponse.Device); err != nil {
		return deviceException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	response.Device = retrieveDeviceResponse.Device

	return nil
}

func (mrh *mongoRecordHandler) ValidateDeleteRequest(request *deviceRecordHandler.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !device.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for device", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Delete(request *deviceRecordHandler.DeleteRequest, response *deviceRecordHandler.DeleteResponse) error {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	deviceCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	if err := deviceCollection.Remove(request.Identifier.ToFilter()); err != nil {
		return err
	}

	return nil
}

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *deviceRecordHandler.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Validate(request *deviceRecordHandler.ValidateRequest, response *deviceRecordHandler.ValidateResponse) error {
	if err := mrh.ValidateValidateRequest(request); err != nil {
		return err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	deviceToValidate := &request.Device

	if (*deviceToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*deviceToValidate).Id,
		})
	}

	if (*deviceToValidate).IMEI == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "imei",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*deviceToValidate).IMEI,
		})
	}

	if (*deviceToValidate).SimCountryCode == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "simCountryCode",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*deviceToValidate).SimCountryCode,
		})
	}

	if (*deviceToValidate).SimNumber == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "simNumber",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*deviceToValidate).SimNumber,
		})
	}

	// owner party type must be set, cannot be blank
	if (*deviceToValidate).OwnerPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "ownerPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*deviceToValidate).OwnerPartyType,
		})
	} else {
		// if it is not blank
		// owner party type must be valid. i.e. must be of a valid type and the party must exist
		switch (*deviceToValidate).OwnerPartyType {
		case party.System:
			// system owner party type means ownerId must be the system id
			rootPartyID := id.Identifier{Id: "root"}
			if (*deviceToValidate).OwnerId != rootPartyID {
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "ownerId",
					Type:  reasonInvalid.MustExist,
					Help:  "owner party must exist",
					Data:  (*deviceToValidate).OwnerId,
				})
			}

		case party.Company:
			// company owner must exist, try and retrieve to confirm
			if err := mrh.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
				Identifier: (*deviceToValidate).OwnerId,
			},
				&companyRecordHandler.RetrieveResponse{});
				err != nil {
				switch err.(type) {
				case companyException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "ownerId",
						Type:  reasonInvalid.MustExist,
						Help:  "owner party must exist",
						Data:  (*deviceToValidate).OwnerId,
					})
				default:
					return brainException.Unexpected{Reasons: []string{"error retrieving company", err.Error()}}
				}
			}

		case party.Client:
			// client owner must exist, try and retrieve to confirm
			if err := mrh.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
				Identifier: (*deviceToValidate).OwnerId,
			},
				&clientRecordHandler.RetrieveResponse{});
				err != nil {
				switch err.(type) {
				case clientException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "ownerId",
						Type:  reasonInvalid.MustExist,
						Help:  "owner party must exist",
						Data:  (*deviceToValidate).OwnerId,
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
				Data:  (*deviceToValidate).OwnerPartyType,
			})
		}
	}

	blankId := id.Identifier{}
	// although assigned party type can be blank, if it is then the assigned id must also be blank
	if ((*deviceToValidate).AssignedPartyType == "" && (*deviceToValidate).AssignedId != blankId) ||
		((*deviceToValidate).AssignedId == blankId && (*deviceToValidate).AssignedPartyType != "") {
		if (*deviceToValidate).AssignedId != blankId {
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "assignedPartyType",
				Type:  reasonInvalid.Invalid,
				Help:  "assigned must be blank if assignedPartyType is",
				Data:  (*deviceToValidate).AssignedPartyType,
			})
			allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
				Field: "assignedId",
				Type:  reasonInvalid.Invalid,
				Help:  "assigned must be blank if assignedPartyType is",
				Data:  (*deviceToValidate).AssignedId,
			})
		}
	} else {
		// neither are blank
		switch (*deviceToValidate).AssignedPartyType {
		case party.System:
			// system owner party type means assignedId must be the system id
			rootPartyID := id.Identifier{Id: "root"}
			if (*deviceToValidate).AssignedId != rootPartyID {
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "ownerId",
					Type:  reasonInvalid.MustExist,
					Help:  "owner party must exist",
					Data:  (*deviceToValidate).AssignedId,
				})
			}

		case party.Company:
			// company assigned must exist, try and retrieve to confirm
			if err := mrh.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
				Identifier: (*deviceToValidate).AssignedId,
			},
				&companyRecordHandler.RetrieveResponse{});
				err != nil {
				switch err.(type) {
				case companyException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.MustExist,
						Help:  "assigned party must exist",
						Data:  (*deviceToValidate).AssignedId,
					})
				default:
					return brainException.Unexpected{Reasons: []string{"error retrieving company", err.Error()}}
				}
			}

		case party.Client:
			// client assigned must exist, try and retrieve to confirm
			if err := mrh.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
				Identifier: (*deviceToValidate).AssignedId,
			},
				&clientRecordHandler.RetrieveResponse{});
				err != nil {
				switch err.(type) {
				case clientException.NotFound:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "assignedId",
						Type:  reasonInvalid.MustExist,
						Help:  "assigned party must exist",
						Data:  (*deviceToValidate).AssignedId,
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
				Data:  (*deviceToValidate).OwnerPartyType,
			})
		}
	}

	returnedReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)

	// Perform additional checks/ignores considering method field
	switch request.Method {
	case deviceRecordHandler.Create:

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

func (mrh *mongoRecordHandler) ValidateCollectRequest(request *deviceRecordHandler.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Collect(request *deviceRecordHandler.CollectRequest, response *deviceRecordHandler.CollectResponse) error {
	if err := mrh.ValidateCollectRequest(request); err != nil {
		return err
	}

	// Build filters from criteria
	filter := bson.M{}
	criteriaFilters := make([]bson.M, 0)
	for criterionIdx := range request.Criteria {
		criteriaFilters = append(criteriaFilters, request.Criteria[criterionIdx].ToFilter())
	}
	if len(criteriaFilters) > 0 {
		filter["$and"] = criteriaFilters
	}

	// Get Device Collection
	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()
	deviceCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Perform Query
	query := deviceCollection.Find(filter)

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
	response.Records = make([]device.Device, 0)
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return err
	}

	return nil
}
