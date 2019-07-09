package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	clientRecordHandler "github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	"github.com/iot-my-world/brain/pkg/party/client/recordHandler/adaptor/jsonRpc"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
)

type recordHandler struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) clientRecordHandler.RecordHandler {
	return &recordHandler{
		jsonRpcClient: jsonRpcClient,
	}
}

func (r *recordHandler) ValidateCollectRequest(request *clientRecordHandler.CollectRequest) error {
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

func (r *recordHandler) Collect(request *clientRecordHandler.CollectRequest) (*clientRecordHandler.CollectResponse, error) {
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

	clientCollectResponse := client.CollectResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		clientRecordHandler.CollectService,
		client.CollectRequest{
			Criteria: criteria,
			Query:    request.Query,
		},
		&clientCollectResponse); err != nil {
		return nil, err
	}

	return &clientRecordHandler.CollectResponse{
		Records: clientCollectResponse.Records,
		Total:   clientCollectResponse.Total,
	}, nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *clientRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Retrieve(request *clientRecordHandler.RetrieveRequest) (*clientRecordHandler.RetrieveResponse, error) {
	if err := r.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	// wrap identifier
	id, err := wrappedIdentifier.Wrap(request.Identifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	clientRetrieveResponse := client.RetrieveResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		clientRecordHandler.RetrieveService,
		client.RetrieveRequest{
			WrappedIdentifier: *id,
		},
		&clientRetrieveResponse); err != nil {
		return nil, err
	}

	return &clientRecordHandler.RetrieveResponse{
		Client: clientRetrieveResponse.Client,
	}, nil
}

func (r *recordHandler) Create(request *clientRecordHandler.CreateRequest) (*clientRecordHandler.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Update(request *clientRecordHandler.UpdateRequest) (*clientRecordHandler.UpdateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Delete(request *clientRecordHandler.DeleteRequest) (*clientRecordHandler.DeleteResponse, error) {
	return nil, brainException.NotImplemented{}
}
