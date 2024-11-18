package internal

import (
	"testing"
	"time"

	ucapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *InternalSuite) Test_LoadControlLimits() {
	var data []ucapi.LoadLimitsPhase
	var err error
	limitType := model.LoadControlLimitTypeTypeMaxValueLimit
	scope := model.ScopeTypeTypeSelfConsumption
	category := model.LoadControlCategoryTypeObligation

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(limitType),
		LimitCategory: util.Ptr(category),
		ScopeType:     util.Ptr(scope),
	}
	data, err = LoadControlLimits(nil, nil, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = LoadControlLimits(s.localEntity, s.mockRemoteEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				LimitType:     util.Ptr(limitType),
				ScopeType:     util.Ptr(scope),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				LimitType:     util.Ptr(limitType),
				ScopeType:     util.Ptr(scope),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(2)),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				LimitType:     util.Ptr(limitType),
				ScopeType:     util.Ptr(scope),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(10)),
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
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(2)),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), 2, len(data))

	permData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(0),
						},
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

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(data))
	assert.Equal(s.T(), 16.0, data[0].Value)
}

func (s *InternalSuite) Test_LoadControlLimits_AudiMobileConnect_1Phase() {
	var data []ucapi.LoadLimitsPhase
	var err error
	limitType := model.LoadControlLimitTypeTypeMaxValueLimit
	scope := model.ScopeTypeTypeSelfConsumption
	category := model.LoadControlCategoryTypeObligation

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(limitType),
		LimitCategory: util.Ptr(category),
		ScopeType:     util.Ptr(scope),
	}
	data, err = LoadControlLimits(nil, nil, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = LoadControlLimits(s.localEntity, s.mockRemoteEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(1)),
				LimitType:      util.Ptr(limitType),
				LimitCategory:  util.Ptr(category),
				LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
				MeasurementId:  util.Ptr(model.MeasurementIdType(1)),
				Unit:           util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:      util.Ptr(scope),
			},
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(2)),
				LimitType:      util.Ptr(limitType),
				LimitCategory:  util.Ptr(model.LoadControlCategoryTypeRecommendation),
				LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
				MeasurementId:  util.Ptr(model.MeasurementIdType(1)),
				Unit:           util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:      util.Ptr(model.ScopeTypeTypeSelfConsumption),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(10)),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(1)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(4)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			},
			{
				ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(3)),
				MeasurementId:           util.Ptr(model.MeasurementIdType(7)),
				VoltageType:             util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:        util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasuredInReferenceTo: util.Ptr(model.ElectricalConnectionPhaseNameTypeNeutral),
				AcMeasurementType:       util.Ptr(model.ElectricalConnectionAcMeasurementTypeTypeReal),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(8)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(0),
			},
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(2)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(0),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	permData := &model.ElectricalConnectionPermittedValueSetListDataType{
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
								Min: model.NewScaledNumberType(6),
								Max: model.NewScaledNumberType(10),
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
								Min: model.NewScaledNumberType(1437),
								Max: model.NewScaledNumberType(2395),
							},
						},
					},
				},
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))
	assert.Equal(s.T(), 10.0, data[0].Value)
	assert.Equal(s.T(), false, data[0].IsActive)
}

func (s *InternalSuite) Test_LoadControlLimits_Bender_1Phase() {
	var data []ucapi.LoadLimitsPhase
	var err error
	limitType := model.LoadControlLimitTypeTypeMaxValueLimit
	scope := model.ScopeTypeTypeOverloadProtection
	category := model.LoadControlCategoryTypeObligation

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(limitType),
		LimitCategory: util.Ptr(category),
		ScopeType:     util.Ptr(scope),
	}
	data, err = LoadControlLimits(nil, nil, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = LoadControlLimits(s.localEntity, s.mockRemoteEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(0)),
				LimitType:     util.Ptr(limitType),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:     util.Ptr(scope),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				LimitType:     util.Ptr(limitType),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:     util.Ptr(scope),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(2)),
				LimitType:     util.Ptr(limitType),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:     util.Ptr(scope),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
				AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
				AcMeasurementVariant:   util.Ptr(model.ElectricalConnectionMeasurandVariantTypeRms),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(3)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(2)),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	// according to OpEV Spec 1.0.1b, page 30: "At least one set of permitted values SHALL be stated."
	// which is not the case here for all elements
	permData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(6),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(3)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(0),
						},
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(1362),
								Max: model.NewScaledNumberType(3632),
							},
						},
					},
				},
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))
	assert.Equal(s.T(), 16.0, data[0].Value)
}

func (s *InternalSuite) Test_LoadControlLimits_Elli_1Phase() {
	var data []ucapi.LoadLimitsPhase
	var err error
	limitType := model.LoadControlLimitTypeTypeMaxValueLimit
	scope := model.ScopeTypeTypeOverloadProtection
	scopeSelf := model.ScopeTypeTypeSelfConsumption
	category := model.LoadControlCategoryTypeObligation
	categoryRec := model.LoadControlCategoryTypeRecommendation

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(limitType),
		LimitCategory: util.Ptr(category),
		ScopeType:     util.Ptr(scope),
	}
	data, err = LoadControlLimits(nil, nil, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = LoadControlLimits(s.localEntity, s.mockRemoteEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(0)),
				LimitType:     util.Ptr(limitType),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:     util.Ptr(scope),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				LimitType:     util.Ptr(limitType),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:     util.Ptr(scope),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(2)),
				LimitType:     util.Ptr(limitType),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:     util.Ptr(scope),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(3)),
				LimitType:     util.Ptr(limitType),
				LimitCategory: util.Ptr(categoryRec),
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:     util.Ptr(scopeSelf),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(4)),
				LimitType:     util.Ptr(limitType),
				LimitCategory: util.Ptr(categoryRec),
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:     util.Ptr(scopeSelf),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(5)),
				LimitType:     util.Ptr(limitType),
				LimitCategory: util.Ptr(categoryRec),
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:     util.Ptr(scopeSelf),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(2)),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(3)),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(4)),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(5)),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	// according to OpEV Spec 1.0.1b, page 30: "At least one set of permitted values SHALL be stated."
	// which is not the case here for all elements
	permData := &model.ElectricalConnectionPermittedValueSetListDataType{
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
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(3)),
			},
		},
	}

	_, fErr = rElFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *InternalSuite) Test_WriteLoadControlLimit() {
	loadLimit := ucapi.LoadLimit{
		Duration: time.Minute * 2,
		IsActive: true,
		Value:    5000,
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}

	msgCounter, err := WriteLoadControlLimit(nil, nil, filter, loadLimit, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = WriteLoadControlLimit(s.localEntity, s.mockRemoteEntity, filter, loadLimit, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = WriteLoadControlLimit(s.localEntity, s.monitoredEntity, filter, loadLimit, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
				MeasurementId:  util.Ptr(model.MeasurementIdType(0)),
				LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
				LimitDirection: util.Ptr(model.EnergyDirectionTypeProduce),
			},
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(1)),
				LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
				MeasurementId:  util.Ptr(model.MeasurementIdType(0)),
				LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
				LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}
	lc := s.monitoredEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	_, _ = lc.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)

	msgCounter, err = WriteLoadControlLimit(s.localEntity, s.monitoredEntity, filter, loadLimit, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	data := &model.LoadControlLimitListDataType{LoadControlLimitData: []model.LoadControlLimitDataType{}}
	_, _ = lc.UpdateData(true, model.FunctionTypeLoadControlLimitListData, data, nil, nil)

	msgCounter, err = WriteLoadControlLimit(s.localEntity, s.monitoredEntity, filter, loadLimit, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	data = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitChangeable: util.Ptr(false),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(0),
			},
		},
	}
	_, _ = lc.UpdateData(true, model.FunctionTypeLoadControlLimitListData, data, nil, nil)

	msgCounter, err = WriteLoadControlLimit(s.localEntity, s.monitoredEntity, filter, loadLimit, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	data = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			},
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(0),
			},
		},
	}
	_, _ = lc.UpdateData(true, model.FunctionTypeLoadControlLimitListData, data, nil, nil)

	s.mux.Lock()
	cbInvoked := false
	s.mux.Unlock()
	cb := func(result model.ResultDataType) {
		s.mux.Lock()
		cbInvoked = true
		s.mux.Unlock()
	}
	msgCounter, err = WriteLoadControlLimit(s.localEntity, s.monitoredEntity, filter, loadLimit, cb)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	lf := s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
	msg := &spineapi.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter:          util.Ptr(model.MsgCounterType(100)),
			MsgCounterReference: msgCounter,
		},
		CmdClassifier: model.CmdClassifierTypeResult,
		Cmd: model.CmdType{
			ResultData: &model.ResultDataType{
				ErrorNumber: util.Ptr(model.ErrorNumberType(1)),
			},
		},
		FeatureRemote: lc,
		EntityRemote:  s.monitoredEntity,
		DeviceRemote:  s.remoteDevice,
	}
	lf.HandleMessage(msg)
	time.Sleep(time.Millisecond * 200)
	s.mux.Lock()
	assert.True(s.T(), cbInvoked)
	s.mux.Unlock()

	loadLimit = ucapi.LoadLimit{
		Duration:       0,
		IsActive:       true,
		Value:          5000,
		DeleteDuration: true,
	}

	msgCounter, err = WriteLoadControlLimit(s.localEntity, s.monitoredEntity, filter, loadLimit, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
}

func (s *InternalSuite) Test_WriteLoadControlLimits() {
	loadLimits := []ucapi.LoadLimitsPhase{}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
		Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
		ScopeType:     util.Ptr(model.ScopeTypeTypeOverloadProtection),
	}

	msgCounter, err := WriteLoadControlPhaseLimits(nil, nil, filter, loadLimits, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = WriteLoadControlPhaseLimits(s.localEntity, s.mockRemoteEntity, filter, loadLimits, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = WriteLoadControlPhaseLimits(s.localEntity, s.monitoredEntity, filter, loadLimits, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(10)),
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
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	msgCounter, err = WriteLoadControlPhaseLimits(s.localEntity, s.monitoredEntity, filter, loadLimits, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	type dataStruct struct {
		phases                 int
		permittedDefaultExists bool
		permittedDefaultValue  float64
		permittedMinValue      float64
		permittedMaxValue      float64
		limits, limitsExpected []float64
	}

	tests := []struct {
		name string
		data []dataStruct
	}{
		{
			"1 Phase ISO15118",
			[]dataStruct{
				{1, true, 0.1, 2, 16, []float64{0}, []float64{0.1}},
				{1, true, 0.1, 2, 16, []float64{2.2}, []float64{2.2}},
				{1, true, 0.1, 2, 16, []float64{10}, []float64{10}},
				{1, true, 0.1, 2, 16, []float64{16}, []float64{16}},
			},
		},
		{
			"3 Phase ISO15118",
			[]dataStruct{
				{3, true, 0.1, 2, 16, []float64{0, 0, 0}, []float64{0.1, 0.1, 0.1}},
				{3, true, 0.1, 2, 16, []float64{2.2, 2.2, 2.2}, []float64{2.2, 2.2, 2.2}},
				{3, true, 0.1, 2, 16, []float64{10, 10, 10}, []float64{10, 10, 10}},
				{3, true, 0.1, 2, 16, []float64{16, 16, 16}, []float64{16, 16, 16}},
			},
		},
		{
			"1 Phase IEC61851",
			[]dataStruct{
				{1, true, 0, 6, 16, []float64{0}, []float64{0}},
				{1, true, 0, 6, 16, []float64{6}, []float64{6}},
				{1, true, 0, 6, 16, []float64{10}, []float64{10}},
				{1, true, 0, 6, 16, []float64{16}, []float64{16}},
			},
		},
		{
			"3 Phase IEC61851",
			[]dataStruct{
				{3, true, 0, 6, 16, []float64{0, 0, 0}, []float64{0, 0, 0}},
				{3, true, 0, 6, 16, []float64{6, 6, 6}, []float64{6, 6, 6}},
				{3, true, 0, 6, 16, []float64{10, 10, 10}, []float64{10, 10, 10}},
				{3, true, 0, 6, 16, []float64{16, 16, 16}, []float64{16, 16, 16}},
			},
		},
		{
			"3 Phase IEC61851 Elli",
			[]dataStruct{
				{3, false, 0, 6, 16, []float64{0, 0, 0}, []float64{0, 0, 0}},
				{3, false, 0, 6, 16, []float64{6, 6, 6}, []float64{6, 6, 6}},
				{3, false, 0, 6, 16, []float64{10, 10, 10}, []float64{10, 10, 10}},
				{3, false, 0, 6, 16, []float64{16, 16, 16}, []float64{16, 16, 16}},
			},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			dataSet := []model.ElectricalConnectionPermittedValueSetDataType{}
			permittedData := []model.ScaledNumberSetType{}
			for _, data := range tc.data {
				// clean up data
				remoteLoadControlF := s.monitoredEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
				assert.NotNil(s.T(), remoteLoadControlF)

				emptyLimits := model.LoadControlLimitListDataType{}
				_, errT := remoteLoadControlF.UpdateData(true, model.FunctionTypeLoadControlLimitListData, &emptyLimits, nil, nil)
				assert.Nil(s.T(), errT)

				for phase := 0; phase < data.phases; phase++ {
					item := model.ScaledNumberSetType{
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(data.permittedMinValue),
								Max: model.NewScaledNumberType(data.permittedMaxValue),
							},
						},
					}
					if data.permittedDefaultExists {
						item.Value = []model.ScaledNumberType{*model.NewScaledNumberType(data.permittedDefaultValue)}
					}
					permittedData = append(permittedData, item)

					permittedItem := model.ElectricalConnectionPermittedValueSetDataType{
						ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
						ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(phase)),
						PermittedValueSet:      permittedData,
					}
					dataSet = append(dataSet, permittedItem)
				}

				permData := &model.ElectricalConnectionPermittedValueSetListDataType{
					ElectricalConnectionPermittedValueSetData: dataSet,
				}

				_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
				assert.Nil(s.T(), fErr)

				msgCounter, err := WriteLoadControlPhaseLimits(s.localEntity, s.monitoredEntity, filter, loadLimits, nil)
				assert.NotNil(t, err)
				assert.Nil(t, msgCounter)

				limitDesc := []model.LoadControlLimitDescriptionDataType{}
				for index := range data.limits {
					id := model.LoadControlLimitIdType(index)
					limitItem := model.LoadControlLimitDescriptionDataType{
						LimitId:       util.Ptr(id),
						LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
						MeasurementId: util.Ptr(model.MeasurementIdType(index)),
						LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
						Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
						ScopeType:     util.Ptr(model.ScopeTypeTypeOverloadProtection),
					}
					limitDesc = append(limitDesc, limitItem)
				}
				add := len(limitDesc)
				for index := range data.limits {
					id := model.LoadControlLimitIdType(index + add)
					limitItem := model.LoadControlLimitDescriptionDataType{
						LimitId:       util.Ptr(id),
						LimitCategory: util.Ptr(model.LoadControlCategoryTypeRecommendation),
						MeasurementId: util.Ptr(model.MeasurementIdType(index)),
						LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
						Unit:          util.Ptr(model.UnitOfMeasurementTypeA),
						ScopeType:     util.Ptr(model.ScopeTypeTypeSelfConsumption),
					}
					limitDesc = append(limitDesc, limitItem)
				}

				descData := &model.LoadControlLimitDescriptionListDataType{
					LoadControlLimitDescriptionData: limitDesc,
				}

				rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
				_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
				assert.Nil(s.T(), fErr)

				msgCounter, err = WriteLoadControlPhaseLimits(s.localEntity, s.monitoredEntity, filter, loadLimits, nil)
				assert.NotNil(t, err)
				assert.Nil(t, msgCounter)

				limitData := []model.LoadControlLimitDataType{}
				for index := range limitDesc {
					limitItem := model.LoadControlLimitDataType{
						LimitId:           util.Ptr(model.LoadControlLimitIdType(index)),
						IsLimitChangeable: util.Ptr(true),
						IsLimitActive:     util.Ptr(false),
						Value:             model.NewScaledNumberType(data.permittedMaxValue),
					}
					limitData = append(limitData, limitItem)
				}

				limitListData := &model.LoadControlLimitListDataType{
					LoadControlLimitData: limitData,
				}

				_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, limitListData, nil, nil)
				assert.Nil(s.T(), fErr)

				msgCounter, err = WriteLoadControlPhaseLimits(s.localEntity, s.monitoredEntity, filter, loadLimits, nil)
				assert.NotNil(t, err)
				assert.Nil(t, msgCounter)

				phaseLimitValues := []ucapi.LoadLimitsPhase{}
				for index, limit := range data.limits {
					phase := ucapi.PhaseNameMapping[index]
					phaseLimitValues = append(phaseLimitValues, ucapi.LoadLimitsPhase{
						Phase:    phase,
						IsActive: true,
						Value:    limit,
					})
				}

				s.mux.Lock()
				cbInvoked := false
				s.mux.Unlock()
				cb := func(result model.ResultDataType) {
					s.mux.Lock()
					cbInvoked = true
					s.mux.Unlock()
				}
				msgCounter, err = WriteLoadControlPhaseLimits(s.localEntity, s.monitoredEntity, filter, phaseLimitValues, cb)
				assert.Nil(t, err)
				assert.NotNil(t, msgCounter)

				msgCounter, err = WriteLoadControlPhaseLimits(s.localEntity, s.monitoredEntity, filter, phaseLimitValues, cb)
				assert.Nil(t, err)
				assert.NotNil(t, msgCounter)

				lf := s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
				msg := &spineapi.Message{
					RequestHeader: &model.HeaderType{
						MsgCounter:          util.Ptr(model.MsgCounterType(100)),
						MsgCounterReference: msgCounter,
					},
					CmdClassifier: model.CmdClassifierTypeResult,
					Cmd: model.CmdType{
						ResultData: &model.ResultDataType{
							ErrorNumber: util.Ptr(model.ErrorNumberType(1)),
						},
					},
					FeatureRemote: rFeature,
					EntityRemote:  s.monitoredEntity,
					DeviceRemote:  s.remoteDevice,
				}
				lf.HandleMessage(msg)
				time.Sleep(time.Millisecond * 200)
				s.mux.Lock()
				assert.True(s.T(), cbInvoked)
				s.mux.Unlock()
			}
		})
	}
}
