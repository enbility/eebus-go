package opev

import (
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *CemOPEVSuite) Test_Public() {
	// The actual tests of the functionality is located in the util package

	_, _, _, err := s.sut.CurrentLimits(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, _, _, err = s.sut.CurrentLimits(s.evEntity)
	assert.NotNil(s.T(), err)

	meas := s.evEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	assert.NotNil(s.T(), meas)

	mData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}
	_, errT := meas.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, mData, nil, nil)
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
