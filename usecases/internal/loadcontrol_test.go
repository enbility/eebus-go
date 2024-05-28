package internal

import (
	"testing"

	ucapi "github.com/enbility/eebus-go/usecases/api"
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
	fErr := rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(data))
	assert.Equal(s.T(), 0.0, data[0].Value)

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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
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

	fErr = rFeature.UpdateData(model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

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

	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.localEntity, s.monitoredEntity, filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(data))
	assert.Equal(s.T(), 16.0, data[0].Value)
}

func (s *InternalSuite) Test_WriteLoadControlLimits() {
	loadLimits := []ucapi.LoadLimitsPhase{}

	category := model.LoadControlCategoryTypeObligation

	msgCounter, err := WriteLoadControlLimits(s.localEntity, s.mockRemoteEntity, category, loadLimits)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = WriteLoadControlLimits(s.localEntity, s.monitoredEntity, category, loadLimits)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	msgCounter, err = WriteLoadControlLimits(s.localEntity, s.monitoredEntity, category, loadLimits)
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
				errT := remoteLoadControlF.UpdateData(model.FunctionTypeLoadControlLimitListData, &emptyLimits, nil, nil)
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

				fErr = rFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
				assert.Nil(s.T(), fErr)

				msgCounter, err := WriteLoadControlLimits(s.localEntity, s.monitoredEntity, category, loadLimits)
				assert.NotNil(t, err)
				assert.Nil(t, msgCounter)

				limitDesc := []model.LoadControlLimitDescriptionDataType{}
				for index := range data.limits {
					id := model.LoadControlLimitIdType(index)
					limitItem := model.LoadControlLimitDescriptionDataType{
						LimitId:       util.Ptr(id),
						LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
						MeasurementId: util.Ptr(model.MeasurementIdType(index)),
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
					}
					limitDesc = append(limitDesc, limitItem)
				}

				descData := &model.LoadControlLimitDescriptionListDataType{
					LoadControlLimitDescriptionData: limitDesc,
				}

				rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
				fErr = rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
				assert.Nil(s.T(), fErr)

				msgCounter, err = WriteLoadControlLimits(s.localEntity, s.monitoredEntity, category, loadLimits)
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

				fErr = rFeature.UpdateData(model.FunctionTypeLoadControlLimitListData, limitListData, nil, nil)
				assert.Nil(s.T(), fErr)

				msgCounter, err = WriteLoadControlLimits(s.localEntity, s.monitoredEntity, category, loadLimits)
				assert.NotNil(t, err)
				assert.Nil(t, msgCounter)

				phaseLimitValues := []ucapi.LoadLimitsPhase{}
				for index, limit := range data.limits {
					phase := PhaseNameMapping[index]
					phaseLimitValues = append(phaseLimitValues, ucapi.LoadLimitsPhase{
						Phase:    phase,
						IsActive: true,
						Value:    limit,
					})
				}

				msgCounter, err = WriteLoadControlLimits(s.localEntity, s.monitoredEntity, category, phaseLimitValues)
				assert.Nil(t, err)
				assert.NotNil(t, msgCounter)

				msgCounter, err = WriteLoadControlLimits(s.localEntity, s.monitoredEntity, category, phaseLimitValues)
				assert.Nil(t, err)
				assert.NotNil(t, msgCounter)
			}
		})
	}
}
