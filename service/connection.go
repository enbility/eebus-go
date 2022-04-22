package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/DerAndereAndi/eebus-go/service/util"
	"github.com/DerAndereAndi/eebus-go/ship"
	"github.com/gorilla/websocket"
)

type ShipRole uint

const (
	writeWait               = 10 * time.Second
	cmiTimeout              = 10 * time.Second // SHIP 4.2
	cmiCloseTimeout         = 100 * time.Millisecond
	tHelloInit              = 60 * time.Second
	tHelloProlongThrInc     = 30 * time.Second
	tHelloProlongWaitingGap = 15 * time.Second
	tHellogProlongMin       = 1 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second // SHIP 4.2: ping interval + pong timeout
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 60 * time.Second // SHIP 4.2: ping interval

	// SHIP 9.2: Set maximum fragment length to 1024 bytes
	maxMessageSize = 1024

	ShipRoleServer ShipRole = 0
	ShipRoleClient ShipRole = 1
)

// A ConnectionHandler handles websocket connections.
type ConnectionHandler struct {
	// The ship connection mode of this connection
	role ShipRole

	// The remote service
	remoteService *ServiceDetails

	// The local service
	localService *ServiceDetails

	// The connection hub handling all service connections
	connectionsHub *connectionsHub

	// The actual websocket connection
	conn *websocket.Conn

	// Is this connection initiated by the local service
	isConnectedFromLocalService bool

	// The read channel for incoming messages
	readChannel chan []byte

	// The write channel for outgoing messages
	writeChannel chan []byte

	// Indicates wether the ship handshake has been completed
	shipHandshakeComplete bool
}

// Connection handler when the service initiates a connection to a remote service
func (c *ConnectionHandler) handleConnection() {
	if len(c.remoteService.SKI) == 0 {
		fmt.Println("SKI is not set")
		c.conn.Close()
		return
	}

	c.startup()
}

func (c *ConnectionHandler) startup() {
	c.readChannel = make(chan []byte, 1)  // Listen to incoming websocket messages
	c.writeChannel = make(chan []byte, 1) // Send outgoing websocket messages

	fmt.Println("Open pumps")
	go c.readPump()
	go c.writePump()

	go func() {
		if err := c.shipHandshake(); err != nil {
			fmt.Println("SHIP handshake error: ", err)
			c.shutdown(false)
			return
		}
	}()
}

// shutdown the connection and all internals
// may only invoked after startup() is invoked!
func (c *ConnectionHandler) shutdown(safeShutdown bool) {
	c.connectionsHub.unregister <- c

	if !isChannelClosed(c.readChannel) {
		close(c.readChannel)
	}

	if !isChannelClosed(c.writeChannel) {
		close(c.writeChannel)
	}

	if c.conn != nil {
		if c.shipHandshakeComplete && safeShutdown {
			// close the SHIP connection according to the SHIP protocol
			c.shipClose()
		}

		c.conn.Close()
	}

}

func (c *ConnectionHandler) shipClose() {
	if c.conn == nil {
		return
	}

	// SHIP 13.4.7: Connection Termination
	closeMessage := ship.ConnectionClose{
		ConnectionClose: ship.ConnectionCloseType{
			Phase: ship.ConnectionClosePhaseTypeAnnounce,
		},
	}

	_ = c.sendModel(ship.MsgTypeControl, closeMessage)
}

func isChannelClosed[T any](ch <-chan T) bool {
	select {
	case <-ch:
		return false
	default:
		return true
	}
}

// writePump pumps messages from the writeChannel to the websocket connection
func (c *ConnectionHandler) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.writeChannel:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The write channel has been closed
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection into the read message channel
func (c *ConnectionHandler) readPump() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		message, err := c.readWebsocketMessage()
		if err != nil {
			fmt.Println("Error reading message: ", err)
			c.shutdown(false)
			return
		}

		c.readChannel <- message
	}
}

// read a message from the websocket connection
func (c *ConnectionHandler) readWebsocketMessage() ([]byte, error) {
	if c.conn == nil {
		return nil, errors.New("Connection is not initialized")
	}

	msgType, b, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	if msgType != websocket.BinaryMessage {
		return nil, errors.New("Message is not a binary message")
	}

	if len(b) < 2 {
		return nil, fmt.Errorf("Invalid ship message length")
	}

	return b, nil
}

// write a message to the websocket connection
func (c *ConnectionHandler) writeWebsocketMessage(message []byte) error {
	if c.conn == nil {
		return errors.New("Connection is not initialized")
	}

	c.writeChannel <- message

	return nil
}

// send a json message for a provided model to the websocket connection
func (c *ConnectionHandler) sendModel(typ byte, model interface{}) error {
	msg, err := json.Marshal(model)
	if err != nil {
		return err
	}

	eebusMsg, err := util.JsonIntoEEBUSJson(msg)
	if err != nil {
		return err
	}

	fmt.Println("Send: ", string(eebusMsg))

	// Wrap the message into a binary message with the ship header
	shipMsg := []byte{typ}
	shipMsg = append(shipMsg, eebusMsg...)

	err = c.writeWebsocketMessage(shipMsg)
	if err != nil {
		return err
	}

	return nil
}

// enable jsonFormat if the return message is expected to be encoded in
// the eebus json format
// return the SHIP message type, the SHIP message and an error
func (c *ConnectionHandler) parseMessage(msg []byte, jsonFormat bool) (byte, []byte) {
	// Extract the SHIP header byte
	shipHeaderByte := msg[0]
	// remove the SHIP header byte from the message
	msg = msg[1:]

	if len(msg) > 1 {
		fmt.Println("Recv: ", string(msg))
	}

	if jsonFormat {
		return shipHeaderByte, util.JsonFromEEBUSJson(msg)
	}

	return shipHeaderByte, msg
}

// read the next message from the websocket connection and
// return an error if the provided timeout is reached
func (c *ConnectionHandler) readNextMessage(duration time.Duration) ([]byte, error) {
	timeout := time.NewTimer(duration)

	select {
	case <-timeout.C:
		return nil, errors.New("Timeout waiting for message")
	case msg := <-c.readChannel:
		return msg, nil
	}
}
