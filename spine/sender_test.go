package spine

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

type WriteMessageHandler struct {
	sentMessage []byte

	mux sync.Mutex
}

var _ SpineDataConnection = (*WriteMessageHandler)(nil)

func (t *WriteMessageHandler) WriteSpineMessage(message []byte) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.sentMessage = message
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

	sentBytes := temp.sentMessage
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}
