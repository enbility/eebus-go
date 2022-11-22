package integrationtests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

const (
	wallbox_detaileddiscoverydata_recv_reply_file_path  = "./testdata/wallbox_detaileddiscoverydata_recv_reply.json"
	wallbox_detaileddiscoverydata_recv_notify_file_path = "./testdata/wallbox_detaileddiscoverydata_recv_notify.json"
)

type WriteMessageHandler struct {
	sentMessages [][]byte
}

var _ spine.WriteMessageI = (*WriteMessageHandler)(nil)

func (t *WriteMessageHandler) WriteMessage(message []byte) {
	t.sentMessages = append(t.sentMessages, message)
}

func (t *WriteMessageHandler) LastMessage() []byte {
	if len(t.sentMessages) == 0 {
		return nil
	}
	return t.sentMessages[len(t.sentMessages)-1]
}

func (t *WriteMessageHandler) MessageWithReference(msgCounterReference *model.MsgCounterType) []byte {
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

func beforeTest(suiteName, testName string, fId uint, ftype model.FeatureTypeType, frole model.RoleType) (*spine.DeviceLocalImpl, string, spine.ReadMessageI, *WriteMessageHandler) {
	sut := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := spine.NewEntityLocalImpl(sut, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	sut.AddEntity(localEntity)
	f := spine.NewFeatureLocalImpl(fId, localEntity, ftype, frole)
	localEntity.AddFeature(f)

	remoteSki := "TestRemoteSki"

	writeHandler := &WriteMessageHandler{}
	remoteDevice := sut.AddRemoteDevice(remoteSki, writeHandler)

	return sut, remoteSki, remoteDevice, writeHandler
}

func initialCommunication(t *testing.T, readHandler spine.ReadMessageI, writeHandler *WriteMessageHandler) {
	// Initial generic communication

	_, _ = readHandler.ReadMessage(loadFileData(t, wallbox_detaileddiscoverydata_recv_reply_file_path))

	// Act
	msgCounter, _ := readHandler.ReadMessage(loadFileData(t, wallbox_detaileddiscoverydata_recv_notify_file_path))
	waitForAck(t, msgCounter, writeHandler)
}

func loadFileData(t *testing.T, fileName string) []byte {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	return fileData
}

func checkSentData(t *testing.T, sendBytes []byte, msgSendFilePrefix string) {
	msgSendExpectedBytes, err := os.ReadFile(msgSendFilePrefix + "_expected.json")
	if err != nil {
		t.Fatal(err)
	}

	msgSendActualFileName := msgSendFilePrefix + "_actual.json"
	equal := jsonDatagramEqual(t, msgSendExpectedBytes, sendBytes)
	if !equal {
		saveJsonToFile(t, sendBytes, msgSendActualFileName)
	}
	assert.Truef(t, equal, "Assert equal failed! Check '%s' ", msgSendActualFileName)
}

func jsonDatagramEqual(t *testing.T, expectedJson, actualJson []byte) bool {
	var actualDatagram model.Datagram
	if err := json.Unmarshal(actualJson, &actualDatagram); err != nil {
		t.Fatal(err)
	}
	var expectedDatagram model.Datagram
	if err := json.Unmarshal(expectedJson, &expectedDatagram); err != nil {
		t.Fatal(err)
	}

	less := func(a, b model.FunctionPropertyType) bool { return string(*a.Function) < string(*b.Function) }
	return cmp.Equal(expectedDatagram, actualDatagram, cmpopts.SortSlices(less))
}

func saveJsonToFile(t *testing.T, data json.RawMessage, fileName string) {
	jsonIndent, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(fileName, jsonIndent, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
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
