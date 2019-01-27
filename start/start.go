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
	"gitlab.com/iotTracker/brain/user"
	"gopkg.in/mgo.v2"
	"time"
	"gitlab.com/iotTracker/brain/security/encrypt"
	"gitlab.com/iotTracker/brain/security/systemRole"
	"gitlab.com/iotTracker/brain/security/apiAuth"
	"gitlab.com/iotTracker/brain/security/token"

	"gitlab.com/iotTracker/brain/security/auth"
)

var ServerPort = "9006"

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

	// Create Record Handlers
	SystemRoleRecordHandler := systemRole.NewMongoRecordHandler(mainMongoSession, databaseName, systemRoleCollection)
	UserRecordHandler := user.NewMongoRecordHandler(mainMongoSession, databaseName, userCollection)

	// Create Services
	SystemRoleService := systemRole.NewService(SystemRoleRecordHandler)
	UserService := user.NewServiceAdaptor(UserRecordHandler)
	AuthService := auth.NewService(UserRecordHandler, rsaPrivateKey)

	// Initialise the APIAuthorizer
	mainAPIAuthorizer.JWTValidator = token.NewJWTValidator(&rsaPrivateKey.PublicKey)
	mainAPIAuthorizer.RoleRecordHandler = SystemRoleRecordHandler

	// Create secureAPIServer
	secureAPIServer := rpc.NewServer()
	secureAPIServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	// Register Services with secureAPIServer
	if err := secureAPIServer.RegisterService(SystemRoleService, "SystemRole"); err != nil {
		log.Fatal("Unable to Register System Role Service")
	}
	if err:= secureAPIServer.RegisterService(AuthService, "Auth"); err != nil {
		log.Fatal("Unable to Register Auth Service Adaptor")
	}
	if err := secureAPIServer.RegisterService(UserService, "User"); err != nil {
		log.Fatal("Unable to Register User Service")
	}

	// Set up Router for secureAPIServer
	secureAPIServerMux := mux.NewRouter()
	secureAPIServerMux.Methods("OPTIONS").HandlerFunc(preFlightHandler)
	secureAPIServerMux.Handle("/api", apiAuthApplier(secureAPIServer)).Methods("POST")
	//secureAPIServerMux.Handle("/api", secureAPIServer).Methods("POST")
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
