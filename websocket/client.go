package websocket

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/gorilla/websocket"
)


var (
	//how long to wait for a pong response
	pongWait = 60 * time.Second
	//how often to ping the client
	pingInterval = (pongWait * 9) / 10
)

type WebsocketClient struct{
	conn *websocket.Conn;
	manager *WebsocketManager;
	//messageQueue to avoid concurrent writing
	messageQueue chan Event;
}

func NewWebsocketClient(conn *websocket.Conn,ws *WebsocketManager) *WebsocketClient{
	return &WebsocketClient{
		conn:conn,
		manager:ws,
		messageQueue:make(chan Event),
	}
}

//gorilla websocket only allows one writer at a time so use a unbuffered channel
//reading
func (c *WebsocketClient) readMessages(){
	defer func(){
		c.conn.Close()
		c.manager.removeClient(c)
	}()
	if err:=c.conn.SetReadDeadline(time.Now().Add(pongWait));err != nil{
		fmt.Println(err)
		return;
	}
	c.conn.SetPongHandler(c.pongHandler)
	c.conn.SetReadLimit(1024)
	for{
		_,payload,err := c.conn.ReadMessage()
		if err != nil {
			//what type of error
			if websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				fmt.Println("error reading messages")
			}
			//will trigger the defer
			break;
		}
		var request Event
		//parsing the payload
		if err := json.Unmarshal(payload,&request); err != nil{
			fmt.Println("error unmarshalling")
			break;
		}
		//routing the events
		if err:= c.manager.routeEvents(request,c);err!=nil{
			fmt.Println("error handling event")
			break;
		}
	}
}

//writing
func (c *WebsocketClient) writeMessages(){
	defer func(){
		c.conn.Close()
		c.manager.removeClient(c)
	}()
	//triger the ping
	ticker := time.NewTicker(pingInterval)
	//only write one message at a time and taken from the channel like a queue
	for{
		select {
		case message,ok := <- c.messageQueue:
			if !ok{
				if err:=c.conn.WriteMessage(websocket.CloseMessage, nil);err != nil{
					fmt.Println("Connection Closed",err)
				}
				return;
			}
			data,err := json.Marshal(message)
			if err != nil{
				fmt.Println("Error marshalling message",err)
			}
			if err:=c.conn.WriteMessage(websocket.TextMessage,data);err != nil{
				fmt.Println("Error writing message",err)
			}
			fmt.Println("message sent!")

		case <- ticker.C:
			fmt.Println("ping")
			if err:=c.conn.WriteMessage(websocket.PingMessage, []byte{});err != nil{
				fmt.Println("Error writing ping",err)
				return;
			}
		}
	}
}

//helper
func (c *WebsocketClient) pongHandler(string) error{
	fmt.Println("pong")
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}