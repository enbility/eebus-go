package ship

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/util"
	"github.com/gorilla/websocket"
)

// Handling of the actual websocket connection to a remote device
type websocketConnection struct {
	// The actual websocket connection
	conn *websocket.Conn

	// The implementation handling message processing
	dataProcessing ShipDataProcessing

	// The connection was closed
	closeChannel chan struct{}

	// The ship write channel for outgoing SHIP messages
	shipWriteChannel chan []byte

	// internal handling of closed connections
	connectionClosed bool

	// the error message received for the closed connection
	connectionClosedError error

	remoteSki string

	muxConnClosed sync.Mutex
	muxShipWrite  sync.Mutex
	muxConWrite   sync.Mutex
	shutdownOnce  sync.Once
}

// create a new websocket based shipDataProcessing implementation
func NewWebsocketConnection(conn *websocket.Conn, remoteSki string) *websocketConnection {
	return &websocketConnection{
		conn:                  conn,
		remoteSki:             remoteSki,
		connectionClosedError: nil,
	}
}

// sets the error message for the closed connection
func (w *websocketConnection) setConnClosedError(err error) {
	w.muxConnClosed.Lock()
	defer w.muxConnClosed.Unlock()

	w.connectionClosed = true

	if err != nil {
		w.connectionClosedError = err
	}
}

func (w *websocketConnection) connClosedError() error {
	w.muxConnClosed.Lock()
	defer w.muxConnClosed.Unlock()

	return w.connectionClosedError
}

// check if the websocket connection is closed
func (w *websocketConnection) isConnClosed() bool {
	w.muxConnClosed.Lock()
	defer w.muxConnClosed.Unlock()

	return w.connectionClosed
}

func (w *websocketConnection) run() {
	w.shipWriteChannel = make(chan []byte, 1) // Send outgoing ship messages
	w.closeChannel = make(chan struct{}, 1)   // Listen to close events

	go w.readShipPump()
	go w.writeShipPump()
}

// writePump pumps messages from the SPINE and SHIP writeChannels to the websocket connection
func (w *websocketConnection) writeShipPump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-w.closeChannel:
			return

		case message, ok := <-w.shipWriteChannel:
			if w.isConnClosed() {
				return
			}

			w.muxConWrite.Lock()
			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w.muxConWrite.Unlock()
			if !ok {
				logging.Log().Debug(w.remoteSki, "ship write channel closed")
				// The write channel has been closed
				_ = w.writeMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := w.writeMessage(websocket.BinaryMessage, message); err != nil {
				// ignore write errors if the connection got closed
				if w.isConnClosed() {
					return
				}

				w.closeWithError(err, "error writing to websocket: ")
				return
			}

			var text string
			if len(message) > 2 {
				text = string(message[1:])
			} else if bytes.Equal(message, shipInit) {
				text = "ship init"
			} else {
				text = "unknown single byte"
			}
			logging.Log().Trace("Send:", w.remoteSki, text)

		case <-ticker.C:
			w.handlePing()
		}
	}
}

func (w *websocketConnection) handlePing() {
	if w.isConnClosed() {
		return
	}

	w.muxConWrite.Lock()
	_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))
	w.muxConWrite.Unlock()
	if err := w.writeMessage(websocket.PingMessage, nil); err != nil {
		w.closeWithError(err, "error writing to websocket: ")
		return
	}
}

func (w *websocketConnection) closeWithError(err error, reason string) {
	logging.Log().Debug(w.remoteSki, reason, err)
	w.setConnClosedError(err)
	w.dataProcessing.ReportConnectionError(err)
}

// readShipPump checks for messages from the websocket connection
func (w *websocketConnection) readShipPump() {
	_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPongHandler(func(string) error { _ = w.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		if w.isConnClosed() {
			return
		}

		message, err := w.readWebsocketMessage()
		// ignore read errors if the connection got closed
		if w.isConnClosed() {
			return
		}

		if err != nil {
			logging.Log().Debug(w.remoteSki, "websocket read error: ", err)
			w.close()
			w.setConnClosedError(err)
			w.dataProcessing.ReportConnectionError(err)
			return
		}

		var text string
		if len(message) > 2 {
			text = string(message[1:])
		} else if bytes.Equal(message, shipInit) {
			text = "ship init"
		} else {
			text = "unknown single byte"
		}
		logging.Log().Trace("Recv:", w.remoteSki, text)

		w.dataProcessing.HandleIncomingShipMessage(message)
	}
}

// read a message from the websocket connection
func (w *websocketConnection) readWebsocketMessage() ([]byte, error) {
	if w.conn == nil {
		return nil, errors.New("connection is not initialized")
	}

	msgType, b, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	if msgType != websocket.BinaryMessage {
		return nil, errors.New("message is not a binary message")
	}

	if len(b) < 2 {
		return nil, fmt.Errorf("invalid ship message length")
	}

	return b, nil
}

// close the current websocket connection
func (w *websocketConnection) close() {
	w.shutdownOnce.Do(func() {
		if w.isConnClosed() {
			return
		}

		w.setConnClosedError(nil)

		w.muxShipWrite.Lock()

		if !util.IsChannelClosed(w.closeChannel) {
			close(w.closeChannel)
			w.closeChannel = nil
		}

		if !util.IsChannelClosed(w.shipWriteChannel) {
			close(w.shipWriteChannel)
			w.shipWriteChannel = nil
		}

		if w.conn != nil {
			w.conn.Close()
		}

		w.muxShipWrite.Unlock()
	})
}

var _ ShipDataConnection = (*websocketConnection)(nil)

func (w *websocketConnection) InitDataProcessing(dataProcessing ShipDataProcessing) {
	w.dataProcessing = dataProcessing

	w.run()
}

// write a message to the websocket connection
func (w *websocketConnection) WriteMessageToDataConnection(message []byte) error {
	if w.isConnClosed() {
		return errors.New("connection is closed")
	}

	w.muxShipWrite.Lock()
	defer w.muxShipWrite.Unlock()

	if w.conn == nil || w.shipWriteChannel == nil {
		return errors.New("connection is closed")
	}

	w.shipWriteChannel <- message
	return nil
}

// make sure websocket Write is only called once at a time
func (w *websocketConnection) writeMessage(messageType int, data []byte) error {
	w.muxConWrite.Lock()
	defer w.muxConWrite.Unlock()

	return w.conn.WriteMessage(messageType, data)
}

// shutdown the connection and all internals
func (w *websocketConnection) CloseDataConnection(closeCode int, reason string) {
	if !w.isConnClosed() {
		if reason != "" {
			_ = w.writeMessage(websocket.CloseMessage, websocket.FormatCloseMessage(closeCode, reason))
		}
		w.setConnClosedError(nil)
		w.close()
	}
}

// return if the connection is closed
func (w *websocketConnection) IsDataConnectionClosed() (bool, error) {
	isClosed := w.isConnClosed()
	err := w.connClosedError()

	if isClosed && err == nil {
		err = errors.New("connection is closed")
	}

	return isClosed, err
}
