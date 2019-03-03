package tracking

import (
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/party/client"
	"gitlab.com/iotTracker/brain/tracker/reading"
)

type Report interface {
	Live(request *LiveRequest, response *LiveResponse) error
	Historical(request *HistoricalRequest, response *HistoricalResponse) error
}

type LiveRequest struct {
	Claims          claims.Claims
	ClientCriteria  []criterion.Criterion
	CompanyCriteria []criterion.Criterion
}

type LiveResponse struct {
	Companies []company.Company
	Clients   []client.Client
	Readings  []reading.Reading
}

type HistoricalRequest struct {
	Claims          claims.Claims
	ClientCriteria  []criterion.Criterion
	CompanyCriteria []criterion.Criterion
}

type HistoricalResponse struct {
	Companies []company.Company
	Clients   []client.Client
	Readings  []reading.Reading
}
