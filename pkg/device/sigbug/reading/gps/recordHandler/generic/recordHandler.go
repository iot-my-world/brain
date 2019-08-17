package sigbugGPSReadingRecordHandler

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	sigbugGPSReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	sigbugGPSReadingRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/recordHandler"
	brainRecordHandler "github.com/iot-my-world/brain/pkg/recordHandler"
	brainRecordHandlerException "github.com/iot-my-world/brain/pkg/recordHandler/exception"
	sigfoxBackendDataCallbackReadingRecordHandlerException "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler/exception"
)

type RecordHandler struct {
	sigbugGPSReadingRecordHandler brainRecordHandler.RecordHandler
}

func New(
	brainReadingRecordHandler brainRecordHandler.RecordHandler,
) sigbugGPSReadingRecordHandler.RecordHandler {

	return &RecordHandler{
		sigbugGPSReadingRecordHandler: brainReadingRecordHandler,
	}
}

type CreateRequest struct {
	Reading sigbugGPSReading.Reading
}

type CreateResponse struct {
	Reading sigbugGPSReading.Reading
}

func (r *RecordHandler) ValidateCreateRequest(request *sigbugGPSReadingRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r *RecordHandler) Create(request *sigbugGPSReadingRecordHandler.CreateRequest) (*sigbugGPSReadingRecordHandler.CreateResponse, error) {
	if err := r.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.sigbugGPSReadingRecordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.Reading,
	}, &createResponse); err != nil {
		return nil, sigfoxBackendDataCallbackReadingRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdReading, ok := createResponse.Entity.(*sigbugGPSReading.Reading)
	if !ok {
		return nil, sigfoxBackendDataCallbackReadingRecordHandlerException.Create{Reasons: []string{"could not cast created entity to message"}}
	}

	return &sigbugGPSReadingRecordHandler.CreateResponse{
		Reading: *createdReading,
	}, nil
}

func (r *RecordHandler) Retrieve(request *sigbugGPSReadingRecordHandler.RetrieveRequest) (*sigbugGPSReadingRecordHandler.RetrieveResponse, error) {
	retrievedReading := sigbugGPSReading.Reading{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedReading,
	}
	if err := r.sigbugGPSReadingRecordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, sigfoxBackendDataCallbackReadingRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &sigbugGPSReadingRecordHandler.RetrieveResponse{
		Reading: retrievedReading,
	}, nil
}

func (r *RecordHandler) Update(request *sigbugGPSReadingRecordHandler.UpdateRequest) (*sigbugGPSReadingRecordHandler.UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.sigbugGPSReadingRecordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Reading,
	}, &updateResponse); err != nil {
		return nil, sigfoxBackendDataCallbackReadingRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &sigbugGPSReadingRecordHandler.UpdateResponse{}, nil
}

func (r *RecordHandler) Delete(request *sigbugGPSReadingRecordHandler.DeleteRequest) (*sigbugGPSReadingRecordHandler.DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.sigbugGPSReadingRecordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, sigfoxBackendDataCallbackReadingRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &sigbugGPSReadingRecordHandler.DeleteResponse{}, nil
}

func (r *RecordHandler) Collect(request *sigbugGPSReadingRecordHandler.CollectRequest) (*sigbugGPSReadingRecordHandler.CollectResponse, error) {
	var collectedReading []sigbugGPSReading.Reading
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedReading,
	}
	err := r.sigbugGPSReadingRecordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, sigfoxBackendDataCallbackReadingRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	if collectedReading == nil {
		collectedReading = make([]sigbugGPSReading.Reading, 0)
	}

	return &sigbugGPSReadingRecordHandler.CollectResponse{
		Records: collectedReading,
		Total:   collectResponse.Total,
	}, nil
}
