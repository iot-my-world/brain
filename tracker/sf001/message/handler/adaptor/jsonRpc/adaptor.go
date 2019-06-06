package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"net/http"
)

type Adaptor struct {
}

func New() *Adaptor {
	return &Adaptor{}
}

type TestRequest struct {
}

type TestResponse struct {
}

func (a *Adaptor) Test(r *http.Request, request *TestRequest, response *TestResponse) error {

	log.Info("Test success!")

	return nil
}
