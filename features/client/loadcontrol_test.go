package client

import (
	"testing"
	"time"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestLoadControlSuite(t *testing.T) {
	suite.Run(t, new(LoadControlSuite))
}

type LoadControlSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	loadControl *LoadControl
	sentMessage []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*LoadControlSuite)(nil)

func (s *LoadControlSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *LoadControlSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeLoadControl,
				functions: []model.FunctionType{
					model.FunctionTypeLoadControlLimitDescriptionListData,
					model.FunctionTypeLoadControlLimitConstraintsListData,
					model.FunctionTypeLoadControlLimitListData,
				},
			},
		},
	)

	var err error
	s.loadControl, err = NewLoadControl(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.loadControl)

	s.loadControl, err = NewLoadControl(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.loadControl)
}

func (s *LoadControlSuite) Test_RequestLimitDescription() {
	counter, err := s.loadControl.RequestLimitDescriptions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.loadControl.RequestLimitDescriptions(
		&model.LoadControlLimitDescriptionListDataSelectorsType{},
		&model.LoadControlLimitDescriptionDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_RequestLimitConstraints() {
	counter, err := s.loadControl.RequestLimitConstraints(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.loadControl.RequestLimitConstraints(
		&model.LoadControlLimitConstraintsListDataSelectorsType{},
		&model.LoadControlLimitConstraintsDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_RequestLimits() {
	counter, err := s.loadControl.RequestLimitData(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.loadControl.RequestLimitData(
		&model.LoadControlLimitListDataSelectorsType{},
		&model.LoadControlLimitDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_WriteLimitValues() {
	counter, err := s.loadControl.WriteLimitData(nil, nil, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.LoadControlLimitDataType{}
	counter, err = s.loadControl.WriteLimitData(data, nil, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data = []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			Value:   model.NewScaledNumberType(10),
			TimePeriod: &model.TimePeriodType{
				EndTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(time.Minute * 5),
			},
		},
	}
	counter, err = s.loadControl.WriteLimitData(data, nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	deleteSelectors := &model.LoadControlLimitListDataSelectorsType{
		LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
	}
	deleteElements := &model.LoadControlLimitDataElementsType{
		TimePeriod: &model.TimePeriodElementsType{},
	}
	counter, err = s.loadControl.WriteLimitData(data, deleteSelectors, deleteElements)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
