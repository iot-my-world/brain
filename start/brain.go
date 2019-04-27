package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	gorillaJson "github.com/gorilla/rpc/json"
	"gitlab.com/iotTracker/brain/cors"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/security/encrypt"
	"gitlab.com/iotTracker/brain/security/token"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	authServiceJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/authorization/service/adaptor/jsonRpc"
	apiUserAuthorizationService "gitlab.com/iotTracker/brain/security/authorization/service/user/api"
	humanUserAuthorizationService "gitlab.com/iotTracker/brain/security/authorization/service/user/human"

	humanUserHttpAPIAuthApplier "gitlab.com/iotTracker/brain/security/authorization/api/applier/http/user/human"
	humanUserAPIAuthorizer "gitlab.com/iotTracker/brain/security/authorization/api/authorizer/user/human"

	apiUserHttpAPIAuthApplier "gitlab.com/iotTracker/brain/security/authorization/api/applier/http/user/api"
	apiUserAPIAuthorizer "gitlab.com/iotTracker/brain/security/authorization/api/authorizer/user/api"

	permissionAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/permission/administrator/adaptor/jsonRpc"
	permissionBasicAdministrator "gitlab.com/iotTracker/brain/security/permission/administrator/basic"

	roleMongoRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler/mongo"

	userAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/user/human/administrator/adaptor/jsonRpc"
	userBasicAdministrator "gitlab.com/iotTracker/brain/user/human/administrator/basic"
	userRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/user/human/recordHandler/adaptor/jsonRpc"
	userMongoRecordHandler "gitlab.com/iotTracker/brain/user/human/recordHandler/mongo"
	userValidatorJsonRpcAdaptor "gitlab.com/iotTracker/brain/user/human/validator/adaptor/jsonRpc"
	userBasicValidator "gitlab.com/iotTracker/brain/user/human/validator/basic"

	companyAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/administrator/adaptor/jsonRpc"
	companyBasicAdministrator "gitlab.com/iotTracker/brain/party/company/administrator/basic"
	companyRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/recordHandler/adaptor/jsonRpc"
	companyMongoRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler/mongo"
	companyValidatorJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/company/validator/adaptor/jsonRpc"
	companyBasicValidator "gitlab.com/iotTracker/brain/party/company/validator/basic"

	clientAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/client/administrator/adaptor/jsonRpc"
	clientBasicAdministrator "gitlab.com/iotTracker/brain/party/client/administrator/basic"
	clientRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/client/recordHandler/adaptor/jsonRpc"
	clientMongoRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler/mongo"
	clientValidatorJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/client/validator/adaptor/jsonRpc"
	clientBasicValidator "gitlab.com/iotTracker/brain/party/client/validator/basic"

	systemRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/system/recordHandler/adaptor/jsonRpc"
	systemMongoRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler/mongo"

	tk102DeviceAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/tk102/administrator/adaptor/jsonRpc"
	tk102DeviceBasicAdministrator "gitlab.com/iotTracker/brain/tracker/tk102/administrator/basic"
	tk102ReadingAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/tk102/reading/administrator/adaptor/jsonRpc"
	tk102ReadingBasicAdministrator "gitlab.com/iotTracker/brain/tracker/tk102/reading/administrator/basic"
	tk102ReadingRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/tk102/reading/recordHandler/adaptor/jsonRpc"
	tk102ReadingMongoRecordHandler "gitlab.com/iotTracker/brain/tracker/tk102/reading/recordHandler/mongo"
	tk102ReadingBasicValidator "gitlab.com/iotTracker/brain/tracker/tk102/reading/validator/basic"
	tk102DeviceRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/tk102/recordHandler/adaptor/jsonRpc"
	tk102DeviceMongoRecordHandler "gitlab.com/iotTracker/brain/tracker/tk102/recordHandler/mongo"
	tk102DeviceValidatorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/tk102/validator/adaptor/jsonRpc"
	tk102DeviceBasicValidator "gitlab.com/iotTracker/brain/tracker/tk102/validator/basic"

	zx303DeviceAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/zx303/administrator/adaptor/jsonRpc"
	zx303DeviceBasicAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/administrator/basic"
	zx303DeviceAuthenticatorAdaptorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/zx303/authenticator/adaptor/jsonRpc"
	zx303DeviceBasicAuthenticator "gitlab.com/iotTracker/brain/tracker/zx303/authenticator/basic"
	zx303DeviceRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/zx303/recordHandler/adaptor/jsonRpc"
	zx303DeviceMongoRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/recordHandler/mongo"
	zx303DeviceValidatorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/zx303/validator/adaptor/jsonRpc"
	zx303DeviceBasicValidator "gitlab.com/iotTracker/brain/tracker/zx303/validator/basic"

	apiUserAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/user/api/administrator/adaptor/jsonRpc"
	apiUserBasicAdministrator "gitlab.com/iotTracker/brain/user/api/administrator/basic"
	apiUserBasicPasswordGenerator "gitlab.com/iotTracker/brain/user/api/password/generator/basic"
	apiUserRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/user/api/recordHandler/adaptor/jsonRpc"
	apiUserMongoRecordHandler "gitlab.com/iotTracker/brain/user/api/recordHandler/mongo"
	apiUserValidatorJsonRpcAdaptor "gitlab.com/iotTracker/brain/user/api/validator/adaptor/jsonRpc"
	apiUserBasicValidator "gitlab.com/iotTracker/brain/user/api/validator/basic"

	trackingReportJsonRpcAdaptor "gitlab.com/iotTracker/brain/report/tracking/adaptor/jsonRpc"
	trackingBasicReport "gitlab.com/iotTracker/brain/report/tracking/basic"

	"flag"
	"gitlab.com/iotTracker/brain/communication/email/mailer"
	gmailMailer "gitlab.com/iotTracker/brain/communication/email/mailer/gmail"
	partyBasicRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	partyBasicRegistrar "gitlab.com/iotTracker/brain/party/registrar/basic"

	registrationEmailGenerator "gitlab.com/iotTracker/brain/communication/email/generator/registration"
	setPasswordEmailGenerator "gitlab.com/iotTracker/brain/communication/email/generator/set/password"

	partyAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/administrator/adaptor/jsonRpc"
	partyBasicAdministrator "gitlab.com/iotTracker/brain/party/administrator/basic"

	barcodeScanner "gitlab.com/iotTracker/brain/barcode/scanner"
	barcodeScannerJsonRpcAdaptor "gitlab.com/iotTracker/brain/barcode/scanner/adaptor/jsonRpc"

	zx303GPSReadingHandler "gitlab.com/iotTracker/brain/messaging/message/handler/zx303/reading/gps"
	messageConsumerGroup "gitlab.com/iotTracker/messaging/consumer/group"
	messagingMessageHandler "gitlab.com/iotTracker/messaging/message/handler"

	"gitlab.com/iotTracker/brain/party"
	humanUserLoginClaims "gitlab.com/iotTracker/brain/security/claims/login/user/human"
	"strings"
)

var humanUserAPIServerPort = "9010"
var apiUserAPIServerPort = "9011"

func main() {
	// get the command line args
	mongoNodes := flag.String("mongoNodes", "localhost:27017", "the nodes in the db cluster")
	mongoUser := flag.String("mongoUser", "", "brains mongo db user")
	mongoPassword := flag.String("mongoPassword", "", "passwords for brains mongo db")
	mailRedirectBaseUrl := flag.String("mailRedirectBaseUrl", "http://localhost:3000", "base url for all email invites")
	rootPasswordFileLocation := flag.String("rootPasswordFileLocation", "", "path to file containing root password")
	pathToEmailTemplateFolder := flag.String("pathToEmailTemplateFolder", "communication/email/template", "path to email template files")
	kafkaBrokers := flag.String("kafkaBrokers", "localhost:9092", "ipAddress:port of each kafka broker node (, separated)")

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

	// email generators
	RegistrationEmailGenerator := registrationEmailGenerator.New(
		*pathToEmailTemplateFolder,
	)
	SetPasswordEmailGenerator := setPasswordEmailGenerator.New(
		*pathToEmailTemplateFolder,
	)

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
	)
	UserValidator := userBasicValidator.New(
		UserRecordHandler,
		&systemClaims,
	)
	UserBasicAdministrator := userBasicAdministrator.New(
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

	// Company
	CompanyRecordHandler := companyMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		companyCollection,
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
		clientCollection,
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
		systemCollection,
		*rootPasswordFileLocation,
		PartyBasicRegistrar,
		&systemClaims,
	)

	// Reading
	TK102ReadingRecordHandler := tk102ReadingMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		tk102ReadingCollection,
	)
	ReadingValidator := tk102ReadingBasicValidator.New()
	tk102ReadingAdministrator := tk102ReadingBasicAdministrator.New(
		TK102ReadingRecordHandler,
		ReadingValidator,
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
		apiUserCollection,
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

	// TK102 Device
	TK102DeviceRecordHandler := tk102DeviceMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		tk102DeviceCollection,
	)
	TK102DeviceValidator := tk102DeviceBasicValidator.New(
		PartyBasicAdministrator,
	)
	TK102DeviceAdministrator := tk102DeviceBasicAdministrator.New(
		TK102DeviceRecordHandler,
		CompanyRecordHandler,
		ClientRecordHandler,
		PartyBasicAdministrator,
		TK102ReadingRecordHandler,
		TK102DeviceValidator,
	)

	// Device
	// ZX303 Device
	ZX303DeviceRecordHandler := zx303DeviceMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		zx303DeviceCollection,
	)
	ZX303DeviceValidator := zx303DeviceBasicValidator.New(
		PartyBasicAdministrator,
	)
	ZX303DeviceAdministrator := zx303DeviceBasicAdministrator.New(
		ZX303DeviceValidator,
		ZX303DeviceRecordHandler,
	)
	ZX303DeviceAuthenticator := zx303DeviceBasicAuthenticator.New(
		ZX303DeviceRecordHandler,
	)

	// Report
	TrackingReport := trackingBasicReport.New(
		PartyBasicAdministrator,
		TK102ReadingRecordHandler,
		TK102DeviceRecordHandler,
	)

	// Barcode Scanner
	BarcodeScanner := barcodeScanner.New()

	// Create Service Provider Adaptors
	// User
	UserRecordHandlerAdaptor := userRecordHandlerJsonRpcAdaptor.New(UserRecordHandler)
	UserValidatorAdaptor := userValidatorJsonRpcAdaptor.New(UserValidator)
	UserAdministratorAdaptor := userAdministratorJsonRpcAdaptor.New(UserBasicAdministrator)

	// APIUser
	APIUserRecordHandlerAdaptor := apiUserRecordHandlerJsonRpcAdaptor.New(APIUserRecordHandler)
	APIUserValidatorAdaptor := apiUserValidatorJsonRpcAdaptor.New(APIUserValidator)
	APIUserAdministratorAdaptor := apiUserAdministratorJsonRpcAdaptor.New(APIUserAdministrator)

	// Auth
	HumanUserAuthServiceAdaptor := authServiceJsonRpcAdaptor.New(HumanUserAuthorizationService)
	APIUserAuthServiceAdaptor := authServiceJsonRpcAdaptor.New(APIUserAuthorizationService)

	// Permission
	PermissionHandlerAdaptor := permissionAdministratorJsonRpcAdaptor.New(PermissionBasicHandler)

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

	// TK102 Device
	TK102DeviceRecordHandlerAdaptor := tk102DeviceRecordHandlerJsonRpcAdaptor.New(TK102DeviceRecordHandler)
	TK102DeviceAdministratorAdaptor := tk102DeviceAdministratorJsonRpcAdaptor.New(TK102DeviceAdministrator)
	TK102DeviceValidatorAdaptor := tk102DeviceValidatorJsonRpcAdaptor.New(TK102DeviceValidator)
	tk102ReadingRecordHandlerAdaptor := tk102ReadingRecordHandlerJsonRpcAdaptor.New(TK102ReadingRecordHandler)
	tk102ReadingAdministratorAdaptor := tk102ReadingAdministratorJsonRpcAdaptor.New(tk102ReadingAdministrator)

	// ZX303 Device
	ZX303DeviceRecordHandlerAdaptor := zx303DeviceRecordHandlerJsonRpcAdaptor.New(ZX303DeviceRecordHandler)
	ZX303DeviceAdministratorAdaptor := zx303DeviceAdministratorJsonRpcAdaptor.New(ZX303DeviceAdministrator)
	ZX303DeviceValidatorAdaptor := zx303DeviceValidatorJsonRpcAdaptor.New(ZX303DeviceValidator)
	ZX303DeviceAuthenticatorAdaptor := zx303DeviceAuthenticatorAdaptorJsonRpcAdaptor.New(ZX303DeviceAuthenticator)

	// Report
	TrackingReportAdaptor := trackingReportJsonRpcAdaptor.New(TrackingReport)

	// Barcode Scanner
	BarcodeScannerAdaptor := barcodeScannerJsonRpcAdaptor.New(BarcodeScanner)

	// Create secureHumanUserAPIServer
	secureHumanUserAPIServer := rpc.NewServer()
	secureHumanUserAPIServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	// Register Service Provider Adaptors with secureHumanUserAPIServer
	// User
	if err := secureHumanUserAPIServer.RegisterService(UserRecordHandlerAdaptor, "UserRecordHandler"); err != nil {
		log.Fatal("Unable to Register User Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(UserValidatorAdaptor, "UserValidator"); err != nil {
		log.Fatal("Unable to Register User Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(UserAdministratorAdaptor, "UserAdministrator"); err != nil {
		log.Fatal("Unable to Register User Administrator Service")
	}
	// API User
	if err := secureHumanUserAPIServer.RegisterService(APIUserRecordHandlerAdaptor, "APIUserRecordHandler"); err != nil {
		log.Fatal("Unable to Register API User Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(APIUserValidatorAdaptor, "APIUserValidator"); err != nil {
		log.Fatal("Unable to Register API User Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(APIUserAdministratorAdaptor, "APIUserAdministrator"); err != nil {
		log.Fatal("Unable to Register API User Administrator Service")
	}

	// Auth
	if err := secureHumanUserAPIServer.RegisterService(HumanUserAuthServiceAdaptor, "Auth"); err != nil {
		log.Fatal("Unable to Register Auth Service Adaptor")
	}

	// Permission
	if err := secureHumanUserAPIServer.RegisterService(PermissionHandlerAdaptor, "PermissionHandler"); err != nil {
		log.Fatal("Unable to Register Permission Handler Service Adaptor")
	}

	// Company
	if err := secureHumanUserAPIServer.RegisterService(CompanyRecordHandlerAdaptor, "CompanyRecordHandler"); err != nil {
		log.Fatal("Unable to Register Company Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(CompanyValidatorAdaptor, "CompanyValidator"); err != nil {
		log.Fatal("Unable to Register Company Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(CompanyAdministratorAdaptor, "CompanyAdministrator"); err != nil {
		log.Fatal("Unable to Register Company Administrator Service")
	}

	// Client
	if err := secureHumanUserAPIServer.RegisterService(ClientRecordHandlerAdaptor, "ClientRecordHandler"); err != nil {
		log.Fatal("Unable to Register Client Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(ClientValidatorAdaptor, "ClientValidator"); err != nil {
		log.Fatal("Unable to Register Client Validator Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(ClientAdministratorAdaptor, "ClientAdministrator"); err != nil {
		log.Fatal("Unable to Register Client Administrator Service")
	}

	// Party
	if err := secureHumanUserAPIServer.RegisterService(PartyBasicRegistrarAdaptor, "PartyRegistrar"); err != nil {
		log.Fatal("Unable to Register Party Registrar Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(PartyHandlerAdaptor, "PartyAdministrator"); err != nil {
		log.Fatal("Unable to Register Party Administrator Service")
	}

	// System
	if err := secureHumanUserAPIServer.RegisterService(SystemRecordHandlerAdaptor, "SystemRecordHandler"); err != nil {
		log.Fatal("Unable to Register System Record Handler Service")
	}

	// TK102 Device
	if err := secureHumanUserAPIServer.RegisterService(TK102DeviceRecordHandlerAdaptor, "TK102DeviceRecordHandler"); err != nil {
		log.Fatal("Unable to Register TK102 Device Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(TK102DeviceValidatorAdaptor, "TK102DeviceValidator"); err != nil {
		log.Fatal("Unable to Register TK102 Device Validator")
	}
	if err := secureHumanUserAPIServer.RegisterService(TK102DeviceAdministratorAdaptor, "TK102DeviceAdministrator"); err != nil {
		log.Fatal("Unable to Register TK102 Device Administrator")
	}

	// ZX303 Device
	if err := secureHumanUserAPIServer.RegisterService(ZX303DeviceRecordHandlerAdaptor, "ZX303DeviceRecordHandler"); err != nil {
		log.Fatal("Unable to Register ZX303 Device Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(ZX303DeviceValidatorAdaptor, "ZX303DeviceValidator"); err != nil {
		log.Fatal("Unable to Register ZX303 Device Validator")
	}
	if err := secureHumanUserAPIServer.RegisterService(ZX303DeviceAdministratorAdaptor, "ZX303DeviceAdministrator"); err != nil {
		log.Fatal("Unable to Register ZX303 Device Administrator")
	}

	// Reading
	if err := secureHumanUserAPIServer.RegisterService(tk102ReadingRecordHandlerAdaptor, "TK102ReadingRecordHandler"); err != nil {
		log.Fatal("Unable to Register Reading Record Handler Service")
	}
	if err := secureHumanUserAPIServer.RegisterService(tk102ReadingAdministratorAdaptor, "TK102ReadingAdministrator"); err != nil {
		log.Fatal("Unable to Register Reading Administrator Service")
	}

	// Reports
	if err := secureHumanUserAPIServer.RegisterService(TrackingReportAdaptor, "TrackingReport"); err != nil {
		log.Fatal("Unable to Register Tracking Report Service")
	}

	// Barcode Scanner
	if err := secureHumanUserAPIServer.RegisterService(BarcodeScannerAdaptor, "BarcodeScanner"); err != nil {
		log.Fatal("Unable to Register Barcode Scanner Service")
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
	humanUserSecureAPIServerMux.Handle("/api", HumanUserHttpAPIAuthApplier.ApplyAuth(secureHumanUserAPIServer)).Methods("POST")
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

	// ZX303 Device
	if err := secureAPIUserAPIServer.RegisterService(ZX303DeviceAuthenticatorAdaptor, "ZX303DeviceAuthenticator"); err != nil {
		log.Fatal("Unable to Register API User ZX303 Device Authenticator Service Adaptor")
	}

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
	apiUserSecureAPIServerMux.Handle("/api", APIUserHttpAPIAuthApplier.ApplyAuth(secureAPIUserAPIServer)).Methods("POST")
	// Start secureAPIUserAPIServer
	log.Info("Starting API User Secure API Server on port " + apiUserAPIServerPort)
	go func() {
		err := http.ListenAndServe(":"+apiUserAPIServerPort, apiUserSecureAPIServerMux)
		log.Error("apiUserSecureAPIServerMux stopped: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}()

	// set up kafka messaging
	kafkaBrokerNodes := strings.Split(*kafkaBrokers, ",")
	MessageConsumerGroup := messageConsumerGroup.New(
		kafkaBrokerNodes,
		[]string{"brainQueue"},
		"brain",
		[]messagingMessageHandler.Handler{
			zx303GPSReadingHandler.New(),
		},
	)
	go func() {
		err := MessageConsumerGroup.Start()
		log.Error(err.Error())
		os.Exit(1)
	}()

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
