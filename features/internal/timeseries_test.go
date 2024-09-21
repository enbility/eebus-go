package internal_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/features/internal"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestTimeSeriesSuite(t *testing.T) {
	suite.Run(t, new(TimeSeriesSuite))
}

type TimeSeriesSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.TimeSeriesCommon
}

func (s *TimeSeriesSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
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

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalTimeSeries(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteTimeSeries(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *TimeSeriesSuite) Test_GetData() {
	filter := model.TimeSeriesDescriptionDataType{}
	data, err := s.localSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addDescription()
	s.addData()

	data, err = s.localSut.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

func (s *TimeSeriesSuite) Test_GetDataForId() {
	filter := model.TimeSeriesDescriptionDataType{
		TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
	}
	data, err := s.localSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.localSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	filter.TimeSeriesType = util.Ptr(model.TimeSeriesTypeTypePlan)
	data, err = s.localSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *TimeSeriesSuite) Test_GetDescriptions() {
	filter := model.TimeSeriesDescriptionDataType{}
	data, err := s.localSut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addDescription()

	data, err = s.localSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

func (s *TimeSeriesSuite) Test_GetDescriptionsForId() {
	filter := model.TimeSeriesDescriptionDataType{
		TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
	}
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

	filter.TimeSeriesId = util.Ptr(model.TimeSeriesIdType(1))
	data, err = s.localSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *TimeSeriesSuite) Test_GetDescriptionForType() {
	filter := model.TimeSeriesDescriptionDataType{
		TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
	}
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

func (s *TimeSeriesSuite) Test_GetConstraints() {
	data, err := s.localSut.GetConstraints()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))
	data, err = s.remoteSut.GetConstraints()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addConstraints()

	data, err = s.localSut.GetConstraints()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
	data, err = s.remoteSut.GetConstraints()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

// helpers

func (s *TimeSeriesSuite) addData() {
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
	_ = s.localFeature.UpdateData(model.FunctionTypeTimeSeriesListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeTimeSeriesListData, fData, nil, nil)
}

func (s *TimeSeriesSuite) addDescription() {
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
	_ = s.localFeature.UpdateData(model.FunctionTypeTimeSeriesDescriptionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeTimeSeriesDescriptionListData, fData, nil, nil)
}

func (s *TimeSeriesSuite) addConstraints() {
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
	_ = s.localFeature.UpdateData(model.FunctionTypeTimeSeriesConstraintsListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeTimeSeriesConstraintsListData, fData, nil, nil)
}
