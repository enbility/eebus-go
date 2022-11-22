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
	sut *spine.DeviceLocalImpl

	remoteSki string

	readHandler  spine.ReadMessageI
	writeHandler *WriteMessageHandler
}

func (s *DeviceDiagnosisSuite) SetupSuite() {
}

func (s *DeviceDiagnosisSuite) BeforeTest(suiteName, testName string) {
	s.sut, s.remoteSki, s.readHandler, s.writeHandler = beforeTest(suiteName, testName, 1, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)

	// f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)

	initialCommunication(s.T(), s.readHandler, s.writeHandler)
}

func (s *DeviceDiagnosisSuite) AfterTest(suiteName, testName string) {
}

func (s *DeviceDiagnosisSuite) TestHeartbeatSubscription_RecvNotify() {
	// Act
	msgCounter, _ := s.readHandler.ReadMessage(loadFileData(s.T(), dd_subscriptionRequestCall_recv_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)

	ddFeature := remoteDevice.FeatureByEntityTypeAndRole(
		remoteDevice.Entity(spine.NewAddressEntityType([]uint{1})),
		model.FeatureTypeTypeDeviceDiagnosis,
		model.RoleTypeClient)
	assert.NotNil(s.T(), ddFeature)
}
