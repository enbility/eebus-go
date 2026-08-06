package internal_test

import (
	"testing"

	"github.com/enbility/eebus-go/features/internal"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestHvacSuite(t *testing.T) {
	suite.Run(t, new(HvacSuite))
}

type HvacSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.HvacCommon
}

func (s *HvacSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeHvac,
				functions: []model.FunctionType{
					model.FunctionTypeHvacSystemFunctionDescriptionListData,
					model.FunctionTypeHvacSystemFunctionListData,
					model.FunctionTypeHvacOperationModeDescriptionListData,
					model.FunctionTypeHvacSystemFunctionOperationModeRelationListData,
					model.FunctionTypeHvacSystemFunctionSetPointRelationListData,
				},
			},
		},
	)

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeHvac, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalHvac(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeHvac, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteHvac(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *HvacSuite) Test_GetHvacSystemFunctionDescriptions() {
	filter := model.HvacSystemFunctionDescriptionDataType{}
	data, err := s.localSut.GetHvacSystemFunctionDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetHvacSystemFunctionDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	desc, err := s.localSut.GetHvacSystemFunctionDescriptionForId(model.HvacSystemFunctionIdType(1))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)

	s.addDescriptions()

	data, err = s.localSut.GetHvacSystemFunctionDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))
	data, err = s.remoteSut.GetHvacSystemFunctionDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))

	filter.SystemFunctionType = util.Ptr(model.HvacSystemFunctionTypeTypeHeating)
	data, err = s.localSut.GetHvacSystemFunctionDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))

	desc, err = s.localSut.GetHvacSystemFunctionDescriptionForId(model.HvacSystemFunctionIdType(1))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
	desc, err = s.remoteSut.GetHvacSystemFunctionDescriptionForId(model.HvacSystemFunctionIdType(10))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)
}

func (s *HvacSuite) Test_GetHvacSystemFunctions() {
	filter := model.HvacSystemFunctionDataType{}
	data, err := s.localSut.GetHvacSystemFunctionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetHvacSystemFunctionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	function, err := s.localSut.GetHvacSystemFunctionForId(model.HvacSystemFunctionIdType(1))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), function)

	s.addData()

	data, err = s.localSut.GetHvacSystemFunctionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))
	data, err = s.remoteSut.GetHvacSystemFunctionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))

	function, err = s.localSut.GetHvacSystemFunctionForId(model.HvacSystemFunctionIdType(1))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), function)
	assert.Equal(s.T(), model.HvacOperationModeIdType(2), *function.CurrentOperationModeId)

	function, err = s.remoteSut.GetHvacSystemFunctionForId(model.HvacSystemFunctionIdType(10))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), function)
}

func (s *HvacSuite) Test_GetHvacOperationModeDescriptions() {
	filter := model.HvacOperationModeDescriptionDataType{}
	data, err := s.localSut.GetHvacOperationModeDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetHvacOperationModeDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	desc, err := s.localSut.GetHvacOperationModeDescriptionForId(model.HvacOperationModeIdType(1))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)

	s.addOperationModeDescriptions()

	data, err = s.localSut.GetHvacOperationModeDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))
	data, err = s.remoteSut.GetHvacOperationModeDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))

	filter.OperationModeType = util.Ptr(model.HvacOperationModeTypeTypeAuto)
	data, err = s.localSut.GetHvacOperationModeDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))

	desc, err = s.localSut.GetHvacOperationModeDescriptionForId(model.HvacOperationModeIdType(1))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
	desc, err = s.remoteSut.GetHvacOperationModeDescriptionForId(model.HvacOperationModeIdType(10))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)
}

func (s *HvacSuite) Test_GetHvacSystemFunctionOperationModeRelations() {
	filter := model.HvacSystemFunctionOperationModeRelationDataType{}
	data, err := s.localSut.GetHvacSystemFunctionOperationModeRelationsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetHvacSystemFunctionOperationModeRelationsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addOperationModeRelations()

	data, err = s.localSut.GetHvacSystemFunctionOperationModeRelationsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))
	data, err = s.remoteSut.GetHvacSystemFunctionOperationModeRelationsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))

	filter.SystemFunctionId = util.Ptr(model.HvacSystemFunctionIdType(1))
	data, err = s.localSut.GetHvacSystemFunctionOperationModeRelationsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))
}

func (s *HvacSuite) Test_GetHvacSystemFunctionSetpointRelations() {
	filter := model.HvacSystemFunctionSetpointRelationDataType{}
	data, err := s.localSut.GetHvacSystemFunctionSetpointRelationsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetHvacSystemFunctionSetpointRelationsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addSetpointRelations()

	data, err = s.localSut.GetHvacSystemFunctionSetpointRelationsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))
	data, err = s.remoteSut.GetHvacSystemFunctionSetpointRelationsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))

	filter.OperationModeId = util.Ptr(model.HvacOperationModeIdType(1))
	data, err = s.localSut.GetHvacSystemFunctionSetpointRelationsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))
}

func (s *HvacSuite) Test_CheckEventPayloadDataForFilter() {
	filter := model.HvacSystemFunctionDataType{
		SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(1)),
	}

	exists := s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	temp := true
	exists = s.localSut.CheckEventPayloadDataForFilter(temp, filter)
	assert.False(s.T(), exists)

	data := &model.HvacSystemFunctionListDataType{}
	exists = s.localSut.CheckEventPayloadDataForFilter(data, filter)
	assert.False(s.T(), exists)

	data = &model.HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []model.HvacSystemFunctionDataType{
			{
				SystemFunctionId:       util.Ptr(model.HvacSystemFunctionIdType(1)),
				CurrentOperationModeId: util.Ptr(model.HvacOperationModeIdType(2)),
			},
		},
	}
	exists = s.localSut.CheckEventPayloadDataForFilter(data, filter)
	assert.True(s.T(), exists)

	filter.SystemFunctionId = util.Ptr(model.HvacSystemFunctionIdType(10))
	exists = s.localSut.CheckEventPayloadDataForFilter(data, filter)
	assert.False(s.T(), exists)
}

// helpers

func (s *HvacSuite) addDescriptions() {
	fData := &model.HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []model.HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId:   util.Ptr(model.HvacSystemFunctionIdType(1)),
				SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeHeating),
			},
			{
				SystemFunctionId:   util.Ptr(model.HvacSystemFunctionIdType(2)),
				SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeDhw),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeHvacSystemFunctionDescriptionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionDescriptionListData, fData, nil, nil)
}

func (s *HvacSuite) addData() {
	fData := &model.HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []model.HvacSystemFunctionDataType{
			{
				SystemFunctionId:            util.Ptr(model.HvacSystemFunctionIdType(1)),
				CurrentOperationModeId:      util.Ptr(model.HvacOperationModeIdType(2)),
				IsOperationModeIdChangeable: util.Ptr(true),
				CurrentSetpointId:           util.Ptr(model.SetpointIdType(1)),
				IsSetpointIdChangeable:      util.Ptr(true),
				IsOverrunActive:             util.Ptr(false),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeHvacSystemFunctionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionListData, fData, nil, nil)
}

func (s *HvacSuite) addOperationModeDescriptions() {
	fData := &model.HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: []model.HvacOperationModeDescriptionDataType{
			{
				OperationModeId:   util.Ptr(model.HvacOperationModeIdType(1)),
				OperationModeType: util.Ptr(model.HvacOperationModeTypeTypeAuto),
			},
			{
				OperationModeId:   util.Ptr(model.HvacOperationModeIdType(2)),
				OperationModeType: util.Ptr(model.HvacOperationModeTypeTypeOff),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeHvacOperationModeDescriptionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeHvacOperationModeDescriptionListData, fData, nil, nil)
}

func (s *HvacSuite) addOperationModeRelations() {
	fData := &model.HvacSystemFunctionOperationModeRelationListDataType{
		HvacSystemFunctionOperationModeRelationData: []model.HvacSystemFunctionOperationModeRelationDataType{
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(1)),
				OperationModeId:  []model.HvacOperationModeIdType{1},
			},
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(2)),
				OperationModeId:  []model.HvacOperationModeIdType{2},
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeHvacSystemFunctionOperationModeRelationListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionOperationModeRelationListData, fData, nil, nil)
}

func (s *HvacSuite) addSetpointRelations() {
	fData := &model.HvacSystemFunctionSetpointRelationListDataType{
		HvacSystemFunctionSetpointRelationData: []model.HvacSystemFunctionSetpointRelationDataType{
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(1)),
				OperationModeId:  util.Ptr(model.HvacOperationModeIdType(1)),
				SetpointId:       []model.SetpointIdType{1},
			},
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(1)),
				OperationModeId:  util.Ptr(model.HvacOperationModeIdType(2)),
				SetpointId:       []model.SetpointIdType{2},
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeHvacSystemFunctionSetPointRelationListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionSetPointRelationListData, fData, nil, nil)
}
