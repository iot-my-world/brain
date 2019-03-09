package mongo

import (
	"fmt"
	"github.com/satori/go.uuid"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/party/company"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	companyRecordHandlerException "gitlab.com/iotTracker/brain/party/company/recordHandler/exception"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	userRecordHandlerException "gitlab.com/iotTracker/brain/party/user/recordHandler/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier/adminEmailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/api"
)

type mongoRecordHandler struct {
	mongoSession      *mgo.Session
	database          string
	collection        string
	userRecordHandler userRecordHandler.RecordHandler
	ignoredReasons    map[api.Method]reasonInvalid.IgnoredReasonsInvalid
}

func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
	userRecordHandler userRecordHandler.RecordHandler,
) *mongoRecordHandler {

	setupIndices(mongoSession, database, collection)

	ignoredReasons := map[api.Method]reasonInvalid.IgnoredReasonsInvalid{
		companyRecordHandler.Create: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
			},
		},
	}

	newCompanyMongoRecordHandler := mongoRecordHandler{
		mongoSession:      mongoSession,
		database:          database,
		collection:        collection,
		ignoredReasons:    ignoredReasons,
		userRecordHandler: userRecordHandler,
	}

	return &newCompanyMongoRecordHandler
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise Company collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	companyCollection := mgoSesh.DB(database).C(collection)

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := companyCollection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}

	// Ensure admin email uniqueness
	adminEmailUnique := mgo.Index{
		Key:    []string{"adminEmailAddress"},
		Unique: true,
	}
	if err := companyCollection.EnsureIndex(adminEmailUnique); err != nil {
		log.Fatal("Could not ensure admin email uniqueness: ", err)
	}

}

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *companyRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	// A new company can only be made by root
	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "nil claims")
	} else {
		if request.Claims.PartyDetails().PartyType != party.System {
			reasonsInvalid = append(reasonsInvalid, "only system party can make a new company")
		}
	}

	// Validate the new company
	companyValidateResponse := companyRecordHandler.ValidateResponse{}

	if err := mrh.Validate(&companyRecordHandler.ValidateRequest{
		Claims:  request.Claims,
		Company: request.Company,
		Method:  companyRecordHandler.Create},
		&companyValidateResponse); err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate newCompany")
	} else {
		for _, reason := range companyValidateResponse.ReasonsInvalid {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Create(request *companyRecordHandler.CreateRequest, response *companyRecordHandler.CreateResponse) error {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.Company.Id = newId.String()

	if err := companyCollection.Insert(request.Company); err != nil {
		return companyRecordHandlerException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	response.Company = request.Company
	return nil
}

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *companyRecordHandler.RetrieveRequest) error {
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
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Retrieve(request *companyRecordHandler.RetrieveRequest, response *companyRecordHandler.RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var companyRecord company.Company

	filter := request.Identifier.ToFilter()
	filter = company.ContextualiseFilter(filter, request.Claims)

	if err := companyCollection.Find(filter).One(&companyRecord); err != nil {
		if err == mgo.ErrNotFound {
			return companyRecordHandlerException.NotFound{}
		} else {
			return brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.Company = companyRecord
	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *companyRecordHandler.UpdateRequest) error {
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
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Update(request *companyRecordHandler.UpdateRequest, response *companyRecordHandler.UpdateResponse) error {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve Company
	retrieveCompanyResponse := companyRecordHandler.RetrieveResponse{}
	if err := mrh.Retrieve(&companyRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveCompanyResponse); err != nil {
		return companyRecordHandlerException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveCompanyResponse.Company.Id = request.Company.Id // cannot update ever
	retrieveCompanyResponse.Company.Name = request.Company.Name

	filter := request.Identifier.ToFilter()
	filter = company.ContextualiseFilter(filter, request.Claims)
	if err := companyCollection.Update(filter, retrieveCompanyResponse.Company); err != nil {
		return companyRecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	response.Company = retrieveCompanyResponse.Company

	return nil
}

func (mrh *mongoRecordHandler) ValidateDeleteRequest(request *companyRecordHandler.DeleteRequest) error {
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
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Delete(request *companyRecordHandler.DeleteRequest, response *companyRecordHandler.DeleteResponse) error {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	filter := request.Identifier.ToFilter()
	filter = company.ContextualiseFilter(filter, request.Claims)
	if err := companyCollection.Remove(filter); err != nil {
		return err
	}

	return nil
}

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *companyRecordHandler.ValidateRequest) error {
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

func (mrh *mongoRecordHandler) Validate(request *companyRecordHandler.ValidateRequest, response *companyRecordHandler.ValidateResponse) error {
	if err := mrh.ValidateValidateRequest(request); err != nil {
		return err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	companyToValidate := &request.Company

	if (*companyToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*companyToValidate).Id,
		})
	}

	if (*companyToValidate).Name == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*companyToValidate).Name,
		})
	}

	if (*companyToValidate).AdminEmailAddress == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "adminEmailAddress",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*companyToValidate).AdminEmailAddress,
		})
	}

	if (*companyToValidate).ParentPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*companyToValidate).ParentPartyType,
		})
	}

	blankIdIdentifier := id.Identifier{}
	if (*companyToValidate).ParentId == blankIdIdentifier {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*companyToValidate).ParentId,
		})
	}

	returnedReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)

	// Perform additional checks/ignores considering method field
	switch request.Method {
	case companyRecordHandler.Create:

		// Check if there is another client that is already using the same admin email address
		if (*companyToValidate).AdminEmailAddress != "" {
			if err := mrh.Retrieve(&companyRecordHandler.RetrieveRequest{
				Claims: request.Claims,
				Identifier: adminEmailAddress.Identifier{
					AdminEmailAddress: (*companyToValidate).AdminEmailAddress,
				},
			},
				&companyRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case companyRecordHandlerException.NotFound:
					// this is what we want, do nothing
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "adminEmailAddress",
						Type:  reasonInvalid.Unknown,
						Help:  "unknown error",
						Data:  (*companyToValidate).AdminEmailAddress,
					})
				}
			} else {
				// there was no error, this email is already in database
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "adminEmailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*companyToValidate).AdminEmailAddress,
				})
			}

			// check if there any users with this email address
			if err := mrh.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
				Claims: request.Claims,
				Identifier: emailAddress.Identifier{
					EmailAddress: (*companyToValidate).AdminEmailAddress,
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
						Data:  (*companyToValidate).AdminEmailAddress,
					})
				}
			} else {
				// there was no error, this email is already in database
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "adminEmailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*companyToValidate).AdminEmailAddress,
				})
			}
		}
	}

	// Ignore reasons applicable to method if relevant
	if mrh.ignoredReasons[request.Method].ReasonsInvalid != nil {
		for _, reason := range allReasonsInvalid {
			if !mrh.ignoredReasons[request.Method].CanIgnore(reason) {
				returnedReasonsInvalid = append(returnedReasonsInvalid, reason)
			}
		}
	}

	response.ReasonsInvalid = returnedReasonsInvalid
	return nil
}

func (mrh *mongoRecordHandler) ValidateCollectRequest(request *companyRecordHandler.CollectRequest) error {
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

func (mrh *mongoRecordHandler) Collect(request *companyRecordHandler.CollectRequest, response *companyRecordHandler.CollectResponse) error {
	if err := mrh.ValidateCollectRequest(request); err != nil {
		return err
	}

	filter := criterion.CriteriaToFilter(request.Criteria)
	filter = company.ContextualiseFilter(filter, request.Claims)

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
		return err
	}

	// Apply limit if applicable
	if request.Query.Limit > 0 {
		query.Limit(request.Query.Limit)
	}

	// Determine the Sort Order
	mongoSortOrder := request.Query.ToMongoSortFormat()

	// Populate records
	response.Records = make([]company.Company, 0)
	if err := query.
		Skip(request.Query.Offset).
		Sort(mongoSortOrder...).
		All(&response.Records); err != nil {
		return err
	}

	return nil
}
