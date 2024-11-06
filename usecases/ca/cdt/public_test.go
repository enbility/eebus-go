package cdt

import (
	"github.com/enbility/ship-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

// Test_Setpoints verifies the retrieval of setpoints from a remote entity.
func (s *CaCDTSuite) Test_Setpoints() {
	// Test case: No setpoints available for mock remote entity
	data, err := s.sut.Setpoints(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	// Test case: No setpoints available for CDT entity
	data, err = s.sut.Setpoints(s.cdtEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	// Prepare setpoint data
	setpointsData := &model.SetpointListDataType{}
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.cdtEntity, model.FeatureTypeTypeSetpoint, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSetpointListData, setpointsData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test case: No setpoints available after updating data
	data, err = s.sut.Setpoints(s.cdtEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	// Add setpoint data
	setpointsData.SetpointData = []model.SetpointDataType{
		{
			SetpointId: util.Ptr(model.SetpointIdType(1)),
			Value: util.Ptr(model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(34)),
				Scale:  util.Ptr(model.ScaleType(0)),
			}),
			ValueMin: util.Ptr(model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(20)),
				Scale:  util.Ptr(model.ScaleType(0)),
			}),
			ValueMax: util.Ptr(model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(70)),
				Scale:  util.Ptr(model.ScaleType(0)),
			}),
			ValueToleranceAbsolute: util.Ptr(model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(0)),
				Scale:  util.Ptr(model.ScaleType(0)),
			}),
			ValueTolerancePercentage: util.Ptr(model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(0)),
				Scale:  util.Ptr(model.ScaleType(0)),
			}),
			IsSetpointChangeable: util.Ptr(true),
			IsSetpointActive:     util.Ptr(true),
			TimePeriod:           util.Ptr(model.TimePeriodType{}),
		},
	}

	// Update data with new setpoint data
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeSetpointListData, setpointsData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test case: Setpoints available after updating data
	data, err = s.sut.Setpoints(s.cdtEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Len(s.T(), data, 1)
	assert.Equal(s.T(), uint(1), data[0].Id)
	assert.Equal(s.T(), 34.0, data[0].Value)
	assert.Equal(s.T(), 20.0, data[0].MinValue)
	assert.Equal(s.T(), 70.0, data[0].MaxValue)
	assert.True(s.T(), data[0].IsActive)
	assert.True(s.T(), data[0].IsChangeable)
	assert.Equal(s.T(), model.TimePeriodType{}, data[0].TimePeriod)
}

// Test_SetpointConstraints verifies the retrieval of setpoint constraints from a remote entity.
func (s *CaCDTSuite) Test_SetpointConstraints() {
	// Test case: No setpoint constraints available for mock remote entity
	data, err := s.sut.SetpointConstraints(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	// Test case: No setpoint constraints available for CDT entity
	data, err = s.sut.SetpointConstraints(s.cdtEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	// Prepare setpoint constraints data
	constraintsData := &model.SetpointConstraintsListDataType{}
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.cdtEntity, model.FeatureTypeTypeSetpoint, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSetpointConstraintsListData, constraintsData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test case: No setpoint constraints available after updating data
	data, err = s.sut.SetpointConstraints(s.cdtEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	// Add setpoint constraints data
	constraintsData.SetpointConstraintsData = []model.SetpointConstraintsDataType{
		{
			SetpointId: util.Ptr(model.SetpointIdType(1)),
			SetpointRangeMin: util.Ptr(model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(20)),
				Scale:  util.Ptr(model.ScaleType(0)),
			}),
			SetpointRangeMax: util.Ptr(model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(70)),
				Scale:  util.Ptr(model.ScaleType(0)),
			}),
			SetpointStepSize: util.Ptr(model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(1)),
				Scale:  util.Ptr(model.ScaleType(0)),
			}),
		},
	}

	// Update data with new setpoint constraints data
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeSetpointConstraintsListData, constraintsData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test case: Setpoint constraints available after updating data
	data, err = s.sut.SetpointConstraints(s.cdtEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Len(s.T(), data, 1)
	assert.Equal(s.T(), uint(1), data[0].Id)
	assert.Equal(s.T(), 20.0, data[0].MinValue)
	assert.Equal(s.T(), 70.0, data[0].MaxValue)
	assert.Equal(s.T(), 1.0, data[0].StepSize)
}

// Test_WriteSetpoint verifies the functionality of writing a setpoint to a remote entity.
func (s *CaCDTSuite) Test_WriteSetpoint() {
	// Test case: No setpoints available for mock remote entity
	err := s.sut.WriteSetpoint(s.mockRemoteEntity, model.HvacOperationModeTypeTypeOn, 35.0)
	assert.NotNil(s.T(), err)

	// Create a setpoint
	setpoint := model.SetpointDataType{
		SetpointId: util.Ptr(model.SetpointIdType(1)),
		Value: util.Ptr(model.ScaledNumberType{
			Number: util.Ptr(model.NumberType(34)),
			Scale:  util.Ptr(model.ScaleType(0)),
		}),
		ValueMin: util.Ptr(model.ScaledNumberType{
			Number: util.Ptr(model.NumberType(20)),
			Scale:  util.Ptr(model.ScaleType(0)),
		}),
		ValueMax: util.Ptr(model.ScaledNumberType{
			Number: util.Ptr(model.NumberType(70)),
			Scale:  util.Ptr(model.ScaleType(0)),
		}),
		IsSetpointChangeable: util.Ptr(true),
		IsSetpointActive:     util.Ptr(true),
	}

	setpoints := &model.SetpointListDataType{
		SetpointData: []model.SetpointDataType{setpoint},
	}

	// Update the remote feature with the new setpoint data
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.cdtEntity, model.FeatureTypeTypeSetpoint, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSetpointListData, setpoints, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test case: No mapping of operation modes to setpoints available
	err = s.sut.WriteSetpoint(s.cdtEntity, model.HvacOperationModeTypeTypeOn, 35.0)
	assert.NotNil(s.T(), err)

	// Create a mapping of operation modes to setpoints
	s.sut.modes = map[model.HvacOperationModeTypeType]model.SetpointIdType{
		model.HvacOperationModeTypeTypeOn: 1,
	}

	// Test case: Setpoint and operation mode mapping available - the write should succeed
	err = s.sut.WriteSetpoint(s.cdtEntity, model.HvacOperationModeTypeTypeOn, 35.0)
	assert.Nil(s.T(), err)
}
