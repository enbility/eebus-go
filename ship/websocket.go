package ship

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/util"
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
	isConnectionClosed bool

	remoteSki string

	mux          sync.Mutex
	shutdownOnce sync.Once
}

// create a new websocket based shipDataProcessing implementation
func NewWebsocketConnection(conn *websocket.Conn, remoteSki string) *websocketConnection {
	return &websocketConnection{
		conn:      conn,
		remoteSki: remoteSki,
	}
}

// check if the websocket connection is closed
func (w *websocketConnection) isConnClosed() bool {
	w.mux.Lock()
	defer w.mux.Unlock()

	return w.isConnectionClosed
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

			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				logging.Log.Debug(w.remoteSki, "Ship write channel closed")
				// The write channel has been closed
				_ = w.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := w.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				logging.Log.Debug(w.remoteSki, "error writing to websocket: ", err)
				return
			}

			if len(message) > 2 {
				logging.Log.Trace("Send:", w.remoteSki, string(message[1:]))
			}

		case <-ticker.C:
			if w.isConnClosed() {
				return
			}

			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logging.Log.Debug(w.remoteSki, "error writing to websocket: ", err)
				return
			}
		}
	}
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
		if err != nil {
			logging.Log.Error(w.remoteSki, "websocket read error: ", err)
			w.close()
			w.dataProcessing.ReportConnectionError(err)
			return
		}

		if len(message) > 2 {
			logging.Log.Trace("Recv:", w.remoteSki, string(message[1:]))
		}

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
		if w.isConnectionClosed {
			return
		}

		w.mux.Lock()

		if !util.IsChannelClosed(w.closeChannel) {
			close(w.closeChannel)
			w.closeChannel = nil
		}

		if !util.IsChannelClosed(w.shipWriteChannel) {
			close(w.shipWriteChannel)
			w.shipWriteChannel = nil
		}

		w.mux.Unlock()

		if w.conn != nil {
			w.conn.Close()
		}

		w.isConnectionClosed = true
	})
}

var _ ShipDataConnection = (*websocketConnection)(nil)

func (w *websocketConnection) InitDataProcessing(dataProcessing ShipDataProcessing) {
	w.dataProcessing = dataProcessing

	w.run()
}

// write a message to the websocket connection
func (w *websocketConnection) WriteMessageToDataConnection(message []byte) error {
	if w.conn == nil {
		return errors.New("connection is not initialized")
	}

	w.shipWriteChannel <- message
	return nil
}

// shutdown the connection and all internals
func (w *websocketConnection) CloseDataConnection() {
	if !w.isConnClosed() {
		w.close()
	}
}

// return if the connection is closed
func (w *websocketConnection) IsDataConnectionClosed() bool {
	return w.isConnClosed()
}
