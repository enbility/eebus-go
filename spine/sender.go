//go:generate mockery --name=Sender
package spine

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type ComControl interface {
	//Read() (*model.CmiDatagramType, error)

	// This must be connected to the correct remote device !!
	SendSpineMessage(datagram model.DatagramType) error
}

type Sender interface {
	Request(cmdClassifier model.CmdClassifierType, senderAddress, destinationAddress *model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error)
	ResultSuccess(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType) error
	ResultError(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *ErrorType) error
	Reply(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, cmd model.CmdType) error
	Subscribe(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) error
	Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd []model.CmdType) error
}

type SenderImpl struct {
	msgNum uint64 // 64bit values need to be defined on top of the struct to make atomic commands work on 32bit systems
	//log        util.Logger

	writeChannel chan<- []byte
}

var _ Sender = (*SenderImpl)(nil)

func NewSender(writeC chan<- []byte) Sender {
	return &SenderImpl{
		//log:        log,
		writeChannel: writeC,
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

	if c.writeChannel == nil {
		return errors.New("write channel not set")
	}

	if msg == nil {
		return errors.New("message is nil")
	}

	fmt.Printf("%s\n", datagram.PrintMessageOverview(true, "", ""))

	// write to channel
	go func() { c.writeChannel <- msg }()

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
	// TODO where ack handling?

	// if ackRequest {
	// 	_ = c.sendAcknowledgementMessage(nil, featureSource, featureDestination, msgCounterReference)
	// }

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
func (c *SenderImpl) Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd []model.CmdType) error {
	cmdClassifier := model.CmdClassifierTypeNotify

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           c.getMsgCounter(),
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: cmd,
		},
	}

	return c.sendSpineMessage(datagram)
}

// Write sends notification to destination
func (c *SenderImpl) Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd []model.CmdType) error {
	cmdClassifier := model.CmdClassifierTypeWrite
	ackRequest := true

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           c.getMsgCounter(),
			CmdClassifier:        &cmdClassifier,
			AckRequest:           &ackRequest,
		},
		Payload: model.PayloadType{
			Cmd: cmd,
		},
	}

	return c.sendSpineMessage(datagram)
}

// Send a subscription request to a remote server feature
func (c *SenderImpl) Subscribe(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) error {
	cmd := model.CmdType{
		NodeManagementSubscriptionRequestCall: NewNodeManagementSubscriptionRequestCallType(senderAddress, destinationAddress, serverFeatureType),
	}

	// we always send it to the remote NodeManagment feature, which always is at entity:[0],feature:0
	var feature0 model.AddressFeatureType = 0
	remoteAddress := model.FeatureAddressType{
		Entity:  []model.AddressEntityType{0},
		Feature: &feature0,
		Device:  destinationAddress.Device,
	}

	_, err := c.Request(model.CmdClassifierTypeCall, senderAddress, &remoteAddress, true, []model.CmdType{cmd})

	return err
}

func (c *SenderImpl) getMsgCounter() *model.MsgCounterType {
	// TODO:  persistence
	i := model.MsgCounterType(atomic.AddUint64(&c.msgNum, 1))
	return &i
}
