package integrationtests

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	detaileddiscoverydata_send_read_file_prefix  = "./testdata/01_detaileddiscoverydata_send_read"
	detaileddiscoverydata_recv_read_file_path    = "./testdata/01_detaileddiscoverydata_recv_read.json"
	detaileddiscoverydata_send_reply_file_prefix = "./testdata/01_detaileddiscoverydata_send_reply"
)

func TestNodeManagementSuite(t *testing.T) {
	suite.Run(t, new(NodeManagementSuite))
}

type NodeManagementSuite struct {
	suite.Suite
	sut    *spine.DeviceLocalImpl
	readC  chan []byte
	writeC chan []byte
}

func (s *NodeManagementSuite) SetupSuite() {
}

func (s *NodeManagementSuite) BeforeTest(suiteName, testName string) {
	s.sut = spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestDeviceCode",
		"TestSerialNumber", "TestDeviceAddress", model.DeviceTypeTypeChargingStation)

	s.readC = make(chan []byte)
	s.writeC = make(chan []byte)

	s.sut.AddRemoteDevice("", s.readC, s.writeC)
}

func (s *NodeManagementSuite) AfterTest(suiteName, testName string) {
}

func (s *NodeManagementSuite) TestDetailedDiscovery_SendRead() {
	// Act (see BeforeTest)

	// Assert
	sendBytes := <-s.writeC
	checkSentData(s.T(), sendBytes, detaileddiscoverydata_send_read_file_prefix)
}

func (s *NodeManagementSuite) TestDetailedDiscovery_SendReply() {
	// irgnore detaileddiscoverydata_send_read
	<-s.writeC

	// Act
	s.readC <- loadFileData(s.T(), detaileddiscoverydata_recv_read_file_path)

	// Assert
	sendBytes := <-s.writeC
	checkSentData(s.T(), sendBytes, detaileddiscoverydata_send_reply_file_prefix)
}

func (s *NodeManagementSuite) TestDetailedDiscovery_SendReplyWithAcknowledge() {
	// TODO
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

	saveJsonToFile(t, sendBytes, msgSendFilePrefix+"_actual.json")
	equal := jsonDatagramEqual(t, msgSendExpectedBytes, sendBytes)
	assert.True(t, equal, "Assert equal failed for "+msgSendFilePrefix)
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
