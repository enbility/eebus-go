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

func beforeTest(suiteName, testName string, fId uint, ftype model.FeatureTypeType, frole model.RoleType) (*spine.DeviceLocalImpl, string, chan []byte, chan []byte) {
	sut := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := spine.NewEntityLocalImpl(sut, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	sut.AddEntity(localEntity)
	f := spine.NewFeatureLocalImpl(fId, localEntity, ftype, frole)
	localEntity.AddFeature(f)

	remoteSki := "TestRemoteSki"

	readC := make(chan []byte, 1)
	writeC := make(chan []byte, 1)

	sut.AddRemoteDevice(remoteSki, readC, writeC)

	return sut, remoteSki, readC, writeC
}

func initialCommunication(t *testing.T, readC, writeC chan []byte) {
	// Initial generic communication
	<-writeC // ignore NodeManagementDetailedDiscoveryData read

	readC <- loadFileData(t, wallbox_detaileddiscoverydata_recv_reply_file_path)
	<-writeC // ignore NodeManagementSubscriptionRequestCall

	// Act
	readC <- loadFileData(t, wallbox_detaileddiscoverydata_recv_notify_file_path)
	waitForAck(t, writeC)
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

func waitForAck(t *testing.T, writeC chan []byte) {
	var datagram model.Datagram

	maxSentDatagram := 100
	for i := 0; i < maxSentDatagram; i++ {
		sentBytes := <-writeC
		if err := json.Unmarshal(sentBytes, &datagram); err != nil {
			t.Fatal(err)
		}
		cmd := datagram.Datagram.Payload.Cmd[0]
		if cmd.ResultData != nil {
			if cmd.ResultData.ErrorNumber != nil && uint(*cmd.ResultData.ErrorNumber) != uint(model.ErrorNumberTypeNoError) {
				t.Fatal(fmt.Errorf("error '%d' result data received", uint(*cmd.ResultData.ErrorNumber)))
			}
			return
		}
	}

	t.Fatal("acknowledge message was not sent!!")
}
