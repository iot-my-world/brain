package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	gorillaJson "github.com/gorilla/rpc/json"
	"github.com/iot-my-world/brain/cors"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/security/encrypt"
	"github.com/iot-my-world/brain/security/token"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	databaseCollection "github.com/iot-my-world/brain/database/collection"

	authServiceJsonRpcAdaptor "github.com/iot-my-world/brain/security/authorization/service/adaptor/jsonRpc"
	apiUserAuthorizationService "github.com/iot-my-world/brain/security/authorization/service/user/api"
	humanUserAuthorizationService "github.com/iot-my-world/brain/security/authorization/service/user/human"

	humanUserHttpAPIAuthApplier "github.com/iot-my-world/brain/security/authorization/api/applier/http/user/human"
	humanUserAPIAuthorizer "github.com/iot-my-world/brain/security/authorization/api/authorizer/user/human"

	apiUserHttpAPIAuthApplier "github.com/iot-my-world/brain/security/authorization/api/applier/http/user/api"
	apiUserAPIAuthorizer "github.com/iot-my-world/brain/security/authorization/api/authorizer/user/api"

	permissionAdministrator "github.com/iot-my-world/brain/security/permission/administrator"
	permissionAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/security/permission/administrator/adaptor/jsonRpc"
	permissionBasicAdministrator "github.com/iot-my-world/brain/security/permission/administrator/basic"

	roleMongoRecordHandler "github.com/iot-my-world/brain/security/role/recordHandler/mongo"

	humanUserAdministrator "github.com/iot-my-world/brain/user/human/administrator"
	humanUserAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/user/human/administrator/adaptor/jsonRpc"
	humanUserBasicAdministrator "github.com/iot-my-world/brain/user/human/administrator/basic"
	humanUserRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler"
	humanUserRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/user/human/recordHandler/adaptor/jsonRpc"
	humanUserMongoRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler/mongo"
	humanUserValidator "github.com/iot-my-world/brain/user/human/validator"
	humanUserValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/user/human/validator/adaptor/jsonRpc"
	humanUserBasicValidator "github.com/iot-my-world/brain/user/human/validator/basic"

	individualAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/party/individual/administrator/adaptor/jsonRpc"
	individualBasicAdministrator "github.com/iot-my-world/brain/party/individual/administrator/basic"
	individualRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/party/individual/recordHandler/adaptor/jsonRpc"
	individualMongoRecordHandler "github.com/iot-my-world/brain/party/individual/recordHandler/mongo"
	individualValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/party/individual/validator/adaptor/jsonRpc"
	individualBasicValidator "github.com/iot-my-world/brain/party/individual/validator/basic"

	companyAdministrator "github.com/iot-my-world/brain/party/company/administrator"
	companyAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/party/company/administrator/adaptor/jsonRpc"
	companyBasicAdministrator "github.com/iot-my-world/brain/party/company/administrator/basic"
	companyRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler"
	companyRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/party/company/recordHandler/adaptor/jsonRpc"
	companyMongoRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler/mongo"
	companyValidator "github.com/iot-my-world/brain/party/company/validator"
	companyValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/party/company/validator/adaptor/jsonRpc"
	companyBasicValidator "github.com/iot-my-world/brain/party/company/validator/basic"

	clientAdministrator "github.com/iot-my-world/brain/party/client/administrator"
	clientAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/party/client/administrator/adaptor/jsonRpc"
	clientBasicAdministrator "github.com/iot-my-world/brain/party/client/administrator/basic"
	clientRecordHandler "github.com/iot-my-world/brain/party/client/recordHandler"
	clientRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/party/client/recordHandler/adaptor/jsonRpc"
	clientMongoRecordHandler "github.com/iot-my-world/brain/party/client/recordHandler/mongo"
	clientValidator "github.com/iot-my-world/brain/party/client/validator"
	clientValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/party/client/validator/adaptor/jsonRpc"
	clientBasicValidator "github.com/iot-my-world/brain/party/client/validator/basic"

	systemRecordHandler "github.com/iot-my-world/brain/party/system/recordHandler"
	systemRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/party/system/recordHandler/adaptor/jsonRpc"
	systemMongoRecordHandler "github.com/iot-my-world/brain/party/system/recordHandler/mongo"

	sf001AdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/tracker/sf001/administrator/adaptor/jsonRpc"
	sf001TrackerBasicAdministrator "github.com/iot-my-world/brain/tracker/sf001/administrator/basic"
	sf001TrackerRecordHandler "github.com/iot-my-world/brain/tracker/sf001/recordHandler"
	sf001RecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/tracker/sf001/recordHandler/adaptor/jsonRpc"
	sf001TrackerMongoRecordHandler "github.com/iot-my-world/brain/tracker/sf001/recordHandler/mongo"
	sf001ValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/tracker/sf001/validator/adaptor/jsonRpc"
	sf001TrackerBasicValidator "github.com/iot-my-world/brain/tracker/sf001/validator/basic"

	apiUserAdministrator "github.com/iot-my-world/brain/user/api/administrator"
	apiUserAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/user/api/administrator/adaptor/jsonRpc"
	apiUserBasicAdministrator "github.com/iot-my-world/brain/user/api/administrator/basic"
	apiUserBasicPasswordGenerator "github.com/iot-my-world/brain/user/api/password/generator/basic"
	apiUserRecordHandler "github.com/iot-my-world/brain/user/api/recordHandler"
	apiUserRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/user/api/recordHandler/adaptor/jsonRpc"
	apiUserMongoRecordHandler "github.com/iot-my-world/brain/user/api/recordHandler/mongo"
	apiUserValidator "github.com/iot-my-world/brain/user/api/validator"
	apiUserValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/user/api/validator/adaptor/jsonRpc"
	apiUserBasicValidator "github.com/iot-my-world/brain/user/api/validator/basic"

	trackingReport "github.com/iot-my-world/brain/report/tracking"
	trackingReportJsonRpcAdaptor "github.com/iot-my-world/brain/report/tracking/adaptor/jsonRpc"
	trackingBasicReport "github.com/iot-my-world/brain/report/tracking/basic"

	"flag"
	"github.com/iot-my-world/brain/communication/email/mailer"
	gmailMailer "github.com/iot-my-world/brain/communication/email/mailer/gmail"
	partyRegistrar "github.com/iot-my-world/brain/party/registrar"
	partyBasicRegistrarJsonRpcAdaptor "github.com/iot-my-world/brain/party/registrar/adaptor/jsonRpc"
	partyBasicRegistrar "github.com/iot-my-world/brain/party/registrar/basic"

	registrationEmailGenerator "github.com/iot-my-world/brain/communication/email/generator/registration"
	setPasswordEmailGenerator "github.com/iot-my-world/brain/communication/email/generator/set/password"

	partyAdministrator "github.com/iot-my-world/brain/party/administrator"
	partyAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/party/administrator/adaptor/jsonRpc"
	partyBasicAdministrator "github.com/iot-my-world/brain/party/administrator/basic"

	"fmt"
	"github.com/iot-my-world/brain/party"
	humanUserLoginClaims "github.com/iot-my-world/brain/security/claims/login/user/human"
	"path/filepath"
	"strings"
)

var humanUserAPIServerPort = "9010"
var apiUserAPIServerPort = "9011"

func main() {
	// get the command line args
	mongoNodes := flag.String("mongoNodes", "localhost:27016", "the nodes in the db cluster")
	mongoUser := flag.String("mongoUser", "", "brains mongo db user")
	mongoPassword := flag.String("mongoPassword", "", "passwords for brains mongo db")
	mailRedirectBaseUrl := flag.String("mailRedirectBaseUrl", "http://localhost:3000", "base url for all email invites")
	rootPasswordFileLocation := flag.String("rootPasswordFileLocation", "", "path to file containing root password")
	pathToEmailTemplateFolder := flag.String("pathToEmailTemplateFolder", "communication/email/template", "path to email template files")
	keysFilePath := flag.String("keysFilePath", "", "path to pvt and pub keys")
	//kafkaBrokers := flag.String("kafkaBrokers", "localhost:9092", "ipAddress:port of each kafka broker node (, separated)")

	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("fail to get dir:" + err.Error())
	}
	log.Info("brain working directory: " + dir)

	// Connect to database
	databaseName := "brain"
	mongoCluster := strings.Split(*mongoNodes, ",")
	log.Info(fmt.Sprintf("connecting to mongo @ node addresses: [%s]", strings.Join(mongoCluster, ", ")))
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
	rsaPrivateKey := encrypt.FetchPrivateKey(*keysFilePath)

	// Create Mailer
	Mailer := gmailMailer.New(mailer.AuthInfo{
		Identity: "",
		Username: "spotnavza@gmail.com",
		Password: "spotNav123",
		Host:     "smtp.gmail.com",
	})

	// email generators
	RegistrationEmailGenerator := registrationEmailGenerator.New(
		*pathToEmailTemplateFolder,
	)
	SetPasswordEmailGenerator := setPasswordEmailGenerator.New(
		*pathToEmailTemplateFolder,
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
		//UserId          id.Identifier `json:"userId"`
		//IssueTime       int64         `json:"issueTime"`
		//ExpirationTime  int64         `json:"expirationTime"`
		//ParentPartyType party.Type    `json:"parentPartyType"`
		//ParentId        id.Identifier `json:"parentId"`
		PartyType: party.System,
		//PartyId         id.Identifier `json:"partyId"`
	}

	// ________________________________ Create Service Providers ________________________________

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
	UserValidator := humanUserBasicValidator.New(
		UserRecordHandler,
		&systemClaims,
	)
	UserBasicAdministrator := humanUserBasicAdministrator.New(
		UserRecordHandler,
		UserValidator,
		Mailer,
		rsaPrivateKey,
		*mailRedirectBaseUrl,
		&systemClaims,
		SetPasswordEmailGenerator,
	)

	// Auth
	HumanUserAuthorizationService := humanUserAuthorizationService.New(
		UserRecordHandler,
		rsaPrivateKey,
		&systemClaims,
	)

	// Individual
	IndividualRecordHandler := individualMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.Individual,
	)
	IndividualValidator := individualBasicValidator.New()
	IndividualBasicAdministrator := individualBasicAdministrator.New(
		IndividualValidator,
		IndividualRecordHandler,
	)

	// Company
	CompanyRecordHandler := companyMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.Company,
	)

	CompanyValidator := companyBasicValidator.New(
		CompanyRecordHandler,
		UserRecordHandler,
		&systemClaims,
	)
	CompanyBasicAdministrator := companyBasicAdministrator.New(
		CompanyRecordHandler,
		CompanyValidator,
		UserRecordHandler,
	)

	// Client
	ClientRecordHandler := clientMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.Client,
	)
	ClientValidator := clientBasicValidator.New(
		ClientRecordHandler,
		UserRecordHandler,
		&systemClaims,
	)
	ClientBasicAdministrator := clientBasicAdministrator.New(
		ClientRecordHandler,
		ClientValidator,
		UserRecordHandler,
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
		*mailRedirectBaseUrl,
		&systemClaims,
		RegistrationEmailGenerator,
	)

	// System
	SystemRecordHandler := systemMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.System,
		*rootPasswordFileLocation,
		PartyBasicRegistrar,
		&systemClaims,
	)

	// Party
	PartyBasicAdministrator := partyBasicAdministrator.New(
		ClientRecordHandler,
		CompanyRecordHandler,
		SystemRecordHandler,
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

	APIUserAuthorizationService := apiUserAuthorizationService.New(
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

	// SF001 Device
	SF001TrackerMongoRecordHandler := sf001TrackerMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.SF001Tracker,
	)
	SF001TrackerBasicValidator := sf001TrackerBasicValidator.New(
		PartyBasicAdministrator,
	)
	SF001TrackerBasicAdministrator := sf001TrackerBasicAdministrator.New(
		SF001TrackerBasicValidator,
		SF001TrackerMongoRecordHandler,
	)

	// Report
	TrackingReport := trackingBasicReport.New(
		PartyBasicAdministrator,
	)

	// ________________________________ Create Service Provider Adaptors ________________________________

	// User
	HumanUserRecordHandlerAdaptor := humanUserRecordHandlerJsonRpcAdaptor.New(UserRecordHandler)
	HumanUserValidatorAdaptor := humanUserValidatorJsonRpcAdaptor.New(UserValidator)
	HumanUserAdministratorAdaptor := humanUserAdministratorJsonRpcAdaptor.New(UserBasicAdministrator)

	// APIUser
	APIUserRecordHandlerAdaptor := apiUserRecordHandlerJsonRpcAdaptor.New(APIUserRecordHandler)
	APIUserValidatorAdaptor := apiUserValidatorJsonRpcAdaptor.New(APIUserValidator)
	APIUserAdministratorAdaptor := apiUserAdministratorJsonRpcAdaptor.New(APIUserAdministrator)

	// Auth
	HumanUserAuthServiceAdaptor := authServiceJsonRpcAdaptor.New(HumanUserAuthorizationService)
	APIUserAuthServiceAdaptor := authServiceJsonRpcAdaptor.New(APIUserAuthorizationService)

	// Permission
	PermissionAdministratorAdaptor := permissionAdministratorJsonRpcAdaptor.New(PermissionBasicHandler)

	// Individual
	IndividualRecordHandlerAdaptor := individualRecordHandlerJsonRpcAdaptor.New(IndividualRecordHandler)
	IndividualValidatorAdaptor := individualValidatorJsonRpcAdaptor.New(IndividualValidator)
	IndividualAdministratorAdaptor := individualAdministratorJsonRpcAdaptor.New(IndividualBasicAdministrator)

	// Company
	CompanyRecordHandlerAdaptor := companyRecordHandlerJsonRpcAdaptor.New(CompanyRecordHandler)
	CompanyValidatorAdaptor := companyValidatorJsonRpcAdaptor.New(CompanyValidator)
	CompanyAdministratorAdaptor := companyAdministratorJsonRpcAdaptor.New(CompanyBasicAdministrator)

	// Client
	ClientRecordHandlerAdaptor := clientRecordHandlerJsonRpcAdaptor.New(ClientRecordHandler)
	ClientValidatorAdaptor := clientValidatorJsonRpcAdaptor.New(ClientValidator)
	ClientAdministratorAdaptor := clientAdministratorJsonRpcAdaptor.New(ClientBasicAdministrator)

	// Party
	PartyBasicRegistrarAdaptor := partyBasicRegistrarJsonRpcAdaptor.New(PartyBasicRegistrar)
	PartyHandlerAdaptor := partyAdministratorJsonRpcAdaptor.New(PartyBasicAdministrator)

	// System
	SystemRecordHandlerAdaptor := systemRecordHandlerJsonRpcAdaptor.New(SystemRecordHandler)

	// Report
	TrackingReportAdaptor := trackingReportJsonRpcAdaptor.New(TrackingReport)

	// SF001 Tracker
	SF001TrackerRecordHandlerJsonRpcAdaptor := sf001RecordHandlerJsonRpcAdaptor.New(SF001TrackerMongoRecordHandler)
	SF001TrackerValidatorJsonRpcAdaptor := sf001ValidatorJsonRpcAdaptor.New(SF001TrackerBasicValidator)
	SF001TrackerAdministratorJsonRpcAdaptor := sf001AdministratorJsonRpcAdaptor.New(SF001TrackerBasicAdministrator)

	// ________________________________ Register Service Provider Adaptors with secureHumanUserAPIServer ________________________________
	// Create secureHumanUserAPIServer
	secureHumanUserAPIServer := rpc.NewServer()
	secureHumanUserAPIServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	// User
	if err := secureHumanUserAPIServer.RegisterService(HumanUserRecordHandlerAdaptor, humanUserRecordHandler.ServiceProvider); err != nil {
		log.Fatal("Unable to Register User Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(HumanUserValidatorAdaptor, humanUserValidator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register User Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(HumanUserAdministratorAdaptor, humanUserAdministrator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register User Administrator Service")
	}
	// API User
	if err := secureHumanUserAPIServer.RegisterService(APIUserRecordHandlerAdaptor, apiUserRecordHandler.ServiceProvider); err != nil {
		log.Fatal("Unable to Register API User Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(APIUserValidatorAdaptor, apiUserValidator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register API User Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(APIUserAdministratorAdaptor, apiUserAdministrator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register API User Administrator Service")
	}

	// Auth
	if err := secureHumanUserAPIServer.RegisterService(HumanUserAuthServiceAdaptor, "Auth"); err != nil {
		log.Fatal("Unable to Register Auth Service Adaptor")
	}

	// Permission
	if err := secureHumanUserAPIServer.RegisterService(PermissionAdministratorAdaptor, permissionAdministrator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Permission Handler Service Adaptor")
	}

	// Individual
	if err := secureHumanUserAPIServer.RegisterService(IndividualRecordHandlerAdaptor, "IndividualRecordHandler"); err != nil {
		log.Fatal("Unable to Register Individual Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(IndividualValidatorAdaptor, "IndividualValidator"); err != nil {
		log.Fatal("Unable to Register Individual Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(IndividualAdministratorAdaptor, "IndividualAdministrator"); err != nil {
		log.Fatal("Unable to Register Individual Administrator Service")
	}

	// Company
	if err := secureHumanUserAPIServer.RegisterService(CompanyRecordHandlerAdaptor, companyRecordHandler.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Company Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(CompanyValidatorAdaptor, companyValidator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Company Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(CompanyAdministratorAdaptor, companyAdministrator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Company Administrator Service")
	}

	// Client
	if err := secureHumanUserAPIServer.RegisterService(ClientRecordHandlerAdaptor, clientRecordHandler.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Client Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(ClientValidatorAdaptor, clientValidator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Client Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(ClientAdministratorAdaptor, clientAdministrator.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Client Administrator Service")
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

	// Reports
	if err := secureHumanUserAPIServer.RegisterService(TrackingReportAdaptor, trackingReport.ServiceProvider); err != nil {
		log.Fatal("Unable to Register Tracking Report Service")
	}

	// SF001 Tracker
	if err := secureHumanUserAPIServer.RegisterService(SF001TrackerRecordHandlerJsonRpcAdaptor, sf001TrackerRecordHandler.ServiceProvider); err != nil {
		log.Fatal("Unable to Register SF001 Tracker RecordHandler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(SF001TrackerValidatorJsonRpcAdaptor, "SF001TrackerValidator"); err != nil {
		log.Fatal("Unable to Register SF001 Tracker Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(SF001TrackerAdministratorJsonRpcAdaptor, "SF001TrackerAdministrator"); err != nil {
		log.Fatal("Unable to Register SF001 Tracker Administrator Service")
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
	if err := secureAPIUserAPIServer.RegisterService(APIUserAuthServiceAdaptor, "Auth"); err != nil {
		log.Fatal("Unable to Register API User Authorization Service Adaptor")
	}

	// Sigfox Test

	// Set up Secure API User API Server
	APIUserAPIAuthorizer := apiUserAPIAuthorizer.New(
		token.NewJWTValidator(&rsaPrivateKey.PublicKey),
		PermissionBasicHandler,
	)
	APIUserHttpAPIAuthApplier := apiUserHttpAPIAuthApplier.New(
		APIUserAPIAuthorizer,
	)
	apiUserSecureAPIServerMux := mux.NewRouter()
	apiUserSecureAPIServerMux.Methods("OPTIONS").HandlerFunc(APIUserHttpAPIAuthApplier.PreFlightHandler)
	//apiUserSecureAPIServerMux.Handle("/api-2", APIUserHttpAPIAuthApplier.ApplyAuth(secureAPIUserAPIServer)).Methods("POST")
	apiUserSecureAPIServerMux.Handle("/api-2", secureAPIUserAPIServer).Methods("POST")
	// Start secureAPIUserAPIServer
	log.Info("Starting API User Secure API Server on port " + apiUserAPIServerPort)
	go func() {
		err := http.ListenAndServe(":"+apiUserAPIServerPort, apiUserSecureAPIServerMux)
		log.Error("apiUserSecureAPIServerMux stopped: ", err, "\n", string(debug.Stack()))
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
