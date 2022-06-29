package integrationtests

import (
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	ec_detaileddiscoverydata_recv_reply_file_path              = "./testdata/ec_detaileddiscoverydata_recv_reply.json"
	ec_permittedvaluesetlistdata_recv_notify_partial_file_path = "./testdata/ec_permittedvaluesetlistdata_recv_notify_partial.json"
	ec_subscriptionRequestCall_recv_result_file_path           = "./testdata/ec_subscriptionRequestCall_recv_result.json"
)

func TestElectricalConnectionSuite(t *testing.T) {
	suite.Run(t, new(ElectricalConnectionSuite))
}

type ElectricalConnectionSuite struct {
	suite.Suite
	sut       *spine.DeviceLocalImpl
	remoteSki string
	readC     chan []byte
	writeC    chan []byte
}

func (s *ElectricalConnectionSuite) SetupSuite() {
}

func (s *ElectricalConnectionSuite) BeforeTest(suiteName, testName string) {
	s.sut = spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestDeviceCode",
		"TestSerialNumber", "TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := spine.NewEntityLocalImpl(s.sut, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	s.sut.AddEntity(localEntity)
	f := spine.NewFeatureLocalImpl(1, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)

	s.remoteSki = "TestRemoteSki"

	s.readC = make(chan []byte, 1)
	s.writeC = make(chan []byte, 1)

	s.sut.AddRemoteDevice(s.remoteSki, s.readC, s.writeC)
}

func (s *ElectricalConnectionSuite) AfterTest(suiteName, testName string) {
}

func (s *ElectricalConnectionSuite) TestPermittedValueSetListData_RecvNotifyPartial() {
	<-s.writeC // ignore NodeManagementDetailedDiscoveryData read

	// init with detaileddiscoverydata
	s.readC <- loadFileData(s.T(), ec_detaileddiscoverydata_recv_reply_file_path)
	<-s.writeC // ignore NodeManagementSubscriptionRequestCall
	s.readC <- loadFileData(s.T(), ec_subscriptionRequestCall_recv_result_file_path)
	<-s.writeC // ignore NodeManagementUseCaseData read

	// Act
	s.readC <- loadFileData(s.T(), ec_permittedvaluesetlistdata_recv_notify_partial_file_path)
	<-s.writeC // wait for ack

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)

	ecFeature := remoteDevice.FeatureByEntityTypeAndRole(
		remoteDevice.Entity(spine.NewAddressEntityType([]uint{1})),
		model.FeatureTypeTypeElectricalConnection,
		model.RoleTypeServer)
	assert.NotNil(s.T(), ecFeature)

	data := ecFeature.Data(model.FunctionTypeElectricalConnectionPermittedValueSetListData).(*model.ElectricalConnectionPermittedValueSetListDataType)
	if assert.NotNil(s.T(), data) {
		if assert.Equal(s.T(), 3, len(data.ElectricalConnectionPermittedValueSetData)) {
			item1 := data.ElectricalConnectionPermittedValueSetData[0]
			assert.Equal(s.T(), 0, int(*item1.ElectricalConnectionId))
			assert.Equal(s.T(), 1, int(*item1.ParameterId))
			assert.Equal(s.T(), 1, len(item1.PermittedValueSet))
			assert.Equal(s.T(), 1, len(item1.PermittedValueSet[0].Range))
			assert.NotNil(s.T(), item1.PermittedValueSet[0].Range)
			assert.Equal(s.T(), 6, int(*item1.PermittedValueSet[0].Range[0].Min.Number))
			assert.Equal(s.T(), 0, int(*item1.PermittedValueSet[0].Range[0].Min.Scale))
			assert.Equal(s.T(), 16, int(*item1.PermittedValueSet[0].Range[0].Max.Number))
			assert.Equal(s.T(), 0, int(*item1.PermittedValueSet[0].Range[0].Max.Scale))
			assert.Nil(s.T(), item1.PermittedValueSet[0].Value)
		}
	}
}
