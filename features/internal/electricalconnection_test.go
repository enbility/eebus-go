package internal_test

import (
	"testing"

	"github.com/enbility/eebus-go/features/internal"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestElectricalConnectionSuite(t *testing.T) {
	suite.Run(t, new(ElectricalConnectionSuite))
}

type ElectricalConnectionSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature, localMeasFeature   spineapi.FeatureLocalInterface
	remoteFeature, remoteMeasFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.ElectricalConnectionCommon
}

func (s *ElectricalConnectionSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeElectricalConnection,
				functions: []model.FunctionType{
					model.FunctionTypeElectricalConnectionDescriptionListData,
					model.FunctionTypeElectricalConnectionParameterDescriptionListData,
					model.FunctionTypeElectricalConnectionPermittedValueSetListData,
					model.FunctionTypeElectricalConnectionCharacteristicListData,
				},
			},
			{
				featureType: model.FeatureTypeTypeMeasurement,
				functions: []model.FunctionType{
					model.FunctionTypeMeasurementDescriptionListData,
				},
			},
		},
	)

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalElectricalConnection(s.localFeature)
	assert.NotNil(s.T(), s.localSut)
	s.localMeasFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localMeasFeature)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteElectricalConnection(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
	s.remoteMeasFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteMeasFeature)
}

func (s *ElectricalConnectionSuite) Test_CheckEventPayloadDataForFilter() {
	filter := model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
	}
	exists := s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	exists = s.localSut.CheckEventPayloadDataForFilter(filter, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(filter, filter)
	assert.False(s.T(), exists)

	descData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
		},
	}

	fErr := s.localFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)
	_, fErr = s.remoteFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	exists = s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	permittedData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{},
	}

	exists = s.localSut.CheckEventPayloadDataForFilter(permittedData, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(permittedData, filter)
	assert.False(s.T(), exists)

	permittedData = &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{*model.NewScaledNumberType(0)},
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(6),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}

	exists = s.localSut.CheckEventPayloadDataForFilter(permittedData, filter)
	assert.True(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(permittedData, filter)
	assert.True(s.T(), exists)
}

func (s *ElectricalConnectionSuite) Test_GetDescriptions() {
	filter := model.ElectricalConnectionDescriptionDataType{}
	data, err := s.localSut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetDescriptionForMeasurementId() {
	filter := model.ElectricalConnectionParameterDescriptionDataType{
		MeasurementId: util.Ptr(model.MeasurementIdType(1)),
	}
	data, err := s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetDescriptionForParameterDescriptionFilter() {
	filter := model.ElectricalConnectionParameterDescriptionDataType{
		ParameterId: util.Ptr(model.ElectricalConnectionParameterIdType(1)),
	}
	data, err := s.localSut.GetDescriptionForParameterDescriptionFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDescriptionForParameterDescriptionFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetDescriptionForParameterDescriptionFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDescriptionForParameterDescriptionFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.localSut.GetDescriptionForParameterDescriptionFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetDescriptionForParameterDescriptionFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetParameterDescriptions() {
	filter := model.ElectricalConnectionParameterDescriptionDataType{}
	data, err := s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetParameterDescriptionForScope() {
	filter := model.ElectricalConnectionParameterDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}

	data, err := s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionPower()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	filter.ScopeType = util.Ptr(model.ScopeTypeTypeACCurrent)
	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
	filter.ScopeType = util.Ptr(model.ScopeTypeTypeACCurrent)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetParameterDescriptionForParameterId() {
	filter := model.ElectricalConnectionParameterDescriptionDataType{
		ParameterId: util.Ptr(model.ElectricalConnectionParameterIdType(1)),
	}

	data, err := s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	filter.ParameterId = util.Ptr(model.ElectricalConnectionParameterIdType(10))
	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetParameterDescriptionForMeasurementId() {
	filter := model.ElectricalConnectionParameterDescriptionDataType{
		MeasurementId: util.Ptr(model.MeasurementIdType(1)),
	}
	data, err := s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	filter.MeasurementId = util.Ptr(model.MeasurementIdType(10))
	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetParameterDescriptionForMeasuredPhase() {
	filter := model.ElectricalConnectionParameterDescriptionDataType{
		AcMeasuredPhases: util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
	}
	data, err := s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	filter.AcMeasuredPhases = util.Ptr(model.ElectricalConnectionPhaseNameTypeBc)
	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetPermittedValueSets() {
	filter := model.ElectricalConnectionPermittedValueSetDataType{}
	data, err := s.localSut.GetPermittedValueSetForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetPermittedValueSetForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addPermittedValueSet()

	data, err = s.localSut.GetPermittedValueSetForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetPermittedValueSetForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	s.addParamDescriptionPower()

	data, err = s.localSut.GetPermittedValueSetForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetPermittedValueSetForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetPermittedValueSetsEmptyElli() {
	filter := model.ElectricalConnectionPermittedValueSetDataType{}
	data, err := s.localSut.GetPermittedValueSetForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetPermittedValueSetForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addPermittedValueSetEmptyElli()

	data, err = s.localSut.GetPermittedValueSetForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetPermittedValueSetForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetPermittedValueSetForParameterId() {
	filter := model.ElectricalConnectionPermittedValueSetDataType{
		ParameterId: util.Ptr(model.ElectricalConnectionParameterIdType(1)),
	}

	data, err := s.localSut.GetPermittedValueSetForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetPermittedValueSetForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addPermittedValueSet()

	data, err = s.localSut.GetPermittedValueSetForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetPermittedValueSetForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	filter.ParameterId = util.Ptr(model.ElectricalConnectionParameterIdType(10))
	data, err = s.localSut.GetPermittedValueSetForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetPermittedValueSetForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetPermittedValueSetForMeasurementId() {
	filter := model.ElectricalConnectionParameterDescriptionDataType{
		MeasurementId: util.Ptr(model.MeasurementIdType(1)),
	}
	data, err := s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addPermittedValueSet()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addParamDescriptionCurrents()

	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	filter.MeasurementId = util.Ptr(model.MeasurementIdType(10))
	data, err = s.localSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetLimitsForParameterId() {
	filter := model.ElectricalConnectionPermittedValueSetDataType{
		ParameterId: util.Ptr(model.ElectricalConnectionParameterIdType(1)),
	}
	minV, maxV, defaultV, err := s.localSut.GetPermittedValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), minV, 0.0)
	assert.Equal(s.T(), maxV, 0.0)
	assert.Equal(s.T(), defaultV, 0.0)
	minV, maxV, defaultV, err = s.remoteSut.GetPermittedValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), minV, 0.0)
	assert.Equal(s.T(), maxV, 0.0)
	assert.Equal(s.T(), defaultV, 0.0)

	s.addPermittedValueSet()
	s.addParamDescriptionCurrents()

	minV, maxV, defaultV, err = s.localSut.GetPermittedValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), minV, 2.0)
	assert.Equal(s.T(), maxV, 16.0)
	assert.Equal(s.T(), defaultV, 0.1)
	minV, maxV, defaultV, err = s.remoteSut.GetPermittedValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), minV, 2.0)
	assert.Equal(s.T(), maxV, 16.0)
	assert.Equal(s.T(), defaultV, 0.1)
}

func (s *ElectricalConnectionSuite) Test_GetLimitsForParameterId_Elli_1Phase() {
	filter := model.ElectricalConnectionPermittedValueSetDataType{
		ParameterId: util.Ptr(model.ElectricalConnectionParameterIdType(0)),
	}
	minV, maxV, defaultV, err := s.localSut.GetPermittedValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), minV, 0.0)
	assert.Equal(s.T(), maxV, 0.0)
	assert.Equal(s.T(), defaultV, 0.0)
	minV, maxV, defaultV, err = s.remoteSut.GetPermittedValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), minV, 0.0)
	assert.Equal(s.T(), maxV, 0.0)
	assert.Equal(s.T(), defaultV, 0.0)

	s.addParamDescription_Elli_1Phase()
	s.addPermittedValueSetEmptyElli()

	minV, maxV, defaultV, err = s.localSut.GetPermittedValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), minV, 0.0)
	assert.Equal(s.T(), maxV, 0.0)
	assert.Equal(s.T(), defaultV, 0.0)
	minV, maxV, defaultV, err = s.remoteSut.GetPermittedValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), minV, 0.0)
	assert.Equal(s.T(), maxV, 0.0)
	assert.Equal(s.T(), defaultV, 0.0)
}

func (s *ElectricalConnectionSuite) Test_AdjustValueToBeWithinPermittedValuesForParameter() {
	parameterId := model.ElectricalConnectionParameterIdType(1)
	s.addPermittedValueSet()
	s.addParamDescriptionCurrents()

	value := s.localSut.AdjustValueToBeWithinPermittedValuesForParameterId(20, parameterId)
	assert.Equal(s.T(), 16.0, value)
	value = s.remoteSut.AdjustValueToBeWithinPermittedValuesForParameterId(20, parameterId)
	assert.Equal(s.T(), 16.0, value)

	value = s.localSut.AdjustValueToBeWithinPermittedValuesForParameterId(2, parameterId)
	assert.Equal(s.T(), 2.0, value)
	value = s.remoteSut.AdjustValueToBeWithinPermittedValuesForParameterId(2, parameterId)
	assert.Equal(s.T(), 2.0, value)

	value = s.localSut.AdjustValueToBeWithinPermittedValuesForParameterId(1, parameterId)
	assert.Equal(s.T(), 0.1, value)
	value = s.remoteSut.AdjustValueToBeWithinPermittedValuesForParameterId(1, parameterId)
	assert.Equal(s.T(), 0.1, value)
}

func (s *ElectricalConnectionSuite) Test_AdjustValueToBeWithinPermittedValuesForParameter_Elli_1Phase() {
	parameterId := model.ElectricalConnectionParameterIdType(1)

	s.addPermittedValueSetEmptyElli()
	s.addParamDescription_Elli_1Phase()

	value := s.localSut.AdjustValueToBeWithinPermittedValuesForParameterId(20, parameterId)
	assert.Equal(s.T(), value, 20.0)
	value = s.remoteSut.AdjustValueToBeWithinPermittedValuesForParameterId(20, parameterId)
	assert.Equal(s.T(), value, 20.0)

	value = s.localSut.AdjustValueToBeWithinPermittedValuesForParameterId(2, parameterId)
	assert.Equal(s.T(), value, 2.0)
	value = s.remoteSut.AdjustValueToBeWithinPermittedValuesForParameterId(2, parameterId)
	assert.Equal(s.T(), value, 2.0)

	value = s.localSut.AdjustValueToBeWithinPermittedValuesForParameterId(1, parameterId)
	assert.Equal(s.T(), value, 1.0)
	value = s.remoteSut.AdjustValueToBeWithinPermittedValuesForParameterId(1, parameterId)
	assert.Equal(s.T(), value, 1.0)
}

func (s *ElectricalConnectionSuite) Test_GetCharacteristics() {
	filter := model.ElectricalConnectionCharacteristicDataType{}
	data, err := s.localSut.GetCharacteristicsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetCharacteristicsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addCharacteristics()

	data, err = s.localSut.GetCharacteristicsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetCharacteristicsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetCharacteristicForContextType() {
	filter := model.ElectricalConnectionCharacteristicDataType{
		CharacteristicContext: util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:    util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeEnergyCapacityNominalMax),
	}

	data, err := s.localSut.GetCharacteristicsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetCharacteristicsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addCharacteristics()

	filter.CharacteristicType = util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeApparentPowerConsumptionNominalMax)
	data, err = s.localSut.GetCharacteristicsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetCharacteristicsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)

	filter.CharacteristicType = util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeEnergyCapacityNominalMax)
	data, err = s.localSut.GetCharacteristicsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetCharacteristicsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_EVCurrentLimits() {
	minData, maxData, defaultData, err := s.localSut.GetPhaseCurrentLimits(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)
	minData, maxData, defaultData, err = s.remoteSut.GetPhaseCurrentLimits(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	minData, maxData, defaultData, err = s.localSut.GetPhaseCurrentLimits(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)
	minData, maxData, defaultData, err = s.remoteSut.GetPhaseCurrentLimits(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	minData, maxData, defaultData, err = s.localSut.GetPhaseCurrentLimits(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)
	minData, maxData, defaultData, err = s.remoteSut.GetPhaseCurrentLimits(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
		},
	}

	fErr := s.localFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)
	_, fErr = s.remoteFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	minData, maxData, defaultData, err = s.localSut.GetPhaseCurrentLimits(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	minData, maxData, defaultData, err = s.remoteSut.GetPhaseCurrentLimits(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	measData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent)},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(2)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent)},
		},
	}
	fErr = s.localMeasFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)
	_, fErr = s.remoteMeasFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	type permittedStruct struct {
		defaultExists                      bool
		defaultValue, expectedDefaultValue float64
		minValue, expectedMinValue         float64
		maxValue, expectedMaxValue         float64
	}

	tests := []struct {
		name      string
		permitted []permittedStruct
	}{
		{
			"1 Phase ISO15118",
			[]permittedStruct{
				{true, 0.1, 0.1, 2, 2, 16, 16},
			},
		},
		{
			"1 Phase IEC61851",
			[]permittedStruct{
				{true, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
		{
			"1 Phase IEC61851 Elli",
			[]permittedStruct{
				{false, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
		{
			"3 Phase ISO15118",
			[]permittedStruct{
				{true, 0.1, 0.1, 2, 2, 16, 16},
				{true, 0.1, 0.1, 2, 2, 16, 16},
				{true, 0.1, 0.1, 2, 2, 16, 16},
			},
		},
		{
			"3 Phase IEC61851",
			[]permittedStruct{
				{true, 0.0, 0.0, 6, 6, 16, 16},
				{true, 0.0, 0.0, 6, 6, 16, 16},
				{true, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
		{
			"3 Phase IEC61851 Elli",
			[]permittedStruct{
				{false, 0.0, 0.0, 6, 6, 16, 16},
				{false, 0.0, 0.0, 6, 6, 16, 16},
				{false, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			dataSet := []model.ElectricalConnectionPermittedValueSetDataType{}
			permittedData := []model.ScaledNumberSetType{}
			for index, data := range tc.permitted {
				item := model.ScaledNumberSetType{
					Range: []model.ScaledNumberRangeType{
						{
							Min: model.NewScaledNumberType(data.minValue),
							Max: model.NewScaledNumberType(data.maxValue),
						},
					},
				}
				if data.defaultExists {
					item.Value = []model.ScaledNumberType{*model.NewScaledNumberType(data.defaultValue)}
				}
				permittedData = append(permittedData, item)

				permittedItem := model.ElectricalConnectionPermittedValueSetDataType{
					ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(index)),
					PermittedValueSet:      permittedData,
				}
				dataSet = append(dataSet, permittedItem)
			}

			permData := &model.ElectricalConnectionPermittedValueSetListDataType{
				ElectricalConnectionPermittedValueSetData: dataSet,
			}

			fErr := s.localFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
			assert.Nil(s.T(), fErr)
			_, fErr = s.remoteFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
			assert.Nil(s.T(), fErr)

			minData, maxData, defaultData, err = s.localSut.GetPhaseCurrentLimits(measData.MeasurementDescriptionData)
			assert.Nil(s.T(), err)

			assert.Equal(s.T(), len(tc.permitted), len(minData))
			assert.Equal(s.T(), len(tc.permitted), len(maxData))
			assert.Equal(s.T(), len(tc.permitted), len(defaultData))
			for index, item := range tc.permitted {
				assert.Equal(s.T(), item.expectedMinValue, minData[index])
				assert.Equal(s.T(), item.expectedMaxValue, maxData[index])
				assert.Equal(s.T(), item.expectedDefaultValue, defaultData[index])
			}

			minData, maxData, defaultData, err = s.remoteSut.GetPhaseCurrentLimits(measData.MeasurementDescriptionData)
			assert.Nil(s.T(), err)

			assert.Equal(s.T(), len(tc.permitted), len(minData))
			assert.Equal(s.T(), len(tc.permitted), len(maxData))
			assert.Equal(s.T(), len(tc.permitted), len(defaultData))
			for index, item := range tc.permitted {
				assert.Equal(s.T(), item.expectedMinValue, minData[index])
				assert.Equal(s.T(), item.expectedMaxValue, maxData[index])
				assert.Equal(s.T(), item.expectedDefaultValue, defaultData[index])
			}
		})
	}
}

// helper

func (s *ElectricalConnectionSuite) addDescription() {
	fData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				PowerSupplyType: util.Ptr(model.ElectricalConnectionVoltageTypeTypeDc),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				PowerSupplyType:         util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcConnectedPhases:       util.Ptr(uint(3)),
				PositiveEnergyDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addCharacteristics() {
	fData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeEnergyCapacityNominalMax),
				Value:                  model.NewScaledNumberType(98),
				Unit:                   util.Ptr(model.UnitOfMeasurementTypeWh),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeElectricalConnectionCharacteristicListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addParamDescriptionCurrents() {
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
	_ = s.localFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addParamDescriptionPower() {
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
	_ = s.localFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addParamDescription_Elli_1Phase() {
	fData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
				AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(3)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
				AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(4)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(3)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(5)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(4)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(6)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(5)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(7)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(6)),
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
				AcMeasurementType:      util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addPermittedValueSet() {
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
	_ = s.localFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, fData, nil, nil)
}

func (s *ElectricalConnectionSuite) addPermittedValueSetEmptyElli() {
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
	_ = s.localFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, fData, nil, nil)
}
