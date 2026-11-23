package mpc

import (
	"slices"
	"time"

	"github.com/enbility/eebus-go/features/server"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
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
		9999, cert, time.Second*4, nil, nil)

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

	assert.Nil(s.T(), s.sut.AddFeatures())
	s.sut.AddUseCase()
}

func (s *MuMPCSuite) measurementPhaseSpecificDataForFilter(
	measurementFilter model.MeasurementDescriptionDataType,
	energyDirection model.EnergyDirectionType,
	validPhaseNameTypes []model.ElectricalConnectionPhaseNameType,
) ([]float64, error) {
	measurements, err := server.NewMeasurement(s.sut.LocalEntity)
	if err != nil {
		return nil, err
	}

	electricalConnection, err := server.NewElectricalConnection(s.sut.LocalEntity)
	if err != nil {
		return nil, err
	}

	data, err := measurements.GetDataForFilter(measurementFilter)
	if err != nil || len(data) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	var result []float64

	for _, item := range data {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		if validPhaseNameTypes != nil {
			filter := model.ElectricalConnectionParameterDescriptionDataType{
				MeasurementId: item.MeasurementId,
			}
			param, err := electricalConnection.GetParameterDescriptionsForFilter(filter)
			if err != nil || len(param) == 0 ||
				param[0].AcMeasuredPhases == nil ||
				!slices.Contains(validPhaseNameTypes, *param[0].AcMeasuredPhases) {
				continue
			}
		}

		if energyDirection != "" {
			filter := model.ElectricalConnectionParameterDescriptionDataType{
				MeasurementId: item.MeasurementId,
			}
			desc, err := electricalConnection.GetDescriptionForParameterDescriptionFilter(filter)
			if err != nil || desc == nil {
				continue
			}

			// if energy direction is not consume
			if desc.PositiveEnergyDirection == nil || *desc.PositiveEnergyDirection != energyDirection {
				return nil, err
			}
		}

		value := item.Value.GetValue()

		result = append(result, value)
	}

	return result, nil
}
