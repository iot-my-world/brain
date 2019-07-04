package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	wrappedCriterion "github.com/iot-my-world/brain/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	humanUserRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler"
	humanUserRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/user/human/recordHandler/adaptor/jsonRpc"
)

type recordHandler struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) humanUserRecordHandler.RecordHandler {
	return &recordHandler{
		jsonRpcClient: jsonRpcClient,
	}
}

func (r *recordHandler) ValidateCollectRequest(request *humanUserRecordHandler.CollectRequest) error {
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

func (r *recordHandler) Collect(request *humanUserRecordHandler.CollectRequest) (*humanUserRecordHandler.CollectResponse, error) {
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

	clientCollectResponse := humanUserRecordHandlerJsonRpcAdaptor.CollectResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		humanUserRecordHandler.CollectService,
		humanUserRecordHandlerJsonRpcAdaptor.CollectRequest{
			Criteria: criteria,
			Query:    request.Query,
		},
		&clientCollectResponse); err != nil {
		return nil, err
	}

	return &humanUserRecordHandler.CollectResponse{
		Records: clientCollectResponse.Records,
		Total:   clientCollectResponse.Total,
	}, nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *humanUserRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Retrieve(request *humanUserRecordHandler.RetrieveRequest) (*humanUserRecordHandler.RetrieveResponse, error) {
	if err := r.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	// wrap identifier
	id, err := wrappedIdentifier.Wrap(request.Identifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	clientRetrieveResponse := humanUserRecordHandlerJsonRpcAdaptor.RetrieveResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		humanUserRecordHandler.RetrieveService,
		humanUserRecordHandlerJsonRpcAdaptor.RetrieveRequest{
			WrappedIdentifier: *id,
		},
		&clientRetrieveResponse); err != nil {
		return nil, err
	}

	return &humanUserRecordHandler.RetrieveResponse{
		User: clientRetrieveResponse.User,
	}, nil
}

func (r *recordHandler) Create(request *humanUserRecordHandler.CreateRequest) (*humanUserRecordHandler.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Update(request *humanUserRecordHandler.UpdateRequest) (*humanUserRecordHandler.UpdateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Delete(request *humanUserRecordHandler.DeleteRequest) (*humanUserRecordHandler.DeleteResponse, error) {
	return nil, brainException.NotImplemented{}
}
