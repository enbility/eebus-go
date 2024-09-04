package cevc

import (
	"testing"
	"time"

	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemCEVCSuite) Test_TimeSlotConstraints() {
	constraints, err := s.sut.TimeSlotConstraints(s.mockRemoteEntity)
	assert.Equal(s.T(), uint(0), constraints.MinSlots)
	assert.Equal(s.T(), uint(0), constraints.MaxSlots)
	assert.Equal(s.T(), time.Duration(0), constraints.MinSlotDuration)
	assert.Equal(s.T(), time.Duration(0), constraints.MaxSlotDuration)
	assert.Equal(s.T(), time.Duration(0), constraints.SlotDurationStepSize)
	assert.NotEqual(s.T(), err, nil)

	constraints, err = s.sut.TimeSlotConstraints(s.evEntity)
	assert.Equal(s.T(), uint(0), constraints.MinSlots)
	assert.Equal(s.T(), uint(0), constraints.MaxSlots)
	assert.Equal(s.T(), time.Duration(0), constraints.MinSlotDuration)
	assert.Equal(s.T(), time.Duration(0), constraints.MaxSlotDuration)
	assert.Equal(s.T(), time.Duration(0), constraints.SlotDurationStepSize)
	assert.NotEqual(s.T(), err, nil)

	constData := &model.TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []model.TimeSeriesConstraintsDataType{
			{
				TimeSeriesId:         util.Ptr(model.TimeSeriesIdType(0)),
				SlotCountMin:         util.Ptr(model.TimeSeriesSlotCountType(1)),
				SlotCountMax:         util.Ptr(model.TimeSeriesSlotCountType(10)),
				SlotDurationMin:      model.NewDurationType(1 * time.Minute),
				SlotDurationMax:      model.NewDurationType(60 * time.Minute),
				SlotDurationStepSize: model.NewDurationType(1 * time.Minute),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeTimeSeriesConstraintsListData, constData, nil, nil)
	assert.Nil(s.T(), fErr)

	constraints, err = s.sut.TimeSlotConstraints(s.evEntity)
	assert.Equal(s.T(), uint(1), constraints.MinSlots)
	assert.Equal(s.T(), uint(10), constraints.MaxSlots)
	assert.Equal(s.T(), time.Duration(1*time.Minute), constraints.MinSlotDuration)
	assert.Equal(s.T(), time.Duration(1*time.Hour), constraints.MaxSlotDuration)
	assert.Equal(s.T(), time.Duration(1*time.Minute), constraints.SlotDurationStepSize)
	assert.Equal(s.T(), err, nil)
}

func (s *CemCEVCSuite) Test_WritePowerLimits() {
	data := []ucapi.DurationSlotValue{}

	err := s.sut.WritePowerLimits(s.mockRemoteEntity, data)
	assert.NotNil(s.T(), err)

	err = s.sut.WritePowerLimits(s.evEntity, data)
	assert.NotNil(s.T(), err)

	elParamDesc := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	err = s.sut.WritePowerLimits(s.evEntity, data)
	assert.NotNil(s.T(), err)

	elPermDesc := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, elPermDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	err = s.sut.WritePowerLimits(s.evEntity, data)
	assert.NotNil(s.T(), err)

	elPermDesc = &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, elPermDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	err = s.sut.WritePowerLimits(s.evEntity, data)
	assert.NotNil(s.T(), err)

	descData := &model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:   util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeConstraints),
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeTimeSeriesDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	err = s.sut.WritePowerLimits(s.evEntity, data)
	assert.NotNil(s.T(), err)

	type dataStruct struct {
		error              bool
		minSlots, maxSlots uint
		slots              []ucapi.DurationSlotValue
	}

	tests := []struct {
		name string
		data []dataStruct
	}{
		{
			"too few slots",
			[]dataStruct{
				{
					true, 2, 2,
					[]ucapi.DurationSlotValue{
						{Duration: time.Hour, Value: 11000},
					},
				},
			},
		}, {
			"too many slots",
			[]dataStruct{
				{
					true, 1, 1,
					[]ucapi.DurationSlotValue{
						{Duration: time.Hour, Value: 11000},
						{Duration: time.Hour, Value: 11000},
					},
				},
			},
		},
		{
			"1 slot",
			[]dataStruct{
				{
					false, 1, 1,
					[]ucapi.DurationSlotValue{
						{Duration: time.Hour, Value: 11000},
					},
				},
			},
		},
		{
			"2 slots",
			[]dataStruct{
				{
					false, 1, 2,
					[]ucapi.DurationSlotValue{
						{Duration: time.Hour, Value: 11000},
						{Duration: 30 * time.Minute, Value: 5000},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			for _, data := range tc.data {
				constData := &model.TimeSeriesConstraintsListDataType{
					TimeSeriesConstraintsData: []model.TimeSeriesConstraintsDataType{
						{
							TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
							SlotCountMin: util.Ptr(model.TimeSeriesSlotCountType(data.minSlots)),
							SlotCountMax: util.Ptr(model.TimeSeriesSlotCountType(data.maxSlots)),
						},
					},
				}

				_, fErr := rFeature.UpdateData(true, model.FunctionTypeTimeSeriesConstraintsListData, constData, nil, nil)
				assert.Nil(s.T(), fErr)

				err = s.sut.WritePowerLimits(s.evEntity, data.slots)
				if data.error {
					assert.NotNil(t, err)
					continue
				}

				assert.Nil(t, err)
			}
		})
	}
}
