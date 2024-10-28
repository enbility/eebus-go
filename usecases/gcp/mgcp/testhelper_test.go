package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	spinemocks "github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

const remoteSki string = "testremoteski"

func TestMuMPCSuite(t *testing.T) {
	suite.Run(t, new(GcpMpcgSuite))
}

type GcpMpcgSuite struct {
	suite.Suite

	sut *MGCP

	service api.ServiceInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	mockRemoteEntity *spinemocks.EntityRemoteInterface
	monitoredEntity  spineapi.EntityRemoteInterface
	loadControlFeature,
	deviceDiagnosisFeature,
	deviceConfigurationFeature spineapi.FeatureLocalInterface

	eventCalled bool
}

func (s *GcpMpcgSuite) Event(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	s.eventCalled = true
}

func (s *GcpMpcgSuite) BeforeTest(suiteName, testName string) {
	s.eventCalled = false
	cert, _ := cert.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := api.NewConfiguration(
		"test", "test", "test", "test",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeGridConnectionHub},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeInverter},
		9999, cert, time.Second*4)

	serviceHandler := mocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	_ = s.service.Setup()

	mockRemoteDevice := spinemocks.NewDeviceRemoteInterface(s.T())
	s.mockRemoteEntity = spinemocks.NewEntityRemoteInterface(s.T())
	mockRemoteFeature := spinemocks.NewFeatureRemoteInterface(s.T())
	mockRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(mockRemoteFeature).Maybe()
	mockRemoteDevice.EXPECT().Ski().Return(remoteSki).Maybe()
	s.mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice).Maybe()
	s.mockRemoteEntity.EXPECT().EntityType().Return(mock.Anything).Maybe()
	entityAddress := &model.EntityAddressType{}
	s.mockRemoteEntity.EXPECT().Address().Return(entityAddress).Maybe()
	mockRemoteFeature.EXPECT().DataCopy(mock.Anything).Return(mock.Anything).Maybe()
	mockRemoteFeature.EXPECT().Address().Return(&model.FeatureAddressType{}).Maybe()
	mockRemoteFeature.EXPECT().Operations().Return(nil).Maybe()

	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeInverter)
	s.sut = NewMGCP(
		localEntity,
		s.Event,
		&MonitorPvFeedInPowerLimitationFactorConfig{},
		&MonitorPowerConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
		},
		&MonitorEnergyConfig{
			ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorVoltageConfig{
			ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorFrequencyConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
	)

	s.sut.AddFeatures()
	s.sut.AddUseCase()

	//s.remoteDevice, s.monitoredEntity = setupDevices(s.service, s.T())
}
