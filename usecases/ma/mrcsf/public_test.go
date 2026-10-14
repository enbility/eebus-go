package mrcsf

import (
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *MaMRCSFSuite) Test_OperationModes() {
	data, err := s.sut.OperationModes(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.OperationModes(s.hvacRoomEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addHvacData()

	data, err = s.sut.OperationModes(s.hvacRoomEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(data))
	assert.Equal(s.T(), ucapi.HvacOperationModeTypeAuto, data[0])
}

func (s *MaMRCSFSuite) Test_CurrentOperationMode() {
	data, err := s.sut.CurrentOperationMode(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), ucapi.HvacOperationModeType(""), data)

	data, err = s.sut.CurrentOperationMode(s.hvacRoomEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), ucapi.HvacOperationModeType(""), data)

	s.addHvacData()

	data, err = s.sut.CurrentOperationMode(s.hvacRoomEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), ucapi.HvacOperationModeTypeEco, data)
}

// helpers

func (s *MaMRCSFSuite) addHvacData() {
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.hvacRoomEntity, model.FeatureTypeTypeHvac, model.RoleTypeServer)

	descData := &model.HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []model.HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId:   util.Ptr(model.HvacSystemFunctionIdType(1)),
				SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeCooling),
			},
		},
	}
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	modeData := &model.HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: []model.HvacOperationModeDescriptionDataType{
			{
				OperationModeId:   util.Ptr(model.HvacOperationModeIdType(1)),
				OperationModeType: util.Ptr(model.HvacOperationModeTypeTypeAuto),
			},
			{
				OperationModeId:   util.Ptr(model.HvacOperationModeIdType(2)),
				OperationModeType: util.Ptr(model.HvacOperationModeTypeTypeOn),
			},
			{
				OperationModeId:   util.Ptr(model.HvacOperationModeIdType(3)),
				OperationModeType: util.Ptr(model.HvacOperationModeTypeTypeEco),
			},
		},
	}
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeHvacOperationModeDescriptionListData, modeData, nil, nil)
	assert.Nil(s.T(), fErr)

	relationData := &model.HvacSystemFunctionOperationModeRelationListDataType{
		HvacSystemFunctionOperationModeRelationData: []model.HvacSystemFunctionOperationModeRelationDataType{
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(1)),
				OperationModeId:  []model.HvacOperationModeIdType{1, 2, 3},
			},
		},
	}
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionOperationModeRelationListData, relationData, nil, nil)
	assert.Nil(s.T(), fErr)

	functionData := &model.HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []model.HvacSystemFunctionDataType{
			{
				SystemFunctionId:       util.Ptr(model.HvacSystemFunctionIdType(1)),
				CurrentOperationModeId: util.Ptr(model.HvacOperationModeIdType(3)),
			},
		},
	}
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionListData, functionData, nil, nil)
	assert.Nil(s.T(), fErr)
}
