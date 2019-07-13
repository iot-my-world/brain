package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	recordHandler2 "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	"github.com/iot-my-world/brain/pkg/user/human/recordHandler/adaptor/jsonRpc"
)

type recordHandler struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) recordHandler2.RecordHandler {
	return &recordHandler{
		jsonRpcClient: jsonRpcClient,
	}
}

func (r *recordHandler) ValidateCollectRequest(request *recordHandler2.CollectRequest) error {
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

func (r *recordHandler) Collect(request *recordHandler2.CollectRequest) (*recordHandler2.CollectResponse, error) {
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

	clientCollectResponse := user.CollectResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		recordHandler2.CollectService,
		user.CollectRequest{
			Criteria: criteria,
			Query:    request.Query,
		},
		&clientCollectResponse); err != nil {
		return nil, err
	}

	return &recordHandler2.CollectResponse{
		Records: clientCollectResponse.Records,
		Total:   clientCollectResponse.Total,
	}, nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *recordHandler2.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Retrieve(request *recordHandler2.RetrieveRequest) (*recordHandler2.RetrieveResponse, error) {
	if err := r.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	// wrap identifier
	id, err := wrappedIdentifier.Wrap(request.Identifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	clientRetrieveResponse := user.RetrieveResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		recordHandler2.RetrieveService,
		user.RetrieveRequest{
			WrappedIdentifier: *id,
		},
		&clientRetrieveResponse); err != nil {
		return nil, err
	}

	return &recordHandler2.RetrieveResponse{
		User: clientRetrieveResponse.User,
	}, nil
}

func (r *recordHandler) Create(request *recordHandler2.CreateRequest) (*recordHandler2.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Update(request *recordHandler2.UpdateRequest) (*recordHandler2.UpdateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Delete(request *recordHandler2.DeleteRequest) (*recordHandler2.DeleteResponse, error) {
	return nil, brainException.NotImplemented{}
}
