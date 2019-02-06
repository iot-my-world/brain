package mongo

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/log"
	"gopkg.in/mgo.v2/bson"
	globalException "gitlab.com/iotTracker/brain/exception"
	companyException "gitlab.com/iotTracker/brain/party/company/exception"
	"fmt"
	"gitlab.com/iotTracker/brain/validate"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gitlab.com/iotTracker/brain/party/company"
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

	newCompanyMongoRecordHandler := mongoRecordHandler{
		mongoSession:         mongoSession,
		database:             database,
		collection:           collection,
		createIgnoredReasons: createIgnoredReasons,
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
}

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *company.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	// Validate the new company
	companyValidateResponse := company.ValidateResponse{}

	err := mrh.Validate(&company.ValidateRequest{Company: request.Company}, &companyValidateResponse)
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate newCompany")
	} else {
		for _, reason := range companyValidateResponse.ReasonsInvalid {
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

func (mrh *mongoRecordHandler) Create(request *company.CreateRequest, response *company.CreateResponse) error {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	request.Company.Id = bson.NewObjectId().Hex()

	if err := companyCollection.Insert(request.Company); err != nil {
		return companyException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	response.Company = request.Company
	return nil
}

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *company.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !company.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for company", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Retrieve(request *company.RetrieveRequest, response *company.RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var companyRecord company.Company

	if err := companyCollection.Find(request.Identifier.ToFilter()).One(&companyRecord); err != nil {
		if err == mgo.ErrNotFound {
			return companyException.NotFound{}
		} else {
			return globalException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.Company = companyRecord
	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *company.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Update(request *company.UpdateRequest, response *company.UpdateResponse) error {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve Company
	retrieveCompanyResponse := company.RetrieveResponse{}
	if err := mrh.Retrieve(&company.RetrieveRequest{Identifier: request.Identifier}, &retrieveCompanyResponse); err != nil {
		return companyException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveCompanyResponse.Company.Id = request.Company.Id // cannot update ever
	retrieveCompanyResponse.Company.Name = request.Company.Name

	if err := companyCollection.Update(request.Identifier.ToFilter(), retrieveCompanyResponse.Company); err != nil {
		return companyException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	response.Company = retrieveCompanyResponse.Company

	return nil
}

func (mrh *mongoRecordHandler) ValidateDeleteRequest(request *company.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !company.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for company", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Delete(request *company.DeleteRequest, response *company.DeleteResponse) error {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	companyCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	if err := companyCollection.Remove(request.Identifier.ToFilter()); err != nil {
		return err
	}

	return nil
}

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *company.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Validate(request *company.ValidateRequest, response *company.ValidateResponse) error {
	if err := mrh.ValidateValidateRequest(request); err != nil {
		return err
	}

	reasonsInvalid := make([]validate.ReasonInvalid, 0)
	companyToValidate := &request.Company

	if (*companyToValidate).Id == "" {
		reasonsInvalid = append(reasonsInvalid, validate.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*companyToValidate).Id,
		})
	}

	if (*companyToValidate).Name == "" {
		reasonsInvalid = append(reasonsInvalid, validate.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*companyToValidate).Name,
		})
	}

	response.ReasonsInvalid = reasonsInvalid
	return nil
}
