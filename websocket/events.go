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
	var message SendMessageEvent
	if err := json.Unmarshal(event.Payload,&message); err != nil{
		fmt.Println("error unmarshalling message: ", err)
		return err
	}
	var broadcastMessage NewMessageEvent
	broadcastMessage.Message = message.Message
	broadcastMessage.From = message.From
	broadcastMessage.Sent = time.Now()
	data, err := json.Marshal(broadcastMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	//turn back into Event
	var serverResponse Event
	serverResponse.Payload = data
	serverResponse.Type = EventSendMessage
	// Broadcast to all other Clients
	c.manager.Broadcast(serverResponse,c)
	return nil
}


type ChangeRoomEvent struct {
	Name string `json:"name"`
}

func ChangeRoomHandler(event Event, c *WebsocketClient) error {
	var changeRoomEvent ChangeRoomEvent
	if err := json.Unmarshal(event.Payload,&changeRoomEvent); err != nil{
		return err
	}
	c.chatroom = changeRoomEvent.Name
	fmt.Println("now in room:",c.chatroom)
	return nil
}