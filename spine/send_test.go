package spine_test

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

type WriteMessageHandler struct {
	sentMessage []byte

	mux sync.Mutex
}

var _ spine.SpineDataConnection = (*WriteMessageHandler)(nil)

func (t *WriteMessageHandler) WriteSpineMessage(message []byte) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.sentMessage = message
}

func TestSender_Notify_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := spine.NewSender(temp)

	senderAddress := featureAddressType(1, spine.NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, spine.NewEntityAddressType("destination", []uint{1}))
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

func featureAddressType(id uint, entityAddress *model.EntityAddressType) *model.FeatureAddressType {
	res := model.FeatureAddressType{
		Device:  entityAddress.Device,
		Entity:  entityAddress.Entity,
		Feature: util.Ptr(model.AddressFeatureType(id)),
	}

	return &res
}
