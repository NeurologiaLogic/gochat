package websocket

import (
	"encoding/json"
	"fmt"
	"time"
)

//event structure
type Event struct{
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

//contract
type EventHandler func(event Event, client *WebsocketClient) error

//constants
const (
	EventSendMessage = "send_message"
	EventChangeRoom = "change_room"
)

//types of event handlers

//Message Received from payload
type SendMessageEvent struct {
	Message string `json:"message"`
	From string `json:"from"`
}

//MessageResponse from Server to be broadcast
type NewMessageEvent struct{
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

func SendMessageHandler(event Event, c *WebsocketClient) error {
	fmt.Println(event)
	c.manager.Broadcast(event)
	return nil
}


type ChangeRoom struct {
	Name string `json:"name"`
}

func ChangeRoomHandler(event Event, c *WebsocketClient) error {

	return nil
}