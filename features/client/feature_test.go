package client_test

import (
	"testing"

	features "github.com/enbility/eebus-go/features/client"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestFeatureSuite(t *testing.T) {
	suite.Run(t, new(FeatureSuite))
}

type FeatureSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	testFeature *features.Feature
	sentMessage []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*FeatureSuite)(nil)

func (s *FeatureSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *FeatureSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeAlarm,
				functions: []model.FunctionType{
					model.FunctionTypeAlarmListData,
				},
			},
		},
	)

	var err error
	s.testFeature, err = features.NewFeature(model.FeatureTypeTypeAlarm, s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.testFeature)
}

func (s *FeatureSuite) Test_NewFeature() {
	newFeature, err := features.NewFeature(model.FeatureTypeTypeBill, nil, s.remoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), newFeature)

	newFeature, err = features.NewFeature(model.FeatureTypeTypeBill, s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), newFeature)

	newFeature, err = features.NewFeature(model.FeatureTypeTypeBill, s.localEntity, s.remoteEntity)
	assert.NotNil(s.T(), err)
	assert.NotNil(s.T(), newFeature)

	f := spine.NewFeatureLocal(1, s.localEntity, model.FeatureTypeTypeBill, model.RoleTypeClient)
	s.localEntity.AddFeature(f)

	newFeature, err = features.NewFeature(model.FeatureTypeTypeBill, s.localEntity, s.remoteEntity)
	assert.NotNil(s.T(), err)
	assert.NotNil(s.T(), newFeature)
}

func (s *FeatureSuite) Test_Subscription() {
	subscription := s.testFeature.HasSubscription()
	assert.Equal(s.T(), false, subscription)

	counter, err := s.testFeature.Subscribe()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	subscription = s.testFeature.HasSubscription()
	assert.Equal(s.T(), true, subscription)

	counter, err = s.testFeature.Subscribe()
	assert.NotNil(s.T(), counter)
	assert.Nil(s.T(), err)

	counter, err = s.testFeature.Unsubscribe()
	assert.NotNil(s.T(), counter)
	assert.Nil(s.T(), err)

	subscription = s.testFeature.HasSubscription()
	assert.Equal(s.T(), false, subscription)
}

func (s *FeatureSuite) Test_Binding() {
	binding := s.testFeature.HasBinding()
	assert.Equal(s.T(), false, binding)

	counter, err := s.testFeature.Bind()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	binding = s.testFeature.HasBinding()
	assert.Equal(s.T(), true, binding)

	counter, err = s.testFeature.Bind()
	assert.NotNil(s.T(), counter)
	assert.Nil(s.T(), err)

	counter, err = s.testFeature.Unbind()
	assert.NotNil(s.T(), counter)
	assert.Nil(s.T(), err)

	binding = s.testFeature.HasBinding()
	assert.Equal(s.T(), false, binding)
}

func (s *FeatureSuite) Test_ResultCallback() {
	testFct := func(msg spineapi.ResponseMessage) {}
	err := s.testFeature.AddResponseCallback(10, testFct)
	assert.Nil(s.T(), err)

	s.testFeature.AddResultCallback(testFct)
}
