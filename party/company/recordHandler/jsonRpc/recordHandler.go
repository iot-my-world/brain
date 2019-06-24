package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	companyRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler"
	companyRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/party/company/recordHandler/adaptor/jsonRpc"
	wrappedCriterion "github.com/iot-my-world/brain/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
)

type RecordHandler struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) *RecordHandler {
	return &RecordHandler{
		jsonRpcClient: jsonRpcClient,
	}
}

func (r *RecordHandler) ValidateCollectRequest(request *companyRecordHandler.CollectRequest) error {
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

func (r *RecordHandler) Collect(request *companyRecordHandler.CollectRequest) (*companyRecordHandler.CollectResponse, error) {
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
		"CompanyRecordHandler.Collect",
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

func (r *RecordHandler) ValidateRetrieveRequest(request *companyRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (r *RecordHandler) Retrieve(request *companyRecordHandler.RetrieveRequest) (*companyRecordHandler.RetrieveResponse, error) {
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
		"CompanyRecordHandler.Retrieve",
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

func (r *RecordHandler) Create(request *companyRecordHandler.CreateRequest) (*companyRecordHandler.CreateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *RecordHandler) Update(request *companyRecordHandler.UpdateRequest) (*companyRecordHandler.UpdateResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (r *RecordHandler) Delete(request *companyRecordHandler.DeleteRequest) (*companyRecordHandler.DeleteResponse, error) {
	return nil, brainException.NotImplemented{}
}
