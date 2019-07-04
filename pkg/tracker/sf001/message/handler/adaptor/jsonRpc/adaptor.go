package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
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
