package spine

import (
	"encoding/json"
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSender_Notify_MsgCounter(t *testing.T) {
	writeC := make(chan []byte, 1)
	sut := NewSender(writeC)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))
	cmd := model.CmdType{
		ResultData: &model.ResultDataType{ErrorNumber: util.Ptr(model.ErrorNumberType(model.ErrorNumberTypeNoError))},
	}

	assert.NoError(t, sut.Notify(senderAddress, destinationAddress, cmd))
	<-writeC

	// Act
	assert.NoError(t, sut.Notify(senderAddress, destinationAddress, cmd))
	expectedMsgCounter := 2 //because Notify was called twice

	sentBytes := <-writeC
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}
