package cevc

import (
	"time"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemCEVCSuite) Test_CoordinatedChargingScenarios() {
	timeConst := &model.TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []model.TimeSeriesConstraintsDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				SlotCountMax: util.Ptr(model.TimeSeriesSlotCountType(30)),
			},
		},
	}

	rTimeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	_, fErr := rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesConstraintsListData, timeConst, nil, nil)
	assert.Nil(s.T(), fErr)

	timeDesc := &model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(1)),
				TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypeConstraints),
				TimeSeriesWriteable: util.Ptr(true),
				UpdateRequired:      util.Ptr(false),
				Unit:                util.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(2)),
				TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypePlan),
				TimeSeriesWriteable: util.Ptr(false),
				Unit:                util.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(3)),
				TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
				TimeSeriesWriteable: util.Ptr(false),
				Unit:                util.Ptr(model.UnitOfMeasurementTypeWh),
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesDescriptionListData, timeDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	incDesc := &model.IncentiveTableDescriptionDataType{
		IncentiveTableDescription: []model.IncentiveTableDescriptionType{
			{
				TariffDescription: &model.TariffDescriptionDataType{
					TariffId:        util.Ptr(model.TariffIdType(1)),
					TariffWriteable: util.Ptr(true),
					UpdateRequired:  util.Ptr(false),
					ScopeType:       util.Ptr(model.ScopeTypeTypeSimpleIncentiveTable),
				},
			},
		},
	}

	rIncentiveFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer)
	_, fErr = rIncentiveFeature.UpdateData(true, model.FunctionTypeIncentiveTableDescriptionData, incDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	// demand, No Profile No Timer demand

	timeData := &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(3)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Value:            model.NewScaledNumberType(0),
						MaxValue:         model.NewScaledNumberType(74690),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	demand, err := s.sut.EnergyDemand(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0.0, demand.MinDemand)
	assert.Equal(s.T(), 0.0, demand.OptDemand)
	assert.Equal(s.T(), 74690.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), 0.0, demand.DurationUntilEnd)

	// the final plan

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(2)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(model.DurationType("PT18H3M7S")),
						MaxValue:         model.NewScaledNumberType(4163),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(model.DurationType("PT42M")),
						MaxValue:         model.NewScaledNumberType(2736),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(model.DurationType("P1D")),
						MaxValue:         model.NewScaledNumberType(0),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	// demand, profile + timer with 80% target and no climate, minSoC reached

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(3)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(model.DurationType("P2DT4H40M36S")),
						Value:            model.NewScaledNumberType(53400),
						MaxValue:         model.NewScaledNumberType(74690),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(model.DurationType("P1D")),
						MaxValue:         model.NewScaledNumberType(0),
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
	assert.Equal(s.T(), 53400.0, demand.OptDemand)
	assert.Equal(s.T(), 74690.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), time.Duration(time.Hour*52+time.Minute*40+time.Second*36).Seconds(), demand.DurationUntilEnd)

	// the final plan

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(2)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(model.DurationType("P1DT15H24M24S")),
						MaxValue:         model.NewScaledNumberType(0),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(model.DurationType("PT12H35M50S")),
						MaxValue:         model.NewScaledNumberType(4163),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(2)),
						Duration:         util.Ptr(model.DurationType("PT40M22S")),
						MaxValue:         model.NewScaledNumberType(0),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	// demand, profile with 25% min SoC, minSoC not reached, no timer

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
			},
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(2)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(model.DurationType("PT8M42S")),
						MaxValue:         model.NewScaledNumberType(4212),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(model.DurationType("P1D")),
						MaxValue:         model.NewScaledNumberType(0),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(3)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Value:            model.NewScaledNumberType(600),
						MinValue:         model.NewScaledNumberType(600),
						MaxValue:         model.NewScaledNumberType(75600),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	demand, err = s.sut.EnergyDemand(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 600.0, demand.MinDemand)
	assert.Equal(s.T(), 600.0, demand.OptDemand)
	assert.Equal(s.T(), 75600.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), 0.0, demand.DurationUntilEnd)

	// the final plan

	timeData = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(2)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(model.DurationType("PT8M42S")),
						MaxValue:         model.NewScaledNumberType(4212),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(model.DurationType("P1D")),
						MaxValue:         model.NewScaledNumberType(0),
					},
				},
			},
		},
	}

	_, fErr = rTimeFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)
}

/*
func requestIncentiveUpdate(t *testing.T, datagram model.DatagramType, localDevice api.DeviceLocal, remoteDevice api.DeviceRemote) {
	cmd := []model.CmdType{{
		IncentiveTableDescriptionData: &model.IncentiveTableDescriptionDataType{
			IncentiveTableDescription: []model.IncentiveTableDescriptionType{
				{
					TariffDescription: &model.TariffDescriptionDataType{
						TariffId:        util.Ptr(model.TariffIdType(1)),
						TariffWriteable: util.Ptr(true),
						UpdateRequired:  util.Ptr(true),
						ScopeType:       util.Ptr(model.ScopeTypeTypeSimpleIncentiveTable),
					},
				},
			},
		},
	}}

	datagram.Payload.Cmd = cmd

	err := localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)
}

func requestPowerTableUpdate(t *testing.T, datagram model.DatagramType, localDevice api.DeviceLocal, remoteDevice api.DeviceRemote) {
	cmd := []model.CmdType{{
		TimeSeriesDescriptionListData: &model.TimeSeriesDescriptionListDataType{
			TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
				{
					TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(1)),
					TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypeConstraints),
					TimeSeriesWriteable: util.Ptr(true),
					UpdateRequired:      util.Ptr(true),
				},
				{
					TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(2)),
					TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypePlan),
					TimeSeriesWriteable: util.Ptr(false),
					Unit:                util.Ptr(model.UnitOfMeasurementTypeW),
				},
				{
					TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(3)),
					TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypeConstraints),
					TimeSeriesWriteable: util.Ptr(false),
					Unit:                util.Ptr(model.UnitOfMeasurementTypeWh),
				},
			},
		},
	}}

	datagram.Payload.Cmd = cmd

	err := localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)
}
*/
