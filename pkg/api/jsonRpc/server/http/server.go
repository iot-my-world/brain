package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	gorillaJson "github.com/gorilla/rpc/json"
	"github.com/iot-my-world/brain/internal/cors"
	"github.com/iot-my-world/brain/internal/log"
	server2 "github.com/iot-my-world/brain/pkg/api/jsonRpc/server"
	jsonRpcServerAuthoriser "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authoriser"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	"io/ioutil"
	netHttp "net/http"
	"runtime/debug"
	"strings"
)

type server struct {
	path             string
	host             string
	port             string
	rpcServer        *rpc.Server
	authoriser       jsonRpcServerAuthoriser.Authoriser
	serverMux        *mux.Router
	serviceProviders map[jsonRpcServiceProvider.Name]jsonRpcServiceProvider.Provider
}

func New(
	path string,
	host string,
	port string,
	authoriser jsonRpcServerAuthoriser.Authoriser,
) server2.Server {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")
	return &server{
		path:             path,
		host:             host,
		port:             port,
		serverMux:        mux.NewRouter(),
		rpcServer:        rpcServer,
		authoriser:       authoriser,
		serviceProviders: make(map[jsonRpcServiceProvider.Name]jsonRpcServiceProvider.Provider),
	}
}

func (s *server) Start() error {
	s.serverMux.Methods("OPTIONS").HandlerFunc(preFlightHandler)
	s.serverMux.Handle(
		s.path,
		s.rpcServer,
	).Methods("POST")
	if err := netHttp.ListenAndServe(s.host+":"+s.port, s.serverMux); err != nil {
		log.Error("json rpc api server stopped: ", err, "\n", string(debug.Stack()))
	}
	return nil
}

func (s *server) SecureStart() error {
	s.serverMux.Methods("OPTIONS").HandlerFunc(securePreFlightHandler)
	s.serverMux.Handle(
		s.path,
		s.applyAuthorization(s.rpcServer),
	).Methods("POST")
	if err := netHttp.ListenAndServe(s.host+":"+s.port, s.serverMux); err != nil {
		log.Error("json rpc api server stopped: ", err, "\n", string(debug.Stack()))
	}
	return nil
}

func (s *server) RegisterServiceProvider(serviceProvider jsonRpcServiceProvider.Provider) error {
	s.serviceProviders[serviceProvider.Name()] = serviceProvider
	if err := s.rpcServer.RegisterService(serviceProvider, string(serviceProvider.Name())); err != nil {
		err = errors.New(fmt.Sprintf(
			"error registering service %s with json rpc http server: %s",
			string(serviceProvider.Name()),
			err.Error(),
		))
		log.Error(err.Error())
		return err
	}
	return nil
}

func (s *server) RegisterBatchServiceProviders(serviceProviders []jsonRpcServiceProvider.Provider) error {
	for _, serviceProvider := range serviceProviders {
		if err := s.RegisterServiceProvider(serviceProvider); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) applyAuthorization(next netHttp.Handler) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		// Retrieve json rpc service method from request body
		serviceProvider, jsonRpcServiceMethod, err := s.getServiceProvider(r)
		if err != nil {
			// if it can't be retrieved, error 404
			netHttp.Error(w, err.Error(), netHttp.StatusNotFound)
			return
		}

		// check if method requires authorization
		if !serviceProvider.MethodRequiresAuthorization(jsonRpcServiceMethod) {
			next.ServeHTTP(w, r)
			return
		}

		// check if an authorization header was provided
		if r.Header["Authorization"] == nil {
			log.Info("Unauthorised Json RPC access! - No Authorisation header!")
			// unauthorised api access, error 403
			netHttp.Error(w, "Unauthorised", netHttp.StatusForbidden)
			return
		}

		// authorize access to the service
		jwt := r.Header["Authorization"][0]
		if wrappedClaims, err := s.authoriser.AuthoriseServiceMethod(jwt, jsonRpcServiceMethod); err == nil {
			ctx := context.WithValue(r.Context(), "wrappedClaims", wrappedClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		} else {
			log.Warn("Unauthorised Access Attempt", err.Error())
			// unauthorised api access, error 403
			netHttp.Error(w, "Unauthorised", netHttp.StatusForbidden)
			return
		}
	})
}

func (s *server) getServiceProvider(r *netHttp.Request) (jsonRpcServiceProvider.Provider, string, error) {
	// Confirm that body of request has data
	if r.Body == nil {
		return nil, "", errors.New("body is nil")
	}

	// Extract body of http Request
	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(r.Body)

	// Reset body of request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Retrieve id and method of json rpc request
	var req struct {
		// To unmarshal the received json
		Id     string `json:"id"`
		Method string `json:"method"`
	}
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return nil, "", err
	}
	providerAndMethod := strings.Split(req.Method, ".")
	if len(providerAndMethod) != 2 {
		return nil, "", errors.New("invalid rpc method provider string")
	}

	provider, found := s.serviceProviders[jsonRpcServiceProvider.Name(providerAndMethod[0])]
	if !found {
		return nil, "", errors.New("no registered service provider")
	}

	return provider, req.Method, nil
}

func securePreFlightHandler(w netHttp.ResponseWriter, r *netHttp.Request) {
	w.Header().Set(
		"Access-Control-Allow-Origin",
		"*",
	)
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin, Authorization",
	)
	w.WriteHeader(netHttp.StatusOK)
}

func preFlightHandler(w netHttp.ResponseWriter, r *netHttp.Request) {
	w.Header().Set(
		"Access-Control-Allow-Origin",
		"*",
	)
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin",
	)
	w.WriteHeader(netHttp.StatusOK)
}
