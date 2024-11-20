package mpc

import (
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	spinemocks "github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const remoteSki string = "testremoteski"

type MuMPCSuite struct {
	*suite.Suite

	powerConfig     *MonitorPowerConfig
	energyConfig    *MonitorEnergyConfig
	currentConfig   *MonitorCurrentConfig
	voltageConfig   *MonitorVoltageConfig
	frequencyConfig *MonitorFrequencyConfig
	sut             *MPC

	service api.ServiceInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	mockRemoteEntity *spinemocks.EntityRemoteInterface
	monitoredEntity  spineapi.EntityRemoteInterface
	loadControlFeature,
	deviceDiagnosisFeature,
	deviceConfigurationFeature spineapi.FeatureLocalInterface

	eventCalled bool
}

func NewMuMPCSuite(
	suite *suite.Suite,
	powerConfig *MonitorPowerConfig,
	energyConfig *MonitorEnergyConfig,
	currentConfig *MonitorCurrentConfig,
	voltageConfig *MonitorVoltageConfig,
	frequencyConfig *MonitorFrequencyConfig,
) *MuMPCSuite {
	return &MuMPCSuite{
		Suite:           suite,
		powerConfig:     powerConfig,
		energyConfig:    energyConfig,
		currentConfig:   currentConfig,
		voltageConfig:   voltageConfig,
		frequencyConfig: frequencyConfig,
	}
}

func (s *MuMPCSuite) Event(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	s.eventCalled = true
}

func (s *MuMPCSuite) BeforeTest(suiteName, testName string) {
	s.eventCalled = false
	cert, _ := cert.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := api.NewConfiguration(
		"test", "test", "test", "test",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
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
	s.sut, _ = NewMPC(
		localEntity,
		s.Event,
		s.powerConfig,
		s.energyConfig,
		s.currentConfig,
		s.voltageConfig,
		s.frequencyConfig,
	)
	s.sut.AddFeatures()
	s.sut.AddUseCase()

	//s.remoteDevice, s.monitoredEntity = setupDevices(s.service, s.T())
}
