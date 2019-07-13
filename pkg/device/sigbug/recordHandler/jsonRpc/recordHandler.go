package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	sigbugRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/adaptor/jsonRpc"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
)

type recordHandler struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) sigbugRecordHandler.RecordHandler {
	return &recordHandler{
		jsonRpcClient: jsonRpcClient,
	}
}

func (r *recordHandler) Create(request *sigbugRecordHandler.CreateRequest) (*sigbugRecordHandler.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Retrieve(request *sigbugRecordHandler.RetrieveRequest) (*sigbugRecordHandler.RetrieveResponse, error) {
	return nil, brainException.NotImplemented{}
}
func (r *recordHandler) Update(request *sigbugRecordHandler.UpdateRequest) (*sigbugRecordHandler.UpdateResponse, error) {
	return nil, brainException.NotImplemented{}
}
func (r *recordHandler) Delete(request *sigbugRecordHandler.DeleteRequest) (*sigbugRecordHandler.DeleteResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) ValidateCollectRequest(request *sigbugRecordHandler.CollectRequest) error {
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

func (r *recordHandler) Collect(request *sigbugRecordHandler.CollectRequest) (*sigbugRecordHandler.CollectResponse, error) {
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

	collectResponse := sigbugRecordHandlerJsonRpcAdaptor.CollectResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		sigbugRecordHandler.CollectService,
		sigbugRecordHandlerJsonRpcAdaptor.CollectRequest{
			Criteria: criteria,
			Query:    request.Query,
		},
		&collectResponse); err != nil {
		return nil, err
	}

	return &sigbugRecordHandler.CollectResponse{
		Records: collectResponse.Records,
		Total:   collectResponse.Total,
	}, nil
}
