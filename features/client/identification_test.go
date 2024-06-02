package client

import (
	"testing"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestIdentificationSuite(t *testing.T) {
	suite.Run(t, new(IdentificationSuite))
}

type IdentificationSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	identification *Identification
	sentMessage    []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*IdentificationSuite)(nil)

func (s *IdentificationSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *IdentificationSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeIdentification,
				functions: []model.FunctionType{
					model.FunctionTypeIdentificationListData,
				},
			},
		},
	)

	var err error
	s.identification, err = NewIdentification(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.identification)

	s.identification, err = NewIdentification(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.identification)
}

func (s *IdentificationSuite) Test_RequestValues() {
	counter, err := s.identification.RequestValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
