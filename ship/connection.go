package ship

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/ship/model"
	shipUtil "github.com/DerAndereAndi/eebus-go/ship/util"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/util"
)

// implemented by connectionsHub and used by shipConnection
type ConnectionHandler interface {
	HandleClosedConnection(connection *ShipConnection)
}

// A ShipConnection handles the data connection and coordinates SHIP and SPINE messages i/o
type ShipConnection struct {
	// The ship connection mode of this connection
	role shipRole

	// The remote SKI
	RemoteSKI string

	// the remote SHIP Id
	remoteShipID string

	// The local SHIP ID
	localShipID string

	// data provider
	dataProvider ShipServiceDataProvider

	// Where to pass incoming SPINE messages to
	spineDataProcessing spine.SpineDataProcessing

	// the handler for sending messages on the data connection
	DataHandler ShipDataConnection

	// The current SHIP state
	smeState shipMessageExchangeState

	// handles timeouts for the various states
	//
	// WaitForReady SHIP 13.4.4.1.3: The communication partner must send its "READY" state (or request for prolongation") before the timer expires.
	//
	// SendProlongationRequest SHIP 13.4.4.1.3: Local timer to request for prolongation at the communication partner in time (i.e. before the communication partner's Wait-For-Ready-Timer expires).
	//
	// ProlongationRequestReply SHIP 13.4.4.1.3: Detection of response timeout on prolongation request.
	handshakeTimerRunning  bool
	handshakeTimerType     timeoutTimerType
	handshakeTimerStopChan chan struct{}
	handshakeTimerMux      sync.Mutex

	lastReceivedWaitingValue time.Duration // required for Prolong-Request-Reply-Timer

	// the SPINE local device
	spineLocalDevice *spine.DeviceLocalImpl

	shutdownOnce sync.Once

	mux sync.Mutex
}

func NewConnectionHandler(dataProvider ShipServiceDataProvider, dataHandler ShipDataConnection, spineLocalDevice *spine.DeviceLocalImpl, role shipRole, localShipID, remoteSki, remoteShipId string) *ShipConnection {
	ship := &ShipConnection{
		dataProvider:     dataProvider,
		spineLocalDevice: spineLocalDevice,
		role:             role,
		localShipID:      localShipID,
		RemoteSKI:        remoteSki,
		remoteShipID:     remoteShipId,
		DataHandler:      dataHandler,
		smeState:         cmiStateInitStart,
	}

	ship.handshakeTimerStopChan = make(chan struct{})

	dataHandler.InitDataProcessing(ship)

	return ship
}

// start SHIP communication
func (c *ShipConnection) Run() {
	c.handleShipMessage(false, nil)
}

func (c *ShipConnection) removeRemoteDeviceConnection() {
	remoteDevice := c.spineLocalDevice.RemoteDeviceForSki(c.RemoteSKI)

	c.spineLocalDevice.RemoveRemoteDevice(c.RemoteSKI)

	// inform about the disconnection
	payload := spine.EventPayload{
		Ski:        c.RemoteSKI,
		EventType:  spine.EventTypeDeviceChange,
		ChangeType: spine.ElementChangeRemove,
		Device:     remoteDevice,
	}
	spine.Events.Publish(payload)
}

// close this ship connection
func (c *ShipConnection) CloseConnection(safe bool, reason string) {
	c.shutdownOnce.Do(func() {
		c.stopHandshakeTimer()

		c.removeRemoteDeviceConnection()

		if safe && c.smeState > cmiStateInitStart {
			// SHIP 13.4.7: Connection Termination Announce
			closeMessage := model.ConnectionClose{
				ConnectionClose: model.ConnectionCloseType{
					Phase:  model.ConnectionClosePhaseTypeAnnounce,
					Reason: util.Ptr(model.ConnectionCloseReasonType(reason)),
				},
			}

			_ = c.sendShipModel(model.MsgTypeEnd, closeMessage)
			return
		}

		c.DataHandler.CloseDataConnection()
		c.dataProvider.HandleConnectionClosed(c)
	})
}

var _ spine.SpineDataConnection = (*ShipConnection)(nil)

// SpineDataConnection interface implementation
func (c *ShipConnection) WriteSpineMessage(message []byte) {
	if err := c.sendSpineData(message); err != nil {
		logging.Log.Error(c.RemoteSKI, "Error sending spine message: ", err)
		return
	}
}

var _ ShipDataProcessing = (*ShipConnection)(nil)

func (c *ShipConnection) shipModelFromMessage(message []byte) (*model.ShipData, error) {
	_, jsonData := c.parseMessage(message, true)

	// Get the datagram from the message
	data := model.ShipData{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		logging.Log.Error(c.RemoteSKI, "error unmarshalling message: ", err)
		return nil, err
	}

	if data.Data.Payload == nil {
		errorMsg := "received no valid payload"
		logging.Log.Error(c.RemoteSKI, errorMsg)
		return nil, errors.New(errorMsg)
	}

	return &data, nil
}

// route the incoming message to either SHIP or SPINE message handlers
func (c *ShipConnection) HandleIncomingShipMessage(message []byte) {
	if len(message) > 2 {
		logging.Log.Trace("Recv:", c.RemoteSKI, string(message[1:]))
	}

	// Check if this is a SHIP SME or SPINE message
	if !c.hasSpineDatagram(message) {
		c.handleShipMessage(false, message)
		return
	}

	data, err := c.shipModelFromMessage(message)
	if err != nil {
		return
	}

	if c.spineDataProcessing == nil {
		return
	}

	// pass the payload to the SPINE read handler
	_, _ = c.spineDataProcessing.HandleIncomingSpineMesssage([]byte(data.Data.Payload))
}

// checks wether the provided messages is a SHIP message
func (c *ShipConnection) hasSpineDatagram(message []byte) bool {
	return bytes.Contains(message, []byte("datagram"))
}

// the websocket data connection was closed from remote
func (c *ShipConnection) ReportConnectionError(err error) {
	// right now there is nothing to do
}

const payloadPlaceholder = `{"place":"holder"}`

func (c *ShipConnection) transformSpineDataIntoShipJson(data []byte) ([]byte, error) {
	spineMsg, err := shipUtil.JsonIntoEEBUSJson(data)
	if err != nil {
		return nil, err
	}

	payload := json.RawMessage([]byte(spineMsg))

	// Workaround for the fact that SHIP payload is a json.RawMessage
	// which would also be transformed into an array element but it shouldn't
	// hence patching the payload into the message later after the SHIP
	// and SPINE model are transformed independently

	// Create the message
	shipMessage := model.ShipData{
		Data: model.DataType{
			Header: model.HeaderType{
				ProtocolId: model.ShipProtocolId,
			},
			Payload: json.RawMessage([]byte(payloadPlaceholder)),
		},
	}

	msg, err := json.Marshal(shipMessage)
	if err != nil {
		return nil, err
	}

	eebusMsg, err := shipUtil.JsonIntoEEBUSJson(msg)
	if err != nil {
		return nil, err
	}

	eebusMsg = strings.ReplaceAll(eebusMsg, `[`+payloadPlaceholder+`]`, string(payload))

	return []byte(eebusMsg), nil
}

func (c *ShipConnection) sendSpineData(data []byte) error {
	eebusMsg, err := c.transformSpineDataIntoShipJson(data)
	if err != nil {
		return err
	}

	if c.DataHandler.IsDataConnectionClosed() {
		return err
	}

	logging.Log.Trace("Send:", c.RemoteSKI, string(eebusMsg))

	// Wrap the message into a binary message with the ship header
	shipMsg := []byte{model.MsgTypeData}
	shipMsg = append(shipMsg, eebusMsg...)

	err = c.DataHandler.WriteMessageToDataConnection(shipMsg)
	if err != nil {
		logging.Log.Debug("Error sending message: ", err)
		return err
	}

	return nil
}

// send a json message for a provided model to the websocket connection
func (c *ShipConnection) sendShipModel(typ byte, model interface{}) error {
	shipMsg, err := c.shipMessage(typ, model)
	if err != nil {
		return err
	}

	err = c.DataHandler.WriteMessageToDataConnection(shipMsg)
	if err != nil {
		return err
	}

	return nil
}

// Process a SHIP Json message
func (c *ShipConnection) processShipJsonMessage(message []byte, target any) error {
	_, data := c.parseMessage(message, true)

	return json.Unmarshal(data, &target)
}

// transform a SHIP model into EEBUS specific JSON
func (c *ShipConnection) shipMessage(typ byte, model interface{}) ([]byte, error) {
	if model == nil {
		return nil, errors.New("invalid data")
	}

	msg, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	eebusMsg, err := shipUtil.JsonIntoEEBUSJson(msg)
	if err != nil {
		return nil, err
	}

	logging.Log.Trace("Send:", c.RemoteSKI, string(eebusMsg))

	// Wrap the message into a binary message with the ship header
	shipMsg := []byte{typ}
	shipMsg = append(shipMsg, eebusMsg...)

	return shipMsg, nil
}

// return the SHIP message type, the SHIP message and an error
//
// enable jsonFormat if the return message is expected to be encoded in the eebus json format
func (c *ShipConnection) parseMessage(msg []byte, jsonFormat bool) (byte, []byte) {
	if len(msg) == 0 {
		return 0, nil
	}

	// Extract the SHIP header byte
	shipHeaderByte := msg[0]
	// remove the SHIP header byte from the message
	msg = msg[1:]

	if jsonFormat {
		return shipHeaderByte, shipUtil.JsonFromEEBUSJson(msg)
	}

	return shipHeaderByte, msg
}
