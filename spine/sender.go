//go:generate mockery --name=Sender
package spine

import (
	"encoding/json"
	"errors"
	"sync/atomic"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type ComControl interface {
	// This must be connected to the correct remote device !!
	SendSpineMessage(datagram model.DatagramType) error
}

type Sender interface {
	// Sends a read cmd to request some data
	Request(cmdClassifier model.CmdClassifierType, senderAddress, destinationAddress *model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error)
	// Sends a result cmd with no error to indicate that a message was processed successfully
	ResultSuccess(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType) error
	// Sends a result cmd with error information to indicate that a message processing failed
	ResultError(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *ErrorType) error
	// Sends a reply cmd to response to a read cmd
	Reply(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, cmd model.CmdType) error
	// Sends a call cmd with a subscription request
	Subscribe(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error)
	// Sends a call cmd with a binding request
	Bind(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error)
	// Sends a notify cmd to indicate that a subscribed feature changed
	Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error)
	// Sends a write cmd, setting properties of remote features
	Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error)
}

type SenderImpl struct {
	msgNum uint64 // 64bit values need to be defined on top of the struct to make atomic commands work on 32bit systems

	writeHandler WriteMessageI
}

var _ Sender = (*SenderImpl)(nil)

func NewSender(writeI WriteMessageI) Sender {
	return &SenderImpl{
		writeHandler: writeI,
	}
}

func (c *SenderImpl) sendSpineMessage(datagram model.DatagramType) error {
	// pack into datagram
	data := model.Datagram{
		Datagram: datagram,
	}

	// marshal
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if c.writeHandler == nil {
		return errors.New("outgoing interface implementation not set")
	}

	if msg == nil {
		return errors.New("message is nil")
	}

	logging.Log.Debug(datagram.PrintMessageOverview(true, "", ""))

	// write to channel
	c.writeHandler.WriteMessage(msg)

	return nil
}

// Sends request
func (c *SenderImpl) Request(cmdClassifier model.CmdClassifierType, senderAddress, destinationAddress *model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error) {
	msgCounter := c.getMsgCounter()

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           msgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: cmd,
		},
	}

	if ackRequest {
		datagram.Header.AckRequest = &ackRequest
	}

	return msgCounter, c.sendSpineMessage(datagram)
}

func (c *SenderImpl) ResultSuccess(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType) error {
	return c.result(requestHeader, senderAddress, nil)
}

func (c *SenderImpl) ResultError(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *ErrorType) error {
	return c.result(requestHeader, senderAddress, err)
}

// sends a result for a request
func (c *SenderImpl) result(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *ErrorType) error {
	cmdClassifier := model.CmdClassifierTypeResult

	addressSource := *requestHeader.AddressDestination
	addressSource.Device = senderAddress.Device

	var resultData model.ResultDataType
	if err != nil {
		resultData = model.ResultDataType{
			ErrorNumber: &err.ErrorNumber,
			Description: err.Description,
		}
	} else {
		resultData = model.ResultDataType{
			ErrorNumber: util.Ptr(model.ErrorNumberTypeNoError),
		}
	}

	cmd := model.CmdType{
		ResultData: &resultData,
	}

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        &addressSource,
			AddressDestination:   requestHeader.AddressSource,
			MsgCounter:           c.getMsgCounter(),
			MsgCounterReference:  requestHeader.MsgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return c.sendSpineMessage(datagram)
}

// Reply sends reply to original sender
func (c *SenderImpl) Reply(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, cmd model.CmdType) error {
	cmdClassifier := model.CmdClassifierTypeReply

	addressSource := *requestHeader.AddressDestination
	addressSource.Device = senderAddress.Device

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        &addressSource,
			AddressDestination:   requestHeader.AddressSource,
			MsgCounter:           c.getMsgCounter(),
			MsgCounterReference:  requestHeader.MsgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return c.sendSpineMessage(datagram)
}

// Notify sends notification to destination
func (c *SenderImpl) Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error) {
	msgCounter := c.getMsgCounter()

	cmdClassifier := model.CmdClassifierTypeNotify

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           msgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return msgCounter, c.sendSpineMessage(datagram)
}

// Write sends notification to destination
func (c *SenderImpl) Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error) {
	msgCounter := c.getMsgCounter()

	cmdClassifier := model.CmdClassifierTypeWrite
	ackRequest := true

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           msgCounter,
			CmdClassifier:        &cmdClassifier,
			AckRequest:           &ackRequest,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return msgCounter, c.sendSpineMessage(datagram)
}

// Send a subscription request to a remote server feature
func (c *SenderImpl) Subscribe(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error) {

	cmd := model.CmdType{
		NodeManagementSubscriptionRequestCall: NewNodeManagementSubscriptionRequestCallType(senderAddress, destinationAddress, serverFeatureType),
	}

	// we always send it to the remote NodeManagment feature, which always is at entity:[0],feature:0
	localAddress := NodeManagementAddress(senderAddress.Device)
	remoteAddress := NodeManagementAddress(destinationAddress.Device)

	return c.Request(model.CmdClassifierTypeCall, localAddress, remoteAddress, true, []model.CmdType{cmd})
}

// Send a binding request to a remote server feature
func (c *SenderImpl) Bind(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error) {
	cmd := model.CmdType{
		NodeManagementBindingRequestCall: NewNodeManagementBindingRequestCallType(senderAddress, destinationAddress, serverFeatureType),
	}

	// we always send it to the remote NodeManagment feature, which always is at entity:[0],feature:0
	localAddress := NodeManagementAddress(senderAddress.Device)
	remoteAddress := NodeManagementAddress(destinationAddress.Device)

	return c.Request(model.CmdClassifierTypeCall, localAddress, remoteAddress, true, []model.CmdType{cmd})
}

func (c *SenderImpl) getMsgCounter() *model.MsgCounterType {
	// TODO:  persistence
	i := model.MsgCounterType(atomic.AddUint64(&c.msgNum, 1))
	return &i
}
