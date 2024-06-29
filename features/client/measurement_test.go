package client

import (
	"testing"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestMeasurementSuite(t *testing.T) {
	suite.Run(t, new(MeasurementSuite))
}

type MeasurementSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	measurement *Measurement
	sentMessage []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*MeasurementSuite)(nil)

func (s *MeasurementSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *MeasurementSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
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
	s.measurement, err = NewMeasurement(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.measurement)

	s.measurement, err = NewMeasurement(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.measurement)
}

func (s *MeasurementSuite) Test_RequestDescriptions() {
	msgCounter, err := s.measurement.RequestDescriptions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msgCounter, err = s.measurement.RequestDescriptions(
		&model.MeasurementDescriptionListDataSelectorsType{},
		&model.MeasurementDescriptionDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
}

func (s *MeasurementSuite) Test_RequestConstraints() {
	msgCounter, err := s.measurement.RequestConstraints(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msgCounter, err = s.measurement.RequestConstraints(
		&model.MeasurementConstraintsListDataSelectorsType{},
		&model.MeasurementConstraintsDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
}

func (s *MeasurementSuite) Test_RequestData() {
	counter, err := s.measurement.RequestData(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.measurement.RequestData(
		&model.MeasurementListDataSelectorsType{},
		&model.MeasurementDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
