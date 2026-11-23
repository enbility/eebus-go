package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
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
	"slices"
	"testing"
	"time"
)

const remoteSki string = "testremoteski"

func TestMuMPCPublicSuite(t *testing.T) {
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
	localEntity spineapi.EntityLocalInterface
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
		9999, cert, time.Second*4, nil, nil)

	serviceHandler := mocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	_ = s.service.Setup()

	s.localEntity = s.service.LocalDevice().EntityForType(model.EntityTypeTypeInverter)
	s.sut, _ = NewMGCP(
		s.localEntity,
		s.Event,
		&MonitorPvFeedInPowerLimitationFactorConfig{},
		&MonitorPowerConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeCalculatedValue),
			ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(1),
			}),
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

	_ = s.sut.AddFeatures()
	s.sut.AddUseCase()
}

func (s *GcpMpcgSuite) measurementPhaseSpecificDataForFilter(
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
