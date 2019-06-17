package exoWSC

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gitlab.com/iotTracker/brain/exoWSC/message"
	"gitlab.com/iotTracker/brain/log"
	"time"
)

func NewClientHelper(
	Conn *websocket.Conn,
	Hub *Hub,
) *clientHelper {

	return &clientHelper{
		Conn:      Conn,
		MsgToSend: make(chan Message),
		Hub:       Hub,
	}
}

type clientHelper struct {
	// The websocket connection.
	Conn *websocket.Conn
	// Buffered channel of outbound messages.
	MsgToSend chan Message
	// hub
	Hub *Hub
}

func (c *clientHelper) Send(message Message) error {
	sendTimeOutTicker := time.NewTicker(2 * time.Second)
	defer func() {
		sendTimeOutTicker.Stop()
	}()

	select {
	case c.MsgToSend <- message:
	case <-sendTimeOutTicker.C:
		log.Error("Time out on waiting to get message into Client Helper's MsgToSend Channel. Msg: ", message)
	}

	return nil
}

func (c *clientHelper) HandleRX() {
	defer func() {
		//unregister clients here
		log.Info("wsClientReader Connection Closed")
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	c.Conn.SetPongHandler(func(string) error {
		//fmt.Println("Server got pong back")
		c.Conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	// Build message data to send welcome
	messageByteData, err := json.Marshal(struct {
		Msg string `json:"msg"`
	}{
		Msg: "Connected to websocket client helper",
	})
	if err != nil {
		log.Warn("Unable to marshal welcome message data for client")
	}
	messageData := string(messageByteData[:])

	welcomeMessage := Message{
		Type:       message.WelcomeMessage,
		SerialData: messageData,
	}

	if err == nil {
		c.Send(welcomeMessage)
	} else {
		log.Warn("Unable to marshal message for client")
	}

	for {
		_, rawMsgData, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Warn("clientHelper closed ws connection unexpectedly: ", err)
			}
			log.Debug("clientHelper closed ws connection err: ", err)
			break
		}

		receivedMsg := Message{}
		if err := json.Unmarshal(rawMsgData, &receivedMsg); err != nil {
			log.Error("Error unmarshalling received message: " + err.Error())
			continue
		}
		c.Hub.Broadcast <- receivedMsg
	}
}

func (c *clientHelper) HandleTX() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		c.Conn.Close()
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-c.MsgToSend:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				// The hub closed the channel.
				log.Debug("The hub closed the channel")
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			rawMsgData, err := json.Marshal(message)
			if err != nil {
				log.Error("Unable to marshall ws message to send to client")
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, rawMsgData); err != nil {
				log.Warn("Could not write to websocket client: ", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			//fmt.Println("Ping the Client!")
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Warn("Could not send ping: ", err)
				return
			}
		}
	}
}
