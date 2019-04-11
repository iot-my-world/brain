package administrator

import (
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/device"
)

type Administrator interface {
	Collect(request *CollectRequest) (*CollectResponse, error)
	Create(request *CreateRequest) (*CreateResponse, error)
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []device.Device
	Total   int
}

type CreateRequest struct {
	Claims claims.Claims
	Device device.Device
}

type CreateResponse struct {
	Device device.Device
}
