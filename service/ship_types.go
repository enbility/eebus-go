package service

import (
	"time"

	"github.com/DerAndereAndi/eebus-go/ship"
)

type shipRole string

const (
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second // SHIP 4.2: ping interval + pong timeout
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 50 * time.Second // SHIP 4.2: ping interval

	// SHIP 9.2: Set maximum fragment length to 1024 bytes
	maxMessageSize = 1024

	ShipRoleServer shipRole = "server"
	ShipRoleClient shipRole = "client"
)

const (
	cmiTimeout              = 10 * time.Second       // SHIP 4.2
	cmiCloseTimeout         = 100 * time.Millisecond //nolint
	tHelloInit              = 60 * time.Second
	tHelloProlongThrInc     = 30 * time.Second
	tHelloProlongWaitingGap = 15 * time.Second //nolint
	tHellogProlongMin       = 1 * time.Second  //nolint
)

type shipMessageExchangeState uint

const (
	// Connection Mode Initialisation (CMI) SHIP 13.4.3
	cmiStateInitStart shipMessageExchangeState = iota
	cmiStateClientSend
	cmiStateClientWait
	cmiStateClientEvaluate
	cmiStateServerWait
	cmiStateServerEvaluate
	// Connection Data Preparation SHIP 13.4.4
	smeHelloState
	smeHelloStateReady //nolint
	smeHelloStateReadyInit
	smeHelloStateReadyListen
	smeHelloStateReadyTimeout
	smeHelloStatePending //nolint
	smeHelloStatePendingInit
	smeHelloStatePendingListen
	smeHelloStatePendingTimeout
	smeHelloStateOk
	smeHelloStateAbort
	// Connection State Protocol Handhsake SHIP 13.4.4.2
	smeProtHStateServerInit           //nolint
	smeProtHStateClientInit           //nolint
	smeProtHStateServerListenProposal //nolint
	smeProtHStateServerListenConfirm  //nolint
	smeProtHStateClientListenChoice   //nolint
	smeProtHStateTimeout              //nolint
	smeProtHStateClientOk             //nolint
	smeProtHStateServerOk             //nolint
	// Connection PIN State 13.4.5
	smePinStateCheckInit     //nolint
	smePinStateCheckListen   //nolint
	smePinStateCheckError    //nolint
	smePinStateCheckBusyInit //nolint
	smePinStateCheckBusyWait //nolint
	smePinStateCheckOk       //nolint
	smePinStateAskInit       //nolint
	smePinStateAskProcess    //nolint
	smePinStateAskRestricted //nolint
	smePinStateAskOk         //nolint
	// ConnectionAccess Methods Identification 13.4.6
	smeAccessMethodsRequest

	// Handshake approved on both ends
	smeApproved

	// Handshake process is successfully completed
	smeComplete

	// Handshake Error
	smeHandshakeError
)

var shipInit []byte = []byte{ship.MsgTypeInit, 0x00}
