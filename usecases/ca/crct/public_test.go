package crct

import (
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CaCRCTSuite) Test_Setpoints() {
	data, err := s.sut.Setpoints(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.Setpoints(s.hvacRoomEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addHvacData()

	data, err = s.sut.Setpoints(s.hvacRoomEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addSetpointData()

	data, err = s.sut.Setpoints(s.hvacRoomEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))
	assert.Equal(s.T(), uint(1), data[0].Id)
	assert.Equal(s.T(), 21.0, data[0].Value)
	assert.True(s.T(), data[0].IsChangeable)
	assert.True(s.T(), data[1].IsActive)
}

func (s *CaCRCTSuite) Test_SetpointConstraints() {
	data, err := s.sut.SetpointConstraints(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.SetpointConstraints(s.hvacRoomEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addHvacData()
	s.addSetpointData()

	data, err = s.sut.SetpointConstraints(s.hvacRoomEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))
	assert.Equal(s.T(), uint(1), data[0].Id)
	assert.Equal(s.T(), 16.0, data[0].MinValue)
	assert.Equal(s.T(), 25.0, data[0].MaxValue)
	assert.Equal(s.T(), 0.5, data[0].StepSize)
}

func (s *CaCRCTSuite) Test_WriteSetpoint() {
	_, err := s.sut.WriteSetpoint(s.mockRemoteEntity, ucapi.HvacOperationModeTypeEco, 19, nil)
	assert.NotNil(s.T(), err)

	// the setpoints of the auto mode cannot be written
	_, err = s.sut.WriteSetpoint(s.hvacRoomEntity, ucapi.HvacOperationModeTypeAuto, 19, nil)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteSetpoint(s.hvacRoomEntity, ucapi.HvacOperationModeTypeEco, 19, nil)
	assert.NotNil(s.T(), err)

	s.addHvacData()

	_, err = s.sut.WriteSetpoint(s.hvacRoomEntity, ucapi.HvacOperationModeTypeEco, 19, nil)
	assert.Nil(s.T(), err)

	// the off mode has no setpoint in the test data
	_, err = s.sut.WriteSetpoint(s.hvacRoomEntity, ucapi.HvacOperationModeTypeOff, 19, nil)
	assert.NotNil(s.T(), err)

	// a setpoint marked not changeable cannot be written
	s.addSetpointData()

	_, err = s.sut.WriteSetpoint(s.hvacRoomEntity, ucapi.HvacOperationModeTypeOn, 22, nil)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteSetpoint(s.hvacRoomEntity, ucapi.HvacOperationModeTypeEco, 19, nil)
	assert.Nil(s.T(), err)
}

func (s *CaCRCTSuite) Test_WriteSetpoint_AmbiguousSystemFunction() {
	s.addHvacData()

	// two cooling system functions make the id ambiguous, so the write must fail
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.hvacRoomEntity, model.FeatureTypeTypeHvac, model.RoleTypeServer)
	descData := &model.HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []model.HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId:   util.Ptr(model.HvacSystemFunctionIdType(1)),
				SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeCooling),
			},
			{
				SystemFunctionId:   util.Ptr(model.HvacSystemFunctionIdType(2)),
				SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeCooling),
			},
		},
	}
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err := s.sut.WriteSetpoint(s.hvacRoomEntity, ucapi.HvacOperationModeTypeEco, 19, nil)
	assert.NotNil(s.T(), err)
}

// helpers

func (s *CaCRCTSuite) addHvacData() {
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
				OperationModeType: util.Ptr(model.HvacOperationModeTypeTypeEco),
			},
			{
				OperationModeId:   util.Ptr(model.HvacOperationModeIdType(2)),
				OperationModeType: util.Ptr(model.HvacOperationModeTypeTypeOn),
			},
			{
				OperationModeId:   util.Ptr(model.HvacOperationModeIdType(3)),
				OperationModeType: util.Ptr(model.HvacOperationModeTypeTypeOff),
			},
		},
	}
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeHvacOperationModeDescriptionListData, modeData, nil, nil)
	assert.Nil(s.T(), fErr)

	relationData := &model.HvacSystemFunctionSetpointRelationListDataType{
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
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionSetPointRelationListData, relationData, nil, nil)
	assert.Nil(s.T(), fErr)
}

func (s *CaCRCTSuite) addSetpointData() {
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.hvacRoomEntity, model.FeatureTypeTypeSetpoint, model.RoleTypeServer)

	spData := &model.SetpointListDataType{
		SetpointData: []model.SetpointDataType{
			{
				SetpointId:           util.Ptr(model.SetpointIdType(1)),
				Value:                model.NewScaledNumberType(21),
				ValueMin:             model.NewScaledNumberType(16),
				ValueMax:             model.NewScaledNumberType(25),
				IsSetpointChangeable: util.Ptr(true),
			},
			{
				SetpointId:           util.Ptr(model.SetpointIdType(2)),
				Value:                model.NewScaledNumberType(23),
				IsSetpointActive:     util.Ptr(true),
				IsSetpointChangeable: util.Ptr(false),
			},
		},
	}
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSetpointListData, spData, nil, nil)
	assert.Nil(s.T(), fErr)

	constraintsData := &model.SetpointConstraintsListDataType{
		SetpointConstraintsData: []model.SetpointConstraintsDataType{
			{
				SetpointId:       util.Ptr(model.SetpointIdType(1)),
				SetpointRangeMin: model.NewScaledNumberType(16),
				SetpointRangeMax: model.NewScaledNumberType(25),
				SetpointStepSize: model.NewScaledNumberType(0.5),
			},
		},
	}
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeSetpointConstraintsListData, constraintsData, nil, nil)
	assert.Nil(s.T(), fErr)
}
