package recordHandler

import (
	"github.com/iot-my-world/brain/pkg/party/company"
	"github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	"github.com/iot-my-world/brain/pkg/party/company/recordHandler/exception"
	brainRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/pkg/recordHandler/exception"
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

func (r *RecordHandler) Create(request *recordHandler.CreateRequest) (*recordHandler.CreateResponse, error) {
	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.Company,
	}, &createResponse); err != nil {
		return nil, exception.Create{Reasons: []string{err.Error()}}
	}
	createdCompany, ok := createResponse.Entity.(*company.Company)
	if !ok {
		return nil, exception.Create{Reasons: []string{"could not cast created entity to company"}}
	}

	return &recordHandler.CreateResponse{
		Company: *createdCompany,
	}, nil
}

func (r *RecordHandler) Retrieve(request *recordHandler.RetrieveRequest) (*recordHandler.RetrieveResponse, error) {
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
			return nil, exception.NotFound{}
		default:
			return nil, err
		}
	}

	return &recordHandler.RetrieveResponse{
		Company: retrievedCompany,
	}, nil
}

func (r *RecordHandler) Update(request *recordHandler.UpdateRequest) (*recordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Company,
	}, &updateResponse); err != nil {
		return nil, exception.Update{Reasons: []string{err.Error()}}
	}

	return &recordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *recordHandler.DeleteRequest) (*recordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, exception.Delete{Reasons: []string{err.Error()}}
	}

	return &recordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *recordHandler.CollectRequest) (*recordHandler.CollectResponse, error) {
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
		return nil, exception.Collect{Reasons: []string{err.Error()}}
	}

	if collectedCompanies == nil {
		collectedCompanies = make([]company.Company, 0)
	}

	return &recordHandler.CollectResponse{
		Records: collectedCompanies,
		Total:   collectResponse.Total,
	}, nil
}
