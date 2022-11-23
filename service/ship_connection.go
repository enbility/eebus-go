package service

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/service/util"
	"github.com/DerAndereAndi/eebus-go/ship"
	"github.com/DerAndereAndi/eebus-go/spine"
)

// interface for SHIP and SPINE interactions
// implemented by service
type interactionShipSpine interface {
	// (dis-)approval of a connection trust is needed
	requestUserTrustForService(ski string)

	// inform about the SHIP identifier, which needs to be persisted
	shipIDUpdateForService(details *ServiceDetails)

	// new spine connection established, inform SPINE
	addRemoteDeviceConnection(ski string, writeI spine.WriteMessageI) spine.ReadMessageI

	// remove an existing connection from SPINE
	removeRemoteDeviceConnection(ski string)
}

// interface for handling the actual remote device data connection
//
// implemented by the websocket connection
type shipDataConnection interface {
	// initialize data processing
	InitDataProcessing(shipDataProcessing)

	// send data via the connection to the remote device
	WriteMessageToDataConnection([]byte) error

	// close the data connection
	CloseDataConnection()
}

// interface for handling incoming data
//
// implemented by shipConnection
type shipDataProcessing interface {
	// called for each incoming message
	HandleIncomingMessage([]byte)

	// called if the data connection is closed unsafe
	// e.g. due to connection issues
	ReportConnectionError(error)
}

// A shipConnection handles the data connection and coordinates SHIP and SPINE messages i/o
type shipConnection struct {
	// The ship connection mode of this connection
	role shipRole

	// The remote service
	remoteService *ServiceDetails

	// The local service
	localService *ServiceDetails

	// Where to pass incoming SPINE messages to
	readHandler spine.ReadMessageI

	// the handler for sending messages on the data connection
	dataHandler shipDataConnection

	// The current SHIP state
	smeState shipMessageExchangeState

	// contains the error message if smeState is in state smeHandshakeError
	handshakeError error

	// handles timeouts for the current smeState
	handshakeTimer        *time.Timer
	handshakeTimerRunning bool

	// stores if the connection should be trusted right away
	connectionTrustPending bool

	hubHandler         connectionHubHandler
	interactionHandler interactionShipSpine

	shutdownOnce sync.Once
}

func newConnectionHandler(hubHandler connectionHubHandler, connectionDelegate interactionShipSpine, role shipRole, localService, remoteService *ServiceDetails, dataHandler shipDataConnection) *shipConnection {
	ship := &shipConnection{
		hubHandler:         hubHandler,
		interactionHandler: connectionDelegate,
		role:               role,
		localService:       localService,
		remoteService:      remoteService,
		dataHandler:        dataHandler,
	}
	dataHandler.InitDataProcessing(ship)
	ship.startup()

	return ship
}

// start SHIP communication
func (c *shipConnection) startup() {
	c.handshakeTimer = time.NewTimer(time.Hour * 1)
	if !c.handshakeTimer.Stop() {
		<-c.handshakeTimer.C
	}

	// if the user trusted this connection e.g. via the UI or if we already have a stored SHIP ID for this SKI
	if !c.remoteService.userTrust && len(c.remoteService.ShipID) > 0 {
		c.remoteService.userTrust = true
	}

	c.smeState = cmiStateInitStart
	c.handleShipMessage(false, nil)
}

// close this ship connection
func (c *shipConnection) CloseConnection(safe bool) {
	c.shutdownOnce.Do(func() {
		c.interactionHandler.removeRemoteDeviceConnection(c.remoteService.SKI)
		c.hubHandler.HandleConnectionClosing(c)

		if safe {
			// SHIP 13.4.7: Connection Termination
			closeMessage := ship.ConnectionClose{
				ConnectionClose: ship.ConnectionCloseType{
					Phase: ship.ConnectionClosePhaseTypeAnnounce,
				},
			}

			_ = c.sendShipModel(ship.MsgTypeControl, closeMessage)

			// TODO: finish safe close implementation
		}

		c.dataHandler.CloseDataConnection()
	})
}

var _ connectionInteraction = (*shipConnection)(nil)

// handle an incoming trust result
func (c *shipConnection) ReportUserTrust(trust bool) {
	c.remoteService.userTrust = trust
	c.connectionTrustPending = false
	c.handleShipMessage(false, nil)
}

var _ spine.WriteMessageI = (*shipConnection)(nil)

// WriteMessageI interface implementation
func (c *shipConnection) WriteMessage(message []byte) {
	if err := c.sendSpineData(message); err != nil {
		logging.Log.Error(c.remoteService.SKI, "Error sending spine message: ", err)
		return
	}
}

var _ shipDataProcessing = (*shipConnection)(nil)

// route the incoming message to either SHIP or SPINE message handlers
func (c *shipConnection) HandleIncomingMessage(message []byte) {
	// Check if this is a SHIP SME or SPINE message
	if c.isShipMessage(message) {
		c.handleShipMessage(false, message)
		return
	}

	_, jsonData := c.parseMessage(message, true)

	// Get the datagram from the message
	data := ship.ShipData{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		logging.Log.Error(c.remoteService.SKI, "Error unmarshalling message: ", err)
		return
	}

	if data.Data.Payload == nil {
		logging.Log.Error(c.remoteService.SKI, "Received no valid payload")
		return
	}

	if c.readHandler == nil {
		return
	}

	// pass the payload to the SPINE read handler
	_, _ = c.readHandler.ReadMessage([]byte(data.Data.Payload))
}

// checks wether the provided messages is a SHIP message
func (c *shipConnection) isShipMessage(message []byte) bool {
	return !bytes.Contains(message, []byte("datagram"))
}

func (c *shipConnection) ReportConnectionError(err error) {
	c.CloseConnection(false)
}

const payloadPlaceholder = `{"place":"holder"}`

func (c *shipConnection) transformSpineDataIntoShipJson(data []byte) ([]byte, error) {
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

func (c *shipConnection) sendSpineData(data []byte) error {
	eebusMsg, err := c.transformSpineDataIntoShipJson(data)
	if err != nil {
		return err
	}

	logging.Log.Trace("Send:", c.remoteService.SKI, string(eebusMsg))

	// Wrap the message into a binary message with the ship header
	shipMsg := []byte{ship.MsgTypeData}
	shipMsg = append(shipMsg, eebusMsg...)

	err = c.dataHandler.WriteMessageToDataConnection(shipMsg)
	if err != nil {
		logging.Log.Error("Error sending message: ", err)
		return err
	}

	return nil
}

// send a json message for a provided model to the websocket connection
func (c *shipConnection) sendShipModel(typ byte, model interface{}) error {
	shipMsg, err := c.shipMessage(typ, model)
	if err != nil {
		return err
	}

	err = c.dataHandler.WriteMessageToDataConnection(shipMsg)
	if err != nil {
		return err
	}

	return nil
}

// transform a SHIP model into EEBUS specific JSON
func (c *shipConnection) shipMessage(typ byte, model interface{}) ([]byte, error) {
	msg, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	eebusMsg, err := util.JsonIntoEEBUSJson(msg)
	if err != nil {
		return nil, err
	}

	logging.Log.Trace("Send:", c.remoteService.SKI, string(eebusMsg))

	// Wrap the message into a binary message with the ship header
	shipMsg := []byte{typ}
	shipMsg = append(shipMsg, eebusMsg...)

	return shipMsg, nil
}

// return the SHIP message type, the SHIP message and an error
//
// enable jsonFormat if the return message is expected to be encoded in the eebus json format
func (c *shipConnection) parseMessage(msg []byte, jsonFormat bool) (byte, []byte) {
	// Extract the SHIP header byte
	shipHeaderByte := msg[0]
	// remove the SHIP header byte from the message
	msg = msg[1:]

	if len(msg) > 1 && c.smeState == smeComplete {
		logging.Log.Trace("Recv:", c.remoteService.SKI, string(msg))
	}

	if jsonFormat {
		return shipHeaderByte, util.JsonFromEEBUSJson(msg)
	}

	return shipHeaderByte, msg
}
