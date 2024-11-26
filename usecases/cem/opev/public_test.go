package opev

import (
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *CemOPEVSuite) Test_Public() {
	// The actual tests of the functionality is located in the internal package

	_, _, _, err := s.sut.CurrentLimits(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, _, _, err = s.sut.CurrentLimits(s.evEntity)
	assert.NotNil(s.T(), err)

	lc := s.evEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	assert.NotNil(s.T(), lc)

	meas := s.evEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	assert.NotNil(s.T(), meas)

	lData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
		},
	}

	_,errT := lc.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, lData, nil, nil)
	assert.Nil(s.T(), errT)

	_, _, _, err = s.sut.CurrentLimits(s.evEntity)
	assert.NotNil(s.T(), err)

	lData = &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
				LimitType: util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
				Unit: util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType: util.Ptr(model.ScopeTypeTypeOverloadProtection),
			},
		},
	}

	_,errT = lc.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, lData, nil, nil)
	assert.Nil(s.T(), errT)

	_, _, _, err = s.sut.CurrentLimits(s.evEntity)
	assert.NotNil(s.T(), err)

	lData = &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
			},
		},
	}

	_,errT = lc.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, lData, &model.FilterType{}, nil)
	assert.Nil(s.T(), errT)

	_, _, _, err = s.sut.CurrentLimits(s.evEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.LoadControlLimits(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.LoadControlLimits(s.evEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteLoadControlLimits(s.mockRemoteEntity, []ucapi.LoadLimitsPhase{}, nil)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteLoadControlLimits(s.evEntity, []ucapi.LoadLimitsPhase{}, nil)
	assert.NotNil(s.T(), err)

	s.sut.StopHeartbeat()
	s.sut.StartHeartbeat()

	err = s.sut.SetOperatingState(true)
	assert.Nil(s.T(), err)
}
