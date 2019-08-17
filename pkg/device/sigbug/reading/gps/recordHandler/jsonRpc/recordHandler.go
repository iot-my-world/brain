package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	sigbugGPSReadingRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/recordHandler"
	sigbugGPSReadingRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/recordHandler/adaptor/jsonRpc"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
)

type recordHandler struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) sigbugGPSReadingRecordHandler.RecordHandler {
	return &recordHandler{
		jsonRpcClient: jsonRpcClient,
	}
}

func (r *recordHandler) Create(request *sigbugGPSReadingRecordHandler.CreateRequest) (*sigbugGPSReadingRecordHandler.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Retrieve(request *sigbugGPSReadingRecordHandler.RetrieveRequest) (*sigbugGPSReadingRecordHandler.RetrieveResponse, error) {
	return nil, brainException.NotImplemented{}
}
func (r *recordHandler) Update(request *sigbugGPSReadingRecordHandler.UpdateRequest) (*sigbugGPSReadingRecordHandler.UpdateResponse, error) {
	return nil, brainException.NotImplemented{}
}
func (r *recordHandler) Delete(request *sigbugGPSReadingRecordHandler.DeleteRequest) (*sigbugGPSReadingRecordHandler.DeleteResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) ValidateCollectRequest(request *sigbugGPSReadingRecordHandler.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Criteria == nil {
		reasonsInvalid = append(reasonsInvalid, "criteria is nil")
	} else {
		for _, crit := range request.Criteria {
			if crit == nil {
				reasonsInvalid = append(reasonsInvalid, "a criterion is nil")
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Collect(request *sigbugGPSReadingRecordHandler.CollectRequest) (*sigbugGPSReadingRecordHandler.CollectResponse, error) {
	if err := r.ValidateCollectRequest(request); err != nil {
		return nil, err
	}

	// wrap criteria
	criteria := make([]wrappedCriterion.Wrapped, 0)
	for _, crit := range request.Criteria {
		wrapped, err := wrappedCriterion.Wrap(crit)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		criteria = append(criteria, *wrapped)
	}

	collectResponse := sigbugGPSReadingRecordHandlerJsonRpcAdaptor.CollectResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		sigbugGPSReadingRecordHandler.CollectService,
		sigbugGPSReadingRecordHandlerJsonRpcAdaptor.CollectRequest{
			Criteria: criteria,
			Query:    request.Query,
		},
		&collectResponse); err != nil {
		return nil, err
	}

	return &sigbugGPSReadingRecordHandler.CollectResponse{
		Records: collectResponse.Records,
		Total:   collectResponse.Total,
	}, nil
}
