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

	systemMongoRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler/mongo"
	systemRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/system/recordHandler/adaptor/jsonRpc"

	readingRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/reading/recordHandler/adaptor/jsonRpc"
	readingMongoRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler/mongo"
	tk102DeviceServer "gitlab.com/iotTracker/brain/tracker/device/tk102/server"

	tk102DeviceRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler/adaptor/jsonRpc"
	tk102DeviceMongoRecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler/mongo"
	tk102DeviceBasicAdministrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator/basic"
	tk102DeviceAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator/adaptor/jsonRpc"

	trackingBasicReport "gitlab.com/iotTracker/brain/report/tracking/basic"
	trackingReportJsonRpcAdaptor "gitlab.com/iotTracker/brain/report/tracking/adaptor/jsonRpc"

	"flag"
	"gitlab.com/iotTracker/brain/email/mailer"
	gmailMailer "gitlab.com/iotTracker/brain/email/mailer/gmail"
	partyBasicRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	partyBasicRegistrar "gitlab.com/iotTracker/brain/party/registrar/basic"
	"strings"
)

var ServerPort = "9010"

var mainAPIAuthorizer = apiAuth.APIAuthorizer{}

func main() {
	// get the command line args
	mongoNodes := flag.String("mongoNodes", "localhost:27017", "the nodes in the db cluster")
	mongoUser := flag.String("mongoUser", "", "brains mongo db user")
	mongoPassword := flag.String("mongoPassword", "", "passwords for brains mongo db")
	mailRedirectBaseUrl := flag.String("mailRedirectBaseUrl", "http://localhost:3000", "base url for all email invites")
	rootPasswordFileLocation := flag.String("rootPasswordFileLocation", "", "path to file containing root password")

	flag.Parse()

	// Connect to database
	log.Info("connecting to mongo...")
	databaseName := "brain"
	mongoCluster := strings.Split(*mongoNodes, ",")
	dialInfo := mgo.DialInfo{
		Addrs:     mongoCluster,
		Username:  *mongoUser,
		Password:  *mongoPassword,
		Mechanism: "SCRAM-SHA-1",
		Timeout:   10 * time.Second,
		Source:    "admin",
		Database:  databaseName,
	}
	mainMongoSession, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		log.Error("Could not connect to Mongo cluster: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}
	log.Info("Connected to Mongo!")
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
	RoleRecordHandler := roleMongoRecordHandler.New(mainMongoSession, databaseName, roleCollection)
	UserRecordHandler := userMongoRecordHandler.New(mainMongoSession, databaseName, userCollection)
	PermissionBasicHandler := permissionBasicHandler.New(UserRecordHandler, RoleRecordHandler)
	AuthService := authBasicService.New(UserRecordHandler, rsaPrivateKey)
	CompanyRecordHandler := companyMongoRecordHandler.New(mainMongoSession, databaseName, companyCollection, UserRecordHandler)
	ClientRecordHandler := clientMongoRecordHandler.New(mainMongoSession, databaseName, clientCollection, UserRecordHandler)
	PartyBasicRegistrar := partyBasicRegistrar.New(CompanyRecordHandler, UserRecordHandler, ClientRecordHandler, Mailer, rsaPrivateKey, *mailRedirectBaseUrl)
	SystemRecordHandler := systemMongoRecordHandler.New(mainMongoSession, databaseName, systemCollection, *rootPasswordFileLocation, PartyBasicRegistrar)
	TK102DeviceRecordHandler := tk102DeviceMongoRecordHandler.New(mainMongoSession, databaseName, tk102DeviceCollection, SystemRecordHandler, CompanyRecordHandler, ClientRecordHandler)
	TK102DeviceAdministrator := tk102DeviceBasicAdministrator.New(TK102DeviceRecordHandler, CompanyRecordHandler, ClientRecordHandler)
	ReadingRecordHandler := readingMongoRecordHandler.New(mainMongoSession, databaseName, readingCollection)
	TrackingReport := trackingBasicReport.New(CompanyRecordHandler, ClientRecordHandler, ReadingRecordHandler)

	// Create Service Provider Adaptors
	RoleRecordHandlerAdaptor := roleRecordHandlerJsonRpcAdaptor.New(RoleRecordHandler)
	UserRecordHandlerAdaptor := userRecordHandlerJsonRpcAdaptor.New(UserRecordHandler)
	AuthServiceAdaptor := authServiceJsonRpcAdaptor.New(AuthService)
	PermissionHandlerAdaptor := permissionServiceJsonRpcAdaptor.New(PermissionBasicHandler)
	CompanyRecordHandlerAdaptor := companyRecordHandlerJsonRpcAdaptor.New(CompanyRecordHandler)
	ClientRecordHandlerAdaptor := clientRecordHandlerJsonRpcAdaptor.New(ClientRecordHandler)
	PartyBasicRegistrarAdaptor := partyBasicRegistrarJsonRpcAdaptor.New(PartyBasicRegistrar)
	SystemRecordHandlerAdaptor := systemRecordHandlerJsonRpcAdaptor.New(SystemRecordHandler)
	TK102DeviceRecordHandlerAdaptor := tk102DeviceRecordHandlerJsonRpcAdaptor.New(TK102DeviceRecordHandler)
	TK102DeviceAdministratorAdaptor := tk102DeviceAdministratorJsonRpcAdaptor.New(TK102DeviceAdministrator)
	ReadingRecordHandlerAdaptor := readingRecordHandlerJsonRpcAdaptor.New(ReadingRecordHandler)
	TrackingReportAdaptor := trackingReportJsonRpcAdaptor.New(TrackingReport)

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
	if err := secureAPIServer.RegisterService(SystemRecordHandlerAdaptor, "SystemRecordHandler"); err != nil {
		log.Fatal("Unable to Register System Record Handler Service")
	}
	if err := secureAPIServer.RegisterService(TK102DeviceRecordHandlerAdaptor, "TK102DeviceRecordHandler"); err != nil {
		log.Fatal("Unable to Register TK102 Device Record Handler Service")
	}
	if err := secureAPIServer.RegisterService(TK102DeviceAdministratorAdaptor, "TK102DeviceAdministrator"); err != nil {
		log.Fatal("Unable to Register TK102 Device Administrator")
	}
	if err := secureAPIServer.RegisterService(ReadingRecordHandlerAdaptor, "ReadingRecordHandler"); err != nil {
		log.Fatal("Unable to Register Reading Record Handler Service")
	}
	if err := secureAPIServer.RegisterService(TrackingReportAdaptor, "TrackingReport"); err != nil {
		log.Fatal("Unable to Register Tracking Report Service")
	}

	// Set up Router for secureAPIServer
	secureAPIServerMux := mux.NewRouter()
	secureAPIServerMux.Methods("OPTIONS").HandlerFunc(preFlightHandler)
	secureAPIServerMux.Handle("/api", apiAuthApplier(secureAPIServer)).Methods("POST")
	// Start secureAPIServer
	log.Info("Starting secureAPIServer on port " + ServerPort)
	go func() {
		err := http.ListenAndServe(":"+ServerPort, secureAPIServerMux)
		log.Error("secureAPIServer stopped: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}()

	// Set up tracker tcp server
	tk102DeviceServerInstance := tk102DeviceServer.New(ReadingRecordHandler, TK102DeviceRecordHandler, "0.0.0.0", "7018")
	log.Info("Starting TK102 Device Server")
	go func() {
		err := tk102DeviceServerInstance.Start()
		log.Error("tcp server stopped: ", err)
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
