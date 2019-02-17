package jsonRpc

import (
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
)

type adaptor struct {
	RecordHandler readingRecordHandler.RecordHandler
}

func New(recordHandler readingRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

