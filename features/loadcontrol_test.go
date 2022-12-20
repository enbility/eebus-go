package features

import (
	"testing"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestLoadControlSuite(t *testing.T) {
	suite.Run(t, new(LoadControlSuite))
}

type LoadControlSuite struct {
	suite.Suite

	localDevice  *spine.DeviceLocalImpl
	remoteEntity *spine.EntityRemoteImpl

	loadControl *LoadControl
	sentMessage []byte
}

var _ spine.SpineDataConnection = (*LoadControlSuite)(nil)

func (s *LoadControlSuite) WriteSpineMessage(message []byte) {
	s.sentMessage = message
}

func (s *LoadControlSuite) BeforeTest(suiteName, testName string) {
	s.localDevice, s.remoteEntity = setupFeatures(
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
	s.loadControl, err = NewLoadControl(model.RoleTypeServer, model.RoleTypeClient, s.localDevice, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.loadControl)
}

func (s *LoadControlSuite) Test_RequestLimitDescription() {
	err := s.loadControl.RequestLimitDescription()
	assert.Nil(s.T(), err)
}

func (s *LoadControlSuite) Test_RequestLimitConstraints() {
	err := s.loadControl.RequestLimitConstraints()
	assert.Nil(s.T(), err)
}

func (s *LoadControlSuite) Test_RequestLimits() {
	counter, err := s.loadControl.RequestLimits()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_GetLimitDescriptions() {
	data, err := s.loadControl.GetLimitDescriptions()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.loadControl.GetLimitDescriptions()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *LoadControlSuite) Test_GetLimitDescriptionsForCategory() {
	data, err := s.loadControl.GetLimitDescriptionsForCategory(model.LoadControlCategoryTypeObligation)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.loadControl.GetLimitDescriptionsForCategory(model.LoadControlCategoryTypeOptimization)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.loadControl.GetLimitDescriptionsForCategory(model.LoadControlCategoryTypeObligation)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *LoadControlSuite) Test_GetLimitDescriptionsForMeasurementId() {
	measurementId := model.MeasurementIdType(0)
	data, err := s.loadControl.GetLimitDescriptionsForMeasurementId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.loadControl.GetLimitDescriptionsForMeasurementId(measurementId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	measurementId = model.MeasurementIdType(10)
	data, err = s.loadControl.GetLimitDescriptionsForMeasurementId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *LoadControlSuite) Test_WriteLimitValues() {
	counter, err := s.loadControl.WriteLimitValues(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.LoadControlLimitDataType{}
	counter, err = s.loadControl.WriteLimitValues(data)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data = []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			Value:   model.NewScaledNumberType(10),
		},
	}
	counter, err = s.loadControl.WriteLimitValues(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_GetLimitData() {
	data, err := s.loadControl.GetLimitData()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.loadControl.GetLimitData()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.loadControl.GetLimitData()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *LoadControlSuite) Test_GetLimitDataForLimitId() {
	limitId := model.LoadControlLimitIdType(0)
	data, err := s.loadControl.GetLimitDataForLimitId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.loadControl.GetLimitDataForLimitId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.loadControl.GetLimitDataForLimitId(limitId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	limitId = model.LoadControlLimitIdType(10)
	data, err = s.loadControl.GetLimitDataForLimitId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

// helper

func (s *LoadControlSuite) addDescription() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(0)),
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, fData, nil, nil)
}

func (s *LoadControlSuite) addData() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(true),
				Value:             model.NewScaledNumberType(12),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeLoadControlLimitListData, fData, nil, nil)
}
