package main

import (
	"net/http"
	"errors"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"gitlab.com/iotTracker/brain/log"
)

type JsonRpcReq struct{
	// To unmarshal the received json
	Id string `json:"id"`
	Method string `json:"method"`
}

func apiAuthApplier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Retrieve json rpc service method from request body
		jsonRpcServiceMethod, err := getJsonRpcServiceMethod(r)
		if err != nil {
			// if it can't be retrieved, error 404
			http.Error(w, err.Error(), http.StatusNotFound)
			return
			// TODO: Unauthorised access attempts like this should be getting tracked more formally. Could indicate an attack.
		}

		// Check if Authorization Header is present in request indicating that the user is logged in
		if r.Header["Authorization"] == nil {
			// Header not present, this is only allowed for a few jsonRpcServiceMethods
			switch {
			case jsonRpcServiceMethod == "Auth.Login":
				next.ServeHTTP(w, r)
			default:
				log.Info("Unauthorised Json RPC access! - No Authorisation header!")
				// unauthorised api access, error 403
				http.Error(w, "Unauthorised", http.StatusForbidden)
				// TODO: Unauthorised access attempts like this should be getting tracked more formally. Could indicate an attack.
			}
			return
		}

		// Validate the jwt and confirm that user has appropriate claims to access given jsonrpc service
		jwt := r.Header["Authorization"][0]
		if err := mainAPIAuthorizer.AuthorizeAPIReq(jwt, jsonRpcServiceMethod); err != nil {
			log.Warn("Unauthorised Access Attempt for Method." + jsonRpcServiceMethod, err.Error())
			// unauthorised api access, error 403
			http.Error(w, "Unauthorised", http.StatusForbidden)
			// TODO: Unauthorised access attempts like this should be getting tracked more formally. Could indicate an attack.
			return
		}

		next.ServeHTTP(w, r)
		return
	})
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
	req := JsonRpcReq{}
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return "", err
	}
	return req.Method, nil
}

func preFlightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers",
		"Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin, Authorization")

	w.WriteHeader(http.StatusOK)
}

func healthcheck() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	}
}