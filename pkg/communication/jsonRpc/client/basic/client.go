package basic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	client2 "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	authorizationAdministrator "github.com/iot-my-world/brain/pkg/security/authorization/administrator"
	authorizationAdministratorJsonRpc "github.com/iot-my-world/brain/pkg/security/authorization/administrator/jsonRpc"
	"github.com/iot-my-world/brain/pkg/security/claims"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/satori/go.uuid"
	"gopkg.in/square/go-jose.v2"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type client struct {
	url                        string
	jwt                        string
	claims                     claims.Claims
	loggedIn                   bool
	loginRequest               authorizationAdministrator.LoginRequest
	authorizationAdministrator authorizationAdministrator.Administrator
}

// Create New basic json rpc client
func New(
	url string,
) client2.Client {
	newJsonRpcClient := client{
		url: url,
	}
	newJsonRpcClient.authorizationAdministrator = authorizationAdministratorJsonRpc.New(&newJsonRpcClient)

	return &newJsonRpcClient
}

func (c *client) LoggedIn() bool {
	return c.loggedIn
}

func (c *client) Post(request *client2.Request) (*client2.Response, error) {
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
	defer func() {
		if err := postResponse.Body.Close(); err != nil {
			log.Error(err.Error())
		}
	}()
	if err != nil {
		return nil, errors.New("error reading post response body bytes " + err.Error())
	}

	// check for an rpc error
	if strings.Contains(string(postResponseBytes), "rpc: can't find service") {
		return nil, errors.New("rpc error: method not found")
	}

	// unmarshal the body into the response
	response := client2.Response{}
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

	jsonRpcRequest := client2.NewRequest(id.String(), method, [1]interface{}{request})

	jsonRpcResponse, err := c.Post(&jsonRpcRequest)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonRpcResponse.Result, response); err != nil {
		return err
	}

	return nil
}

func (c *client) Login(loginRequest authorizationAdministrator.LoginRequest) error {
	loginResponse, err := c.authorizationAdministrator.Login(&loginRequest)
	if err != nil {
		log.Error(err)
		return err
	}

	// save the login request for maintain/refresh login
	c.loginRequest = loginRequest

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

func (c *client) RefreshLogin() error {
	if err := c.Login(c.loginRequest); err != nil {
		return err
	}
	return nil
}

func (c *client) MaintainLogin() error {
	refreshTokenTimer := time.NewTimer(c.claims.TimeToExpiry() - 10*time.Second)
	for {
		select {
		case <-refreshTokenTimer.C:
			log.Info("refresh json rpc client login")
			if err := c.RefreshLogin(); err != nil {
				return err
			}
			refreshTokenTimer.Reset(c.claims.TimeToExpiry() - 1*time.Second)
		}
	}
}

func (c *client) GetURL() string {
	return c.url
}

func (c *client) SetURL(url string) {
	c.url = url
}
