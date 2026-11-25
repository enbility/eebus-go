package mdsf

import (
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *MaMDSFSuite) Test_OperationModes() {
	data, err := s.sut.OperationModes(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.OperationModes(s.dhwCircuitEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addHvacData()

	data, err = s.sut.OperationModes(s.dhwCircuitEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(data))
	assert.Equal(s.T(), ucapi.HvacOperationModeTypeAuto, data[0])
}

func (s *MaMDSFSuite) Test_CurrentOperationMode() {
	data, err := s.sut.CurrentOperationMode(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), ucapi.HvacOperationModeType(""), data)

	data, err = s.sut.CurrentOperationMode(s.dhwCircuitEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), ucapi.HvacOperationModeType(""), data)

	s.addHvacData()

	data, err = s.sut.CurrentOperationMode(s.dhwCircuitEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), ucapi.HvacOperationModeTypeEco, data)
}

func (s *MaMDSFSuite) Test_IsOverrunActive() {
	data, err := s.sut.IsOverrunActive(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.False(s.T(), data)

	data, err = s.sut.IsOverrunActive(s.dhwCircuitEntity)
	assert.NotNil(s.T(), err)
	assert.False(s.T(), data)

	s.addHvacData()
	s.addOverrunData()

	data, err = s.sut.IsOverrunActive(s.dhwCircuitEntity)
	assert.Nil(s.T(), err)
	assert.True(s.T(), data)
}

func (s *MaMDSFSuite) Test_OverrunStatus() {
	data, err := s.sut.OverrunStatus(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), model.HvacOverrunStatusType(""), data)

	data, err = s.sut.OverrunStatus(s.dhwCircuitEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), model.HvacOverrunStatusType(""), data)

	s.addHvacData()

	data, err = s.sut.OverrunStatus(s.dhwCircuitEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), model.HvacOverrunStatusType(""), data)

	s.addOverrunData()

	data, err = s.sut.OverrunStatus(s.dhwCircuitEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), model.HvacOverrunStatusTypeRunning, data)
}

// helpers

func (s *MaMDSFSuite) addOverrunData() {
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.dhwCircuitEntity, model.FeatureTypeTypeHvac, model.RoleTypeServer)

	descData := &model.HvacOverrunDescriptionListDataType{
		HvacOverrunDescriptionData: []model.HvacOverrunDescriptionDataType{
			{
				OverrunId:                util.Ptr(model.HvacOverrunIdType(1)),
				OverrunType:              util.Ptr(model.HvacOverrunTypeTypeOneTimeDhw),
				AffectedSystemFunctionId: []model.HvacSystemFunctionIdType{1},
			},
		},
	}
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeHvacOverrunDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	overrunData := &model.HvacOverrunListDataType{
		HvacOverrunData: []model.HvacOverrunDataType{
			{
				OverrunId:     util.Ptr(model.HvacOverrunIdType(1)),
				OverrunStatus: util.Ptr(model.HvacOverrunStatusTypeRunning),
			},
		},
	}
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeHvacOverrunListData, overrunData, nil, nil)
	assert.Nil(s.T(), fErr)

	functionData := &model.HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []model.HvacSystemFunctionDataType{
			{
				SystemFunctionId:       util.Ptr(model.HvacSystemFunctionIdType(1)),
				CurrentOperationModeId: util.Ptr(model.HvacOperationModeIdType(3)),
				IsOverrunActive:        util.Ptr(true),
			},
		},
	}
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionListData, functionData, nil, nil)
	assert.Nil(s.T(), fErr)
}

func (s *MaMDSFSuite) addHvacData() {
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.dhwCircuitEntity, model.FeatureTypeTypeHvac, model.RoleTypeServer)

	descData := &model.HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []model.HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId:   util.Ptr(model.HvacSystemFunctionIdType(1)),
				SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeDhw),
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
