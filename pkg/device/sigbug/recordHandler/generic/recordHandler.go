package sigbugRecordHandler

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	sigbugRecordHandlerException "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/exception"
	brainRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/pkg/recordHandler/exception"
)

type RecordHandler struct {
	sigbugRecordHandler brainRecordHandler.RecordHandler
}

func New(
	brainSigbugRecordHandler brainRecordHandler.RecordHandler,
) sigbugRecordHandler.RecordHandler {

	return &RecordHandler{
		sigbugRecordHandler: brainSigbugRecordHandler,
	}
}

type CreateRequest struct {
	Sigbug sigbug.Sigbug
}

type CreateResponse struct {
	Sigbug sigbug.Sigbug
}

func (r *RecordHandler) ValidateCreateRequest(request *sigbugRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *RecordHandler) Create(request *sigbugRecordHandler.CreateRequest) (*sigbugRecordHandler.CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.sigbugRecordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.Sigbug,
	}, &createResponse); err != nil {
		return nil, sigbugRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdDevice, ok := createResponse.Entity.(*sigbug.Sigbug)
	if !ok {
		return nil, sigbugRecordHandlerException.Create{Reasons: []string{"could not cast created entity to sigbug"}}
	}

	return &sigbugRecordHandler.CreateResponse{
		Sigbug: *createdDevice,
	}, nil
}

func (r *RecordHandler) Retrieve(request *sigbugRecordHandler.RetrieveRequest) (*sigbugRecordHandler.RetrieveResponse, error) {
	retrievedSigbug := sigbug.Sigbug{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedSigbug,
	}
	if err := r.sigbugRecordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, sigbugRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &sigbugRecordHandler.RetrieveResponse{
		Sigbug: retrievedSigbug,
	}, nil
}

func (r *RecordHandler) Update(request *sigbugRecordHandler.UpdateRequest) (*sigbugRecordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.sigbugRecordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Sigbug,
	}, &updateResponse); err != nil {
		return nil, sigbugRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &sigbugRecordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *sigbugRecordHandler.DeleteRequest) (*sigbugRecordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.sigbugRecordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, sigbugRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &sigbugRecordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *sigbugRecordHandler.CollectRequest) (*sigbugRecordHandler.CollectResponse, error) {
	var collectedSigbug []sigbug.Sigbug
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedSigbug,
	}
	err := r.sigbugRecordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, sigbugRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedSigbug == nil {
		collectedSigbug = make([]sigbug.Sigbug, 0)
	}

	return &sigbugRecordHandler.CollectResponse{
		Records: collectedSigbug,
		Total:   collectResponse.Total,
	}, nil
}
