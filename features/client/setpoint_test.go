package client

import (
	"testing"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSetpointSuite(t *testing.T) {
	suite.Run(t, new(SetpointSuite))
}

type SetpointSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	setpoint    *Setpoint
	sentMessage []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*SetpointSuite)(nil)

func (s *SetpointSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *SetpointSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeSetpoint,
				functions: []model.FunctionType{
					model.FunctionTypeSetpointDescriptionListData,
					model.FunctionTypeSetpointConstraintsListData,
					model.FunctionTypeSetpointListData,
				},
			},
		},
	)

	var err error
	s.setpoint, err = NewSetpoint(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.setpoint)

	s.setpoint, err = NewSetpoint(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.setpoint)
}

func (s *SetpointSuite) Test_RequestSetpointDescriptions() {
	counter, err := s.setpoint.RequestSetpointDescriptions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *SetpointSuite) Test_RequestSetpointConstraints() {
	counter, err := s.setpoint.RequestSetpointConstraints(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *SetpointSuite) Test_RequestSetpoints() {
	counter, err := s.setpoint.RequestSetpoints(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *SetpointSuite) Test_WriteSetpointListData() {
	counter, err := s.setpoint.WriteSetpointListData(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.SetpointDataType{}
	counter, err = s.setpoint.WriteSetpointListData(data)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data = []model.SetpointDataType{
		{
			SetpointId: util.Ptr(model.SetpointIdType(1)),
			Value:      model.NewScaledNumberType(21),
		},
	}
	counter, err = s.setpoint.WriteSetpointListData(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
