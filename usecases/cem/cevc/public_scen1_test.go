package cevc

import (
	"time"

	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemCEVCSuite) Test_ChargeStrategy() {
	data := s.sut.ChargeStrategy(s.mockRemoteEntity)
	assert.Equal(s.T(), ucapi.EVChargeStrategyTypeUnknown, data)

	data = s.sut.ChargeStrategy(s.evEntity)
	assert.Equal(s.T(), ucapi.EVChargeStrategyTypeUnknown, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeCommunicationsStandard),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data = s.sut.ChargeStrategy(s.evEntity)
	assert.Equal(s.T(), ucapi.EVChargeStrategyTypeUnknown, data)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					String: util.Ptr(model.DeviceConfigurationKeyValueStringType(model.DeviceConfigurationKeyValueStringTypeISO151182ED2)),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data = s.sut.ChargeStrategy(s.evEntity)
	assert.Equal(s.T(), ucapi.EVChargeStrategyTypeUnknown, data)

	timeDescData := &model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:   util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
			},
		},
	}

	rTimeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesDescriptionListData, timeDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	timeData := &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	data = s.sut.ChargeStrategy(s.evEntity)
	assert.Equal(s.T(), ucapi.EVChargeStrategyTypeUnknown, data)

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	data = s.sut.ChargeStrategy(s.evEntity)
	assert.Equal(s.T(), ucapi.EVChargeStrategyTypeNoDemand, data)

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(model.DurationType("PT0S")),
						Value:            model.NewScaledNumberType(0),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	data = s.sut.ChargeStrategy(s.evEntity)
	assert.Equal(s.T(), ucapi.EVChargeStrategyTypeNoDemand, data)

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Value:            model.NewScaledNumberType(10000),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	data = s.sut.ChargeStrategy(s.evEntity)
	assert.Equal(s.T(), ucapi.EVChargeStrategyTypeDirectCharging, data)

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Value:            model.NewScaledNumberType(10000),
						Duration:         model.NewDurationType(2 * time.Hour),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	data = s.sut.ChargeStrategy(s.evEntity)
	assert.Equal(s.T(), ucapi.EVChargeStrategyTypeTimedCharging, data)
}

func (s *CemCEVCSuite) Test_EnergySingleDemand() {
	demand, err := s.sut.EnergyDemand(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, demand.MinDemand)
	assert.Equal(s.T(), 0.0, demand.OptDemand)
	assert.Equal(s.T(), 0.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), 0.0, demand.DurationUntilEnd)

	demand, err = s.sut.EnergyDemand(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, demand.MinDemand)
	assert.Equal(s.T(), 0.0, demand.OptDemand)
	assert.Equal(s.T(), 0.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), 0.0, demand.DurationUntilEnd)

	timeDescData := &model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:   util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
			},
		},
	}

	rTimeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	_, fErr := rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesDescriptionListData, timeDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	timeData := &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	demand, err = s.sut.EnergyDemand(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, demand.MinDemand)
	assert.Equal(s.T(), 0.0, demand.OptDemand)
	assert.Equal(s.T(), 0.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), 0.0, demand.DurationUntilEnd)

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						TimePeriod: &model.TimePeriodType{
							StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
						},
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	demand, err = s.sut.EnergyDemand(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0.0, demand.MinDemand)
	assert.Equal(s.T(), 0.0, demand.OptDemand)
	assert.Equal(s.T(), 0.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), 0.0, demand.DurationUntilEnd)

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						MinValue:         model.NewScaledNumberType(1000),
						Value:            model.NewScaledNumberType(10000),
						MaxValue:         model.NewScaledNumberType(100000),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	demand, err = s.sut.EnergyDemand(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1000.0, demand.MinDemand)
	assert.Equal(s.T(), 10000.0, demand.OptDemand)
	assert.Equal(s.T(), 100000.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), 0.0, demand.DurationUntilEnd)

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Value:            model.NewScaledNumberType(10000),
						Duration:         model.NewDurationType(2 * time.Hour),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	demand, err = s.sut.EnergyDemand(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0.0, demand.MinDemand)
	assert.Equal(s.T(), 10000.0, demand.OptDemand)
	assert.Equal(s.T(), 0.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), time.Duration(2*time.Hour).Seconds(), demand.DurationUntilEnd)
}
