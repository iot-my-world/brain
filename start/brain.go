package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	gorillaJson "github.com/gorilla/rpc/json"
	"gitlab.com/iotTracker/brain/cors"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/security/apiAuth"
	"gitlab.com/iotTracker/brain/security/encrypt"
	"gitlab.com/iotTracker/brain/security/token"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	authServiceJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	authBasicService "gitlab.com/iotTracker/brain/security/auth/service/basic"

	permissionServiceJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/permission/handler/adaptor/jsonRpc"
	permissionBasicHandler "gitlab.com/iotTracker/brain/security/permission/handler/basic"

	roleRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/role/recordHandler/adaptor/jsonRpc"
	roleMongoRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler/mongo"

	userRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/user/recordHandler/adaptor/jsonRpc"
	userMongoRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler/mongo"

	companyRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/recordHandler/adaptor/jsonRpc"
	companyMongoRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler/mongo"

	clientRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/client/recordHandler/adaptor/jsonRpc"
	clientMongoRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler/mongo"

	partyBasicRegistrar "gitlab.com/iotTracker/brain/party/registrar/basic"
	partyBasicRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/email/mailer"
	gmailMailer "gitlab.com/iotTracker/brain/email/mailer/gmail"
)

var ServerPort = "9010"

var mainAPIAuthorizer = apiAuth.APIAuthorizer{}

func main() {

	// Connect to database
	databaseName := "brain"
	dialInfo := mgo.DialInfo{
		Addrs:    []string{"localhost:27017"},
		Username: "",
		Password: "",
		Timeout:  10 * time.Second,
		Source:   "admin",
		Database: databaseName,
	}
	mainMongoSession, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		log.Error("Could not connect to Mongo cluster: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}
	log.Debug("Connected to Mongo!")
	defer mainMongoSession.Close()

	// spotNav123
	// spotnavza@gmail.com

	// Get or Generate RSA Key Pair
	rsaPrivateKey := encrypt.FetchPrivateKey("./")

	// Create Mailer
	Mailer := gmailMailer.New(mailer.AuthInfo{
		Identity: "",
		Username: "spotnavza@gmail.com",
		Password: "spotNav123",
		Host:     "smtp.gmail.com",
	})

	// Create Service Providers
	RoleRecordHandler := roleMongoRecordHandler.New(mainMongoSession, databaseName, systemRoleCollection)
	UserRecordHandler := userMongoRecordHandler.New(mainMongoSession, databaseName, userCollection)
	PermissionBasicHandler := permissionBasicHandler.New(UserRecordHandler, RoleRecordHandler)
	AuthService := authBasicService.New(UserRecordHandler, rsaPrivateKey)
	CompanyRecordHandler := companyMongoRecordHandler.New(mainMongoSession, databaseName, companyCollection)
	ClientRecordHandler := clientMongoRecordHandler.New(mainMongoSession, databaseName, clientCollection)
	PartyBasicRegistrar := partyBasicRegistrar.New(CompanyRecordHandler, UserRecordHandler, ClientRecordHandler, Mailer, rsaPrivateKey)

	// Create Service Provider Adaptors
	RoleRecordHandlerAdaptor := roleRecordHandlerJsonRpcAdaptor.New(RoleRecordHandler)
	UserRecordHandlerAdaptor := userRecordHandlerJsonRpcAdaptor.New(UserRecordHandler)
	AuthServiceAdaptor := authServiceJsonRpcAdaptor.New(AuthService)
	PermissionHandlerAdaptor := permissionServiceJsonRpcAdaptor.New(PermissionBasicHandler)
	CompanyRecordHandlerAdaptor := companyRecordHandlerJsonRpcAdaptor.New(CompanyRecordHandler)
	ClientRecordHandlerAdaptor := clientRecordHandlerJsonRpcAdaptor.New(ClientRecordHandler)
	PartyBasicRegistrarAdaptor := partyBasicRegistrarJsonRpcAdaptor.New(PartyBasicRegistrar)

	// Initialise the APIAuthorizer
	mainAPIAuthorizer.JWTValidator = token.NewJWTValidator(&rsaPrivateKey.PublicKey)
	mainAPIAuthorizer.PermissionHandler = PermissionBasicHandler

	// Create secureAPIServer
	secureAPIServer := rpc.NewServer()
	secureAPIServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	// Register Service Provider Adaptors with secureAPIServer
	if err := secureAPIServer.RegisterService(RoleRecordHandlerAdaptor, "RoleRecordHandler"); err != nil {
		log.Fatal("Unable to Register Role Record Handler Service")
	}
	if err := secureAPIServer.RegisterService(UserRecordHandlerAdaptor, "UserRecordHandler"); err != nil {
		log.Fatal("Unable to Register User Record Handler Service")
	}
	if err := secureAPIServer.RegisterService(AuthServiceAdaptor, "Auth"); err != nil {
		log.Fatal("Unable to Register Auth Service Adaptor")
	}
	if err := secureAPIServer.RegisterService(PermissionHandlerAdaptor, "PermissionHandler"); err != nil {
		log.Fatal("Unable to Register Permission Handler Service Adaptor")
	}
	if err := secureAPIServer.RegisterService(CompanyRecordHandlerAdaptor, "CompanyRecordHandler"); err != nil {
		log.Fatal("Unable to Register Company Record Handler Service")
	}
	if err := secureAPIServer.RegisterService(ClientRecordHandlerAdaptor, "ClientRecordHandler"); err != nil {
		log.Fatal("Unable to Register Client Record Handler Service")
	}
	if err := secureAPIServer.RegisterService(PartyBasicRegistrarAdaptor, "PartyRegistrar"); err != nil {
		log.Fatal("Unable to Register Party Registrar Service")
	}

	// Set up Router for secureAPIServer
	secureAPIServerMux := mux.NewRouter()
	secureAPIServerMux.Methods("OPTIONS").HandlerFunc(preFlightHandler)
	secureAPIServerMux.Handle("/api", apiAuthApplier(secureAPIServer)).Methods("POST")
	// Start secureAPIServer
	log.Info("Starting secureAPIServer on port " + ServerPort)
	go func() {
		err := http.ListenAndServe("0.0.0.0:"+ServerPort, secureAPIServerMux)
		log.Error("secureAPIServer stopped: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}()

	//Wait for interrupt signal
	systemSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(systemSignalsChannel)
	for {
		select {
		case s := <-systemSignalsChannel:
			log.Info("Application is shutting down.. ( ", s, " )")
			return
		}
	}
}
