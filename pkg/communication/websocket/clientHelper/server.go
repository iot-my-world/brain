package clientHelper

import (
	"github.com/gorilla/websocket"
	"github.com/iot-my-world/brain/internal/log"
	websocket2 "github.com/iot-my-world/brain/pkg/communication/websocket"
	"net/http"
)

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, hub *websocket2.Hub) {
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
	newClientHelper := websocket2.NewClientHelper(conn, hub)

	log.Info("Registering Client With Hub")
	// Register ClientHelper with hub
	hub.Register <- newClientHelper

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines
	go newClientHelper.HandleRX()
	go newClientHelper.HandleTX()
}
