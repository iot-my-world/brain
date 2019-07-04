package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	recordHandler2 "github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	"github.com/iot-my-world/brain/pkg/party/company/recordHandler/adaptor/jsonRpc"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
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

	companyCollectResponse := company.CollectResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		recordHandler2.CollectService,
		company.CollectRequest{
			Criteria: criteria,
			Query:    request.Query,
		},
		&companyCollectResponse); err != nil {
		return nil, err
	}

	return &recordHandler2.CollectResponse{
		Records: companyCollectResponse.Records,
		Total:   companyCollectResponse.Total,
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

	companyRetrieveResponse := company.RetrieveResponse{}
	if err := r.jsonRpcClient.JsonRpcRequest(
		recordHandler2.RetrieveService,
		company.RetrieveRequest{
			WrappedIdentifier: *id,
		},
		&companyRetrieveResponse); err != nil {
		return nil, err
	}

	return &recordHandler2.RetrieveResponse{
		Company: companyRetrieveResponse.Company,
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
