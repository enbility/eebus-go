package model

import "encoding/json"

const (
	MsgTypeInit    byte = 0
	MsgTypeControl byte = 1
	MsgTypeData    byte = 2
	MsgTypeEnd     byte = 3
)

const (
	ShipProtocolId = "ee1.0"
)

type ConnectionHelloPhaseType string

const (
	ConnectionHelloPhaseTypePending ConnectionHelloPhaseType = "pending"
	ConnectionHelloPhaseTypeReady   ConnectionHelloPhaseType = "ready"
	ConnectionHelloPhaseTypeAborted ConnectionHelloPhaseType = "aborted"
)

type ConnectionHello struct {
	ConnectionHello ConnectionHelloType `json:"connectionHello"`
}

type ConnectionHelloType struct {
	Phase               ConnectionHelloPhaseType `json:"phase"`
	Waiting             *uint                    `json:"waiting,omitempty"`
	ProlongationRequest *bool                    `json:"prolongationRequest,omitempty"`
}

type MessageProtocolFormatType string

const (
	MessageProtocolFormatTypeUTF8  MessageProtocolFormatType = "JSON-UTF8"
	MessageProtocolFormatTypeUTF16 MessageProtocolFormatType = "JSON-UTF16"
)

type MessageProtocolFormatsType struct {
	Format []MessageProtocolFormatType `json:"format"`
}

type ProtocolHandshakeTypeType string

const (
	ProtocolHandshakeTypeTypeAnnounceMax ProtocolHandshakeTypeType = "announceMax"
	ProtocolHandshakeTypeTypeSelect      ProtocolHandshakeTypeType = "select"
)

type Version struct {
	Major uint8 `json:"major"`
	Minor uint8 `json:"minor"`
}

type MessageProtocolHandshakeType struct {
	HandshakeType ProtocolHandshakeTypeType  `json:"handshakeType"`
	Version       Version                    `json:"version"`
	Formats       MessageProtocolFormatsType `json:"formats"`
}

type MessageProtocolHandshake struct {
	MessageProtocolHandshake MessageProtocolHandshakeType `json:"messageProtocolHandshake"`
}

type MessageProtocolHandshakeErrorErrorType uint8

const (
	MessageProtocolHandshakeErrorErrorTypeRFU               MessageProtocolHandshakeErrorErrorType = 0
	MessageProtocolHandshakeErrorErrorTypeTimeout           MessageProtocolHandshakeErrorErrorType = 1
	MessageProtocolHandshakeErrorErrorTypeUnexpectedMessage MessageProtocolHandshakeErrorErrorType = 2
	MessageProtocolHandshakeErrorErrorTypeSelectionMismatch MessageProtocolHandshakeErrorErrorType = 3
)

type MessageProtocolHandshakeErrorType struct {
	Error MessageProtocolHandshakeErrorErrorType `json:"error"`
}

type PinStateType string

const (
	PinStateTypeRequired PinStateType = "required"
	PinStateTypeOptional PinStateType = "optional"
	PinStateTypePinOk    PinStateType = "pinOk"
	PinStateTypeNone     PinStateType = "none"
)

type PinInputPermissionType string

const (
	PinInputPermissionTypeBusy PinInputPermissionType = "busy"
	PinInputPermissionTypeOk   PinInputPermissionType = "ok"
)

type MessageProtocolHandshakeError struct {
	Error MessageProtocolHandshakeErrorErrorType `json:"error"`
}

type ConnectionPinStateType struct {
	PinState        PinStateType            `json:"pinState"`
	InputPermission *PinInputPermissionType `json:"inputPermission,omitempty"`
}

type ConnectionPinState struct {
	ConnectionPinState ConnectionPinStateType `json:"connectionPinState"`
}

type PinValueType string

type ConnectionPinInputType struct {
	Pin PinValueType `json:"pin"`
}

type ConnectionPinErrorErrorType uint8

type ConnectionPinErrorType struct {
	Error ConnectionPinErrorErrorType `json:"error"`
}

type ProtocolIdType string

type HeaderType struct {
	ProtocolId ProtocolIdType `json:"protocolId"`
}

type ExtensionType struct {
	ExtensionId *string `json:"extensionId,omitempty"`
	Binary      *byte   `json:"binary,omitempty"` // HexBinary
	String      *string `json:"string,omitempty"`
}

type ShipData struct {
	Data DataType `json:"data"`
}

type DataType struct {
	Header    HeaderType      `json:"header"`
	Payload   json.RawMessage `json:"payload"`
	Extension *ExtensionType  `json:"extension,omitempty"`
}

type ConnectionClosePhaseType string

const (
	ConnectionClosePhaseTypeAnnounce ConnectionClosePhaseType = "announce"
	ConnectionClosePhaseTypeConfirm  ConnectionClosePhaseType = "confirm"
)

type ConnectionCloseReasonType string

const (
	ConnectionCloseReasonTypeUnspecific        ConnectionCloseReasonType = "unspecific"
	ConnectionCloseReasonTypeRemovedconnection ConnectionCloseReasonType = "removedConnection"
)

type ConnectionClose struct {
	ConnectionClose ConnectionCloseType `json:"connectionClose"`
}

type ConnectionCloseType struct {
	Phase   ConnectionClosePhaseType   `json:"phase"`
	MaxTime *uint                      `json:"maxTime,omitempty"`
	Reason  *ConnectionCloseReasonType `json:"reason,omitempty"`
}

type AccessMethodsRequest struct {
	AccessMethodsRequest AccessMethodsRequestType `json:"accessMethodsRequest"`
}

type AccessMethodsRequestType struct{}

type Dns struct {
	Uri string `json:"uri"`
}

type DnsSdMDns struct {
}

type AccessMethods struct {
	AccessMethods AccessMethodsType `json:"accessMethods"`
}

type AccessMethodsType struct {
	Id        *string    `json:"id"`
	DnsSdMDns *DnsSdMDns `json:"dnsSd_mDns,omitempty"`
	// According to the Spec Dns should be of type *Dns, but the SHM 2.0 only uses a string and would cause a crash
	Dns *string `json:"dns,omitempty"`
}
