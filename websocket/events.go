package websocket

import (
	"encoding/json"
	"fmt"
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
)

//broadcast message
type SendMessageEvent struct {
	Message string `json:"message"`
	From string `json:"from"`
}


//types of event handlers
//event functions
func SendMessage(event Event, c *WebsocketClient) error {
	fmt.Println(event)
	c.manager.Broadcast(event)
	return nil
}