package svc

import (
	"github.com/hongyuxuan/lizardcd/server/internal/types"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*WsClient]bool

	// Inbound messages from the clients.
	broadcast chan types.WsMessage

	// Register requests from the clients.
	register chan *WsClient

	// Unregister requests from clients.
	unregister chan *WsClient
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan types.WsMessage),
		register:   make(chan *WsClient),
		unregister: make(chan *WsClient),
		clients:    make(map[*WsClient]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msg := <-h.broadcast:
			for client := range h.clients {
				if client.id == msg.MessageTo {
					select {
					case client.send <- msg:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
