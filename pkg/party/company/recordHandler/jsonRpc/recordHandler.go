package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	companyRecordHandler "github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	companyRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/company/recordHandler/adaptor/jsonRpc"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
)

type recordHandler struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) companyRecordHandler.RecordHandler {
	return &recordHandler{
		jsonRpcClient: jsonRpcClient,
	}
}

func (r *recordHandler) ValidateCollectRequest(request *companyRecordHandler.CollectRequest) error {
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

func (r *recordHandler) Collect(request *companyRecordHandler.CollectRequest) (*companyRecordHandler.CollectResponse, error) {
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

	companyCollectResponse := companyRecordHandlerJsonRpcAdaptor.CollectResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		companyRecordHandler.CollectService,
		companyRecordHandlerJsonRpcAdaptor.CollectRequest{
			Criteria: criteria,
			Query:    request.Query,
		},
		&companyCollectResponse); err != nil {
		return nil, err
	}

	return &companyRecordHandler.CollectResponse{
		Records: companyCollectResponse.Records,
		Total:   companyCollectResponse.Total,
	}, nil
}

func (r *recordHandler) ValidateRetrieveRequest(request *companyRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *recordHandler) Retrieve(request *companyRecordHandler.RetrieveRequest) (*companyRecordHandler.RetrieveResponse, error) {
	if err := r.ValidateRetrieveRequest(request); err != nil {
		return nil, err
	}

	// wrap identifier
	id, err := wrappedIdentifier.Wrap(request.Identifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	companyRetrieveResponse := companyRecordHandlerJsonRpcAdaptor.RetrieveResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		companyRecordHandler.RetrieveService,
		companyRecordHandlerJsonRpcAdaptor.RetrieveRequest{
			WrappedIdentifier: *id,
		},
		&companyRetrieveResponse); err != nil {
		return nil, err
	}

	return &companyRecordHandler.RetrieveResponse{
		Company: companyRetrieveResponse.Company,
	}, nil
}

func (r *recordHandler) Create(request *companyRecordHandler.CreateRequest) (*companyRecordHandler.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Update(request *companyRecordHandler.UpdateRequest) (*companyRecordHandler.UpdateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *recordHandler) Delete(request *companyRecordHandler.DeleteRequest) (*companyRecordHandler.DeleteResponse, error) {
	return nil, brainException.NotImplemented{}
}
