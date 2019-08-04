package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	"github.com/iot-my-world/brain/pkg/search/query"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	sigfoxBackendDataCallbackMessageRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler"
	"net/http"
)

type adaptor struct {
	RecordHandler sigfoxBackendDataCallbackMessageRecordHandler.RecordHandler
}

func New(recordHandler sigfoxBackendDataCallbackMessageRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(sigfoxBackendDataCallbackMessageRecordHandler.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type RetrieveRequest struct {
	WrappedIdentifier wrappedIdentifier.Wrapped `json:"identifier"`
}

type RetrieveResponse struct {
	Message sigfoxBackendDataCallbackMessage.Message `json:"message"`
}

func (a *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	retrieveMessageResponse, err := a.RecordHandler.Retrieve(
		&sigfoxBackendDataCallbackMessageRecordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: request.WrappedIdentifier.Identifier,
		})
	if err != nil {
		return err
	}

	response.Message = retrieveMessageResponse.Message

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.Wrapped `json:"criteria"`
	Query    query.Query                `json:"query"`
}

type CollectResponse struct {
	Records []sigfoxBackendDataCallbackMessage.Message `json:"records"`
	Total   int                                        `json:"total"`
}

func (a *adaptor) Collect(r *http.Request, request *CollectRequest, response *CollectResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	criteria := make([]criterion.Criterion, 0)
	for criterionIdx := range request.Criteria {
		if c, err := request.Criteria[criterionIdx].UnWrap(); err == nil {
			criteria = append(criteria, c)
		} else {
			return err
		}
	}

	collectMessageResponse, err := a.RecordHandler.Collect(&sigfoxBackendDataCallbackMessageRecordHandler.CollectRequest{
		Claims:   claims,
		Criteria: criteria,
		Query:    request.Query,
	})
	if err != nil {
		return err
	}

	response.Records = collectMessageResponse.Records
	response.Total = collectMessageResponse.Total
	return nil
}
