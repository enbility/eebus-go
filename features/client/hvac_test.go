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

func TestHvacSuite(t *testing.T) {
	suite.Run(t, new(HvacSuite))
}

type HvacSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	hvac        *Hvac
	sentMessage []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*HvacSuite)(nil)

func (s *HvacSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *HvacSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeHvac,
				functions: []model.FunctionType{
					model.FunctionTypeHvacSystemFunctionDescriptionListData,
					model.FunctionTypeHvacSystemFunctionListData,
					model.FunctionTypeHvacOperationModeDescriptionListData,
					model.FunctionTypeHvacSystemFunctionOperationModeRelationListData,
					model.FunctionTypeHvacSystemFunctionSetPointRelationListData,
					model.FunctionTypeHvacOverrunDescriptionListData,
					model.FunctionTypeHvacOverrunListData,
				},
			},
		},
	)

	var err error
	s.hvac, err = NewHvac(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.hvac)

	s.hvac, err = NewHvac(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.hvac)
}

func (s *HvacSuite) Test_RequestHvacSystemFunctionDescriptions() {
	counter, err := s.hvac.RequestHvacSystemFunctionDescriptions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *HvacSuite) Test_RequestHvacSystemFunctions() {
	counter, err := s.hvac.RequestHvacSystemFunctions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *HvacSuite) Test_RequestHvacOperationModeDescriptions() {
	counter, err := s.hvac.RequestHvacOperationModeDescriptions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *HvacSuite) Test_RequestHvacSystemFunctionOperationModeRelations() {
	counter, err := s.hvac.RequestHvacSystemFunctionOperationModeRelations(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *HvacSuite) Test_RequestHvacSystemFunctionSetpointRelations() {
	counter, err := s.hvac.RequestHvacSystemFunctionSetpointRelations(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *HvacSuite) Test_WriteHvacSystemFunctionListData() {
	counter, err := s.hvac.WriteHvacSystemFunctionListData(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.HvacSystemFunctionDataType{}
	counter, err = s.hvac.WriteHvacSystemFunctionListData(data)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data = []model.HvacSystemFunctionDataType{
		{
			SystemFunctionId:       util.Ptr(model.HvacSystemFunctionIdType(1)),
			CurrentOperationModeId: util.Ptr(model.HvacOperationModeIdType(2)),
		},
	}
	counter, err = s.hvac.WriteHvacSystemFunctionListData(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *HvacSuite) Test_RequestHvacOverrunDescriptions() {
	counter, err := s.hvac.RequestHvacOverrunDescriptions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *HvacSuite) Test_RequestHvacOverruns() {
	counter, err := s.hvac.RequestHvacOverruns(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *HvacSuite) Test_WriteHvacOverrunListData() {
	counter, err := s.hvac.WriteHvacOverrunListData(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.HvacOverrunDataType{
		{
			OverrunId:     util.Ptr(model.HvacOverrunIdType(1)),
			OverrunStatus: util.Ptr(model.HvacOverrunStatusTypeActive),
		},
	}
	counter, err = s.hvac.WriteHvacOverrunListData(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
