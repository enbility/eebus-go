package integrationtests

import (
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	dd_subscriptionRequestCall_recv_file_path        = "./testdata/dd_subscriptionRequestCall_recv.json"
	dd_subscriptionRequestCall_recv_result_file_path = "./testdata/ec_subscriptionRequestCall_recv_result.json"
)

func TestDeviceDiagnosisSuite(t *testing.T) {
	suite.Run(t, new(DeviceDiagnosisSuite))
}

type DeviceDiagnosisSuite struct {
	suite.Suite
	sut       *spine.DeviceLocalImpl
	remoteSki string
	readC     chan []byte
	writeC    chan []byte
}

func (s *DeviceDiagnosisSuite) SetupSuite() {
}

func (s *DeviceDiagnosisSuite) BeforeTest(suiteName, testName string) {
	s.sut = spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := spine.NewEntityLocalImpl(s.sut, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	s.sut.AddEntity(localEntity)
	f := spine.NewFeatureLocalImpl(1, localEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	localEntity.AddFeature(f)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)

	s.remoteSki = "TestRemoteSki"

	s.readC = make(chan []byte, 1)
	s.writeC = make(chan []byte, 1)

	s.sut.AddRemoteDevice(s.remoteSki, s.readC, s.writeC)
}

func (s *DeviceDiagnosisSuite) AfterTest(suiteName, testName string) {
}

func (s *DeviceDiagnosisSuite) TestHeartbeatSubscription_RecvNotify() {
	<-s.writeC // ignore NodeManagementDetailedDiscoveryData read

	// init with detaileddiscoverydata
	s.readC <- loadFileData(s.T(), wallbox_detaileddiscoverydata_recv_reply_file_path)
	<-s.writeC // ignore NodeManagementSubscriptionRequestCall
	<-s.writeC // ignore NodeManagementUseCaseData read

	// Act
	s.readC <- loadFileData(s.T(), dd_subscriptionRequestCall_recv_file_path)
	waitForAck(s.T(), s.writeC)

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)

	ddFeature := remoteDevice.FeatureByEntityTypeAndRole(
		remoteDevice.Entity(spine.NewAddressEntityType([]uint{1})),
		model.FeatureTypeTypeDeviceDiagnosis,
		model.RoleTypeClient)
	assert.NotNil(s.T(), ddFeature)
}
