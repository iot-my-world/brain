package exoWSC

import (
	"encoding/json"
	"fmt"
	"github.com/iot-my-world/brain/exoWSC/message"
	"github.com/iot-my-world/brain/log"
)

type Hub struct {
	/*
	   A central hub will receive all incoming messages and broadcast them
	   to all registered "Subscriber"s
	   (i.e. the Subscriber structures in the clients map)
	*/
	Clients    map[Subscriber]bool
	Broadcast  chan Message
	Register   chan Subscriber
	Unregister chan Subscriber
	content    Message
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan Subscriber),
		Unregister: make(chan Subscriber),
		Clients:    make(map[Subscriber]bool),
	}
}

func (h *Hub) Run() {
	/*
	   Channels are like FIFO Stacks
	   A Subscriber will store a message in one of the 3 channels
	   A go routine will unstack them as soon as possible by arrival
	   time.
	*/
	log.Info("Starting websocket connection hub")
	for {
		select {
		case c := <-h.Register:
			log.Info("Subscriber Registered to hub")
			h.Clients[c] = true

			// Build message data to send welcome
			messageByteData, err := json.Marshal(struct {
				Msg string `json:"msg"`
			}{
				Msg: "Welcome to the hub!",
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
			break

		case c := <-h.Unregister:
			_, ok := h.Clients[c]
			if ok {
				delete(h.Clients, c)
			}
			break

		case m := <-h.Broadcast:
			h.content = m
			h.broadcastMessage()
			break
		}
	}
}

func (h *Hub) broadcastMessage() {
	fmt.Println("hub is broadcasting message:", h.content)
	for c := range h.Clients {
		if err := c.Send(h.content); err != nil {
			log.Error("Error Sending msg to client: " + err.Error())
			// TODO: Unregister Client here?
		}
	}

}
