package features

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestTimeSeriesSuite(t *testing.T) {
	suite.Run(t, new(TimeSeriesSuite))
}

type TimeSeriesSuite struct {
	suite.Suite

	localDevice  *spine.DeviceLocalImpl
	remoteEntity *spine.EntityRemoteImpl

	timeSeries  *TimeSeries
	sentMessage []byte
}

var _ spine.SpineDataConnection = (*TimeSeriesSuite)(nil)

func (s *TimeSeriesSuite) WriteSpineMessage(message []byte) {
	s.sentMessage = message
}

func (s *TimeSeriesSuite) BeforeTest(suiteName, testName string) {
	s.localDevice, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeTimeSeries,
				functions: []model.FunctionType{
					model.FunctionTypeTimeSeriesConstraintsListData,
					model.FunctionTypeTimeSeriesDescriptionListData,
					model.FunctionTypeTimeSeriesListData,
				},
			},
		},
	)

	var err error
	s.timeSeries, err = NewTimeSeries(model.RoleTypeServer, model.RoleTypeClient, s.localDevice, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.timeSeries)
}

func (s *TimeSeriesSuite) Test_RequestDescription() {
	err := s.timeSeries.RequestDescriptions()
	assert.Nil(s.T(), err)
}

func (s *TimeSeriesSuite) Test_RequestConstraints() {
	err := s.timeSeries.RequestConstraints()
	assert.Nil(s.T(), err)
}

func (s *TimeSeriesSuite) Test_RequestValues() {
	counter, err := s.timeSeries.RequestValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *TimeSeriesSuite) Test_GetValues() {
	data, err := s.timeSeries.GetValues()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addData()

	data, err = s.timeSeries.GetValues()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

func (s *TimeSeriesSuite) Test_GetDescriptions() {
	data, err := s.timeSeries.GetDescriptions()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addDescription()

	data, err = s.timeSeries.GetDescriptions()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

func (s *TimeSeriesSuite) Test_GetConstraints() {
	data, err := s.timeSeries.GetConstraints()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addConstraints()

	data, err = s.timeSeries.GetConstraints()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

// helpers

func (s *TimeSeriesSuite) addData() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))

	fData := &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
					EndTime:   model.NewAbsoluteOrRelativeTimeType("PT4H"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						Value:            model.NewScaledNumberType(10),
						MinValue:         model.NewScaledNumberType(6),
						MaxValue:         model.NewScaledNumberType(16),
						TimePeriod: &model.TimePeriodType{
							StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
							EndTime:   model.NewAbsoluteOrRelativeTimeType("PT1H"),
						},
						Duration: model.NewDurationType(1 * time.Hour),
					},
				},
			},
		},
	}
	rF.UpdateData(model.FunctionTypeTimeSeriesListData, fData, nil, nil)
}

func (s *TimeSeriesSuite) addDescription() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
				MeasurementId:       util.Ptr(model.MeasurementIdType(0)),
				TimeSeriesWriteable: util.Ptr(false),
				UpdateRequired:      util.Ptr(false),
				Unit:                util.Ptr(model.UnitOfMeasurementTypeWh),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeTimeSeriesDescriptionListData, fData, nil, nil)
}

func (s *TimeSeriesSuite) addConstraints() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []model.TimeSeriesConstraintsDataType{
			{
				TimeSeriesId:                util.Ptr(model.TimeSeriesIdType(0)),
				SlotCountMin:                util.Ptr(model.TimeSeriesSlotCountType(1)),
				SlotCountMax:                util.Ptr(model.TimeSeriesSlotCountType(24)),
				SlotDurationMin:             model.NewDurationType(1 * time.Hour),
				SlotDurationMax:             model.NewDurationType(2 * time.Hour),
				SlotDurationStepSize:        model.NewDurationType(15 * time.Minute),
				EarliestTimeSeriesStartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				LatestTimeSeriesEndTime:     model.NewAbsoluteOrRelativeTimeType("PT1H"),
				SlotValueMin:                model.NewScaledNumberType(2),
				SlotValueMax:                model.NewScaledNumberType(16),
				SlotValueStepSize:           model.NewScaledNumberType(0.1),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeTimeSeriesConstraintsListData, fData, nil, nil)
}
