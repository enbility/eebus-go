package cevc

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemCEVCSuite) Test_ChargePlanConstaints() {
	_, err := s.sut.ChargePlanConstraints(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.ChargePlanConstraints(s.evEntity)
	assert.NotNil(s.T(), err)

	descData := &model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:   util.Ptr(model.TimeSeriesIdType(1)),
				TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeConstraints),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeTimeSeriesDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.ChargePlanConstraints(s.evEntity)
	assert.NotNil(s.T(), err)

	data := &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.ChargePlanConstraints(s.evEntity)
	assert.NotNil(s.T(), err)

	data = &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(model.DurationType("PT5M36S")),
						MaxValue:         model.NewScaledNumberType(4201),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						TimePeriod: &model.TimePeriodType{
							StartTime: model.NewAbsoluteOrRelativeTimeType("PT30S"),
							EndTime:   model.NewAbsoluteOrRelativeTimeType("PT1M"),
						},
						MaxValue: model.NewScaledNumberType(4201),
					},
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.ChargePlanConstraints(s.evEntity)
	assert.Nil(s.T(), err)
}

func (s *CemCEVCSuite) Test_ChargePlan() {
	_, err := s.sut.ChargePlan(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.ChargePlan(s.evEntity)
	assert.NotNil(s.T(), err)

	descData := &model.TimeSeriesDescriptionListDataType{
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeTimeSeriesDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.ChargePlan(s.evEntity)
	assert.NotNil(s.T(), err)

	timeData := &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				TimePeriod:   &model.TimePeriodType{},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(model.DurationType("PT5M36S")),
						MaxValue:         model.NewScaledNumberType(4201),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						TimePeriod: &model.TimePeriodType{
							StartTime: model.NewAbsoluteOrRelativeTimeType("PT30S"),
							EndTime:   model.NewAbsoluteOrRelativeTimeType("PT1M"),
						},
						Value:    model.NewScaledNumberType(5),
						MinValue: model.NewScaledNumberType(0),
						MaxValue: model.NewScaledNumberType(10),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(2)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(model.DurationType("PT5M36S")),
						MaxValue:         model.NewScaledNumberType(4201),
					},
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
						TimePeriod: &model.TimePeriodType{
							StartTime: model.NewAbsoluteOrRelativeTimeType("PT30S"),
							EndTime:   model.NewAbsoluteOrRelativeTimeType("PT1M"),
						},
						Value:    model.NewScaledNumberType(5),
						MinValue: model.NewScaledNumberType(0),
						MaxValue: model.NewScaledNumberType(10),
					},
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.ChargePlan(s.evEntity)
	assert.Nil(s.T(), err)
}
