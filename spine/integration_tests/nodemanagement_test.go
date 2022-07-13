package integrationtests

import (
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	nm_detaileddiscoverydata_send_read_file_prefix     = "./testdata/nm_detaileddiscoverydata_send_read"
	hems_detaileddiscoverydata_recv_reply_file_path    = "./testdata/hems_detaileddiscoverydata_recv_reply.json"
	nm_detaileddiscoverydata_recv_read_file_path       = "./testdata/nm_detaileddiscoverydata_recv_read.json"
	nm_detaileddiscoverydata_send_reply_file_prefix    = "./testdata/nm_detaileddiscoverydata_send_reply"
	nm_detaileddiscoverydata_recv_read_ack_file_path   = "./testdata/nm_detaileddiscoverydata_recv_read_ack.json"
	nm_detaileddiscoverydata_send_result_file_prefix   = "./testdata/nm_detaileddiscoverydata_send_result"
	nm_subscriptionRequestCall_recv_call_file_path     = "./testdata/nm_subscriptionRequestCall_recv_call.json"
	nm_subscriptionRequestCall_send_result_file_prefix = "./testdata/nm_subscriptionRequestCall_send_result"
	nm_destinationListData_recv_read_file_path         = "./testdata/nm_destinationListData_recv_read.json"
	nm_destinationListData_send_reply_file_prefix      = "./testdata/nm_destinationListData_send_reply"
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
		"TestSerialNumber", "TestDeviceAddress", model.DeviceTypeTypeChargingStation, model.NetworkManagementFeatureSetTypeSmart)
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
	checkSentData(s.T(), sendBytes, nm_detaileddiscoverydata_send_read_file_prefix)
}

func (s *NodeManagementSuite) TestDetailedDiscovery_SendReply() {
	// ignore NodeManagementDetailedDiscoveryData read
	<-s.writeC

	// Act
	s.readC <- loadFileData(s.T(), nm_detaileddiscoverydata_recv_read_file_path)

	// Assert
	sendBytes := <-s.writeC
	checkSentData(s.T(), sendBytes, nm_detaileddiscoverydata_send_reply_file_prefix)
}

func (s *NodeManagementSuite) TestDetailedDiscovery_RecvReply() {
	// ignore NodeManagementDetailedDiscoveryData read
	<-s.writeC

	// Act
	s.readC <- loadFileData(s.T(), hems_detaileddiscoverydata_recv_reply_file_path)

	<-s.writeC // to wait until the datagram is processed

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)
	assert.Equal(s.T(), model.DeviceTypeTypeEnergyManagementSystem, *remoteDevice.DeviceType())
	assert.Equal(s.T(), model.NetworkManagementFeatureSetTypeSmart, *remoteDevice.FeatureSet())

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
	// ignore NodeManagementDetailedDiscoveryData read
	<-s.writeC

	// Act
	s.readC <- loadFileData(s.T(), nm_detaileddiscoverydata_recv_read_ack_file_path)

	// Assert
	sentReply := <-s.writeC
	checkSentData(s.T(), sentReply, nm_detaileddiscoverydata_send_reply_file_prefix)
	sentResult := <-s.writeC
	checkSentData(s.T(), sentResult, nm_detaileddiscoverydata_send_result_file_prefix)
}

func (s *NodeManagementSuite) TestSubscriptionRequestCall_BeforeDetailedDiscovery() {
	// ignore NodeManagementDetailedDiscoveryData read
	<-s.writeC

	// Act
	s.readC <- loadFileData(s.T(), nm_subscriptionRequestCall_recv_call_file_path)

	// Assert
	sentResult := <-s.writeC
	checkSentData(s.T(), sentResult, nm_subscriptionRequestCall_send_result_file_prefix)

	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	subscriptionsForDevice := s.sut.SubscriptionManager().Subscriptions(remoteDevice)
	assert.Equal(s.T(), 1, len(subscriptionsForDevice))
	subscriptionsOnFeature := s.sut.SubscriptionManager().SubscriptionsOnFeature(*spine.NodeManagementAddress(s.sut.Address()))
	assert.Equal(s.T(), 1, len(subscriptionsOnFeature))
}

func (s *NodeManagementSuite) TestDestinationList_SendReply() {
	// ignore NodeManagementDetailedDiscoveryData read
	<-s.writeC

	// Act
	s.readC <- loadFileData(s.T(), nm_destinationListData_recv_read_file_path)

	// Assert
	sendBytes := <-s.writeC
	checkSentData(s.T(), sendBytes, nm_destinationListData_send_reply_file_prefix)
}
