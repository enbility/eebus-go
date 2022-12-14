package features

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestMeasurementSuite(t *testing.T) {
	suite.Run(t, new(MeasurementSuite))
}

type MeasurementSuite struct {
	suite.Suite

	localDevice  *spine.DeviceLocalImpl
	remoteEntity *spine.EntityRemoteImpl

	measurement          *Measurement
	electricalConnection *ElectricalConnection
	sentMessage          []byte
}

var _ spine.SpineDataConnection = (*MeasurementSuite)(nil)

func (s *MeasurementSuite) WriteSpineMessage(message []byte) {
	s.sentMessage = message
}

func (s *MeasurementSuite) BeforeTest(suiteName, testName string) {
	s.localDevice, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeMeasurement,
				functions: []model.FunctionType{
					model.FunctionTypeMeasurementDescriptionListData,
					model.FunctionTypeMeasurementConstraintsListData,
					model.FunctionTypeMeasurementListData,
				},
			},
			{
				featureType: model.FeatureTypeTypeElectricalConnection,
				functions: []model.FunctionType{
					model.FunctionTypeElectricalConnectionDescriptionListData,
					model.FunctionTypeElectricalConnectionParameterDescriptionListData,
					model.FunctionTypeElectricalConnectionPermittedValueSetListData,
				},
			},
		},
	)

	var err error
	s.measurement, err = NewMeasurement(model.RoleTypeServer, model.RoleTypeClient, s.localDevice, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.measurement)

	s.electricalConnection, err = NewElectricalConnection(model.RoleTypeServer, model.RoleTypeClient, s.localDevice, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.electricalConnection)
}

func (s *MeasurementSuite) Test_RequestLimitDescription() {
	err := s.measurement.RequestDescription()
	assert.Nil(s.T(), err)
}

func (s *MeasurementSuite) Test_RequestConstraints() {
	err := s.measurement.RequestConstraints()
	assert.Nil(s.T(), err)
}

func (s *MeasurementSuite) Test_Request() {
	counter, err := s.measurement.Request()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *MeasurementSuite) Test_GetValueForScope() {
	measurement := model.MeasurementTypeTypeCurrent
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACCurrent

	data, mId, err := s.measurement.GetValueForTypeCommodityScope(measurement, commodity, scope)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
	assert.Equal(s.T(), uint(0), uint(mId))

	s.addDescription()

	data, mId, err = s.measurement.GetValueForTypeCommodityScope(measurement, commodity, scope)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)
	assert.Equal(s.T(), uint(0), uint(mId))

	s.addData()

	data, mId, err = s.measurement.GetValueForTypeCommodityScope(measurement, commodity, scope)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), 0.0, data)
	assert.Equal(s.T(), uint(0), uint(mId))
}

func (s *MeasurementSuite) Test_GetValuesPerPhaseForScope() {
	measurement := model.MeasurementTypeTypeCurrent
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACCurrent

	data, _, err := s.measurement.GetValuesPerPhaseForTypeCommodityScope(measurement, commodity, scope, s.electricalConnection)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, _, err = s.measurement.GetValuesPerPhaseForTypeCommodityScope(measurement, commodity, scope, s.electricalConnection)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addElectricalParamDescription()

	data, _, err = s.measurement.GetValuesPerPhaseForTypeCommodityScope(measurement, commodity, scope, s.electricalConnection)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, _, err = s.measurement.GetValuesPerPhaseForTypeCommodityScope(measurement, commodity, scope, s.electricalConnection)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetDescriptionForScope() {
	data, err := s.measurement.GetDescriptionForScope(model.ScopeTypeTypeACCurrent)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.measurement.GetDescriptionForScope(model.ScopeTypeTypeACCurrent)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetSoC() {
	data, err := s.measurement.GetSoC()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	s.addDescription()

	data, err = s.measurement.GetSoC()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	s.addData()

	data, err = s.measurement.GetSoC()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), 0, data)
}

func (s *MeasurementSuite) Test_GetConstraints() {
	data, err := s.measurement.GetConstraints()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addConstraints()

	data, err = s.measurement.GetConstraints()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetValues() {
	data, err := s.measurement.GetValues()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addConstraints()

	s.addDescription()

	data, err = s.measurement.GetValues()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.measurement.GetValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

// helper

func (s *MeasurementSuite) addDescription() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePercentage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeStateOfCharge),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeMeasurementDescriptionListData, fData, nil, nil)
}

func (s *MeasurementSuite) addConstraints() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.MeasurementConstraintsListDataType{
		MeasurementConstraintsData: []model.MeasurementConstraintsDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ValueRangeMin: model.NewScaledNumberType(2),
				ValueRangeMax: model.NewScaledNumberType(16),
				ValueStepSize: model.NewScaledNumberType(0.1),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(0.1),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeMeasurementConstraintsListData, fData, nil, nil)
}

func (s *MeasurementSuite) addElectricalParamDescription() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(2)))
	fData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeAbc),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, fData, nil, nil)
}

func (s *MeasurementSuite) addData() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))

	t := time.Now()
	fData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(9),
				Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(t),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(9),
				Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(t),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeMeasurementListData, fData, nil, nil)
}
