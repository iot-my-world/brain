package basic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/satori/go.uuid"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	brainException "gitlab.com/iotTracker/brain/exception"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/authorization/service/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/security/claims"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gopkg.in/square/go-jose.v2"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type client struct {
	url      string
	jwt      string
	claims   claims.Claims
	loggedIn bool
}

// Create New basic json rpc client
func New(
	url string,
) jsonRpcClient.Client {
	return &client{
		url: url,
	}
}

func (c *client) LoggedIn() bool {
	return c.loggedIn
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
		return nil, errors.New("error unmarshalling response bytes into json rpc response: " + err.Error())
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

func (c *client) Login(loginRequest authJsonRpcAdaptor.LoginRequest) error {
	loginResponse := authJsonRpcAdaptor.LoginResponse{}

	if err := c.JsonRpcRequest(
		"Auth.Login",
		loginRequest,
		&loginResponse,
	); err != nil {
		return err
	}

	// save the token
	c.jwt = loginResponse.Jwt

	object, err := jose.ParseSigned(c.jwt)
	if err != nil {
		return errors.New("error parsing jwt " + err.Error())
	}

	// Access Underlying payload without verification
	fv := reflect.ValueOf(object).Elem().FieldByName("payload")

	wrapped := wrappedClaims.Wrapped{}
	if err := json.Unmarshal(fv.Bytes(), &wrapped); err != nil {
		return errors.New("error unmarshalling claims " + err.Error())
	}

	unwrappedClaims, err := wrapped.Unwrap()
	if err != nil {
		return errors.New("error unwrapping claims " + err.Error())
	}

	c.claims = unwrappedClaims
	c.loggedIn = true

	return nil
}

func (c *client) Logout() {
	c.jwt = ""
	c.claims = nil
	c.loggedIn = false
}

func (c *client) Claims() claims.Claims {
	return c.claims
}

func (c *client) SetJWT(jwt string) error {
	// save the token
	c.jwt = jwt

	object, err := jose.ParseSigned(c.jwt)
	if err != nil {
		return errors.New("error parsing jwt " + err.Error())
	}

	// Access Underlying payload without verification
	fv := reflect.ValueOf(object).Elem().FieldByName("payload")

	wrapped := wrappedClaims.Wrapped{}
	if err := json.Unmarshal(fv.Bytes(), &wrapped); err != nil {
		return errors.New("error unmarshalling claims " + err.Error())
	}

	unwrappedClaims, err := wrapped.Unwrap()
	if err != nil {
		return errors.New("error unwrapping claims " + err.Error())
	}

	c.claims = unwrappedClaims

	return nil
}

func (c *client) GetJWT() string {
	return c.jwt
}
