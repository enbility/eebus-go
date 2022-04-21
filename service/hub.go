package service

import (
	"fmt"
	"sync"
)

type connectionsHub struct {
	connections map[string]*ConnectionHandler

	// Register reuqests from a new connection
	register chan *ConnectionHandler

	// Unregister requests from a closing connection
	unregister chan *ConnectionHandler

	localService ServiceDetails

	mux sync.Mutex
}

func newConnectionsHub(localService ServiceDetails) *connectionsHub {
	return &connectionsHub{
		connections:  make(map[string]*ConnectionHandler),
		register:     make(chan *ConnectionHandler),
		unregister:   make(chan *ConnectionHandler),
		localService: localService,
	}
}

func (h *connectionsHub) run() {
	for {
		select {
		case c := <-h.register:

			// SHIP 12.2.2 recommends that the connection initiated with the higher SKI should retain the connection
			existingC := h.connectionForSKI(c.remoteService.SKI)
			if existingC != nil {
				fmt.Println("Connection already exists for SKI: ", c.remoteService.SKI)

				// If the connection is initiated by the local service and the local SKI is higher than the remote SKI
				// then the existing connection should be closed
				if c.isConnectedFromLocalService && c.localService.SKI < c.remoteService.SKI {
					c.conn.Close()
					continue
				} else {
					if existingC.conn != nil {
						existingC.conn.Close()
					}

					h.mux.Lock()
					delete(h.connections, c.remoteService.SKI)
					h.mux.Unlock()
				}
			}

			h.mux.Lock()
			h.connections[c.remoteService.SKI] = c
			h.mux.Unlock()

			c.handleConnection()

			// TODO: shutdown mDNS if this is not a CEM and all registered services are connected
		case c := <-h.unregister:
			if chRegistered, ok := h.connections[c.remoteService.SKI]; ok {
				if chRegistered.conn == c.conn {
					h.mux.Lock()
					delete(h.connections, c.remoteService.SKI)
					h.mux.Unlock()
				}
			}
			// TODO: startup mDNS if this is not a CEM and no registered service is connected
		}
	}
}

// return the connection for a specific SKI
func (h *connectionsHub) connectionForSKI(ski string) *ConnectionHandler {
	h.mux.Lock()
	defer h.mux.Unlock()

	return h.connections[ski]
}

// close all connections
func (h *connectionsHub) shutdown() {
	for _, c := range h.connections {
		c.shutdown(true)
	}
}

// return if there is a connection for a SKI
func (h *connectionsHub) isSkiConnected(ski string) bool {
	h.mux.Lock()
	defer h.mux.Unlock()

	// The connection with the higher SKI should retain the connection
	_, ok := h.connections[ski]
	return ok
}
