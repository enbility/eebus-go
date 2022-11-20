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
	s.sut = spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
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
	s.readC <- loadFileData(s.T(), wallbox_detaileddiscoverydata_recv_reply_file_path)
	<-s.writeC // ignore NodeManagementSubscriptionRequestCall

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)
	assert.Equal(s.T(), model.DeviceTypeTypeChargingStation, *remoteDevice.DeviceType())
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

	evse := rEntities[1]
	assert.NotNil(s.T(), evse)
	assert.Equal(s.T(), model.EntityTypeTypeEVSE, evse.EntityType())

	evseFeatures := evse.Features()
	assert.Equal(s.T(), 3, len(evseFeatures))

	evsedc := evseFeatures[0]
	assert.Equal(s.T(), 1, int(*evsedc.Address().Feature))
	assert.Equal(s.T(), model.FeatureTypeTypeDeviceClassification, evsedc.Type())
	assert.Equal(s.T(), model.RoleTypeClient, evsedc.Role())
	assert.Equal(s.T(), 0, len(evsedc.Operations()))

	evsedd := evseFeatures[1]
	assert.Equal(s.T(), 2, int(*evsedd.Address().Feature))
	assert.Equal(s.T(), model.FeatureTypeTypeDeviceDiagnosis, evsedd.Type())
	assert.Equal(s.T(), model.RoleTypeClient, evsedd.Role())
	assert.Equal(s.T(), 0, len(evsedd.Operations()))

	evseec := evseFeatures[2]
	assert.Equal(s.T(), 3, int(*evseec.Address().Feature))
	assert.Equal(s.T(), model.FeatureTypeTypeElectricalConnection, evseec.Type())
	assert.Equal(s.T(), model.RoleTypeServer, evseec.Role())
	assert.Equal(s.T(), 0, len(evseec.Operations()))

}

func (s *NodeManagementSuite) TestDetailedDiscovery_RecvNotifyAdded() {
	// ignore NodeManagementDetailedDiscoveryData read
	<-s.writeC

	s.readC <- loadFileData(s.T(), wallbox_detaileddiscoverydata_recv_reply_file_path)
	<-s.writeC // ignore NodeManagementSubscriptionRequestCall

	// Act
	s.readC <- loadFileData(s.T(), wallbox_detaileddiscoverydata_recv_notify_file_path)
	waitForAck(s.T(), s.writeC)

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)
	assert.Equal(s.T(), model.DeviceTypeTypeChargingStation, *remoteDevice.DeviceType())
	assert.Equal(s.T(), model.NetworkManagementFeatureSetTypeSmart, *remoteDevice.FeatureSet())

	rEntities := remoteDevice.Entities()
	if assert.Equal(s.T(), 3, len(rEntities)) {
		{
			di := rEntities[spine.DeviceInformationEntityId]
			assert.NotNil(s.T(), di)
			assert.Equal(s.T(), model.EntityTypeTypeDeviceInformation, di.EntityType())
			assert.Equal(s.T(), 2, len(di.Features()))
		}
		{
			evse := rEntities[1]
			assert.NotNil(s.T(), evse)
			assert.Equal(s.T(), model.EntityTypeTypeEVSE, evse.EntityType())
			assert.Equal(s.T(), 3, len(evse.Features()))
		}
		{
			ev := rEntities[2]
			assert.NotNil(s.T(), ev)
			assert.Equal(s.T(), model.EntityTypeTypeEV, ev.EntityType())
			assert.Equal(s.T(), 10, len(ev.Features()))
		}
	}
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
