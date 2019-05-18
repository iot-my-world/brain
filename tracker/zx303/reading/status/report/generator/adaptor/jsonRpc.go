package adaptor

import (
	zx303StatusReadingReportGenerator "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/report/generator"
	"net/http"
)

type Adaptor struct {
	zx303StatusReadingReportGenerator zx303StatusReadingReportGenerator.Generator
}

func New(
	zx303StatusReadingReportGenerator zx303StatusReadingReportGenerator.Generator,
) *Adaptor {
	return &Adaptor{
		zx303StatusReadingReportGenerator: zx303StatusReadingReportGenerator,
	}
}

func (a *Adaptor) Battery(r *http.Request, request *CreateRequest, response *CreateResponse) error {

}
