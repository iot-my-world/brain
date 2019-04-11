package recordHandler

import (
	brainRecordHandler "gitlab.com/iotTracker/brain/recordHandler"
	"gitlab.com/iotTracker/brain/tracker/device"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainCompanyRecordHandler brainRecordHandler.RecordHandler,
) *RecordHandler {

	return &RecordHandler{
		recordHandler: brainCompanyRecordHandler,
	}
}

type CreateRequest struct {
	Device device.Device
}

type CreateResponse struct {
	Device device.Device
}

func (r *RecordHandler) Create(request *CreateRequest) (*CreateResponse, error) {

	return &CreateResponse{}, nil
}
