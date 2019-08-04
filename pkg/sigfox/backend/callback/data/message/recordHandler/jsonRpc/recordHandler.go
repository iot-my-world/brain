package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
	sigfoxBackendDataCallbackMessageRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler"
	messageRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler/adaptor/jsonRpc"
)

type recordHandler struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) sigfoxBackendDataCallbackMessageRecordHandler.RecordHandler {
	return &recordHandler{
		jsonRpcClient: jsonRpcClient,
	}
}

func (r *recordHandler) Create(request *sigfoxBackendDataCallbackMessageRecordHandler.CreateRequest) (*sigfoxBackendDataCallbackMessageRecordHandler.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Retrieve(request *sigfoxBackendDataCallbackMessageRecordHandler.RetrieveRequest) (*sigfoxBackendDataCallbackMessageRecordHandler.RetrieveResponse, error) {
	return nil, brainException.NotImplemented{}
}
func (r *recordHandler) Update(request *sigfoxBackendDataCallbackMessageRecordHandler.UpdateRequest) (*sigfoxBackendDataCallbackMessageRecordHandler.UpdateResponse, error) {
	return nil, brainException.NotImplemented{}
}
func (r *recordHandler) Delete(request *sigfoxBackendDataCallbackMessageRecordHandler.DeleteRequest) (*sigfoxBackendDataCallbackMessageRecordHandler.DeleteResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) ValidateCollectRequest(request *sigfoxBackendDataCallbackMessageRecordHandler.CollectRequest) error {
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

func (r *recordHandler) Collect(request *sigfoxBackendDataCallbackMessageRecordHandler.CollectRequest) (*sigfoxBackendDataCallbackMessageRecordHandler.CollectResponse, error) {
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

	collectResponse := messageRecordHandlerJsonRpcAdaptor.CollectResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		sigfoxBackendDataCallbackMessageRecordHandler.CollectService,
		messageRecordHandlerJsonRpcAdaptor.CollectRequest{
			Criteria: criteria,
			Query:    request.Query,
		},
		&collectResponse); err != nil {
		return nil, err
	}

	return &sigfoxBackendDataCallbackMessageRecordHandler.CollectResponse{
		Records: collectResponse.Records,
		Total:   collectResponse.Total,
	}, nil
}
