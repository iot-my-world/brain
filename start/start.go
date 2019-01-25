package main

import (
	"github.com/gorilla/rpc"
	"bitbucket.org/gotimekeeper/cors"
	gorillaJson "github.com/gorilla/rpc/json"
	"github.com/gorilla/mux"
	"bitbucket.org/gotimekeeper/log"
	"net/http"
	"runtime/debug"
	"os"
	"os/signal"
	"bitbucket.org/gotimekeeper/user"
	"gopkg.in/mgo.v2"
	"time"
	"bitbucket.org/gotimekeeper/security/encrypt"
	"bitbucket.org/gotimekeeper/security/systemRole"
	"bitbucket.org/gotimekeeper/security/apiAuth"
	"bitbucket.org/gotimekeeper/security/token"
	"bitbucket.org/gotimekeeper/rfId/tagEvent"
	"bitbucket.org/gotimekeeper/business/businessRole"

	websocketServer "bitbucket.org/gotimekeeper/exoWSC/clientHelper"
	"bitbucket.org/gotimekeeper/business/employee"
	"bitbucket.org/gotimekeeper/business/ship"
	businessDayConfig "bitbucket.org/gotimekeeper/business/businessDay/config"
	"bitbucket.org/gotimekeeper/business/businessDay"
	"bitbucket.org/gotimekeeper/rfId"
	"bitbucket.org/gotimekeeper/exoWSC"
)

var ServerPort = "9006"
var RFIDServerPort = "9004"
var webSocketServerPort = "9008"

var mainAPIAuthorizer = apiAuth.APIAuthorizer{}

func main(){

	// Connect to database
	databaseName := "timeKeeper"
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

	// Set Up Hub for websocket connections
	hub := exoWSC.NewHub()


	// Create Record Handlers
	SystemRoleRecordHandler := systemRole.NewMongoRecordHandler(mainMongoSession, databaseName, systemRoleCollection)
	BusinessRoleRecordHandler := businessRole.NewMongoRecordHandler(mainMongoSession, databaseName, businessRoleCollection)
	UserRecordHandler := user.NewMongoRecordHandler(mainMongoSession, databaseName, userCollection)
	ShipRecordHandler := ship.NewMongoRecordHandler(mainMongoSession, databaseName, shipCollection)
	BusinessDayConfigRecordHandler := businessDayConfig.NewMongoRecordHandler(
		mainMongoSession,
		databaseName,
		businessDayConfigCollection,
		)
	BusinessDayRecordHandler := businessDay.NewMongoRecordHandler(
		mainMongoSession,
		databaseName,
		businessDayCollection,
		BusinessDayConfigRecordHandler,
		BusinessRoleRecordHandler,
		employeeCollection,
		)
	EmployeeRecordHandler := employee.NewMongoRecordHandler(
		mainMongoSession,
		databaseName,
		employeeCollection,
		BusinessDayRecordHandler,
		)
	TagEventRecordHandler := tagEvent.NewMongoRecordHandler(
		mainMongoSession,
		databaseName,
		tagEventCollection,
		hub,
		EmployeeRecordHandler,
		BusinessDayRecordHandler,
		)




	// Create Services
	SystemRoleService := systemRole.NewService(SystemRoleRecordHandler)
	BusinessRoleService := businessRole.NewService(BusinessRoleRecordHandler)
	UserService := user.NewServiceAdaptor(UserRecordHandler)
	//AuthService := auth.NewService(UserRecordHandler, rsaPrivateKey)
	TagEventService := tagEvent.NewServiceAdaptor(TagEventRecordHandler)
	EmployeeService := employee.NewService(EmployeeRecordHandler)
	ShipService := ship.NewService(ShipRecordHandler)
	BusinessDayConfigService := businessDayConfig.NewService(BusinessDayConfigRecordHandler)
	BusinessDayService := businessDay.NewService(BusinessDayRecordHandler)


	// Initialise the APIAuthorizer
	mainAPIAuthorizer.JWTValidator = token.NewJWTValidator(&rsaPrivateKey.PublicKey)
	mainAPIAuthorizer.RoleRecordHandler = SystemRoleRecordHandler

	// Create rfIDAPIServer API Server
	rfIDAPIServer := rpc.NewServer()
	rfIDAPIServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	// Register Services on rfIDAPIServer API Server
	if err := rfIDAPIServer.RegisterService(TagEventService, "TagEvent"); err != nil {
		log.Fatal("Unable to Register Tag Event Service Adaptor")
	}


	// Create secureAPIServer
	secureAPIServer := rpc.NewServer()
	secureAPIServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	// Register Services with secureAPIServer
	if err := secureAPIServer.RegisterService(SystemRoleService, "SystemRole"); err != nil {
		log.Fatal("Unable to Register System Role Service")
	}
	if err := secureAPIServer.RegisterService(BusinessRoleService, "BusinessRole"); err != nil {
		log.Fatal("Unable to Register Business Role Service")
	}
	//if err:= secureAPIServer.RegisterService(AuthService, "Auth"); err != nil {
	//	log.Fatal("Unable to Register Auth Service Adaptor")
	//}
	if err := secureAPIServer.RegisterService(UserService, "User"); err != nil {
		log.Fatal("Unable to Register User Service")
	}
	if err := secureAPIServer.RegisterService(EmployeeService, "Employee"); err != nil {
		log.Fatal("Unable to Register Employee Service")
	}
	if err := secureAPIServer.RegisterService(ShipService, "Ship"); err != nil {
		log.Fatal("Unable to Register Ship Service")
	}
	if err := secureAPIServer.RegisterService(BusinessDayConfigService, "BusinessDayConfig"); err != nil {
		log.Fatal("Unable to Register BusinessDayConfig Service")
	}
	if err := secureAPIServer.RegisterService(BusinessDayService, "BusinessDay"); err != nil {
		log.Fatal("Unable to Register BusinessDayService Service")
	}

	// Set up Router for secureAPIServer
	secureAPIServerMux := mux.NewRouter()
	secureAPIServerMux.Methods("OPTIONS").HandlerFunc(preFlightHandler)
	// secureAPIServerMux.Handle("/api", apiAuthApplier(secureAPIServer)).Methods("POST")
	secureAPIServerMux.Handle("/api", secureAPIServer).Methods("POST")
	// Start secureAPIServer
	log.Info("Starting secureAPIServer on port " + ServerPort)
	go func() {
		err := http.ListenAndServe("0.0.0.0:" + ServerPort, secureAPIServerMux)
		log.Error("secureAPIServer stopped: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}()

	// Set up Router for rfIDAPIServer
	rfIDAPIServerMux := mux.NewRouter()
	rfIDAPIServerMux.Methods("OPTIONS").HandlerFunc(preFlightHandler)
	rfIDAPIServerMux.Handle("/api", rfIDAPIServer).Methods("POST")
	// Start rfIDAPIServer
	log.Info("Starting rfIDAPIServer on port " + RFIDServerPort)
	go func() {
		err := http.ListenAndServe("0.0.0.0:" + RFIDServerPort, rfIDAPIServerMux)
		log.Error("rfIDAPIServer stopped: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}()


	// Starting websocket connection hub
	go hub.Run()

	// Set up Router for websSocketServer
	webSocketServerMux := mux.NewRouter()
	webSocketServerMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { websocketServer.ServeWs(w, r, hub)})
	webSocketServerMux.HandleFunc("/healthcheck", healthcheck()).Methods("GET")
	// Start 
	log.Info("Starting webSocketServer on port " + webSocketServerPort)
	go func() {
		err := http.ListenAndServe("0.0.0.0:" + webSocketServerPort, webSocketServerMux)
		log.Error("webSocketServer server stopped: ", err, "\n", string(debug.Stack()))
		os.Exit(1)
	}()

	time.Sleep(1 * time.Second)

	// Set Up rfID Context Provider
	rfIDContextProvider := rfId.NewContextProvider(hub)
	if err != nil {
		log.Fatal("unable to create new context provider: " + err.Error())
	}

	// Register Context Provider with hub
	hub.Register <- rfIDContextProvider

	// Register Tag Event Record Handler with the hub
	hub.Register <- TagEventRecordHandler

	go func() {
		if err := rfIDContextProvider.Run(); err != nil {
			log.Error("Context Provider has stopped with error: " + err.Error())
		}
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
