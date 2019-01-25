package rfId

import (
	"fmt"
	"github.com/gorilla/websocket"
	"errors"
	"bitbucket.org/gotimekeeper/exoWSC"
	"bitbucket.org/gotimekeeper/log"
	"time"
	"bitbucket.org/gotimekeeper/exoWSC/message"
	"encoding/json"
)

func NewContextProvider(
	Hub *exoWSC.Hub,
	) (*contextProvider) {

	return &contextProvider{
		Hub: Hub,
		MsgFromHub: make(chan exoWSC.Message, 1000),
	}
}

type contextProvider struct {
	// The websocket contextProvider
	Conn *websocket.Conn
	// Buffered channel of outbound messages
	MsgFromHub chan exoWSC.Message
	Stop      chan bool
	// hub
	Hub *exoWSC.Hub
}

func (cp *contextProvider) Send(message exoWSC.Message) error {
	sendTimeOutTicker := time.NewTicker(2 *time.Second)
	defer func () {
		sendTimeOutTicker.Stop()
	}()

	select {
	case cp.MsgFromHub <- message:
	case <- sendTimeOutTicker.C:
		log.Error("Time out on waiting to get message into context provider's MsgFromHub Channel. Msg: ", message)
	}

	return nil
}

func (cp *contextProvider) Run() error {
	log.Info("Context Provider is starting")
	for {
		select {
		case message, ok := <-cp.MsgFromHub:
			if !ok {
				log.Error("somehow the context provider's internal MsgFromHub Channel has closed")
				return errors.New("somehow the context provider's internal MsgFromHub Channel has closed")
			}
			cp.RespondToMessage(&message)

		case <- cp.Stop:
			fmt.Println("context Provider is stopping")
			return nil
		}

	}
}

type GetServiceContextRequest struct {
	Event string `json:"event"`
}

type GetServiceContextResponse struct {
	Context string `json:"context"`
}

func (cp *contextProvider) RespondToMessage(rxedMsg *exoWSC.Message) error {
	switch rxedMsg.Type {
	case message.CreateServiceContextRequest:
		//fmt.Println("CreateServiceContextRequest")
	case message.CreateServiceContextResponse:
		// Cannot Transmit Back here, will start inf. loop
	case message.GetServiceContextRequest:
		responseMsgData, err := json.Marshal(GetServiceContextResponse{
			//Context: "assignToEmployee",
			Context: "default",
		})
		if err != nil {
			log.Error("error marshalling rfid service context request: " + err.Error())
			return errors.New("error marshalling rfid service context request: " + err.Error())
		}
		responseMsg := exoWSC.Message{
			Type: message.GetServiceContextResponse,
			SerialData: string(responseMsgData[:]),
		}

		timeOutTicker := time.NewTicker(1 * time.Second)
		defer func(){
			timeOutTicker.Stop()
		}()
		select {
		case cp.Hub.Broadcast <- responseMsg:
			break
		case <- timeOutTicker.C:
			log.Error("time out trying to put message into hub's broadcast channel")
			return errors.New("time out trying to put message into hub's broadcast channel")
		}
	default:
	}
	return nil
}
