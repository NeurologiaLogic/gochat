package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//type
type WebsocketManager struct{
	clients map[*WebsocketClient]bool
	upgrader websocket.Upgrader
	sync.Mutex
	eventHandlers map[string]EventHandler
}

//factory
func NewWebsocketManager() *WebsocketManager {
	m:= &WebsocketManager{
				clients: make(map[*WebsocketClient]bool),
				eventHandlers: make(map[string]EventHandler),
				upgrader:websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				},
		}
	m.applyEventHandlers()
	return m
}

//event handlers
func (ws *WebsocketManager) applyEventHandlers(){
	ws.eventHandlers[EventSendMessage] = SendMessageHandler
	ws.eventHandlers[EventChangeRoom] = ChangeRoomHandler

}

//event routing
func (ws *WebsocketManager) routeEvents(event Event, c *WebsocketClient) error {
	handler, ok := ws.eventHandlers[event.Type]
	if !ok {
		fmt.Println("No event handler for event type %s", event.Type)
		//think how to return a nil
		return nil
	}
	//executing the handler
	if err := handler(event, c); err != nil {
		return err
	}
	return nil
}


//http handler
func (ws *WebsocketManager) Handler(c *gin.Context){
	fmt.Println("New Connection")
	ws.upgrader.CheckOrigin = func (r *http.Request) bool {return true}
	conn, err := ws.upgrader.Upgrade(c.Writer, c.Request, nil);
	if err != nil {
		fmt.Println(err)
	}
	client := NewWebsocketClient(conn,ws)
	ws.addClient(client)
	//reader and writer process
	go client.readMessages()
	go client.writeMessages()
}

//helper function
func (ws *WebsocketManager) Broadcast(event Event,c *WebsocketClient){
	for client := range ws.clients{
		if client.chatroom == c.chatroom{
			client.messageQueue <- event
		}
	}
}

func (ws *WebsocketManager) addClient(client *WebsocketClient){
	ws.Lock()
	defer ws.Unlock()
	ws.clients[client] = true
}
func (ws *WebsocketManager) removeClient(client *WebsocketClient){
	ws.Lock()
	defer ws.Unlock()
	//check if exist
	if _,ok := ws.clients[client];ok{
		client.conn.Close()
		delete(ws.clients,client)
	}
}
