package recordHandler

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	brainRecordHandler "gitlab.com/iotTracker/brain/recordHandler"
	brainRecordHandlerException "gitlab.com/iotTracker/brain/recordHandler/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/device/zx303"
	zx303RecordHandlerException "gitlab.com/iotTracker/brain/tracker/device/zx303/recordHandler/exception"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainZX303RecordHandler brainRecordHandler.RecordHandler,
) *RecordHandler {

	return &RecordHandler{
		recordHandler: brainZX303RecordHandler,
	}
}

type CreateRequest struct {
	ZX303 zx303.ZX303
}

type CreateResponse struct {
	ZX303 zx303.ZX303
}

func (r *RecordHandler) ValidateCreateRequest(request *CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *RecordHandler) Create(request *CreateRequest) (*CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.ZX303,
	}, &createResponse); err != nil {
		return nil, zx303RecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdDevice, ok := createResponse.Entity.(*zx303.ZX303)
	if !ok {
		return nil, zx303RecordHandlerException.Create{Reasons: []string{"could not cast created entity to zx303"}}
	}

	return &CreateResponse{
		ZX303: *createdDevice,
	}, nil
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	ZX303 zx303.ZX303
}

func (r *RecordHandler) Retrieve(request *RetrieveRequest) (*RetrieveResponse, error) {
	retrievedZX303 := zx303.ZX303{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedZX303,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, zx303RecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &RetrieveResponse{
		ZX303: retrievedZX303,
	}, nil
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	ZX303      zx303.ZX303
}

type UpdateResponse struct{}

func (r *RecordHandler) Update(request *UpdateRequest) (*UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.ZX303,
	}, &updateResponse); err != nil {
		return nil, zx303RecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &UpdateResponse{}, nil
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
}

func (r *RecordHandler) Delete(request *DeleteRequest) (*DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, zx303RecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &DeleteResponse{}, nil
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []zx303.ZX303
	Total   int
}

func (r *RecordHandler) Collect(request *CollectRequest) (*CollectResponse, error) {
	var collectedZX303 []zx303.ZX303
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedZX303,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, zx303RecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedZX303 == nil {
		collectedZX303 = make([]zx303.ZX303, 0)
	}

	return &CollectResponse{
		Records: collectedZX303,
		Total:   collectResponse.Total,
	}, nil
}
