package clientHelper

import (
	"net/http"
	"github.com/gorilla/websocket"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/exoWSC"
)

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, hub *exoWSC.Hub) {
	log.Info("New Websocket Client Connected")

	//Allow any origin to connect
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	//Upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Could not upgrade ws connection: ", err)
		w.WriteHeader(500)
		return
	}

	//Construct a new client helper
	newClientHelper := exoWSC.NewClientHelper(conn, hub)

	log.Info("Registering Client With Hub")
	// Register ClientHelper with hub
	hub.Register <- newClientHelper


	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines
	go newClientHelper.HandleRX()
	go newClientHelper.HandleTX()
}
