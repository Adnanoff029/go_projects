package main

import (
	"encoding/json"
	"log"
	"time"
	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second
	// Lower than the pongWait
	pingInterval = (pongWait * 9) / 10 // 90% of the pongWait
)

type ClientList map[*Client]bool

type Client struct {
	Connection *websocket.Conn
	Manager    *Manager
	Egress     chan Event
}

func (c *Client) ReadMessages() {
	defer func() {
		c.Manager.RemoveClient(c)
	}()

	if err := c.Connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	c.Connection.SetReadLimit(512)
	c.Connection.SetPongHandler(c.PongHandler)

	for {
		_, payload, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message : %v", err)
			}
			break
		}

		// Broadcasting the message to all clients to check if Egress is working or not.
		// for wsclient := range c.Manager.Clients {
		// 	wsclient.Egress <- payload
		// }
		// log.Println(messageType)
		// log.Println(string(payload))

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("Error marsharling the event : %v", err)
		}
		if err := c.Manager.RouteEvent(request, c); err != nil {
			log.Println("Error handling the event :", err)
		}
	}
}

func (c *Client) WriteMessages() {
	defer func() {
		c.Manager.RemoveClient(c)
	}()
	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case message, ok := <-c.Egress:
			if !ok {
				if err := c.Connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection Closed:", err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println("Error Marshaling the message : ", err)
				return
			}
			if err := c.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("Failed to send message:", err)
			}
			log.Println("Message sent")
		case <-ticker.C:
			log.Println("Ping")
			if err := c.Connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("Write Message Error: ", err)
				return
			}
		}

	}
}

func (c *Client) PongHandler(pongMessage string) error {
	log.Println("Pong")
	return c.Connection.SetReadDeadline(time.Now().Add(pongWait))
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		Connection: conn,
		Manager:    manager,
		Egress:     make(chan Event),
	}
}
