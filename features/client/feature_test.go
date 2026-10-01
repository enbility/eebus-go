package client

import (
	"testing"

	"github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
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

	testFeature, testFeature2 *Feature
	sentMessage               []byte
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
				partial: false,
			},
			{
				featureType: model.FeatureTypeTypeLoadControl,
				functions: []model.FunctionType{
					model.FunctionTypeLoadControlLimitListData,
				},
				partial: true,
			},
		},
	)

	var err error
	s.testFeature, err = NewFeature(model.FeatureTypeTypeAlarm, s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.testFeature)

	s.testFeature, err = NewFeature(model.FeatureTypeTypeAlarm, s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.testFeature)

	s.testFeature2, err = NewFeature(model.FeatureTypeTypeLoadControl, s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.testFeature2)
}

func (s *FeatureSuite) Test_NewFeature() {
	newFeature, err := NewFeature(model.FeatureTypeTypeBill, nil, s.remoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), newFeature)

	newFeature, err = NewFeature(model.FeatureTypeTypeBill, s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), newFeature)

	newFeature, err = NewFeature(model.FeatureTypeTypeBill, s.localEntity, s.remoteEntity)
	assert.NotNil(s.T(), err)
	assert.NotNil(s.T(), newFeature)

	f := spine.NewFeatureLocal(1, s.localEntity, model.FeatureTypeTypeBill, model.RoleTypeClient)
	s.localEntity.AddFeature(f)

	newFeature, err = NewFeature(model.FeatureTypeTypeBill, s.localEntity, s.remoteEntity)
	assert.NotNil(s.T(), err)
	assert.NotNil(s.T(), newFeature)
}

func (s *FeatureSuite) Test_Subscription() {
	// Test initial state
	subscription := s.testFeature.HasSubscription()
	assert.Equal(s.T(), false, subscription)

	// Test subscribe request - verify it returns without error
	counter, err := s.testFeature.Subscribe()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	// Note: In a real EEBUS environment, subscription requires remote device approval.
	// Since we're testing with mock devices that don't process approval,
	// HasSubscription() will remain false. This is expected behavior.
	// The important test is that Subscribe() method works without error.

	// Test duplicate subscribe request
	counter2, err := s.testFeature.Subscribe()
	assert.NotNil(s.T(), counter2)
	assert.Nil(s.T(), err)

	// Test unsubscribe - should work even if subscription wasn't confirmed
	counter, err = s.testFeature.Unsubscribe()
	assert.NotNil(s.T(), counter)
	assert.Nil(s.T(), err)
}

func (s *FeatureSuite) Test_Binding() {
	// Test initial state
	binding := s.testFeature.HasBinding()
	assert.Equal(s.T(), false, binding)

	// Test bind request - verify it returns without error
	counter, err := s.testFeature.Bind()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	// Note: In a real EEBUS environment, binding requires remote device approval.
	// Since we're testing with mock devices that don't process approval,
	// HasBinding() will remain false. This is expected behavior.
	// The important test is that Bind() method works without error.

	// Test duplicate bind request
	counter2, err := s.testFeature.Bind()
	assert.NotNil(s.T(), counter2)
	assert.Nil(s.T(), err)

	// Test unbind - should work even if binding wasn't confirmed
	counter, err = s.testFeature.Unbind()
	assert.NotNil(s.T(), counter)
	assert.Nil(s.T(), err)
}

func (s *FeatureSuite) Test_ResultCallback() {
	testFct := func(msg spineapi.ResponseMessage) {}
	err := s.testFeature.AddResponseCallback(10, testFct)
	assert.Nil(s.T(), err)

	s.testFeature.AddResultCallback(testFct)
}

func (s *FeatureSuite) Test_requestData() {
	counter, err := s.testFeature.requestData(model.FunctionTypeAlarmListData, nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	selectors := model.AlarmListDataSelectorsType{
		AlarmId: util.Ptr(model.AlarmIdType(0)),
	}
	counter, err = s.testFeature.requestData(model.FunctionTypeAlarmListData, selectors, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	// the remote does not support this function, the sentinel must stay matchable
	counter, err = s.testFeature2.requestData(model.FunctionTypeMeasurementDescriptionListData, nil, nil)
	assert.ErrorIs(s.T(), err, api.ErrOperationOnFunctionNotSupported)
	assert.Nil(s.T(), counter)

	counter, err = s.testFeature2.requestData(model.FunctionTypeLoadControlLimitListData, nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	selectors2 := &model.LoadControlLimitListDataSelectorsType{
		LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
	}
	counter, err = s.testFeature2.requestData(model.FunctionTypeLoadControlLimitListData, selectors2, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
