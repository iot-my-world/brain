package recordHandler

import (
	"gitlab.com/iotTracker/brain/tracker/reading"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/query"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Collect(request *CollectRequest, response *CollectResponse) error
}

type CreateRequest struct {
	Reading reading.Reading
}

type CreateResponse struct {
	Reading reading.Reading
}

type CollectRequest struct {
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []reading.Reading
	Total   int
}
