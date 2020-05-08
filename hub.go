package main

import (
	"sync/atomic"
)

// Hub connection hub for managing clients
type Hub struct {
	count     int64            // number of clients connected
	clients   map[*Client]bool // map of connected clients
	broadcast chan []byte      // channel to broadcast messages to all clients
	add       chan *Client     // channel for adding clients
	remove    chan *Client     // channel for removing clients
}

// create new instance of hub
func newHub() *Hub {
	return &Hub{
		count:     0,
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
		add:       make(chan *Client),
		remove:    make(chan *Client),
	}
}

// run hub and manage clients
func (h *Hub) run() {
	for {
		select {
		// add new client and update counter
		case client := <-h.add:
			h.clients[client] = true
			atomic.AddInt64(&h.count, 1)
			// remove exiting client and update counter
		case client := <-h.remove:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				atomic.AddInt64(&h.count, -1)
			}
			// broadcast to all clients
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				// push message to client send channel
				case client.send <- message:
					// close channel when buffer is full
					// delete client and update counter
				default:
					close(client.send)
					delete(h.clients, client)
					atomic.AddInt64(&h.count, -1)
				}
			}
		}
	}
}
