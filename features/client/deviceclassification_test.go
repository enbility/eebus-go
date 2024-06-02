package client

import (
	"testing"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDeviceClassificationSuite(t *testing.T) {
	suite.Run(t, new(DeviceClassificationSuite))
}

type DeviceClassificationSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	deviceClassification *DeviceClassification
	sentMessage          []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*DeviceClassificationSuite)(nil)

func (s *DeviceClassificationSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *DeviceClassificationSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeDeviceClassification,
				functions: []model.FunctionType{
					model.FunctionTypeDeviceClassificationManufacturerData,
				},
			},
		},
	)

	var err error
	s.deviceClassification, err = NewDeviceClassification(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.deviceClassification)

	s.deviceClassification, err = NewDeviceClassification(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.deviceClassification)
}

func (s *DeviceClassificationSuite) Test_RequestManufacturerDetails() {
	counter, err := s.deviceClassification.RequestManufacturerDetails()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
