package client

import (
	"testing"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestElectricalConnectionSuite(t *testing.T) {
	suite.Run(t, new(ElectricalConnectionSuite))
}

type ElectricalConnectionSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	electricalConnection *ElectricalConnection
	sentMessage          []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*ElectricalConnectionSuite)(nil)

func (s *ElectricalConnectionSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *ElectricalConnectionSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeElectricalConnection,
				functions: []model.FunctionType{
					model.FunctionTypeElectricalConnectionDescriptionListData,
					model.FunctionTypeElectricalConnectionParameterDescriptionListData,
					model.FunctionTypeElectricalConnectionPermittedValueSetListData,
					model.FunctionTypeElectricalConnectionCharacteristicListData,
				},
			},
		},
	)

	var err error
	s.electricalConnection, err = NewElectricalConnection(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.electricalConnection)

	s.electricalConnection, err = NewElectricalConnection(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.electricalConnection)
}

func (s *ElectricalConnectionSuite) Test_RequestDescriptions() {
	counter, err := s.electricalConnection.RequestDescriptions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.electricalConnection.RequestDescriptions(
		&model.ElectricalConnectionDescriptionListDataSelectorsType{},
		&model.ElectricalConnectionDescriptionDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *ElectricalConnectionSuite) Test_RequestParameterDescriptions() {
	counter, err := s.electricalConnection.RequestParameterDescriptions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.electricalConnection.RequestParameterDescriptions(
		&model.ElectricalConnectionParameterDescriptionListDataSelectorsType{},
		&model.ElectricalConnectionParameterDescriptionDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *ElectricalConnectionSuite) Test_RequestPermittedValueSets() {
	counter, err := s.electricalConnection.RequestPermittedValueSets(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.electricalConnection.RequestPermittedValueSets(
		&model.ElectricalConnectionPermittedValueSetListDataSelectorsType{},
		&model.ElectricalConnectionPermittedValueSetDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *ElectricalConnectionSuite) Test_RequestCharacteristics() {
	counter, err := s.electricalConnection.RequestCharacteristics(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.electricalConnection.RequestCharacteristics(
		&model.ElectricalConnectionCharacteristicListDataSelectorsType{},
		&model.ElectricalConnectionCharacteristicDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
