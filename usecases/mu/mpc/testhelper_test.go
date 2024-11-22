package mpc

import (
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
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

func (s *MuMPCSuite) Event(_ string, _ spineapi.DeviceRemoteInterface, _ spineapi.EntityRemoteInterface, _ api.EventType) {
}

func (s *MuMPCSuite) BeforeTest(_, _ string) {
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
}
