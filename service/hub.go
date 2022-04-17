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
			fmt.Printf("Registered connections: %#v\n", h.connections)
		case c := <-h.unregister:
			fmt.Println("Connection unregistered: ", c.SKI)
			if _, ok := h.connections[c.SKI]; ok {
				delete(h.connections, c.SKI)
				fmt.Println("Connection unregistered: ", c.SKI)
			}
			fmt.Printf("Remaining connections: %#v\n", h.connections)
		}
	}
}

func (h *ConnectionsHub) ConnectionForSKI(ski string) *ConnectionHandler {
	return h.connections[ski]
}
