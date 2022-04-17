package service

import "fmt"

type ConnectionsHub struct {
	connections map[string]*ConnectionHandler

	// Register reuqests from a new connection
	register chan *ConnectionHandler

	// Unregister requests from a closing connection
	unregister chan *ConnectionHandler
}

func newConnectionsHub() *ConnectionsHub {
	return &ConnectionsHub{
		connections: make(map[string]*ConnectionHandler),
		register:    make(chan *ConnectionHandler),
		unregister:  make(chan *ConnectionHandler),
	}
}

func (h *ConnectionsHub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c.SKI] = c
			fmt.Println("Connection registered: ", c.SKI)
		case c := <-h.unregister:
			if _, ok := h.connections[c.SKI]; ok {
				delete(h.connections, c.SKI)
				fmt.Println("Connection unregistered: ", c.SKI)
			}
		}
	}
}

func (h *ConnectionsHub) ConnectionForSKI(ski string) *ConnectionHandler {
	return h.connections[ski]
}
