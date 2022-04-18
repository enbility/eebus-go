package service

import (
	"crypto/x509"
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
	Role           ShipRole
	SKI            string
	ConnectionsHub *ConnectionsHub

	conn *websocket.Conn

	readChannel chan []byte
}

// Connection handler when the service initiates a connection to a remote service
func (c *ConnectionHandler) handleConnection(conn *websocket.Conn) {
	c.conn = conn

	if len(c.SKI) == 0 {
		fmt.Println("SKI is not set")
		c.shutdown()
		return
	}

	if c.ConnectionsHub.ConnectionForSKI(c.SKI) != nil {
		fmt.Println("Client with SKI already connected")
		c.shutdown()
		return
	}

	c.ConnectionsHub.register <- c

	c.startup()
}

func (c *ConnectionHandler) skiFromX509Certificate(cert *x509.Certificate) (string, error) {
	if len(cert.SubjectKeyId) == 0 {
		return "", errors.New("Client certificate does not provide a SKI")
	}

	return fmt.Sprintf("%0x", cert.SubjectKeyId), nil
}

func (c *ConnectionHandler) startup() {
	c.readChannel = make(chan []byte, 1) // Listen to incoming websocket messages

	go c.readPump()

	go func() {
		// perform ship handshake
		if err := c.shipHandshake(); err != nil {
			fmt.Println("SHIP handshake error: ", err)
			c.shutdown()
		}
	}()
}

// shutdown the connection and all internals
func (c *ConnectionHandler) shutdown() {
	fmt.Println("shutting down connection ", c.SKI)

	if c.conn != nil {
		fmt.Println("closing websocket connection")
		c.conn.Close()
	}

	if !isChannelClosed(c.readChannel) {
		fmt.Println("closing read channel")
		close(c.readChannel)
	}

	fmt.Println("unregistering connection")
	c.ConnectionsHub.unregister <- c
}

func isChannelClosed[T any](ch <-chan T) bool {
	select {
	case <-ch:
		return false
	default:
		return true
	}
}

// readPump pumps messages from the websocket connection into the read message channel
func (c *ConnectionHandler) readPump() {
	defer func() {
		c.shutdown()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		message, err := c.readWebsocketMessage()
		if err != nil {
			fmt.Println("Error reading message: ", err)
			break
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
		fmt.Println("ReadMessage error: ", err)
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

	c.conn.SetWriteDeadline((time.Now().Add(writeWait)))

	err := c.conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		fmt.Println("WriteMessage error: ", err)
		c.shutdown()
		return err
	}

	return nil
}

// send a json message for a provided model to the websocket connection
func (c *ConnectionHandler) sendModel(model interface{}) error {
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
	shipMsg := []byte{ship.MsgTypeControl}
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
