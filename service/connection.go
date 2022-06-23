package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/service/util"
	"github.com/DerAndereAndi/eebus-go/ship"
	"github.com/gorilla/websocket"
)

type ShipRole string

const (
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second // SHIP 4.2: ping interval + pong timeout
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 50 * time.Second // SHIP 4.2: ping interval

	// SHIP 9.2: Set maximum fragment length to 1024 bytes
	maxMessageSize = 1024

	ShipRoleServer ShipRole = "server"
	ShipRoleClient ShipRole = "client"
)

type ConnectionHandlerDelegate interface {
	requestUserTrustForService(service *ServiceDetails)
	shipIDUpdateForService(service *ServiceDetails)

	// new spine connection established, inform SPINE
	addRemoteDeviceConnection(ski string, readC <-chan []byte, writeC chan<- []byte)

	// remove an existing connection from SPINE
	removeRemoteDeviceConnection(ski string)
}

// A ConnectionHandler handles websocket connections.
type ConnectionHandler struct {
	// The ship connection mode of this connection
	role ShipRole

	// The remote service
	remoteService *ServiceDetails

	// The local service
	localService *ServiceDetails

	// The actual websocket connection
	conn *websocket.Conn

	// The read channel for incoming messages
	readChannel chan []byte

	// The write channel for outgoing messages
	writeChannel chan []byte

	// The ship read channel for incoming messages
	shipReadChannel chan []byte

	// The ship write channel for outgoing messages
	shipWriteChannel chan []byte

	// The connection was closed
	closeChannel chan struct{}

	// The channel for handling local trust state update for the remote service
	shipTrustChannel chan bool

	// The current SHIP state
	smeState shipMessageExchangeState

	unregisterChannel  chan<- *ConnectionHandler
	connectionDelegate ConnectionHandlerDelegate

	// internal handling of closed connections
	isConnectionClosed bool

	mux sync.Mutex
}

func newConnectionHandler(unregisterChannel chan<- *ConnectionHandler, connectionDelegate ConnectionHandlerDelegate, role ShipRole, localService, remoteService *ServiceDetails, conn *websocket.Conn) *ConnectionHandler {
	return &ConnectionHandler{
		unregisterChannel:  unregisterChannel,
		connectionDelegate: connectionDelegate,
		role:               role,
		localService:       localService,
		remoteService:      remoteService,
		conn:               conn,
	}
}

func (c *ConnectionHandler) startup() {
	c.readChannel = make(chan []byte, 1)      // Listen to incoming websocket messages
	c.writeChannel = make(chan []byte, 1)     // Send outgoing websocket messages
	c.shipReadChannel = make(chan []byte, 1)  // Listen to incoming ship messages
	c.shipWriteChannel = make(chan []byte, 1) // Send outgoing ship messages
	c.shipTrustChannel = make(chan bool, 1)   // Listen to trust state update
	c.closeChannel = make(chan struct{}, 1)   // Listen to close events

	go c.readShipPump()
	go c.writePump()
	go c.writeShipPump()

	go func() {
		if err := c.shipHandshake(c.remoteService.userTrust || len(c.remoteService.ShipID) > 0); err != nil {
			fmt.Println("SHIP handshake error: ", err)
			c.shutdown(false)
			return
		}

		// Report to SPINE local device about this remote device connection
		c.connectionDelegate.addRemoteDeviceConnection(c.remoteService.SKI, c.readChannel, c.writeChannel)
		c.shipMessageHandler()
	}()
}

// shutdown the connection and all internals
// may only invoked after startup() is invoked!
func (c *ConnectionHandler) shutdown(safeShutdown bool) {
	if c.isConnectionClosed {
		return
	}

	if c.getSmeState() == smeComplete {
		c.connectionDelegate.removeRemoteDeviceConnection(c.remoteService.SKI)
	}

	c.unregisterChannel <- c

	if !util.IsChannelClosed(c.readChannel) {
		close(c.readChannel)
		c.readChannel = nil
	}

	if !util.IsChannelClosed(c.writeChannel) {
		close(c.writeChannel)
		c.writeChannel = nil
	}

	if !util.IsChannelClosed(c.shipReadChannel) {
		close(c.shipReadChannel)
		c.shipReadChannel = nil
	}

	if !util.IsChannelClosed(c.shipWriteChannel) {
		close(c.shipWriteChannel)
		c.shipWriteChannel = nil
	}

	if !util.IsChannelClosed(c.shipTrustChannel) {
		close(c.shipTrustChannel)
		c.shipTrustChannel = nil
	}

	if !util.IsChannelClosed(c.closeChannel) {
		close(c.closeChannel)
		c.closeChannel = nil
	}

	if c.conn != nil {
		if c.getSmeState() == smeComplete && safeShutdown {
			// close the SHIP connection according to the SHIP protocol
			c.shipClose()
		}

		c.conn.Close()
	}

	c.isConnectionClosed = true
}

// writePump pumps messages from the writeChannel to the writeShipChannel
func (c *ConnectionHandler) writePump() {
	for {
		select {
		case <-c.closeChannel:
			return
		case message, ok := <-c.writeChannel:
			if !ok {
				// The write channel is closed
				return
			}
			if c.isConnectionClosed {
				return
			}

			if err := c.sendSpineData(message); err != nil {
				fmt.Println("Error sending spine message: ", err)
				return
			}
		}
	}
}

// writePump pumps messages from the writeChannel to the websocket connection
func (c *ConnectionHandler) writeShipPump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.shutdown(false)
	}()

	for {
		select {
		case <-c.closeChannel:
			return
		case message, ok := <-c.shipWriteChannel:
			if c.isConnectionClosed {
				return
			}

			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				fmt.Println("Ship write channel closed")
				// The write channel has been closed
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				fmt.Println("Error writing to websocket: ", err)
				return
			}
		case <-ticker.C:
			if c.isConnectionClosed {
				return
			}
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("Error writing to websocket: ", err)
				return
			}
		}
	}
}

// readShipPump pumps messages from the websocket connection into the read message channel
func (c *ConnectionHandler) readShipPump() {
	defer func() {
		c.shutdown(false)
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		select {
		case <-c.closeChannel:
			return
		default:
			if c.isConnectionClosed {
				return
			}

			message, err := c.readWebsocketMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Println("Error reading message: ", err)
				}

				if c.isConnectionClosed {
					return
				}

				fmt.Println("Websocket read error: ", err)
				c.shutdown(false)
				return
			}

			// Check if this is a SHIP SME or SPINE message
			isShipMessage := false
			if c.getSmeState() != smeComplete {
				isShipMessage = true
			} else {
				isShipMessage = bytes.Contains([]byte("datagram:"), message)
			}

			if isShipMessage {
				c.shipReadChannel <- message
			} else {
				_, jsonData := c.parseMessage(message, true)

				// Get the datagram from the message
				data := ship.ShipData{}
				if err := json.Unmarshal(jsonData, &data); err != nil {
					fmt.Println("Error unmarshalling message: ", err)
					continue
				}

				if data.Data.Payload == nil {
					fmt.Println("Received no valid payload")
					continue
				}
				go func() {
					if c.readChannel == nil {
						return
					}
					c.readChannel <- []byte(data.Data.Payload)
				}()
			}
		}
	}
}

// handles incoming ship specific messages outside of the handshake process
func (c *ConnectionHandler) shipMessageHandler() {
	for {
		select {
		case <-c.closeChannel:
			return
		case msg := <-c.shipReadChannel:
			// TODO: implement this
			// This should only be a close/abort message, right?
			fmt.Println(string(msg))
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

	c.shipWriteChannel <- message
	return nil
}

const payloadPlaceholder = `{"place":"holder"}`

func (c *ConnectionHandler) transformSpineDataIntoShipJson(data []byte) ([]byte, error) {
	spineMsg, err := util.JsonIntoEEBUSJson(data)
	if err != nil {
		return nil, err
	}

	payload := json.RawMessage([]byte(spineMsg))

	// Workaround for the fact that SHIP payload is a json.RawMessage
	// which would also be transformed into an array element but it shouldn't
	// hence patching the payload into the message later after the SHIP
	// and SPINE model are transformed independently

	// Create the message
	shipMessage := ship.ShipData{
		Data: ship.DataType{
			Header: ship.HeaderType{
				ProtocolId: ship.ShipProtocolId,
			},
			Payload: json.RawMessage([]byte(payloadPlaceholder)),
		},
	}

	msg, err := json.Marshal(shipMessage)
	if err != nil {
		return nil, err
	}

	eebusMsg, err := util.JsonIntoEEBUSJson(msg)
	if err != nil {
		return nil, err
	}

	eebusMsg = strings.ReplaceAll(eebusMsg, `[`+payloadPlaceholder+`]`, string(payload))

	return []byte(eebusMsg), nil
}

func (c *ConnectionHandler) sendSpineData(data []byte) error {
	eebusMsg, err := c.transformSpineDataIntoShipJson(data)
	if err != nil {
		return err
	}

	fmt.Println("Send: ", string(eebusMsg))

	// Wrap the message into a binary message with the ship header
	shipMsg := []byte{ship.MsgTypeData}
	shipMsg = append(shipMsg, eebusMsg...)

	err = c.writeWebsocketMessage(shipMsg)
	if err != nil {
		fmt.Println("Error sending message: ", err)
		return err
	}

	return nil
}

// send a json message for a provided model to the websocket connection
func (c *ConnectionHandler) sendShipModel(typ byte, model interface{}) error {
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
