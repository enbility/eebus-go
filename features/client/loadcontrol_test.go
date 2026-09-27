package client

import (
	"testing"
	"time"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestLoadControlSuite(t *testing.T) {
	suite.Run(t, new(LoadControlSuite))
}

type LoadControlSuite struct {
	suite.Suite

	localEntity        spineapi.EntityLocalInterface
	localEntityPartial spineapi.EntityLocalInterface

	remoteEntity        spineapi.EntityRemoteInterface
	remoteEntityPartial spineapi.EntityRemoteInterface

	loadControl        *LoadControl
	loadControlPartial *LoadControl

	sentMessage []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*LoadControlSuite)(nil)

func (s *LoadControlSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *LoadControlSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
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
				partial: false,
			},
		},
	)

	s.localEntityPartial, s.remoteEntityPartial = setupFeatures(
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
				partial: true,
			},
		},
	)

	var err error
	s.loadControl, err = NewLoadControl(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.loadControl)

	s.loadControl, err = NewLoadControl(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.loadControl)

	s.loadControlPartial, err = NewLoadControl(s.localEntityPartial, s.remoteEntityPartial)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.loadControl)
}

func (s *LoadControlSuite) Test_RequestLimitDescription() {
	counter, err := s.loadControl.RequestLimitDescriptions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.loadControl.RequestLimitDescriptions(
		&model.LoadControlLimitDescriptionListDataSelectorsType{},
		&model.LoadControlLimitDescriptionDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_RequestLimitConstraints() {
	counter, err := s.loadControl.RequestLimitConstraints(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.loadControl.RequestLimitConstraints(
		&model.LoadControlLimitConstraintsListDataSelectorsType{},
		&model.LoadControlLimitConstraintsDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_RequestLimits() {
	counter, err := s.loadControl.RequestLimitData(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.loadControl.RequestLimitData(
		&model.LoadControlLimitListDataSelectorsType{},
		&model.LoadControlLimitDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *LoadControlSuite) Test_WriteLimitValues() {
	counter, err := s.loadControl.WriteLimitData(nil, nil, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.LoadControlLimitDataType{}
	counter, err = s.loadControl.WriteLimitData(data, nil, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	rF := s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	data1 := rF.DataCopy(model.FunctionTypeLoadControlLimitListData).(*model.LoadControlLimitListDataType)
	assert.Nil(s.T(), data1)

	defaultData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(16),
				TimePeriod:        nil,
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(10),
				TimePeriod: &model.TimePeriodType{
					EndTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(time.Minute * 5),
				},
			},
		},
	}
	_, err1 := rF.UpdateData(true, model.FunctionTypeLoadControlLimitListData, defaultData, nil, nil)
	assert.Nil(s.T(), err1)
	data1 = rF.DataCopy(model.FunctionTypeLoadControlLimitListData).(*model.LoadControlLimitListDataType)
	assert.NotNil(s.T(), data1)
	assert.Equal(s.T(), 2, len(data1.LoadControlLimitData))

	data = []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			Value:   model.NewScaledNumberType(10),
			TimePeriod: &model.TimePeriodType{
				EndTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(time.Minute * 5),
			},
		},
	}
	counter, err = s.loadControl.WriteLimitData(data, nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

// a full write persists the merged list right away and restores it when rejected
func (s *LoadControlSuite) Test_WriteLimitValues_Persist() {
	rF := s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

	_, fErr := rF.UpdateData(true, model.FunctionTypeLoadControlLimitListData, &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{LimitId: util.Ptr(model.LoadControlLimitIdType(1)), IsLimitActive: util.Ptr(false), Value: model.NewScaledNumberType(16)},
			{LimitId: util.Ptr(model.LoadControlLimitIdType(2)), IsLimitActive: util.Ptr(true), Value: model.NewScaledNumberType(6)},
		},
	}, nil, nil)
	assert.Nil(s.T(), fErr)

	counter, err := s.loadControl.WriteLimitData([]model.LoadControlLimitDataType{
		{LimitId: util.Ptr(model.LoadControlLimitIdType(1)), IsLimitActive: util.Ptr(true), Value: model.NewScaledNumberType(10)},
	}, nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	limits := func() []model.LoadControlLimitDataType {
		return rF.DataCopy(model.FunctionTypeLoadControlLimitListData).(*model.LoadControlLimitListDataType).LoadControlLimitData
	}

	// the written limit is known before the remote notifies, the other one is kept
	data := limits()
	assert.Equal(s.T(), 2, len(data))
	assert.True(s.T(), *data[0].IsLimitActive)
	assert.Equal(s.T(), 10.0, data[0].Value.GetValue())
	assert.Equal(s.T(), 6.0, data[1].Value.GetValue())

	// a rejected write restores the previous state
	lF := s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
	lF.HandleMessage(&spineapi.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter:          util.Ptr(model.MsgCounterType(100)),
			MsgCounterReference: counter,
		},
		CmdClassifier: model.CmdClassifierTypeResult,
		Cmd: model.CmdType{
			ResultData: &model.ResultDataType{ErrorNumber: util.Ptr(model.ErrorNumberType(7))},
		},
		FeatureRemote: rF,
		EntityRemote:  s.remoteEntity,
		DeviceRemote:  s.remoteEntity.Device(),
	})
	assert.Eventually(s.T(), func() bool {
		data := limits()
		return !*data[0].IsLimitActive && data[0].Value.GetValue() == 16
	}, time.Second, 10*time.Millisecond)
}

// test with partial support
func (s *LoadControlSuite) Test_WriteLimitValues_Partial() {
	counter, err := s.loadControlPartial.WriteLimitData(nil, nil, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.LoadControlLimitDataType{}
	counter, err = s.loadControlPartial.WriteLimitData(data, nil, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	rF := s.remoteEntityPartial.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	data1 := rF.DataCopy(model.FunctionTypeLoadControlLimitListData).(*model.LoadControlLimitListDataType)
	assert.Nil(s.T(), data1)

	defaultData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(16),
				TimePeriod:        nil,
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(10),
				TimePeriod: &model.TimePeriodType{
					EndTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(time.Minute * 5),
				},
			},
		},
	}

	_, err1 := rF.UpdateData(true, model.FunctionTypeLoadControlLimitListData, defaultData, nil, nil)
	assert.Nil(s.T(), err1)
	data1 = rF.DataCopy(model.FunctionTypeLoadControlLimitListData).(*model.LoadControlLimitListDataType)
	assert.NotNil(s.T(), data1)
	assert.Equal(s.T(), 2, len(data1.LoadControlLimitData))

	data = []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			Value:   model.NewScaledNumberType(10),
			TimePeriod: &model.TimePeriodType{
				EndTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(time.Minute * 5),
			},
		},
	}
	counter, err = s.loadControlPartial.WriteLimitData(data, nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	data1 = rF.DataCopy(model.FunctionTypeLoadControlLimitListData).(*model.LoadControlLimitListDataType)
	assert.NotNil(s.T(), data1)
	assert.Equal(s.T(), 2, len(data1.LoadControlLimitData))

	deleteSelectors := &model.LoadControlLimitListDataSelectorsType{
		LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
	}
	deleteElements := &model.LoadControlLimitDataElementsType{
		TimePeriod: &model.TimePeriodElementsType{},
	}
	counter, err = s.loadControlPartial.WriteLimitData(data, deleteSelectors, deleteElements)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
