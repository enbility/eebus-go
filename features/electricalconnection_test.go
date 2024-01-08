package features_test

import (
	"testing"

	"github.com/enbility/eebus-go/features"
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

	localEntity  spine.EntityLocal
	remoteEntity spine.EntityRemote

	electricalConnection *features.ElectricalConnection
	sentMessage          []byte
}

var _ spine.SpineDataConnection = (*ElectricalConnectionSuite)(nil)

func (s *ElectricalConnectionSuite) WriteSpineMessage(message []byte) {
	s.sentMessage = message
}

func (s *ElectricalConnectionSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
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
	s.electricalConnection, err = features.NewElectricalConnection(model.RoleTypeServer, model.RoleTypeClient, s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.electricalConnection)
}

func (s *ElectricalConnectionSuite) Test_RequestDescriptions() {
	err := s.electricalConnection.RequestDescriptions()
	assert.Nil(s.T(), err)
}

func (s *ElectricalConnectionSuite) Test_RequestParameterDescriptions() {
	err := s.electricalConnection.RequestParameterDescriptions()
	assert.Nil(s.T(), err)
}

func (s *ElectricalConnectionSuite) Test_RequestPermittedValueSets() {
	counter, err := s.electricalConnection.RequestPermittedValueSets()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *ElectricalConnectionSuite) Test_GetDescriptions() {
	data, err := s.electricalConnection.GetDescriptions()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.electricalConnection.GetDescriptions()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetDescriptionForMeasurementId() {
	measurementId := model.MeasurementIdType(1)
	data, err := s.electricalConnection.GetDescriptionForMeasurementId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.electricalConnection.GetDescriptionForMeasurementId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.electricalConnection.GetDescriptionForMeasurementId(measurementId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetParameterDescriptions() {
	data, err := s.electricalConnection.GetParameterDescriptions()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.electricalConnection.GetParameterDescriptions()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetParameterDescriptionForParameterId() {
	parametertId := model.ElectricalConnectionParameterIdType(1)
	data, err := s.electricalConnection.GetParameterDescriptionForParameterId(parametertId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.electricalConnection.GetParameterDescriptionForParameterId(parametertId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.electricalConnection.GetParameterDescriptionForParameterId(parametertId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	parametertId = model.ElectricalConnectionParameterIdType(10)
	data, err = s.electricalConnection.GetParameterDescriptionForParameterId(parametertId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetParameterDescriptionForMeasurementId() {
	measurementId := model.MeasurementIdType(1)
	data, err := s.electricalConnection.GetParameterDescriptionForMeasurementId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.electricalConnection.GetParameterDescriptionForMeasurementId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.electricalConnection.GetParameterDescriptionForMeasurementId(measurementId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	measurementId = model.MeasurementIdType(10)
	data, err = s.electricalConnection.GetParameterDescriptionForMeasurementId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetParameterDescriptionForMeasuredPhase() {
	phase := model.ElectricalConnectionPhaseNameTypeA
	data, err := s.electricalConnection.GetParameterDescriptionForMeasuredPhase(phase)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.electricalConnection.GetParameterDescriptionForMeasuredPhase(phase)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.electricalConnection.GetParameterDescriptionForMeasuredPhase(phase)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	phase = model.ElectricalConnectionPhaseNameTypeBc
	data, err = s.electricalConnection.GetParameterDescriptionForMeasuredPhase(phase)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetPermittedValueSets() {
	data, err := s.electricalConnection.GetPermittedValueSets()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addPermittedValueSet()

	data, err = s.electricalConnection.GetPermittedValueSets()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	s.addParamDescriptionPower()

	data, err = s.electricalConnection.GetPermittedValueSets()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetPermittedValueSetsEmptyElli() {
	data, err := s.electricalConnection.GetPermittedValueSets()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addPermittedValueSetEmptyElli()

	data, err = s.electricalConnection.GetPermittedValueSets()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

}

func (s *ElectricalConnectionSuite) Test_GetPermittedValueSetForParameterId() {
	parametertId := model.ElectricalConnectionParameterIdType(1)
	data, err := s.electricalConnection.GetPermittedValueSetForParameterId(parametertId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addPermittedValueSet()

	data, err = s.electricalConnection.GetPermittedValueSetForParameterId(parametertId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	parametertId = model.ElectricalConnectionParameterIdType(10)
	data, err = s.electricalConnection.GetPermittedValueSetForParameterId(parametertId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetPermittedValueSetForMeasurementId() {
	measurementId := model.MeasurementIdType(1)
	data, err := s.electricalConnection.GetPermittedValueSetForMeasurementId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addPermittedValueSet()

	data, err = s.electricalConnection.GetPermittedValueSetForMeasurementId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.electricalConnection.GetPermittedValueSetForMeasurementId(measurementId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	measurementId = model.MeasurementIdType(10)
	data, err = s.electricalConnection.GetPermittedValueSetForMeasurementId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetLimitsForParameterId() {
	parameterId := model.ElectricalConnectionParameterIdType(1)
	minV, maxV, defaultV, err := s.electricalConnection.GetLimitsForParameterId(parameterId)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), minV, 0.0)
	assert.Equal(s.T(), maxV, 0.0)
	assert.Equal(s.T(), defaultV, 0.0)

	s.addPermittedValueSet()
	s.addParamDescriptionCurrents()

	minV, maxV, defaultV, err = s.electricalConnection.GetLimitsForParameterId(parameterId)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), minV, 2.0)
	assert.Equal(s.T(), maxV, 16.0)
	assert.Equal(s.T(), defaultV, 0.1)
}

func (s *ElectricalConnectionSuite) Test_AdjustValueToBeWithinPermittedValuesForParameter() {
	parameterId := model.ElectricalConnectionParameterIdType(1)
	s.addPermittedValueSet()
	s.addParamDescriptionCurrents()

	value := s.electricalConnection.AdjustValueToBeWithinPermittedValuesForParameter(20, parameterId)
	assert.Equal(s.T(), value, 16.0)
	value = s.electricalConnection.AdjustValueToBeWithinPermittedValuesForParameter(2, parameterId)
	assert.Equal(s.T(), value, 2.0)
	value = s.electricalConnection.AdjustValueToBeWithinPermittedValuesForParameter(1, parameterId)
	assert.Equal(s.T(), value, 0.1)
}

// helper

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

func (s *ElectricalConnectionSuite) addParamDescriptionCurrents() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(1)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(4)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(3)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(2)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(4)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(5)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(5)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(3)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(6)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(6)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(7)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(7)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
				AcMeasurementVariant:    util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addParamDescriptionPower() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(8)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addPermittedValueSet() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(0.1),
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
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(3)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(0.1),
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
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(5)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(0.1),
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
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(8)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(0.1),
						},
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(400),
								Max: model.NewScaledNumberType(113664),
							},
						},
					},
				},
			},
		},
	}
	rF.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addPermittedValueSetEmptyElli() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				PermittedValueSet:      []model.ScaledNumberSetType{},
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, fData, nil, nil)
}
