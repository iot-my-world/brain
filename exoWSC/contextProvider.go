package exoWSC

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.com/iotTracker/brain/log"
	"net/url"
	"time"
)

func NewContextProvider(
	Host string,
	Port string,
	Path string,
) *ContextProvider {

	// Create New Websocket ContextProvider
	return &ContextProvider{
		Host: Host,
		Port: Port,
		Path: Path,
		Send: make(chan []byte, 1000),
	}
}

type ContextProvider struct {
	// The websocket ContextProvider
	Conn *websocket.Conn
	// Buffered channel of outbound messages
	Send  chan []byte
	Host  string
	Port  string
	Path  string
	Close chan bool
}

func (cp *ContextProvider) Run() error {

	// Build websocket url
	wsUrl := url.URL{Scheme: "ws", Host: cp.Host + ":" + cp.Port, Path: cp.Path}

	log.Info("Starting Context Provider Connecting to: " + wsUrl.String())

	// Connect websocket
	var err error
	cp.Conn, _, err = websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		return errors.New("WSC error while connecting: " + err.Error())
	}
	defer cp.Conn.Close()

	// Set necessary parameters
	cp.Conn.SetReadLimit(MaxMessageSize)
	cp.Conn.SetReadDeadline(time.Now().Add(PongWait))
	cp.Conn.SetPongHandler(func(string) error {
		//fmt.Println("Client got pong back")
		cp.Conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	go func() {
		//fmt.Println("Start rx")
		err := cp.StartRX()
		fmt.Println("StartRX error ", err)
		cp.Close <- true
	}()

	go func() {
		//fmt.Println("start tx")
		err := cp.StartTX()
		fmt.Println("StartTX error ", err)
		cp.Close <- true
	}()

	for {
		select {
		case <-cp.Close:
			fmt.Println("websocket closing...")
			return errors.New("some good reason for closing...")
		}
	}
}

func (cp *ContextProvider) StartRX() error {
	for {
		_, msgByteData, err := cp.Conn.ReadMessage()
		if err != nil {
			return errors.New("WSC read error: " + err.Error())
		}

		rxedMsg := Message{}
		if err := json.Unmarshal(msgByteData, &rxedMsg); err != nil {
			log.Error("error unmarshalling rxed WS Message: " + err.Error())
		}
		switch rxedMsg.Type {
		default:
			fmt.Println("Default Case!", rxedMsg)
		}
	}
	return nil
}

func (cp *ContextProvider) StartTX() error {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-cp.Send:
			if ok {
				cp.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
				fmt.Println("Need to send message!", message)
				cp.Conn.WriteMessage(websocket.TextMessage, message)
			} else {
				return errors.New("send channel is closed")
			}
		case <-ticker.C:
			cp.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			//fmt.Println("Ping the server!")
			if err := cp.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return errors.New("Could not send ping: " + err.Error())
			}
		}
	}
	return nil
}
