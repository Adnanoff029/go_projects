package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	WebSocketUpgrader = websocket.Upgrader{
		CheckOrigin:     CheckOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	Clients ClientList
	sync.RWMutex
	Handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		Clients:  make(ClientList),
		Handlers: make(map[string]EventHandler),
	}
	m.SetupEventHandlers()
	return m
}

func (m *Manager) RouteEvent(event Event, c *Client) error {
	if handler, ok := m.Handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("No such event type present.")
	}
}

func (m *Manager) SetupEventHandlers() {
	m.Handlers[EventSendMessage] = SendMessage
}

func (m *Manager) ServerWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")
	conn, err := WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, m)
	m.AddClient(client)
	go client.ReadMessages()
	go client.WriteMessages()
}

func (m *Manager) AddClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.Clients[client] = true
}

func (m *Manager) RemoveClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.Clients[client]; ok {
		client.Connection.Close()
		delete(m.Clients, client)
	}
}

func SendMessage(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}

func CheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://localhost:3000":
		return true
	default:
		return false
	}
}
