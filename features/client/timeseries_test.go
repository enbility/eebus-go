package client_test

import (
	"testing"

	features "github.com/enbility/eebus-go/features/client"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestTimeSeriesSuite(t *testing.T) {
	suite.Run(t, new(TimeSeriesSuite))
}

type TimeSeriesSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	timeSeries  *features.TimeSeries
	sentMessage []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*TimeSeriesSuite)(nil)

func (s *TimeSeriesSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *TimeSeriesSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
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
	s.timeSeries, err = features.NewTimeSeries(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.timeSeries)

	s.timeSeries, err = features.NewTimeSeries(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.timeSeries)
}

func (s *TimeSeriesSuite) Test_RequestDescription() {
	msgCounter, err := s.timeSeries.RequestDescriptions()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
}

func (s *TimeSeriesSuite) Test_RequestConstraints() {
	msgCounter, err := s.timeSeries.RequestConstraints()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
}

func (s *TimeSeriesSuite) Test_RequestData() {
	counter, err := s.timeSeries.RequestData()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *TimeSeriesSuite) Test_WriteData() {
	counter, err := s.timeSeries.WriteData(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.TimeSeriesDataType{}
	counter, err = s.timeSeries.WriteData(data)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data = []model.TimeSeriesDataType{
		{
			TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
		},
	}
	counter, err = s.timeSeries.WriteData(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
