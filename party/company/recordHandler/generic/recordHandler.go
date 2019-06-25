package recordHandler

import (
	"github.com/iot-my-world/brain/party/company"
	companyRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler"
	companyRecordHandlerException "github.com/iot-my-world/brain/party/company/recordHandler/exception"
	brainRecordHandler "github.com/iot-my-world/brain/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainCompanyRecordHandler brainRecordHandler.RecordHandler,
) *RecordHandler {

	return &RecordHandler{
		recordHandler: brainCompanyRecordHandler,
	}
}

func (r *RecordHandler) Create(request *companyRecordHandler.CreateRequest) (*companyRecordHandler.CreateResponse, error) {
	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.Company,
	}, &createResponse); err != nil {
		return nil, companyRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdCompany, ok := createResponse.Entity.(*company.Company)
	if !ok {
		return nil, companyRecordHandlerException.Create{Reasons: []string{"could not cast created entity to company"}}
	}

	return &companyRecordHandler.CreateResponse{
		Company: *createdCompany,
	}, nil
}

func (r *RecordHandler) Retrieve(request *companyRecordHandler.RetrieveRequest) (*companyRecordHandler.RetrieveResponse, error) {
	retrievedCompany := company.Company{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedCompany,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, companyRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &companyRecordHandler.RetrieveResponse{
		Company: retrievedCompany,
	}, nil
}

func (r *RecordHandler) Update(request *companyRecordHandler.UpdateRequest) (*companyRecordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Company,
	}, &updateResponse); err != nil {
		return nil, companyRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &companyRecordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *companyRecordHandler.DeleteRequest) (*companyRecordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, companyRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &companyRecordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *companyRecordHandler.CollectRequest) (*companyRecordHandler.CollectResponse, error) {
	var collectedCompanies []company.Company
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedCompanies,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, companyRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	return &companyRecordHandler.CollectResponse{
		Records: collectedCompanies,
		Total:   collectResponse.Total,
	}, nil
}
