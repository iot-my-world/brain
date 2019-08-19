package main

import (
	"github.com/iot-my-world/brain/internal/config"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/internal/security/encrypt"
	"github.com/iot-my-world/brain/pkg/security/token"
	"gopkg.in/mgo.v2"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	databaseCollection "github.com/iot-my-world/brain/pkg/database/collection"

	humanUserJsonRpcServerAuthenticatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator/adaptor/jsonRpc"
	humanUserJsonRpcServerAuthenticator "github.com/iot-my-world/brain/pkg/user/human/authenticator"

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

	systemRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/system/recordHandler/adaptor/jsonRpc"
	systemMongoRecordHandler "github.com/iot-my-world/brain/pkg/party/system/recordHandler/mongo"

	apiUserAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/user/api/administrator/adaptor/jsonRpc"
	apiUserBasicAdministrator "github.com/iot-my-world/brain/pkg/user/api/administrator/basic"
	apiUserBasicPasswordGenerator "github.com/iot-my-world/brain/pkg/user/api/password/generator/basic"
	apiUserRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/user/api/recordHandler/adaptor/jsonRpc"
	apiUserMongoRecordHandler "github.com/iot-my-world/brain/pkg/user/api/recordHandler/mongo"
	apiUserValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/user/api/validator/adaptor/jsonRpc"
	apiUserBasicValidator "github.com/iot-my-world/brain/pkg/user/api/validator/basic"

	trackingReportJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/report/tracking/adaptor/jsonRpc"
	trackingBasicReport "github.com/iot-my-world/brain/pkg/report/tracking/basic"

	"flag"
	"github.com/iot-my-world/brain/pkg/communication/email/mailer"
	gmailMailer "github.com/iot-my-world/brain/pkg/communication/email/mailer/gmail"
	partyBasicRegistrarJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/registrar/adaptor/jsonRpc"
	partyBasicRegistrar "github.com/iot-my-world/brain/pkg/party/registrar/basic"

	registrationEmailGenerator "github.com/iot-my-world/brain/pkg/communication/email/generator/registration"
	setPasswordEmailGenerator "github.com/iot-my-world/brain/pkg/communication/email/generator/set/password"

	partyAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/administrator/adaptor/jsonRpc"
	partyBasicAdministrator "github.com/iot-my-world/brain/pkg/party/administrator/basic"

	"fmt"
	"github.com/iot-my-world/brain/pkg/party"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	"path/filepath"
	"strings"

	sigbugAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/administrator/adaptor/jsonRpc"
	sigbugBasicAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator/basic"
	sigbugGPSReadingAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/administrator/adaptor/jsonRpc"
	sigbugGPSReadingBasicAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/administrator/basic"
	sigbugGPSReadingRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/recordHandler/adaptor/jsonRpc"
	sigbugGPSReadingMongoRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/recordHandler/mongo"
	sigbugGPSReadingValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/validator/adaptor/jsonRpc"
	sigbugGPSReadingBasicValidator "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/validator/basic"
	sigbugRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/adaptor/jsonRpc"
	sigbugMongoRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/mongo"
	sigbugSigfoxMessageHandler "github.com/iot-my-world/brain/pkg/device/sigbug/sigfox/message/handler"
	sigbugValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/device/sigbug/validator/adaptor/jsonRpc"
	sigbugBasicValidator "github.com/iot-my-world/brain/pkg/device/sigbug/validator/basic"

	jsonRpcHttpServer "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"

	sigfoxBackendAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/sigfox/backend/administrator/adaptor/jsonRpc"
	sigfoxBackendBasicAdministrator "github.com/iot-my-world/brain/pkg/sigfox/backend/administrator/basic"
	sigfoxBackendAuthoriser "github.com/iot-my-world/brain/pkg/sigfox/backend/authoriser"
	sigfoxBackendDataMessageHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/handler"
	sigfoxBackendCallbackServerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server/adaptor/jsonRpc"
	sigfoxBasicBackendCallbackServer "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server/basic"
	sigfoxBackendRecordHandlerJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler/adaptor/jsonRpc"
	sigfoxBackendMongoRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler/mongo"
	sigfoxBackendValidatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/sigfox/backend/validator/adaptor/jsonRpc"
	sigfoxBackendBasicValidator "github.com/iot-my-world/brain/pkg/sigfox/backend/validator/basic"

	sigfoxBackendDataCallbackMessageBasicAdministrator "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/administrator/basic"
	sigfoxBackendDataCallbackMessageMongoRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/recordHandler/mongo"
	sigfoxBackendDataCallbackMessageBasicValidator "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/validator/basic"
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
		SigbugRecordHandler,
		PartyBasicAdministrator,
		&systemClaims,
	)
	SigbugAdministrator := sigbugBasicAdministrator.New(
		SigbugValidator,
		SigbugRecordHandler,
	)
	SigbugGPSReadingRecordHandler := sigbugGPSReadingMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.SigbugGPSReading,
	)
	SigbugGPSReadingValidator := sigbugGPSReadingBasicValidator.New(
		SigbugRecordHandler,
		PartyBasicAdministrator,
		&systemClaims,
	)
	SigbugGPSReadingAdministrator := sigbugGPSReadingBasicAdministrator.New(
		SigbugGPSReadingValidator,
		SigbugGPSReadingRecordHandler,
	)

	// Sigfox Backend
	SigfoxBackendRecordHandler := sigfoxBackendMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.SigfoxBackend,
	)
	SigfoxBackendValidator := sigfoxBackendBasicValidator.New(
		PartyBasicAdministrator,
		SigfoxBackendRecordHandler,
		&systemClaims,
		token.NewJWTValidator(&rsaPrivateKey.PublicKey),
	)
	SigfoxBackendAdministrator := sigfoxBackendBasicAdministrator.New(
		SigfoxBackendValidator,
		SigfoxBackendRecordHandler,
		rsaPrivateKey,
	)
	SigfoxBackendDataCallbackMessageMongoRecordHandler := sigfoxBackendDataCallbackMessageMongoRecordHandler.New(
		mainMongoSession,
		databaseName,
		databaseCollection.SigfoxBackendDataCallbackMessage,
	)
	SigfoxBackendDataCallbackMessageBasicValidator := sigfoxBackendDataCallbackMessageBasicValidator.New()
	SigfoxBackendDataCallbackMessageBasicAdministrator := sigfoxBackendDataCallbackMessageBasicAdministrator.New(
		SigfoxBackendDataCallbackMessageBasicValidator,
		SigfoxBackendDataCallbackMessageMongoRecordHandler,
	)

	// Report
	TrackingReport := trackingBasicReport.New(
		PartyBasicAdministrator,
	)

	HumanUserJsonRpcServerAuthenticator := humanUserJsonRpcServerAuthenticator.New(
		UserRecordHandler,
		rsaPrivateKey,
		&systemClaims,
	)

	// Sigfox Backend Callback Server
	SigfoxBackendCallbackServer := sigfoxBasicBackendCallbackServer.New(
		SigfoxBackendDataCallbackMessageBasicAdministrator,
		[]sigfoxBackendDataMessageHandler.Handler{
			sigbugSigfoxMessageHandler.New(
				SigbugRecordHandler,
				SigbugAdministrator,
				SigbugGPSReadingAdministrator,
			),
		},
	)

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
			humanUserJsonRpcServerAuthenticatorJsonRpcAdaptor.New(HumanUserJsonRpcServerAuthenticator),
			humanUserRecordHandlerJsonRpcAdaptor.New(UserRecordHandler),
			humanUserValidatorJsonRpcAdaptor.New(UserValidator),
			humanUserAdministratorJsonRpcAdaptor.New(UserBasicAdministrator),
			apiUserRecordHandlerJsonRpcAdaptor.New(APIUserRecordHandler),
			apiUserValidatorJsonRpcAdaptor.New(APIUserValidator),
			apiUserAdministratorJsonRpcAdaptor.New(APIUserAdministrator),
			permissionAdministratorJsonRpcAdaptor.New(PermissionBasicHandler),
			companyRecordHandlerJsonRpcAdaptor.New(CompanyRecordHandler),
			companyValidatorJsonRpcAdaptor.New(CompanyValidator),
			companyAdministratorJsonRpcAdaptor.New(CompanyAdministrator),
			clientRecordHandlerJsonRpcAdaptor.New(ClientRecordHandler),
			clientValidatorJsonRpcAdaptor.New(ClientValidator),
			clientAdministratorJsonRpcAdaptor.New(ClientBasicAdministrator),
			partyBasicRegistrarJsonRpcAdaptor.New(PartyBasicRegistrar),
			partyAdministratorJsonRpcAdaptor.New(PartyBasicAdministrator),
			systemRecordHandlerJsonRpcAdaptor.New(SystemRecordHandler),
			sigbugRecordHandlerJsonRpcAdaptor.New(SigbugRecordHandler),
			sigbugValidatorJsonRpcAdaptor.New(SigbugValidator),
			sigbugAdministratorJsonRpcAdaptor.New(SigbugAdministrator),
			sigbugGPSReadingRecordHandlerJsonRpcAdaptor.New(SigbugGPSReadingRecordHandler),
			sigbugGPSReadingValidatorJsonRpcAdaptor.New(SigbugGPSReadingValidator),
			sigbugGPSReadingAdministratorJsonRpcAdaptor.New(SigbugGPSReadingAdministrator),
			trackingReportJsonRpcAdaptor.New(TrackingReport),
			sigfoxBackendRecordHandlerJsonRpcAdaptor.New(SigfoxBackendRecordHandler),
			sigfoxBackendValidatorJsonRpcAdaptor.New(SigfoxBackendValidator),
			sigfoxBackendAdministratorJsonRpcAdaptor.New(SigfoxBackendAdministrator),
		},
	); err != nil {
		log.Fatal(err)
	}
	log.Info("Starting Human User API Server on port: " + humanUserAPIServerPort)
	go func() {
		err := humanUserJsonRpcHttpServer.SecureStart()
		log.Error("human user json rpc http server has stopped: ", err)
		os.Exit(1)
	}()

	// set  up sigfox backend server
	sigfoxBackendJsonRpcHttpServer := jsonRpcHttpServer.New(
		"/api-2",
		"0.0.0.0",
		"9011",
		sigfoxBackendAuthoriser.New(
			token.NewJWTValidator(&rsaPrivateKey.PublicKey),
		),
	)
	if err := sigfoxBackendJsonRpcHttpServer.RegisterBatchServiceProviders([]jsonRpcServiceProvider.Provider{
		sigfoxBackendCallbackServerJsonRpcAdaptor.New(SigfoxBackendCallbackServer),
	}); err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Starting Sigfox Backend secure API Server on port: " + "9011")
	go func() {
		err := sigfoxBackendJsonRpcHttpServer.SecureStart()
		log.Error("sigfox backend json rpc http server has stopped: ", err)
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
