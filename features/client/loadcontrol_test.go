package client_test

import (
	"testing"

	features "github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/util"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
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

	loadControl *features.LoadControl
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
	s.loadControl, err = features.NewLoadControl(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.loadControl)
}

func (s *LoadControlSuite) Test_RequestLimitDescription() {
	counter, err := s.loadControl.RequestLimitDescriptions()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_RequestLimitConstraints() {
	counter, err := s.loadControl.RequestLimitConstraints()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_RequestLimits() {
	counter, err := s.loadControl.RequestLimitData()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_WriteLimitValues() {
	counter, err := s.loadControl.WriteLimitData(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.LoadControlLimitDataType{}
	counter, err = s.loadControl.WriteLimitData(data)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data = []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			Value:   model.NewScaledNumberType(10),
		},
	}
	counter, err = s.loadControl.WriteLimitData(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
