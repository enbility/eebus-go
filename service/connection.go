package service

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
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
	Role ShipRole

	// The remote SKI
	SKI string

	// The ship ID
	ShipID string

	// The connection hub handling all service connections
	ConnectionsHub *ConnectionsHub

	// The actual websocket connection
	conn *websocket.Conn

	// The read channel for incoming messages
	readChannel chan []byte

	// The shutdown channel
	shutdownChannel chan struct{}

	shutdownMux sync.Mutex

	// Indicates wether the ship handshake has been completed
	shipHandshakeComplete bool
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
	c.shutdownChannel = make(chan struct{})

	go c.readPump()

	go func() {
		// perform ship handshake
		if err := c.shipHandshake(); err != nil {
			fmt.Println("SHIP handshake error: ", err)
			c.shutdown()
		}
		if !c.shipHandshakeComplete {
			return
		}

		for {
			select {
			case <-c.shutdownChannel:
				return
			case message := <-c.readChannel:
				_, _ = c.parseMessage(message, true)
			}
		}
	}()
}

// shutdown the connection and all internals
func (c *ConnectionHandler) shutdown() {
	c.shutdownMux.Lock()
	defer c.shutdownMux.Unlock()

	fmt.Println("shutting down connection ", c.SKI)

	c.shutdownChannel <- struct{}{}

	if c.conn != nil {
		// close the SHIP connection according to the SHIP protocol
		fmt.Println("closing SHIP")
		c.shipClose()

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

	_ = c.sendModel(closeMessage)
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
		select {
		case <-c.shutdownChannel:
			return
		default:
			message, err := c.readWebsocketMessage()
			if err != nil {
				fmt.Println("Error reading message: ", err)
				break
			}

			c.readChannel <- message
		}
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
