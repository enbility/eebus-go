package spine

import (
	"encoding/json"
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSender_Reply_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))
	requestHeader := &model.HeaderType{
		AddressSource:      senderAddress,
		AddressDestination: destinationAddress,
		MsgCounter:         util.Ptr(model.MsgCounterType(10)),
	}
	cmd := model.CmdType{
		ResultData: &model.ResultDataType{ErrorNumber: util.Ptr(model.ErrorNumberType(model.ErrorNumberTypeNoError))},
	}

	err := sut.Reply(requestHeader, senderAddress, cmd)
	assert.NoError(t, err)

	// Act
	err = sut.Reply(requestHeader, senderAddress, cmd)
	assert.NoError(t, err)
	expectedMsgCounter := 2 //because Notify was called twice

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}

func TestSender_Notify_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))
	cmd := model.CmdType{
		ResultData: &model.ResultDataType{ErrorNumber: util.Ptr(model.ErrorNumberType(model.ErrorNumberTypeNoError))},
	}

	_, err := sut.Notify(senderAddress, destinationAddress, cmd)
	assert.NoError(t, err)

	// Act
	_, err = sut.Notify(senderAddress, destinationAddress, cmd)
	assert.NoError(t, err)
	expectedMsgCounter := 2 //because Notify was called twice

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))

	_, err = sut.DatagramForMsgCounter(model.MsgCounterType(2))
	assert.NoError(t, err)

	_, err = sut.DatagramForMsgCounter(model.MsgCounterType(3))
	assert.Error(t, err)
}

func TestSender_Write_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))
	cmd := model.CmdType{
		ResultData: &model.ResultDataType{ErrorNumber: util.Ptr(model.ErrorNumberType(model.ErrorNumberTypeNoError))},
	}

	_, err := sut.Write(senderAddress, destinationAddress, cmd)
	assert.NoError(t, err)

	// Act
	_, err = sut.Write(senderAddress, destinationAddress, cmd)
	assert.NoError(t, err)
	expectedMsgCounter := 2 //because Write was called twice

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}

func TestSender_Subscribe_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))

	_, err := sut.Subscribe(senderAddress, destinationAddress, model.FeatureTypeTypeLoadControl)
	assert.NoError(t, err)

	// Act
	_, err = sut.Subscribe(senderAddress, destinationAddress, model.FeatureTypeTypeLoadControl)
	assert.NoError(t, err)
	expectedMsgCounter := 2 //because Write was called twice

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}

func TestSender_Unsubscribe_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))

	_, err := sut.Unsubscribe(senderAddress, destinationAddress)
	assert.NoError(t, err)

	// Act
	_, err = sut.Unsubscribe(senderAddress, destinationAddress)
	assert.NoError(t, err)
	expectedMsgCounter := 2 //because Write was called twice

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}

func TestSender_Bind_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))

	_, err := sut.Bind(senderAddress, destinationAddress, model.FeatureTypeTypeLoadControl)
	assert.NoError(t, err)

	// Act
	_, err = sut.Bind(senderAddress, destinationAddress, model.FeatureTypeTypeLoadControl)
	assert.NoError(t, err)
	expectedMsgCounter := 2 //because Write was called twice

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}

func TestSender_Unbind_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))

	_, err := sut.Unbind(senderAddress, destinationAddress)
	assert.NoError(t, err)

	// Act
	_, err = sut.Unbind(senderAddress, destinationAddress)
	assert.NoError(t, err)
	expectedMsgCounter := 2 //because Write was called twice

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}
