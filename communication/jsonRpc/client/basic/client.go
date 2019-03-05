package basic

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"bytes"
	"net/http"
	"fmt"
	"strings"
	"io/ioutil"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	brainException "gitlab.com/iotTracker/brain/exception"
	"github.com/satori/go.uuid"
)

type client struct {
	url string
	jwt string
}

func New(
	url string,
) *client {
	return &client{
		url: url,
	}
}

func (c *client) Post(request *jsonRpcClient.Request) (*jsonRpcClient.Response, error) {
	// marshal the request message
	marshalledRequest, err := json.Marshal(*request)
	if err != nil {
		return nil, errors.New("error marshalling request " + err.Error())
	}

	// put the bytes of the marshalled request into a buffer
	body := bytes.NewBuffer(marshalledRequest)

	// build the post request
	postRequest, err := http.NewRequest("POST", fmt.Sprintf("%s", c.url), body)
	if err != nil {
		return nil, errors.New("error creating post request " + err.Error())
	}

	// set the required headers on the request
	postRequest.Header.Set("Content-Type", "application/json")
	postRequest.Header.Set("Access-Control-Allow-Origin", "*")
	if c.jwt != "" {
		postRequest.Header.Set("Authorization", c.jwt)
	}

	// create the http client
	httpClient := &http.Client{
		//Timeout: time.Second * 5,
	}

	// perform the request
	postResponse, err := httpClient.Do(postRequest)
	if err != nil {
		return nil, errors.New("error performing post request " + err.Error())
	}

	// read the body bytes of the response
	postResponseBytes, err := ioutil.ReadAll(postResponse.Body)
	defer postResponse.Body.Close()
	if err != nil {
		return nil, errors.New("error reading post response body bytes " + err.Error())
	}

	// check for an rpc error
	if strings.Contains(string(postResponseBytes), "rpc: can't find service") {
		return nil, errors.New("rpc error: method not found")
	}

	// unmarshal the body into the response
	response := jsonRpcClient.Response{}
	err = json.Unmarshal(postResponseBytes, &response)
	if err != nil {
		return nil, errors.New("error unmarshalling response bytes into json rpc response " + err.Error())
	}

	if response.Error != "" {
		return &response, errors.New("json rpc service error " + response.Error)
	}

	return &response, nil
}

func (c *client) JsonRpcRequest(method string, request, response interface{}) error {
	id, err := uuid.NewV4()
	if err != nil {
		return brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}

	jsonRpcRequest := jsonRpcClient.NewRequest(id.String(), method, [1]interface{}{request})

	jsonRpcResponse, err := c.Post(&jsonRpcRequest)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonRpcResponse.Result, response); err != nil {
		return err
	}

	return nil
}
