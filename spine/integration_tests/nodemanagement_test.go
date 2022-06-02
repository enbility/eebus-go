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
	detaileddiscoverydata_recv_reply_file_path   = "./testdata/01_detaileddiscoverydata_recv_reply.json"
	detaileddiscoverydata_recv_read_file_path    = "./testdata/01_detaileddiscoverydata_recv_read.json"
	detaileddiscoverydata_send_reply_file_prefix = "./testdata/01_detaileddiscoverydata_send_reply"
)

func TestNodeManagementSuite(t *testing.T) {
	suite.Run(t, new(NodeManagementSuite))
}

type NodeManagementSuite struct {
	suite.Suite
	sut       *spine.DeviceLocalImpl
	remoteSki string
	readC     chan []byte
	writeC    chan []byte
}

func (s *NodeManagementSuite) SetupSuite() {
}

func (s *NodeManagementSuite) BeforeTest(suiteName, testName string) {
	s.sut = spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestDeviceCode",
		"TestSerialNumber", "TestDeviceAddress", model.DeviceTypeTypeChargingStation)
	s.remoteSki = "TestRemoteSki"

	s.readC = make(chan []byte, 1)
	s.writeC = make(chan []byte, 1)

	s.sut.AddRemoteDevice(s.remoteSki, s.readC, s.writeC)
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

func (s *NodeManagementSuite) TestDetailedDiscovery_RecvReply() {
	// irgnore detaileddiscoverydata_send_read
	<-s.writeC

	// Act
	s.readC <- loadFileData(s.T(), detaileddiscoverydata_recv_reply_file_path)

	<-s.writeC // to wait until the datagram is processed

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)

	rEntities := remoteDevice.Entities()
	assert.Equal(s.T(), 2, len(rEntities))
	di := rEntities[spine.DeviceInformationEntityId]
	assert.NotNil(s.T(), di)
	assert.Equal(s.T(), model.EntityTypeTypeDeviceInformation, di.EntityType())

	diFeatures := di.Features()
	assert.Equal(s.T(), 2, len(diFeatures))

	nm := diFeatures[0]
	assert.Equal(s.T(), spine.NodeManagementFeatureId, uint(*nm.Address().Feature))
	assert.Equal(s.T(), model.FeatureTypeTypeNodeManagement, nm.Type())
	assert.Equal(s.T(), model.RoleTypeSpecial, nm.Role())
	assert.Equal(s.T(), 8, len(nm.Operations()))

	dc := diFeatures[1]
	assert.Equal(s.T(), 1, int(*dc.Address().Feature))
	assert.Equal(s.T(), model.FeatureTypeTypeDeviceClassification, dc.Type())
	assert.Equal(s.T(), model.RoleTypeServer, dc.Role())
	assert.Equal(s.T(), 1, len(dc.Operations()))

	cem := rEntities[1]
	assert.NotNil(s.T(), cem)
	assert.Equal(s.T(), model.EntityTypeTypeCEM, cem.EntityType())

	cemFeatures := cem.Features()
	assert.Equal(s.T(), 2, len(cemFeatures))

	cemdc := cemFeatures[0]
	assert.Equal(s.T(), 1, int(*cemdc.Address().Feature))
	assert.Equal(s.T(), model.FeatureTypeTypeDeviceClassification, cemdc.Type())
	assert.Equal(s.T(), model.RoleTypeClient, cemdc.Role())
	assert.Equal(s.T(), 0, len(cemdc.Operations()))

	cemdd := cemFeatures[1]
	assert.Equal(s.T(), 2, int(*cemdd.Address().Feature))
	assert.Equal(s.T(), model.FeatureTypeTypeDeviceDiagnosis, cemdd.Type())
	assert.Equal(s.T(), model.RoleTypeClient, cemdd.Role())
	assert.Equal(s.T(), 0, len(cemdd.Operations()))

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
