package user

import (
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	"github.com/iot-my-world/brain/pkg/search/query"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/user/human"
	"github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	"net/http"
)

type adaptor struct {
	RecordHandler recordHandler.RecordHandler
}

func New(recordHandler recordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(recordHandler.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type RetrieveRequest struct {
	WrappedIdentifier wrappedIdentifier.Wrapped `json:"identifier"`
}

type RetrieveResponse struct {
	User human.User `json:"user" bson:"user"`
}

func (a *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	retrieveUserResponse, err := a.RecordHandler.Retrieve(
		&recordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: request.WrappedIdentifier.Identifier,
		})
	if err != nil {
		return err
	}

	response.User = retrieveUserResponse.User

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.Wrapped `json:"criteria"`
	Query    query.Query                `json:"query"`
}

type CollectResponse struct {
	Records []human.User `json:"records"`
	Total   int          `json:"total"`
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

	collectCompanyResponse, err := a.RecordHandler.Collect(&recordHandler.CollectRequest{
		Criteria: criteria,
		Query:    request.Query,
		Claims:   claims,
	})
	if err != nil {
		return err
	}

	response.Records = collectCompanyResponse.Records
	response.Total = collectCompanyResponse.Total
	return nil
}
