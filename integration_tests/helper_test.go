package integrationtests

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

const (
	wallbox_detaileddiscoverydata_recv_reply_file_path  = "../spine/testdata/wallbox_detaileddiscoverydata_recv_reply.json"
	wallbox_detaileddiscoverydata_recv_notify_file_path = "../spine/testdata/wallbox_detaileddiscoverydata_recv_notify.json"
)

type WriteMessageHandler struct {
	sentMessages [][]byte

	mux sync.Mutex
}

var _ spine.SpineDataConnection = (*WriteMessageHandler)(nil)

func (t *WriteMessageHandler) WriteSpineMessage(message []byte) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.sentMessages = append(t.sentMessages, message)
}

func (t *WriteMessageHandler) LastMessage() []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	if len(t.sentMessages) == 0 {
		return nil
	}

	return t.sentMessages[len(t.sentMessages)-1]
}

func (t *WriteMessageHandler) MessageWithReference(msgCounterReference *model.MsgCounterType) []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	var datagram model.Datagram

	for _, msg := range t.sentMessages {
		if err := json.Unmarshal(msg, &datagram); err != nil {
			return nil
		}
		if datagram.Datagram.Header.MsgCounterReference == nil {
			continue
		}
		if uint(*datagram.Datagram.Header.MsgCounterReference) != uint(*msgCounterReference) {
			continue
		}
		if datagram.Datagram.Payload.Cmd[0].ResultData != nil {
			continue
		}

		return msg
	}

	return nil
}

func (t *WriteMessageHandler) ResultWithReference(msgCounterReference *model.MsgCounterType) []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	var datagram model.Datagram

	for _, msg := range t.sentMessages {
		if err := json.Unmarshal(msg, &datagram); err != nil {
			return nil
		}
		if datagram.Datagram.Header.MsgCounterReference == nil {
			continue
		}
		if uint(*datagram.Datagram.Header.MsgCounterReference) != uint(*msgCounterReference) {
			continue
		}
		if datagram.Datagram.Payload.Cmd[0].ResultData == nil {
			continue
		}

		return msg
	}

	return nil
}

func beforeTest(suiteName, testName string, fId uint, ftype model.FeatureTypeType, frole model.RoleType) (*spine.DeviceLocalImpl, string, spine.SpineDataProcessing, *WriteMessageHandler) {
	sut := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localEntity := spine.NewEntityLocalImpl(sut, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	sut.AddEntity(localEntity)
	f := spine.NewFeatureLocalImpl(fId, localEntity, ftype, frole)
	localEntity.AddFeature(f)

	remoteSki := "TestRemoteSki"

	writeHandler := &WriteMessageHandler{}
	remoteDevice := sut.AddRemoteDevice(remoteSki, writeHandler)

	return sut, remoteSki, remoteDevice, writeHandler
}

func initialCommunication(t *testing.T, readHandler spine.SpineDataProcessing, writeHandler *WriteMessageHandler) {
	// Initial generic communication

	_, _ = readHandler.HandleIncomingSpineMesssage(loadFileData(t, wallbox_detaileddiscoverydata_recv_reply_file_path))

	// Act
	msgCounter, _ := readHandler.HandleIncomingSpineMesssage(loadFileData(t, wallbox_detaileddiscoverydata_recv_notify_file_path))
	waitForAck(t, msgCounter, writeHandler)
}

func loadFileData(t *testing.T, fileName string) []byte {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	return fileData
}

func waitForAck(t *testing.T, msgCounterReference *model.MsgCounterType, writeHandler *WriteMessageHandler) {
	var datagram model.Datagram

	msg := writeHandler.ResultWithReference(msgCounterReference)
	if msg == nil {
		t.Fatal("acknowledge message was not sent!!")
	}

	if err := json.Unmarshal(msg, &datagram); err != nil {
		t.Fatal(err)
	}

	cmd := datagram.Datagram.Payload.Cmd[0]
	if cmd.ResultData != nil {
		if cmd.ResultData.ErrorNumber != nil && uint(*cmd.ResultData.ErrorNumber) != uint(model.ErrorNumberTypeNoError) {
			t.Fatal(fmt.Errorf("error '%d' result data received", uint(*cmd.ResultData.ErrorNumber)))
		}
	}
}
