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

	permissionAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/permission/administrator/adaptor/jsonRpc"
	permissionBasicAdministrator "gitlab.com/iotTracker/brain/security/permission/administrator/basic"

	roleRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/role/recordHandler/adaptor/jsonRpc"
	roleMongoRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler/mongo"

	userAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/user/administrator/adaptor/jsonRpc"
	userBasicAdministrator "gitlab.com/iotTracker/brain/party/user/administrator/basic"
	userRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/user/recordHandler/adaptor/jsonRpc"
	userMongoRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler/mongo"

	companyRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/recordHandler/adaptor/jsonRpc"
	companyMongoRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler/mongo"

	clientRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/client/recordHandler/adaptor/jsonRpc"
	clientMongoRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler/mongo"

	systemRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/system/recordHandler/adaptor/jsonRpc"
	systemMongoRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler/mongo"

	tk102DeviceServer "gitlab.com/iotTracker/brain/tracker/device/tk102/server"
	readingRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/reading/recordHandler/adaptor/jsonRpc"
	readingMongoRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler/mongo"

	tk102DeviceAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator/adaptor/jsonRpc"
	tk102DeviceBasicAdministrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator/basic"
	tk102DeviceRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler/adaptor/jsonRpc"
	tk102DeviceMongoRecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler/mongo"

	trackingReportJsonRpcAdaptor "gitlab.com/iotTracker/brain/report/tracking/adaptor/jsonRpc"
	trackingBasicReport "gitlab.com/iotTracker/brain/report/tracking/basic"

	"flag"
	"gitlab.com/iotTracker/brain/email/mailer"
	gmailMailer "gitlab.com/iotTracker/brain/email/mailer/gmail"
	partyBasicRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	partyBasicRegistrar "gitlab.com/iotTracker/brain/party/registrar/basic"

	partyAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/administrator/adaptor/jsonRpc"
	partyBasicAdministrator "gitlab.com/iotTracker/brain/party/administrator/basic"

	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"strings"
)

var serverPort = "9010"

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

	// Create system claims for the services that root privileges
	var systemClaims = login.Login{
		//UserId          id.Identifier `json:"userId"`
		//IssueTime       int64         `json:"issueTime"`
		//ExpirationTime  int64         `json:"expirationTime"`
		//ParentPartyType party.Type    `json:"parentPartyType"`
		//ParentId        id.Identifier `json:"parentId"`
		PartyType: party.System,
		//PartyId         id.Identifier `json:"partyId"`
	}

	// Create Service Providers
	RoleRecordHandler := roleMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		roleCollection,
	)
	// User
	UserRecordHandler := userMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		userCollection,
		&systemClaims,
	)
	UserBasicAdministrator := userBasicAdministrator.New(
		UserRecordHandler,
	)

	// Permission
	PermissionBasicHandler := permissionBasicAdministrator.New(
		UserRecordHandler,
		RoleRecordHandler,
	)

	// Auth
	AuthService := authBasicService.New(
		UserRecordHandler,
		rsaPrivateKey,
		&systemClaims,
	)

	// Company
	CompanyRecordHandler := companyMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		companyCollection,
		UserRecordHandler,
	)

	// Client
	ClientRecordHandler := clientMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		clientCollection,
		UserRecordHandler,
	)

	// Party
	PartyBasicRegistrar := partyBasicRegistrar.New(
		CompanyRecordHandler,
		UserRecordHandler,
		ClientRecordHandler,
		Mailer,
		rsaPrivateKey,
		*mailRedirectBaseUrl,
		&systemClaims,
	)

	// System
	SystemRecordHandler := systemMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		systemCollection,
		*rootPasswordFileLocation,
		PartyBasicRegistrar,
		&systemClaims,
	)

	// TK102 Device
	TK102DeviceRecordHandler := tk102DeviceMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		tk102DeviceCollection,
		SystemRecordHandler,
		CompanyRecordHandler,
		ClientRecordHandler,
	)

	// Reading
	ReadingRecordHandler := readingMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		readingCollection,
	)

	// Report
	TrackingReport := trackingBasicReport.New(
		SystemRecordHandler,
		CompanyRecordHandler,
		ClientRecordHandler,
		ReadingRecordHandler,
		TK102DeviceRecordHandler,
	)

	// Party
	PartyBasicHandler := partyBasicAdministrator.New(
		ClientRecordHandler,
		CompanyRecordHandler,
		SystemRecordHandler,
	)

	TK102DeviceAdministrator := tk102DeviceBasicAdministrator.New(
		TK102DeviceRecordHandler,
		CompanyRecordHandler,
		ClientRecordHandler,
		PartyBasicHandler,
		ReadingRecordHandler,
	)

	// Create Service Provider Adaptors
	RoleRecordHandlerAdaptor := roleRecordHandlerJsonRpcAdaptor.New(RoleRecordHandler)
	UserRecordHandlerAdaptor := userRecordHandlerJsonRpcAdaptor.New(UserRecordHandler)
	UserAdministratorAdaptor := userAdministratorJsonRpcAdaptor.New(UserBasicAdministrator)
	AuthServiceAdaptor := authServiceJsonRpcAdaptor.New(AuthService)
	PermissionHandlerAdaptor := permissionAdministratorJsonRpcAdaptor.New(PermissionBasicHandler)
	CompanyRecordHandlerAdaptor := companyRecordHandlerJsonRpcAdaptor.New(CompanyRecordHandler)
	ClientRecordHandlerAdaptor := clientRecordHandlerJsonRpcAdaptor.New(ClientRecordHandler)
	PartyBasicRegistrarAdaptor := partyBasicRegistrarJsonRpcAdaptor.New(PartyBasicRegistrar)
	SystemRecordHandlerAdaptor := systemRecordHandlerJsonRpcAdaptor.New(SystemRecordHandler)
	TK102DeviceRecordHandlerAdaptor := tk102DeviceRecordHandlerJsonRpcAdaptor.New(TK102DeviceRecordHandler)
	TK102DeviceAdministratorAdaptor := tk102DeviceAdministratorJsonRpcAdaptor.New(TK102DeviceAdministrator)
	ReadingRecordHandlerAdaptor := readingRecordHandlerJsonRpcAdaptor.New(ReadingRecordHandler)
	TrackingReportAdaptor := trackingReportJsonRpcAdaptor.New(TrackingReport)
	PartyHandlerAdaptor := partyAdministratorJsonRpcAdaptor.New(PartyBasicHandler)

	// Initialise the APIAuthorizer
	mainAPIAuthorizer.JWTValidator = token.NewJWTValidator(&rsaPrivateKey.PublicKey)
	mainAPIAuthorizer.PermissionHandler = PermissionBasicHandler

	// Create secureAPIServer
	secureAPIServer := rpc.NewServer()
	secureAPIServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	// Register Service Provider Adaptors with secureAPIServer
	// Role
	if err := secureAPIServer.RegisterService(RoleRecordHandlerAdaptor, "RoleRecordHandler"); err != nil {
		log.Fatal("Unable to Register Role Record Handler Service")
	}

	// User
	if err := secureAPIServer.RegisterService(UserRecordHandlerAdaptor, "UserRecordHandler"); err != nil {
		log.Fatal("Unable to Register User Record Handler Service")
	}
	if err := secureAPIServer.RegisterService(UserAdministratorAdaptor, "UserAdministrator"); err != nil {
		log.Fatal("Unable to Register User Record Handler Service")
	}

	// Auth
	if err := secureAPIServer.RegisterService(AuthServiceAdaptor, "Auth"); err != nil {
		log.Fatal("Unable to Register Auth Service Adaptor")
	}

	// Permission
	if err := secureAPIServer.RegisterService(PermissionHandlerAdaptor, "PermissionHandler"); err != nil {
		log.Fatal("Unable to Register Permission Handler Service Adaptor")
	}

	// Company
	if err := secureAPIServer.RegisterService(CompanyRecordHandlerAdaptor, "CompanyRecordHandler"); err != nil {
		log.Fatal("Unable to Register Company Record Handler Service")
	}

	// Client
	if err := secureAPIServer.RegisterService(ClientRecordHandlerAdaptor, "ClientRecordHandler"); err != nil {
		log.Fatal("Unable to Register Client Record Handler Service")
	}

	// Party
	if err := secureAPIServer.RegisterService(PartyBasicRegistrarAdaptor, "PartyRegistrar"); err != nil {
		log.Fatal("Unable to Register Party Registrar Service")
	}

	// System
	if err := secureAPIServer.RegisterService(SystemRecordHandlerAdaptor, "SystemRecordHandler"); err != nil {
		log.Fatal("Unable to Register System Record Handler Service")
	}

	// TK102 Device
	if err := secureAPIServer.RegisterService(TK102DeviceRecordHandlerAdaptor, "TK102DeviceRecordHandler"); err != nil {
		log.Fatal("Unable to Register TK102 Device Record Handler Service")
	}
	if err := secureAPIServer.RegisterService(TK102DeviceAdministratorAdaptor, "TK102DeviceAdministrator"); err != nil {
		log.Fatal("Unable to Register TK102 Device Administrator")
	}

	// Reading
	if err := secureAPIServer.RegisterService(ReadingRecordHandlerAdaptor, "ReadingRecordHandler"); err != nil {
		log.Fatal("Unable to Register Reading Record Handler Service")
	}

	// Reports
	if err := secureAPIServer.RegisterService(TrackingReportAdaptor, "TrackingReport"); err != nil {
		log.Fatal("Unable to Register Tracking Report Service")
	}

	// Party Handler
	if err := secureAPIServer.RegisterService(PartyHandlerAdaptor, "PartyHandler"); err != nil {
		log.Fatal("Unable to Register Party Handler Service")
	}

	// Set up Router for secureAPIServer
	secureAPIServerMux := mux.NewRouter()
	secureAPIServerMux.Methods("OPTIONS").HandlerFunc(preFlightHandler)
	secureAPIServerMux.Handle("/api", apiAuthApplier(secureAPIServer)).Methods("POST")
	// Start secureAPIServer
	log.Info("Starting secureAPIServer on port " + serverPort)
	go func() {
		err := http.ListenAndServe(":"+serverPort, secureAPIServerMux)
		log.Error("secureAPIServer stopped: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}()

	// Set up tracker tcp server
	tk102DeviceServerInstance := tk102DeviceServer.New(ReadingRecordHandler, &systemClaims, TK102DeviceRecordHandler, "0.0.0.0", "7018")
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
