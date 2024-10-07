package server_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	features "github.com/enbility/eebus-go/features/server"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestFeatureSuite(t *testing.T) {
	suite.Run(t, new(FeatureSuite))
}

type FeatureSuite struct {
	suite.Suite

	service api.ServiceInterface

	localEntity spineapi.EntityLocalInterface

	remoteDevice spineapi.DeviceRemoteInterface
	remoteEntity spineapi.EntityRemoteInterface

	testFeature *features.Feature
	sentMessage []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*FeatureSuite)(nil)

func (s *FeatureSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *FeatureSuite) BeforeTest(suiteName, testName string) {
	cert, _ := cert.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := api.NewConfiguration(
		"test", "test", "test", "test",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		9999, cert, time.Second*4)

	serviceHandler := mocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	_ = s.service.Setup()
	s.localEntity = s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	var entities []spineapi.EntityRemoteInterface

	s.remoteDevice, entities = setupFeatures(s.service, s.T())
	s.remoteEntity = entities[1]

	var err error
	s.testFeature, err = features.NewFeature(model.FeatureTypeTypeLoadControl, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.testFeature)

	s.testFeature, err = features.NewFeature(model.FeatureTypeTypeLoadControl, s.localEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.testFeature)
}

func (s *FeatureSuite) Test_NewFeature() {
	newFeature, err := features.NewFeature(model.FeatureTypeTypeBill, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), newFeature)

	newFeature, err = features.NewFeature(model.FeatureTypeTypeBill, s.localEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), newFeature)

	newFeature, err = features.NewFeature(model.FeatureTypeTypeDeviceConfiguration, s.localEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), newFeature)
}
