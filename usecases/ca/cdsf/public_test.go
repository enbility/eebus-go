package cdsf

import (
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CaCDSFSuite) Test_OperationModes() {
	data, err := s.sut.OperationModes(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.OperationModes(s.dhwCircuitEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addHvacData(util.Ptr(true))

	data, err = s.sut.OperationModes(s.dhwCircuitEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(data))
	assert.Equal(s.T(), ucapi.HvacOperationModeTypeAuto, data[0])
}

func (s *CaCDSFSuite) Test_CurrentOperationMode() {
	data, err := s.sut.CurrentOperationMode(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), ucapi.HvacOperationModeType(""), data)

	data, err = s.sut.CurrentOperationMode(s.dhwCircuitEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), ucapi.HvacOperationModeType(""), data)

	s.addHvacData(util.Ptr(true))

	data, err = s.sut.CurrentOperationMode(s.dhwCircuitEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), ucapi.HvacOperationModeTypeEco, data)
}

func (s *CaCDSFSuite) Test_WriteOperationMode() {
	_, err := s.sut.WriteOperationMode(s.mockRemoteEntity, ucapi.HvacOperationModeTypeOn, nil)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteOperationMode(s.dhwCircuitEntity, ucapi.HvacOperationModeTypeOn, nil)
	assert.NotNil(s.T(), err)

	s.addHvacData(util.Ptr(true))

	_, err = s.sut.WriteOperationMode(s.dhwCircuitEntity, ucapi.HvacOperationModeTypeOn, nil)
	assert.Nil(s.T(), err)

	// an unsupported mode cannot be written
	_, err = s.sut.WriteOperationMode(s.dhwCircuitEntity, ucapi.HvacOperationModeType("invalid"), nil)
	assert.NotNil(s.T(), err)
}

func (s *CaCDSFSuite) Test_WriteOperationMode_NotChangeable() {
	s.addHvacData(util.Ptr(false))

	_, err := s.sut.WriteOperationMode(s.dhwCircuitEntity, ucapi.HvacOperationModeTypeOn, nil)
	assert.NotNil(s.T(), err)
}

func (s *CaCDSFSuite) Test_WriteOperationMode_ChangeabilityOmitted() {
	// a device may omit the changeability flag but still accept the write
	s.addHvacData(nil)

	_, err := s.sut.WriteOperationMode(s.dhwCircuitEntity, ucapi.HvacOperationModeTypeOn, nil)
	assert.Nil(s.T(), err)
}

func (s *CaCDSFSuite) Test_WriteOperationMode_UnrelatedMode() {
	s.addHvacData(util.Ptr(true))

	// add a mode that exists globally but is not related to the DHW function
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.dhwCircuitEntity, model.FeatureTypeTypeHvac, model.RoleTypeServer)
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
			{
				OperationModeId:   util.Ptr(model.HvacOperationModeIdType(4)),
				OperationModeType: util.Ptr(model.HvacOperationModeTypeTypeOff),
			},
		},
	}
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeHvacOperationModeDescriptionListData, modeData, nil, nil)
	assert.Nil(s.T(), fErr)

	// the relation only lists modes 1, 2, 3, so the unrelated mode must not be written
	_, err := s.sut.WriteOperationMode(s.dhwCircuitEntity, ucapi.HvacOperationModeTypeOff, nil)
	assert.NotNil(s.T(), err)
}

func (s *CaCDSFSuite) Test_StartStopOneTimeDhw() {
	_, err := s.sut.StartOneTimeDhw(s.mockRemoteEntity, nil)
	assert.NotNil(s.T(), err)

	_, err = s.sut.StartOneTimeDhw(s.dhwCircuitEntity, nil)
	assert.NotNil(s.T(), err)

	s.addHvacData(util.Ptr(true))

	_, err = s.sut.StartOneTimeDhw(s.dhwCircuitEntity, nil)
	assert.NotNil(s.T(), err)

	s.addOverrunData(true)

	_, err = s.sut.StartOneTimeDhw(s.dhwCircuitEntity, nil)
	assert.Nil(s.T(), err)

	_, err = s.sut.StopOneTimeDhw(s.dhwCircuitEntity, nil)
	assert.Nil(s.T(), err)

	// the overrun status marked not changeable cannot be written
	s.addOverrunData(false)

	_, err = s.sut.StartOneTimeDhw(s.dhwCircuitEntity, nil)
	assert.NotNil(s.T(), err)
}

func (s *CaCDSFSuite) Test_WriteOperationMode_AmbiguousSystemFunction() {
	s.addHvacData(util.Ptr(true))

	// two DHW system functions make the id ambiguous, so the write must fail
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.dhwCircuitEntity, model.FeatureTypeTypeHvac, model.RoleTypeServer)
	descData := &model.HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []model.HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId:   util.Ptr(model.HvacSystemFunctionIdType(1)),
				SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeDhw),
			},
			{
				SystemFunctionId:   util.Ptr(model.HvacSystemFunctionIdType(2)),
				SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeDhw),
			},
		},
	}
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err := s.sut.WriteOperationMode(s.dhwCircuitEntity, ucapi.HvacOperationModeTypeOn, nil)
	assert.NotNil(s.T(), err)
}

// helpers

func (s *CaCDSFSuite) addOverrunData(isChangeable bool) {
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
				OverrunId:                 util.Ptr(model.HvacOverrunIdType(1)),
				OverrunStatus:             util.Ptr(model.HvacOverrunStatusTypeInactive),
				IsOverrunStatusChangeable: util.Ptr(isChangeable),
			},
		},
	}
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeHvacOverrunListData, overrunData, nil, nil)
	assert.Nil(s.T(), fErr)
}

func (s *CaCDSFSuite) addHvacData(isChangeable *bool) {
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
				SystemFunctionId:            util.Ptr(model.HvacSystemFunctionIdType(1)),
				CurrentOperationModeId:      util.Ptr(model.HvacOperationModeIdType(3)),
				IsOperationModeIdChangeable: isChangeable,
			},
		},
	}
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeHvacSystemFunctionListData, functionData, nil, nil)
	assert.Nil(s.T(), fErr)
}
