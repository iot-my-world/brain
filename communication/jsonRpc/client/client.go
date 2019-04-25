package client

import (
	"encoding/json"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/authorization/service/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Client interface {
	Post(request *Request) (*Response, error)
	JsonRpcRequest(method string, request, response interface{}) error
	Login(authJsonRpcAdaptor.LoginRequest) error
	Logout()
	Claims() claims.Claims
	SetJWT(jwt string) error
	GetJWT() string
	LoggedIn() bool
}

type Request struct {
	Id      string      `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func NewRequest(id, method string, params interface{}) Request {
	return Request{
		Id:      id,
		Method:  method,
		JsonRpc: "2.0",
		Params:  params,
	}
}

type Response struct {
	Id      string          `json:"id"`
	JsonRpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   string          `json:"error"`
}
