package client

import (
	"bitbucket.org/gopiwsclient/log"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.com/iotTracker/brain/exoWSC"
	"net/url"
	"time"
)

func NewConnection(
	Host, Port, Path string,
) *Connection {
	return &Connection{
		Host: Host,
		Port: Port,
		Path: Path,
		Send: make(chan []byte, 1000),
	}
}

type Connection struct {
	// The websocket Connection
	Conn *websocket.Conn
	// Buffered channel of outbound messages
	Send      chan []byte
	Host      string
	Port      string
	Path      string
	Close     chan bool
	MessageRx chan []byte
}

func (c *Connection) Connect() error {
	// Build websocket url
	wsUrl := url.URL{Scheme: "ws", Host: c.Host + ":" + c.Port, Path: c.Path}

	log.Info("WS Connecting to: " + wsUrl.String())

	// Connect websocket
	var err error
	c.Conn, _, err = websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		return errors.New("WSC error while connecting: " + err.Error())
	}
	defer c.Conn.Close()

	// Set necessary parameters
	c.Conn.SetReadLimit(exoWSC.MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(exoWSC.PongWait))
	c.Conn.SetPongHandler(func(string) error {
		fmt.Println("Client got pong back")
		c.Conn.SetReadDeadline(time.Now().Add(exoWSC.PongWait))
		return nil
	})

	go func() {
		//fmt.Println("Start rx")
		err := c.StartRX()
		fmt.Println("StartRX error ", err)
		c.Close <- true
	}()

	go func() {
		//fmt.Println("start tx")
		err := c.StartTX()
		fmt.Println("StartTX error ", err)
		c.Close <- true
	}()

	for {
		select {
		case <-c.Close:
			fmt.Println("websocket closing...")
			return errors.New("some good reason for closing...")
		}
	}
}

func (c *Connection) StartRX() error {
	for {
		_, msgByteData, err := c.Conn.ReadMessage()
		if err != nil {
			return errors.New("WSC read error: " + err.Error())
		}
		//fmt.Printf("Recieved message: %s\n", msgByteData)
		c.MessageRx <- msgByteData
	}
	return nil
}

func (c *Connection) StartTX() error {
	ticker := time.NewTicker(exoWSC.PingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if ok {
				c.Conn.SetWriteDeadline(time.Now().Add(exoWSC.WriteWait))
				//fmt.Println("Need to send message!", message)
				c.Conn.WriteMessage(websocket.TextMessage, message)
			} else {
				return errors.New("send channel is closed")
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(exoWSC.WriteWait))
			//fmt.Println("Ping the server!")
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return errors.New("Could not send ping: " + err.Error())
			}
		}
	}
	return nil
}
