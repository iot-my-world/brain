package client

import "encoding/json"

type Client interface {
	Post(request *Request) (*Response, error)
}

type Request struct {
	Id      string      `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func NewRequest(id, method string) Request {
	return Request{
		Id:      id,
		Method:  method,
		JsonRpc: "2.0",
	}
}

type Response struct {
	Id      string          `json:"id"`
	JsonRpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   string          `json:"error"`
}
