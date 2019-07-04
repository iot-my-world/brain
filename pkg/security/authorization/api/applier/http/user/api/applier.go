package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/iot-my-world/brain/internal/log"
	http2 "github.com/iot-my-world/brain/pkg/security/authorization/api/applier/http"
	"github.com/iot-my-world/brain/pkg/security/authorization/api/authorizer"
	"io/ioutil"
	"net/http"
)

type applier struct {
	apiAuthorizer authorizer.Authorizer
}

func New(
	apiAuthorizer authorizer.Authorizer,
) http2.Applier {
	return &applier{
		apiAuthorizer: apiAuthorizer,
	}
}

func (a *applier) ApplyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Retrieve json rpc service method from request body
		jsonRpcServiceMethod, err := getJsonRpcServiceMethod(r)
		if err != nil {
			// if it can't be retrieved, error 404
			http.Error(w, err.Error(), http.StatusNotFound)
			return
			// TODO: Unauthorised access attempts like this should be getting tracked more formally. Could indicate an attack.
		}

		// UserAdministrator.ForgotPassword
		// Check if Authorization Header is present in request indicating that the api user is logged in
		switch jsonRpcServiceMethod {
		case "Auth.Login":
			next.ServeHTTP(w, r)
			return
		default:
			if r.Header["Authorization"] == nil {
				log.Info("Unauthorised Json RPC access! - No Authorisation header!")
				// unauthorised api access, error 403
				http.Error(w, "Unauthorised", http.StatusForbidden)
				return
			}
		}

		// Validate the jwt and confirm that user has appropriate claims to access given jsonrpc service
		jwt := r.Header["Authorization"][0]
		if wrappedClaims, err := a.apiAuthorizer.AuthorizeAPIReq(jwt, jsonRpcServiceMethod); err == nil {
			ctx := context.WithValue(r.Context(), "wrappedClaims", wrappedClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		} else {
			log.Warn("Unauthorised Access Attempt", err.Error())
			// unauthorised api access, error 403
			http.Error(w, "Unauthorised", http.StatusForbidden)
			// TODO: Unauthorised access attempts like this should be getting tracked more formally. Could indicate an attack.
			return
		}

		// if it can't be retrieved, error 404
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	})
}

func (a *applier) PreFlightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers",
		"Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin, Authorization")

	w.WriteHeader(http.StatusOK)
}

type JSONRPCReq struct {
	// To unmarshal the received json
	Id     string `json:"id"`
	Method string `json:"method"`
}

func getJsonRpcServiceMethod(r *http.Request) (string, error) {
	// Confirm that body of request has data
	if r.Body == nil {
		return "", errors.New("body is nil")
	}

	// Extract body of http Request
	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(r.Body)

	// Reset body of request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Retrieve id and method of json rpc request
	req := JSONRPCReq{}
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return "", err
	}
	return req.Method, nil
}
