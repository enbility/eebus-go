package features

import (
	"testing"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestElectricalConnectionSuite(t *testing.T) {
	suite.Run(t, new(ElectricalConnectionSuite))
}

type ElectricalConnectionSuite struct {
	suite.Suite

	localDevice  *spine.DeviceLocalImpl
	remoteEntity *spine.EntityRemoteImpl

	electricalConnection *ElectricalConnection
	sentMessage          []byte
}

var _ spine.SpineDataConnection = (*ElectricalConnectionSuite)(nil)

func (s *ElectricalConnectionSuite) WriteSpineMessage(message []byte) {
	s.sentMessage = message
}

func (s *ElectricalConnectionSuite) BeforeTest(suiteName, testName string) {
	s.localDevice, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeElectricalConnection,
				functions: []model.FunctionType{
					model.FunctionTypeElectricalConnectionDescriptionListData,
					model.FunctionTypeElectricalConnectionParameterDescriptionListData,
					model.FunctionTypeElectricalConnectionPermittedValueSetListData,
				},
			},
		},
	)

	var err error
	s.electricalConnection, err = NewElectricalConnection(model.RoleTypeServer, model.RoleTypeClient, s.localDevice, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.electricalConnection)
}

func (s *ElectricalConnectionSuite) Test_RequestDescription() {
	err := s.electricalConnection.RequestDescription()
	assert.Nil(s.T(), err)
}

func (s *ElectricalConnectionSuite) Test_RequestParameterDescription() {
	err := s.electricalConnection.RequestParameterDescription()
	assert.Nil(s.T(), err)
}

func (s *ElectricalConnectionSuite) Test_RequestPermittedValueSet() {
	counter, err := s.electricalConnection.RequestPermittedValueSet()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *ElectricalConnectionSuite) Test_GetParamDescriptionListData() {
	mapMeasurementId, mapParamId, err := s.electricalConnection.GetParamDescriptionListData()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), mapMeasurementId)
	assert.Nil(s.T(), mapParamId)

	s.addParamDescription()

	mapMeasurementId, mapParamId, err = s.electricalConnection.GetParamDescriptionListData()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), mapMeasurementId)
	assert.NotNil(s.T(), mapParamId)
}

func (s *ElectricalConnectionSuite) Test_GetDescription() {
	data, err := s.electricalConnection.GetDescription()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.electricalConnection.GetDescription()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetConnectedPhases() {
	data, err := s.electricalConnection.GetConnectedPhases()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	s.addDescription()

	data, err = s.electricalConnection.GetConnectedPhases()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), uint(0), data)
}

func (s *ElectricalConnectionSuite) Test_GetCurrentsLimits() {
	data1, data2, data3, err := s.electricalConnection.GetCurrentsLimits()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data1)
	assert.Nil(s.T(), data2)
	assert.Nil(s.T(), data3)

	s.addParamDescription()

	data1, data2, data3, err = s.electricalConnection.GetCurrentsLimits()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data1)
	assert.Nil(s.T(), data2)
	assert.Nil(s.T(), data3)

	s.addPermittedValueSet()

	data1, data2, data3, err = s.electricalConnection.GetCurrentsLimits()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data1)
	assert.NotNil(s.T(), data2)
	assert.NotNil(s.T(), data3)
}

func (s *ElectricalConnectionSuite) Test_GetEVLimitValues() {
	data, err := s.electricalConnection.GetEVLimitValues()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescription()

	data, err = s.electricalConnection.GetEVLimitValues()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addPermittedValueSet()

	data, err = s.electricalConnection.GetEVLimitValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	s.addParamDescription2()

	data, err = s.electricalConnection.GetEVLimitValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

// helper

func (s *ElectricalConnectionSuite) addParamDescription() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addParamDescription2() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addDescription() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				PowerSupplyType:         util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcConnectedPhases:       util.Ptr(uint(3)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addPermittedValueSet() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(10),
						},
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(2),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}
	rF.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, fData, nil, nil)
}
