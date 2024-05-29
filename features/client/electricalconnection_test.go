package client_test

import (
	"testing"

	features "github.com/enbility/eebus-go/features/client"
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

	electricalConnection *features.ElectricalConnection
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
	s.electricalConnection, err = features.NewElectricalConnection(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.electricalConnection)

	s.electricalConnection, err = features.NewElectricalConnection(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.electricalConnection)
}

func (s *ElectricalConnectionSuite) Test_RequestDescriptions() {
	counter, err := s.electricalConnection.RequestDescriptions()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *ElectricalConnectionSuite) Test_RequestParameterDescriptions() {
	counter, err := s.electricalConnection.RequestParameterDescriptions()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *ElectricalConnectionSuite) Test_RequestPermittedValueSets() {
	counter, err := s.electricalConnection.RequestPermittedValueSets()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *ElectricalConnectionSuite) Test_RequestCharacteristics() {
	counter, err := s.electricalConnection.RequestCharacteristics()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
