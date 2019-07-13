package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	gorillaJson "github.com/gorilla/rpc/json"
	"github.com/iot-my-world/brain/internal/config"
	"github.com/iot-my-world/brain/internal/cors"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/internal/security/encrypt"
	"github.com/iot-my-world/brain/pkg/security/token"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	databaseCollection "github.com/iot-my-world/brain/pkg/database/collection"

	authorizationAdministrator "github.com/iot-my-world/brain/pkg/security/authorization/administrator"

	authServiceJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/security/authorization/administrator/adaptor/jsonRpc"
	apiUserAuthorizationAdministrator "github.com/iot-my-world/brain/pkg/security/authorization/administrator/user/api"
	humanUserAuthorizationAdministrator "github.com/iot-my-world/brain/pkg/security/authorization/administrator/user/human"

	humanUserHttpAPIAuthApplier "github.com/iot-my-world/brain/pkg/security/authorization/api/applier/http/user/human"
	humanUserAPIAuthorizer "github.com/iot-my-world/brain/pkg/security/authorization/api/authorizer/user/human"

	permissionAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/security/permission/administrator/adaptor/jsonRpc"
	permissionBasicAdministrator "github.com/iot-my-world/brain/pkg/security/permission/administrator/basic"

	roleMongoRecordHandler "github.com/iot-my-world/brain/pkg/security/role/recordHandler/mongo"

	humanUserAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/user/human/administrator/adaptor/jsonRpc"
	humanUserBasicAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator/basic"
	humanUserAuthoriser "github.com/iot-my-world/brain/pkg/user/human/authoriser"
	humanUserRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/user/human/recordHandler/adaptor/jsonRpc"
	humanUserMongoRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler/mongo"
	humanUserValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/user/human/validator/adaptor/jsonRpc"
	humanUserBasicValidator "github.com/iot-my-world/brain/pkg/user/human/validator/basic"

	companyAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/company/administrator/adaptor/jsonRpc"
	companyBasicAdministrator "github.com/iot-my-world/brain/pkg/party/company/administrator/basic"
	companyRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/company/recordHandler/adaptor/jsonRpc"
	companyMongoRecordHandler "github.com/iot-my-world/brain/pkg/party/company/recordHandler/mongo"
	companyValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/company/validator/adaptor/jsonRpc"
	companyBasicValidator "github.com/iot-my-world/brain/pkg/party/company/validator/basic"

	clientAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/client/administrator/adaptor/jsonRpc"
	clientBasicAdministrator "github.com/iot-my-world/brain/pkg/party/client/administrator/basic"
	clientRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/client/recordHandler/adaptor/jsonRpc"
	clientMongoRecordHandler "github.com/iot-my-world/brain/pkg/party/client/recordHandler/mongo"
	clientValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/client/validator/adaptor/jsonRpc"
	clientBasicValidator "github.com/iot-my-world/brain/pkg/party/client/validator/basic"

	systemRecordHandler "github.com/iot-my-world/brain/pkg/party/system/recordHandler"
	systemRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/system/recordHandler/adaptor/jsonRpc"
	systemMongoRecordHandler "github.com/iot-my-world/brain/pkg/party/system/recordHandler/mongo"

	apiUserAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/user/api/administrator/adaptor/jsonRpc"
	apiUserBasicAdministrator "github.com/iot-my-world/brain/pkg/user/api/administrator/basic"
	apiUserBasicPasswordGenerator "github.com/iot-my-world/brain/pkg/user/api/password/generator/basic"
	apiUserRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/user/api/recordHandler/adaptor/jsonRpc"
	apiUserMongoRecordHandler "github.com/iot-my-world/brain/pkg/user/api/recordHandler/mongo"
	apiUserValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/user/api/validator/adaptor/jsonRpc"
	apiUserBasicValidator "github.com/iot-my-world/brain/pkg/user/api/validator/basic"

	trackingReport "github.com/iot-my-world/brain/pkg/report/tracking"
	trackingReportJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/report/tracking/adaptor/jsonRpc"
	trackingBasicReport "github.com/iot-my-world/brain/pkg/report/tracking/basic"

	"flag"
	"github.com/iot-my-world/brain/pkg/communication/email/mailer"
	gmailMailer "github.com/iot-my-world/brain/pkg/communication/email/mailer/gmail"
	partyRegistrar "github.com/iot-my-world/brain/pkg/party/registrar"
	partyBasicRegistrarJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/registrar/adaptor/jsonRpc"
	partyBasicRegistrar "github.com/iot-my-world/brain/pkg/party/registrar/basic"

	registrationEmailGenerator "github.com/iot-my-world/brain/pkg/communication/email/generator/registration"
	setPasswordEmailGenerator "github.com/iot-my-world/brain/pkg/communication/email/generator/set/password"

	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	partyAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/administrator/adaptor/jsonRpc"
	partyBasicAdministrator "github.com/iot-my-world/brain/pkg/party/administrator/basic"

	"fmt"
	"github.com/iot-my-world/brain/pkg/party"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	"path/filepath"
	"strings"

	sigbugAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator"
	sigbugAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/administrator/adaptor/jsonRpc"
	sigbugBasicAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator/basic"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	sigbugRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/adaptor/jsonRpc"
	sigbugMongoRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/mongo"
	sigbugValidator "github.com/iot-my-world/brain/pkg/device/sigbug/validator"
	sigbugValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/validator/adaptor/jsonRpc"
	sigbugBasicValidator "github.com/iot-my-world/brain/pkg/device/sigbug/validator/basic"

	jsonRpcHttpServer "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"

	sigfoxBackendAuthoriser "github.com/iot-my-world/brain/pkg/sigfox/backend/authoriser"
	sigfoxBasicBackendCallbackServerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server/adaptor/jsonRpc"
	sigfoxBasicBackendCallbackServer "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server/basic"
)

var humanUserAPIServerPort = "9010"

func main() {
	pathToConfigFile := flag.String("pathToConfigFile", "configs/config.toml", "brain configuration file")
	//kafkaBrokers := flag.String("kafkaBrokers", "localhost:9092", "ipAddress:port of each kafka broker node (, separated)")
	flag.Parse()

	brainConfig := config.New(*pathToConfigFile)
	log.Info("environment: ", brainConfig.Environment)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("fail to get dir:" + err.Error())
	}
	log.Info("working directory: " + dir)

	// Connect to database
	databaseName := "brain"
	log.Info(fmt.Sprintf("connecting to mongo @ node addresses: [%s]", strings.Join(brainConfig.MongoNodes, ", ")))
	dialInfo := mgo.DialInfo{
		Addrs:     brainConfig.MongoNodes,
		Username:  brainConfig.MongoUser,
		Password:  brainConfig.MongoPassword,
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

	// Get or Generate RSA Key Pair
	rsaPrivateKey := encrypt.FetchPrivateKey(brainConfig.KeyFilePath)

	// Create Mailer
	Mailer := gmailMailer.New(mailer.AuthInfo{
		Identity: "",
		Username: brainConfig.EmailAddress,
		Password: brainConfig.EmailPassword,
		Host:     brainConfig.EmailHost,
	})

	// email generators
	RegistrationEmailGenerator := registrationEmailGenerator.New(
		brainConfig.PathToEmailTemplateFolder,
	)
	SetPasswordEmailGenerator := setPasswordEmailGenerator.New(
		brainConfig.PathToEmailTemplateFolder,
	)

	// create and start nerveBroadcast producer
	//kafkaBrokerNodes := strings.Split(*kafkaBrokers, ",")
	//nerveBroadcastProducer := asyncMessagingProducer.New(
	//	kafkaBrokerNodes,
	//	"nerveBroadcast",
	//)
	//if err := nerveBroadcastProducer.Start(); err != nil {
	//	log.Fatal(err.Error())
	//}

	// Create system claims for the services that root privileges
	var systemClaims = humanUserLoginClaims.Login{
		//UserId          id.Identifier
		//IssueTime       int64
		//ExpirationTime  int64
		//ParentPartyType party.Type
		//ParentId        id.Identifier
		PartyType: party.System,
		//PartyId         id.Identifier `json:"partyId"`
	}

	// ________________________________ Create Service Providers ________________________________
	//

	RoleRecordHandler := roleMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.Role,
	)
	// User
	UserRecordHandler := humanUserMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.User,
	)
	CompanyRecordHandler := companyMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.Company,
	)
	ClientRecordHandler := clientMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.Client,
	)
	UserValidator := humanUserBasicValidator.New(
		UserRecordHandler,
		CompanyRecordHandler,
		ClientRecordHandler,
		&systemClaims,
	)
	UserBasicAdministrator := humanUserBasicAdministrator.New(
		UserRecordHandler,
		UserValidator,
		Mailer,
		rsaPrivateKey,
		brainConfig.MailRedirectBaseUrl,
		&systemClaims,
		SetPasswordEmailGenerator,
		brainConfig.Environment,
	)

	// Auth
	HumanUserAuthorizationService := humanUserAuthorizationAdministrator.New(
		UserRecordHandler,
		rsaPrivateKey,
		&systemClaims,
	)

	// Company
	CompanyValidator := companyBasicValidator.New(
		CompanyRecordHandler,
		UserRecordHandler,
		&systemClaims,
	)
	CompanyAdministrator := companyBasicAdministrator.New(
		CompanyRecordHandler,
		CompanyValidator,
		UserRecordHandler,
		&systemClaims,
	)

	// Client
	ClientValidator := clientBasicValidator.New(
		ClientRecordHandler,
		UserRecordHandler,
		&systemClaims,
	)
	ClientBasicAdministrator := clientBasicAdministrator.New(
		ClientRecordHandler,
		ClientValidator,
		UserRecordHandler,
		&systemClaims,
	)

	// Party
	PartyBasicRegistrar := partyBasicRegistrar.New(
		CompanyRecordHandler,
		UserRecordHandler,
		UserValidator,
		UserBasicAdministrator,
		ClientRecordHandler,
		Mailer,
		rsaPrivateKey,
		brainConfig.MailRedirectBaseUrl,
		&systemClaims,
		RegistrationEmailGenerator,
		brainConfig.Environment,
	)

	// System
	SystemRecordHandler := systemMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.System,
		brainConfig.RootPasswordFileLocation,
		PartyBasicRegistrar,
		&systemClaims,
	)

	// Party
	PartyBasicAdministrator := partyBasicAdministrator.New(
		ClientRecordHandler,
		CompanyRecordHandler,
		SystemRecordHandler,
		&systemClaims,
		CompanyAdministrator,
		ClientBasicAdministrator,
		PartyBasicRegistrar,
	)

	// API User
	APIUserRecordHandler := apiUserMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.APIUser,
	)
	APIUserValidator := apiUserBasicValidator.New(
		PartyBasicAdministrator,
	)
	APIUserPasswordGenerator := apiUserBasicPasswordGenerator.New()
	APIUserAdministrator := apiUserBasicAdministrator.New(
		APIUserValidator,
		APIUserRecordHandler,
		APIUserPasswordGenerator,
	)

	APIUserAuthorizationService := apiUserAuthorizationAdministrator.New(
		APIUserRecordHandler,
		rsaPrivateKey,
		&systemClaims,
	)

	// Permission
	PermissionBasicHandler := permissionBasicAdministrator.New(
		UserRecordHandler,
		RoleRecordHandler,
		APIUserRecordHandler,
	)

	// Sigbug Device
	SigbugRecordHandler := sigbugMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.Sigbug,
	)
	SigbugValidator := sigbugBasicValidator.New(
		PartyBasicAdministrator,
	)
	SigbugAdministrator := sigbugBasicAdministrator.New(
		SigbugValidator,
		SigbugRecordHandler,
	)

	// Report
	TrackingReport := trackingBasicReport.New(
		PartyBasicAdministrator,
	)

	// ________________________________ Create Service Provider Adaptors ________________________________

	// Auth
	APIUserAuthServiceAdaptor := authServiceJsonRpcAdaptor.New(APIUserAuthorizationService)

	// Party
	PartyBasicRegistrarAdaptor := partyBasicRegistrarJsonRpcAdaptor.New(PartyBasicRegistrar)
	PartyHandlerAdaptor := partyAdministratorJsonRpcAdaptor.New(PartyBasicAdministrator)

	// System
	SystemRecordHandlerAdaptor := systemRecordHandlerJsonRpcAdaptor.New(SystemRecordHandler)

	// Sigbug
	SigbugRecordHandlerAdaptor := sigbugRecordHandlerJsonRpcAdaptor.New(SigbugRecordHandler)
	SigbugValidatorAdaptor := sigbugValidatorJsonRpcAdaptor.New(SigbugValidator)
	SigbugAdministratorAdaptor := sigbugAdministratorJsonRpcAdaptor.New(SigbugAdministrator)

	// Report
	TrackingReportAdaptor := trackingReportJsonRpcAdaptor.New(TrackingReport)

	// Sigfox Backend Callback Server
	SigfoxBackendCallbackServer := sigfoxBasicBackendCallbackServer.New()

	// ________________________________ Register Service Provider Adaptors with secureHumanUserAPIServer ________________________________
	// Create secureHumanUserAPIServer
	secureHumanUserAPIServer := rpc.NewServer()
	secureHumanUserAPIServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	humanUserJsonRpcHttpServer := jsonRpcHttpServer.New(
		"/api-1",
		"0.0.0.0",
		humanUserAPIServerPort,
		humanUserAuthoriser.New(
			token.NewJWTValidator(&rsaPrivateKey.PublicKey),
			PermissionBasicHandler,
		),
	)
	if err := humanUserJsonRpcHttpServer.RegisterBatchServiceProviders(
		[]jsonRpcServiceProvider.Provider{
			humanUserRecordHandlerJsonRpcAdaptor.New(UserRecordHandler),
			humanUserValidatorJsonRpcAdaptor.New(UserValidator),
			humanUserAdministratorJsonRpcAdaptor.New(UserBasicAdministrator),
			apiUserRecordHandlerJsonRpcAdaptor.New(APIUserRecordHandler),
			apiUserValidatorJsonRpcAdaptor.New(APIUserValidator),
			apiUserAdministratorJsonRpcAdaptor.New(APIUserAdministrator),
			authServiceJsonRpcAdaptor.New(HumanUserAuthorizationService),
			permissionAdministratorJsonRpcAdaptor.New(PermissionBasicHandler),
			companyRecordHandlerJsonRpcAdaptor.New(CompanyRecordHandler),
			companyValidatorJsonRpcAdaptor.New(CompanyValidator),
			companyAdministratorJsonRpcAdaptor.New(CompanyAdministrator),
			clientRecordHandlerJsonRpcAdaptor.New(ClientRecordHandler),
			clientValidatorJsonRpcAdaptor.New(ClientValidator),
			clientAdministratorJsonRpcAdaptor.New(ClientBasicAdministrator),
		},
	); err != nil {
		log.Fatal(err)
	}

	// Party
	if err := secureHumanUserAPIServer.RegisterService(PartyBasicRegistrarAdaptor, partyRegistrar.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Party Registrar Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(PartyHandlerAdaptor, partyAdministrator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Party Administrator Service")
	}

	// System
	if err := secureHumanUserAPIServer.RegisterService(SystemRecordHandlerAdaptor, systemRecordHandler.ServiceProvider); err != nil {
		log.Fatal("Unable to Register System Record Handler Service")
	}

	// Sigbug
	if err := secureHumanUserAPIServer.RegisterService(SigbugRecordHandlerAdaptor, sigbugRecordHandler.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Sigbug Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(SigbugValidatorAdaptor, sigbugValidator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Sigbug Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(SigbugAdministratorAdaptor, sigbugAdministrator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Sigbug Administrator Service")
	}

	// Reports
	if err := secureHumanUserAPIServer.RegisterService(TrackingReportAdaptor, trackingReport.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Tracking Report Service")
	}

	// Set up Secure Human User API Server i.e. the Portal API Server
	HumanUserAPIAuthorizer := humanUserAPIAuthorizer.New(
		token.NewJWTValidator(&rsaPrivateKey.PublicKey),
		PermissionBasicHandler,
	)
	HumanUserHttpAPIAuthApplier := humanUserHttpAPIAuthApplier.New(
		HumanUserAPIAuthorizer,
	)
	humanUserSecureAPIServerMux := mux.NewRouter()
	humanUserSecureAPIServerMux.Methods("OPTIONS").HandlerFunc(HumanUserHttpAPIAuthApplier.PreFlightHandler)
	humanUserSecureAPIServerMux.Handle("/api-1", HumanUserHttpAPIAuthApplier.ApplyAuth(secureHumanUserAPIServer)).Methods("POST")
	// Start secureHumanUserAPIServer
	log.Info("Starting Human User secure API Server on port " + humanUserAPIServerPort)
	go func() {
		err := http.ListenAndServe(":"+humanUserAPIServerPort, humanUserSecureAPIServerMux)
		log.Error("secureHumanUserAPIServer stopped: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}()

	// Create secureAPIUserAPIServer
	secureAPIUserAPIServer := rpc.NewServer()
	secureAPIUserAPIServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	// Auth
	if err := secureAPIUserAPIServer.RegisterService(APIUserAuthServiceAdaptor, authorizationAdministrator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register API User Authorization Service Adaptor")
	}

	// Sigfox Test
	// set  up sigfox backend server
	sigfoxBackendJsonRpcHttpServer := jsonRpcHttpServer.New(
		"/api-sigfox",
		"0.0.0.0",
		"9011",
		sigfoxBackendAuthoriser.New(),
	)

	// register service providers
	if err := sigfoxBackendJsonRpcHttpServer.RegisterServiceProvider(sigfoxBasicBackendCallbackServerJsonRpcAdaptor.New(SigfoxBackendCallbackServer)); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Starting Sigfox Backend secure API Server on port " + "9011")
	go func() {
		err := sigfoxBackendJsonRpcHttpServer.SecureStart()
		log.Error("sigfox backend json rpc http server", err)
		os.Exit(1)
	}()

	//// set up kafka messaging
	//MessageConsumerGroup := messageConsumerGroup.New(
	//	kafkaBrokerNodes,
	//	[]string{"brainQueue"},
	//	"brain",
	//	[]messagingMessageHandler.Handler{
	//		zx303GPSReadingMessageHandler.New(
	//			&systemClaims,
	//			ZX303GPSReadingAdministrator,
	//		),
	//		zx303StatusReadingMessageHandler.New(
	//			&systemClaims,
	//			ZX303StatusReadingAdministrator,
	//		),
	//	},
	//)
	//go func() {
	//	if err := MessageConsumerGroup.Start(); err != nil {
	//		log.Error(err.Error())
	//		os.Exit(1)
	//	}
	//}()

	//Wait for interrupt signal
	systemSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(systemSignalsChannel, os.Interrupt)
	for {
		select {
		case s := <-systemSignalsChannel:
			log.Info("Application is shutting down.. ( ", s, " )")
			return
		}
	}
}
