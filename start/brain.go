package main

import (
	"github.com/gorilla/rpc"
	"gitlab.com/iotTracker/brain/cors"
	gorillaJson "github.com/gorilla/rpc/json"
	"github.com/gorilla/mux"
	"gitlab.com/iotTracker/brain/log"
	"net/http"
	"runtime/debug"
	"os"
	"os/signal"
	"gopkg.in/mgo.v2"
	"time"
	"gitlab.com/iotTracker/brain/security/encrypt"
	"gitlab.com/iotTracker/brain/security/apiAuth"
	"gitlab.com/iotTracker/brain/security/token"

	authBasicService "gitlab.com/iotTracker/brain/security/auth/service/basic"
	authServiceJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	roleMongoRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler/mongo"
	roleRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/role/recordHandler/adaptor/jsonRpc"
	permissionBasicHandler "gitlab.com/iotTracker/brain/security/permission/handler/basic"
	permissionServiceJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/permission/handler/adaptor/jsonRpc"
	userMongoRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler/mongo"
	userRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/user/recordHandler/adaptor/jsonRpc"
)

var ServerPort = "9010"

var mainAPIAuthorizer = apiAuth.APIAuthorizer{}

func main(){

	// Connect to database
	databaseName := "brain"
	dialInfo := mgo.DialInfo{
		Addrs: []string{"localhost:27017"},
		Username: "",
		Password: "",
		Timeout: 10*time.Second,
		Source: "admin",
		Database: databaseName,
	}
	mainMongoSession, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		log.Error("Could not connect to Mongo cluster: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}
	log.Debug("Connected to Mongo!")
	defer mainMongoSession.Close()

	// Get or Generate RSA Key Pair
	rsaPrivateKey := encrypt.FetchPrivateKey("./")

	// Create Service Providers
	RoleRecordHandler := roleMongoRecordHandler.New(mainMongoSession, databaseName, systemRoleCollection)
	UserRecordHandler := userMongoRecordHandler.New(mainMongoSession, databaseName, userCollection)
	PermissionBasicHandler := permissionBasicHandler.New(UserRecordHandler, RoleRecordHandler)
	AuthService := authBasicService.New(UserRecordHandler, rsaPrivateKey)

	// Create Service Provider Adaptors
	RoleRecordHandlerAdaptor := roleRecordHandlerJsonRpcAdaptor.New(RoleRecordHandler)
	UserRecordHandlerAdaptor := userRecordHandlerJsonRpcAdaptor.New(UserRecordHandler)
	AuthServiceAdaptor := authServiceJsonRpcAdaptor.New(AuthService)
	PermissionHandlerAdaptor := permissionServiceJsonRpcAdaptor.New(PermissionBasicHandler)

	// Initialise the APIAuthorizer
	mainAPIAuthorizer.JWTValidator = token.NewJWTValidator(&rsaPrivateKey.PublicKey)
	mainAPIAuthorizer.PermissionHandler = PermissionBasicHandler

	// Create secureAPIServer
	secureAPIServer := rpc.NewServer()
	secureAPIServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	// Register Service Provider Adaptors with secureAPIServer
	if err := secureAPIServer.RegisterService(RoleRecordHandlerAdaptor, "Role"); err != nil {
		log.Fatal("Unable to Register System Role Service")
	}
	if err := secureAPIServer.RegisterService(UserRecordHandlerAdaptor, "User"); err != nil {
		log.Fatal("Unable to Register User Service")
	}
	if err:= secureAPIServer.RegisterService(AuthServiceAdaptor, "Auth"); err != nil {
		log.Fatal("Unable to Register Auth Service Adaptor")
	}
	if err:= secureAPIServer.RegisterService(PermissionHandlerAdaptor, "Permission"); err != nil {
		log.Fatal("Unable to Register Permission Handler Service Adaptor")
	}

	// Set up Router for secureAPIServer
	secureAPIServerMux := mux.NewRouter()
	secureAPIServerMux.Methods("OPTIONS").HandlerFunc(preFlightHandler)
	secureAPIServerMux.Handle("/api", apiAuthApplier(secureAPIServer)).Methods("POST")
	// Start secureAPIServer
	log.Info("Starting secureAPIServer on port " + ServerPort)
	go func() {
		err := http.ListenAndServe("0.0.0.0:" + ServerPort, secureAPIServerMux)
		log.Error("secureAPIServer stopped: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}()

	//Wait for interrupt signal
	systemSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(systemSignalsChannel)
	for {
		select {
		case s := <- systemSignalsChannel:
			log.Info("Application is shutting down.. ( ", s, " )")
			return
		}
	}
}
