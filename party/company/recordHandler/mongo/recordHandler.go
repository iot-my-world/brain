package mongo

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/genRecordHandler"
	imp "gitlab.com/iotTracker/brain/genRecordHandler/mongo"
	"gitlab.com/iotTracker/brain/party/company"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	"gopkg.in/mgo.v2"
)

type mongoRecordHandler struct {
	recordHandler.RecordHandler
}

// New mongo record handler
func New(
	mongoSession *mgo.Session,
	database string,
) companyRecordHandler.RecordHandler {

	//setupIndices(mongoSession, database, collection)

	CompGen := imp.New(
		mongoSession,
		database,
		"company",
		[]mgo.Index{
			{
				Key:    []string{"id"},
				Unique: true,
			},
			{
				Key:    []string{"adminEmailAddress"},
				Unique: true,
			},
		},
	)

	newCompanyMongoRecordHandler := mongoRecordHandler{
		RecordHandler: CompGen,
	}

	newCompanyMongoRecordHandler.RecordHandler.Start()

	return &newCompanyMongoRecordHandler
}

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *companyRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (mrh *mongoRecordHandler) Create(request *companyRecordHandler.CreateRequest) (*companyRecordHandler.CreateResponse, error) {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return nil, err
	} //TODO CRUD validation can be generic
	resp, _ := mrh.GCreate(&recordHandler.CreateRequest{Entity: request.Company})
	comp, _ := resp.Entity.(company.Company)
	return &companyRecordHandler.CreateResponse{Company: comp}, nil
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
	}
	return nil
}

func (mrh *mongoRecordHandler) Retrieve(request *companyRecordHandler.RetrieveRequest) (*companyRecordHandler.RetrieveResponse, error) {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}
	return nil, nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *companyRecordHandler.UpdateRequest) error {
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

func (mrh *mongoRecordHandler) Update(request *companyRecordHandler.UpdateRequest) (*companyRecordHandler.UpdateResponse, error) {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return nil, err
	}

	return nil, nil
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
	}
	return nil
}

func (mrh *mongoRecordHandler) Delete(request *companyRecordHandler.DeleteRequest) (*companyRecordHandler.DeleteResponse, error) {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return nil, err
	}

	return nil, nil
}

func (mrh *mongoRecordHandler) ValidateCollectRequest(request *companyRecordHandler.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (mrh *mongoRecordHandler) Collect(request *companyRecordHandler.CollectRequest) (*companyRecordHandler.CollectResponse, error) {
	if err := mrh.ValidateCollectRequest(request); err != nil {
		return nil, err
	} //TODO use generic validation
	resp, _ := mrh.GCollect(&recordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	})

	compResp := make([]company.Company, 0)
	for _, c := range resp.Records {
		comp, _ := c.(company.Company)
		compResp = append(compResp, comp)
	}
	return &companyRecordHandler.CollectResponse{
		Records: compResp,
		Total:   resp.Total,
	}, nil
}
